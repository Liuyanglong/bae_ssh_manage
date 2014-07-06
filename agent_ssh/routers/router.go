package routers

import (
	"github.com/astaxie/beego"
	"ssh_proxy_manage/agent_ssh/controllers"
)

func init() {
	beego.Router("/updateContainer", &controllers.SshAgentController{}, "get:UpdateContainerRull")
	beego.Router("/deleteContainer", &controllers.SshAgentController{}, "get:DeleteContainerRull")
}
