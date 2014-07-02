package logs

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"ssh_proxy_manage/mail"
	"strings"
)

var normal_log *logs.BeeLogger
var error_log *logs.BeeLogger

func init() {
	normal_log = logs.NewLogger(10000)
	normal_log.EnableFuncCallDepth(true)
	access_log_file := beego.AppConfig.String("access_log")
	normal_log.SetLogger("file", `{"filename":"`+access_log_file+`"}`)
	error_log = logs.NewLogger(10000)
	error_log.EnableFuncCallDepth(true)
	error_log_file := beego.AppConfig.String("error_log")
	error_log.SetLogger("file", `{"filename":"`+error_log_file+`"}`)
}

func Normal(v ...interface{}) {
	normal_log.Info("%v", v)
}
func Error(v ...interface{}) {
	error_log.Error("%v", v)
	err := mail.SendMail("ssh service error!", fmt.Sprintf("[MAIL] "+"%v", v...))
	if err != nil {
		fmt.Println("send mail error!", err)
	}
}
func generateFmtStr(n int) string {
	return strings.Repeat("%v", n)
}
