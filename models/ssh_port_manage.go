package models

import (
	"database/sql"
	"ssh_proxy_manage/logs"
	"time"
)

type SshPort struct {
	PublicPort      int    `json:"port,int"`
	ProxyHost       string `json:"host"`
	ContainerNumber int    `json:"number,int"`
	IsAvail         bool   `json:"avail,string"`
}

type SshPortManage struct {
	DbConn
}

func (u *SshPortManage) TableName() string {
	return "ssh_port_manage"
}

func (u *SshPortManage) Insert(conn *sql.DB, sshPortOb SshPort, logid int64) error {
	logs.Normal("insert sshPortOb:", sshPortOb, "logid:", logid)
	tableName := u.TableName()
	sqli := "INSERT INTO " + tableName + " (publicPort,containerNumber,proxyHost,isAvail,addTime) VALUES (?,?,?,?,?)"
	stmt, sterr := conn.Prepare(sqli)
	if sterr != nil {
		logs.Error("comm prepare error:", sterr, "logid:", logid)
		return sterr
	}
	if _, err := stmt.Exec(sshPortOb.PublicPort, sshPortOb.ContainerNumber, sshPortOb.ProxyHost, sshPortOb.IsAvail, time.Now().Format("2006-01-02 15:04:05")); err != nil {
		logs.Error("sql exec error!", err, sqli, sshPortOb, "logid:", logid)
		return err
	}
	return nil
}

func (u *SshPortManage) Query(conn *sql.DB, pubport string, logid int64) (SshPort, error) {
	tableName := u.TableName()
	sql := "select publicPort,proxyHost, containerNumber,isAvail from " + tableName + " where publicPort = " + pubport
	logs.Normal("query sql:", sql, "logid:", logid)
	rows, err := conn.Query(sql)
	if err != nil {
		logs.Error("query error!", err, sql, "logid:", logid)
		return SshPort{}, err
	}
	sshPortOb := SshPort{}
	for rows.Next() {
		if err = rows.Scan(&sshPortOb.PublicPort, &sshPortOb.ProxyHost, &sshPortOb.ContainerNumber, &sshPortOb.IsAvail); err != nil {
			logs.Error("rows scan error:", err, "logid:", logid)
			return SshPort{}, err
		}
	}

	return sshPortOb, nil
}

func (u *SshPortManage) QueryAll(conn *sql.DB, logid int64) ([]SshPort, error) {
	tableName := u.TableName()
	sql := "select publicPort,proxyHost, containerNumber,isAvail from " + tableName
	logs.Normal("query sql:", sql, "logid:", logid)
	rows, err := conn.Query(sql)
	if err != nil {
		logs.Error("query error!", err, sql, "logid:", logid)
		return []SshPort{}, err
	}
	sshPortObs := make([]SshPort, 0)
	for rows.Next() {
		sshPortOb := SshPort{}
		if err = rows.Scan(&sshPortOb.PublicPort, &sshPortOb.ProxyHost, &sshPortOb.ContainerNumber, &sshPortOb.IsAvail); err != nil {
			logs.Error("rows scan error:", err, "logid:", logid)
			return []SshPort{}, err
		}
		sshPortObs = append(sshPortObs, sshPortOb)
	}

	return sshPortObs, nil
}

func (u *SshPortManage) Delete(conn *sql.DB, pubport string, logid int64) error {
	tableName := u.TableName()
	sqld := "delete from " + tableName + " where publicPort = " + pubport
	logs.Normal("delete sql:", sqld, "logid:", logid)
	if _, err := conn.Exec(sqld); err != nil {
		logs.Error("conn delete sql exec error:", err, "logid:", logid)
		return err
	}
	return nil
}

func (u *SshPortManage) Update(conn *sql.DB, pubport, cnum string, isplus bool, logid int64) error {
	tableName := u.TableName()
	sqld := ""
	if isplus == true {
		sqld = "update " + tableName + " set containerNumber=containerNumber+" + cnum + " where publicPort = " + pubport
	} else {
		sqld = "update " + tableName + " set containerNumber=containerNumber-" + cnum + " where publicPort = " + pubport

	}
	logs.Normal("update sql is:", sqld, "logid:", logid)
	if _, err := conn.Exec(sqld); err != nil {
		logs.Error("update sql exec error:", err, sqld, "logid:", logid)
		return err
	}
	return nil
}

func (u *SshPortManage) GetBestUsePort(conn *sql.DB, logid int64) (SshPort, error) {
	var sshPortOb SshPort
	tableName := u.TableName()
	sqls := "select publicPort,proxyHost,containerNumber,isAvail from " + tableName + " where isAvail = 1 order by containerNumber asc limit 1"
	logs.Normal("get best use port sql:", sqls, "logid:", logid)
	rows, err := conn.Query(sqls)
	if err != nil {
		logs.Error("sql query error:", err, "logid:", logid)
		return SshPort{}, err
	}

	for rows.Next() {
		if err = rows.Scan(&sshPortOb.PublicPort, &sshPortOb.ProxyHost, &sshPortOb.ContainerNumber, &sshPortOb.IsAvail); err != nil {
			logs.Error("rows scan error:", err, "logid:", logid)
			return SshPort{}, err
		}
		break
	}

	return sshPortOb, nil
}
