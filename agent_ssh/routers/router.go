package routers

import (
	"agent_ssh/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.RESTRouter("/sshAgent", &controllers.SshAgentController{})
}
