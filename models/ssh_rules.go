package models

import (
	"database/sql"
	"errors"
	"ssh_proxy_manage/logs"
	"strconv"
	"time"
)

type SshRule struct {
	ContainerName string `json:"name"`
	PublicPort    int    `json:"port"`
	ProxyHost     string `json:"host"`
	UiIpPort      string `json:"rule"`
	Uid           int
}

type SshRuleManage struct {
	DbConn
}

func (u *SshRuleManage) TableName() string {
	return "ssh_proxy_rules"
}

func (u *SshRuleManage) Insert(conn *sql.DB, sshRuleOb SshRule, logid int64) error {
	tableName := u.TableName()
	sqli := "INSERT INTO " + tableName + " (containerName,sshPort,proxyHost,uiIpPort,uid,addTime) VALUES (?,?,?,?,?,?)"
	logs.Normal("insert sql:", sqli, sshRuleOb, "logid:", logid)
	stmt, sterr := conn.Prepare(sqli)
	if sterr != nil {
		logs.Error("conn prepare error:", sterr, sqli, "logid:", logid)
		return sterr
	}
	if _, err := stmt.Exec(sshRuleOb.ContainerName, sshRuleOb.PublicPort, sshRuleOb.ProxyHost, sshRuleOb.UiIpPort, sshRuleOb.Uid, time.Now().Format("2006-01-02 15:04:05")); err != nil {
		logs.Error("sql exec error:", err, sqli, "logid:", logid)
		return err
	}
	return nil
}

func (u *SshRuleManage) Query(conn *sql.DB, uid, cname string, logid int64) (SshRule, error) {
	tableName := u.TableName()
	sql := "select containerName,sshPort,proxyHost,uiIpPort,uid from " + tableName + " where uid = " + uid + " and containerName = '" + cname + "'"
	logs.Normal("sql query:", sql, "logid:", logid)
	rows, err := conn.Query(sql)
	if err != nil {
		logs.Error("conn query error:", err, sql, "logid:", logid)
		return SshRule{}, err
	}
	sshRuleOb := SshRule{}
	for rows.Next() {
		if err = rows.Scan(&sshRuleOb.ContainerName, &sshRuleOb.PublicPort, &sshRuleOb.ProxyHost, &sshRuleOb.UiIpPort, &sshRuleOb.Uid); err != nil {
			logs.Error("row scan error:", err, "logid:", logid)
			return SshRule{}, err
		}
	}

	return sshRuleOb, nil
}

func (u *SshRuleManage) QueryByUid(conn *sql.DB, uid string, logid int64) ([]SshRule, error) {
	tableName := u.TableName()
	rulelist := make([]SshRule, 0)
	sql := "select containerName,sshPort,proxyHost,uiIpPort,uid from " + tableName + " where uid= " + uid
	logs.Normal("sql query by uid:", sql, "logid:", logid)
	rows, err := conn.Query(sql)
	if err != nil {
		logs.Error("conn query error:", err, sql, "logid:", logid)
		return rulelist, err
	}
	for rows.Next() {
		sshRuleOb := SshRule{}
		if err = rows.Scan(&sshRuleOb.ContainerName, &sshRuleOb.PublicPort, &sshRuleOb.ProxyHost, &sshRuleOb.UiIpPort, &sshRuleOb.Uid); err == nil {
			logs.Error("rows scan error:", err, "logid:", logid)
			rulelist = append(rulelist, sshRuleOb)
		}
	}

	return rulelist, nil
}

func (u *SshRuleManage) Delete(conn *sql.DB, uid, cname string, logid int64) error {
	tableName := u.TableName()
	sqld := "delete from " + tableName + " where uid= " + uid + " and containerName = '" + cname + "'"
	logs.Normal("delete sql:", sqld, "logid:", logid)
	if _, err := conn.Exec(sqld); err != nil {
		logs.Error("conn exec error:", err, "logid:", logid)
		return err
	}
	return nil
}

func (u *SshRuleManage) DeleteByUid(conn *sql.DB, uid string, logid int64) error {
	tableName := u.TableName()
	sqld := "delete from " + tableName + " where uid = " + uid
	logs.Normal("delete by uid sql:", sqld, "logid:", logid)
	if _, err := conn.Exec(sqld); err != nil {
		logs.Error("conn exec error:", err, "logid:", logid)
		return err
	}
	return nil
}

func (u *SshRuleManage) Update(conn *sql.DB, ruleOb SshRule, logid int64) error {
	tableName := u.TableName()
	if ruleOb.ContainerName == "" {
		return errors.New("ssh rule object is not good! miss primary key!")
	}
	param := make([]interface{}, 0)
	param = append(param, ruleOb.ContainerName)
	sqlu := "update " + tableName + " set containerName = '" + ruleOb.ContainerName + "'"
	if ruleOb.ProxyHost != "" {
		sqlu += " ,proxyHost='" + ruleOb.ProxyHost + "'"
		param = append(param, ruleOb.ProxyHost)
	}
	if ruleOb.PublicPort != 0 {
		portstr := strconv.Itoa(ruleOb.PublicPort)
		sqlu += " ,sshPort=" + portstr
		param = append(param, ruleOb.PublicPort)
	}
	if ruleOb.UiIpPort != "" {
		sqlu += " ,uiIpPort='" + ruleOb.UiIpPort + "'"
		param = append(param, ruleOb.UiIpPort)
	}

	uidstr := strconv.Itoa(ruleOb.Uid)
	sqlu += " where containerName = '" + ruleOb.ContainerName + "' and uid = " + uidstr
	logs.Normal("sql update:", sqlu, "logid:", logid)
	param = append(param, ruleOb.ContainerName)
	_, err := conn.Exec(sqlu)
	if err != nil {
		logs.Error("conn exec error:", err, "logid:", logid)
	}
	return err
}
