package models

import (
	"database/sql"
	"errors"
	"ssh_proxy_manage/logs"
	"strconv"
	"time"
)

type SshKeys struct {
	Uid       int
	KeyName   string `json:"keyname"`
	PublicKey string `json:"publickey"`
}

type SshKeysManage struct {
	DbConn
}

func (u *SshKeysManage) TableName() string {
	return "ssh_user_keys"
}

func (u *SshKeysManage) Insert(conn *sql.DB, sshKeysOb SshKeys, logid int64) error {
	tableName := u.TableName()
	sqli := "INSERT INTO " + tableName + " (uid,keyName,publicKey,addTime) VALUES (?,?,?,?)"
	stmt, sterr := conn.Prepare(sqli)
	if sterr != nil {
		logs.Error("sql prepare error:", sterr, "logid:", logid)
		return sterr
	}
	if _, err := stmt.Exec(sshKeysOb.Uid, sshKeysOb.KeyName, sshKeysOb.PublicKey, time.Now().Format("2006-01-02 15:04:05")); err != nil {
		logs.Error("sql exec error:", sterr, " the sql is:", sqli, sshKeysOb, "logid:", logid)
		return err
	}
	logs.Normal("sql insert OK!", sqli, sshKeysOb, "logid:", logid)
	return nil
}

func (u *SshKeysManage) Query(conn *sql.DB, uid string, keyname string, logid int64) (SshKeys, error) {
	tableName := u.TableName()
	sql := "select uid,keyName,publicKey from " + tableName + " where uid = " + uid + " and keyname = '" + keyname + "'"
	rows, err := conn.Query(sql)
	if err != nil {
		logs.Error("sql query error!", err, sql, "logid:", logid)
		return SshKeys{}, err
	}
	sshKeysOb := SshKeys{}
	for rows.Next() {
		if err = rows.Scan(&sshKeysOb.Uid, &sshKeysOb.KeyName, &sshKeysOb.PublicKey); err != nil {
			logs.Error("rows scan error!", err, sql, "logid:", logid)
			return SshKeys{}, err
		}
	}
	logs.Normal("query OK!", sql, sshKeysOb)
	return sshKeysOb, nil
}

func (u *SshKeysManage) GetAll(conn *sql.DB, uid string, logid int64) (map[string]string, error) {
	tableName := u.TableName()
	keyMaps := make(map[string]string)
	sql := "select keyName,publicKey from " + tableName + " where uid = " + uid
	rows, err := conn.Query(sql)
	if err != nil {
		logs.Error("sql query error!", err, sql, "logid:", logid)
		return keyMaps, err
	}

	var keyname, publickey string
	for rows.Next() {
		if err = rows.Scan(&keyname, &publickey); err == nil {
			keyMaps[keyname] = publickey
		}
	}
	logs.Normal("getAll OK!", sql, keyMaps, "logid:", logid)
	return keyMaps, nil
}

func (u *SshKeysManage) Delete(conn *sql.DB, uid string, logid int64) error {
	tableName := u.TableName()
	sqld := "delete from " + tableName + " where uid = " + uid
	if _, err := conn.Exec(sqld); err != nil {
		logs.Error("sql delete error!", err, sqld, "logid:", logid)
		return err
	}
	logs.Normal("sql delete OK!", sqld, "logid:", logid)
	return nil
}

func (u *SshKeysManage) DeleteByKey(conn *sql.DB, uid string, key string, logid int64) error {
	tableName := u.TableName()
	sqld := "delete from " + tableName + " where uid = " + uid + " and keyName='" + key + "'"
	if _, err := conn.Exec(sqld); err != nil {
		logs.Error("sql deleteByKey error!", err, sqld, "logid:", logid)
		return err
	}
	logs.Normal("sql deleteByKey OK!", sqld, "logid:", logid)
	return nil
}

func (u *SshKeysManage) Update(conn *sql.DB, keysOb SshKeys, logid int64) error {
	tableName := u.TableName()
	if keysOb.Uid == 0 {
		return errors.New("ssh keys object is not good! miss primary key!")
	}
	param := make([]interface{}, 0)
	param = append(param, keysOb.Uid)
	sqlu := "update " + tableName + " set uid = " + strconv.Itoa(keysOb.Uid)
	if keysOb.KeyName != "" {
		sqlu += " ,keyName='" + keysOb.KeyName + "'"
		param = append(param, keysOb.KeyName)
	}
	if keysOb.PublicKey != "" {
		sqlu += " ,publicKey='" + keysOb.PublicKey + "'"
		param = append(param, keysOb.PublicKey)
	}
	sqlu += " where uid = " + strconv.Itoa(keysOb.Uid)
	logs.Normal("sql update sql is:", sqlu, "logid:", logid)
	param = append(param, keysOb.Uid)
	_, err := conn.Exec(sqlu)
	return err
}
