package main

import (
	"github.com/astaxie/beego"
	"ssh_proxy_manage/logs"
	_ "ssh_proxy_manage/routers"
)

func main() {
	logs.Normal("It is a New Begin!")
	go checkProxyHostAvail()
	beego.Run()
}

//检查proxy server是否可用
func checkProxyHostAvail() {

}
