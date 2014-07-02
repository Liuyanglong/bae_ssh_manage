package models

import (
	"database/sql"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"ssh_proxy_manage/logs"
)

//初始化数据库
func InitDbConn(logid int64) (*sql.DB, error) {
	mysql_conn := beego.AppConfig.String("mysql_conn")
	dbconn, err := sql.Open("mysql", mysql_conn)
	if err != nil {
		logs.Error("db connect Error!", err, "logid:", logid)
	}
	return dbconn, err
}

type DbConn struct {
}

func (u *DbConn) TableName() string {
	return ""
}
