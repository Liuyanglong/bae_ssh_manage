package main

import (
	"github.com/astaxie/beego"
	_ "ssh_proxy_manage/agent_ssh/routers"
)

func main() {
	beego.Run()
}
