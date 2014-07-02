package models

import (
	"database/sql"
	"ssh_proxy_manage/logs"
	"time"
)

type SshProxyMachine struct {
	MachineHost string `json:"host"`
	UserNumber  int    `json:"number,int"`
	IsAvail     bool   `json:"avail,string"`
}

type SshProxyManage struct {
	DbConn
}

func (u *SshProxyManage) TableName() string {
	return "ssh_proxy_machine"
}

func (u *SshProxyManage) Insert(conn *sql.DB, sshProxyOb SshProxyMachine, logid int64) error {
	tableName := u.TableName()
	sqli := "INSERT INTO " + tableName + " (machineHost,userNumber,isAvail,addTime) VALUES (?,?,?,?)"
	logs.Normal("insert:", sqli, sshProxyOb, "logid:", logid)
	stmt, sterr := conn.Prepare(sqli)
	if sterr != nil {
		logs.Error("conn prepare error:", sterr, "logid:", logid)
		return sterr
	}
	if _, err := stmt.Exec(sshProxyOb.MachineHost, sshProxyOb.UserNumber, sshProxyOb.IsAvail, time.Now().Format("2006-01-02 15:04:05")); err != nil {
		logs.Error("stmt exec error:", err, "logid:", logid)
		return err
	}
	return nil
}

func (u *SshProxyManage) Query(conn *sql.DB, host string, logid int64) (SshProxyMachine, error) {
	tableName := u.TableName()
	sql := "select machineHost,userNumber,isAvail from " + tableName + " where machineHost = '" + host + "'"
	logs.Normal("query sql:", sql, "logid:", logid)
	rows, err := conn.Query(sql)
	if err != nil {
		logs.Error("query error:", err, sql, "logid:", logid)
		return SshProxyMachine{}, err
	}
	sshProxyOb := SshProxyMachine{}
	for rows.Next() {
		if err = rows.Scan(&sshProxyOb.MachineHost, &sshProxyOb.UserNumber, &sshProxyOb.IsAvail); err != nil {
			logs.Error("row scan error:", err, "logid:", logid)
			return SshProxyMachine{}, err
		}
	}

	return sshProxyOb, nil
}

func (u *SshProxyManage) Delete(conn *sql.DB, host string, logid int64) error {
	tableName := u.TableName()
	sqld := "delete from " + tableName + " where machineHost = '" + host + "'"
	logs.Normal("sql delete:", sqld, "logid:", logid)
	if _, err := conn.Exec(sqld); err != nil {
		logs.Error("conn exec error:", err, "logid:", logid)
		return err
	}
	return nil
}
