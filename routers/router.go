package routers

import (
	"github.com/astaxie/beego"
	"ssh_proxy_manage/controllers"
)

func init() {
	beego.RESTRouter("/object", &controllers.ObjectController{})
	beego.RESTRouter("/sshPort", &controllers.SshPortController{})
	beego.RESTRouter("/sshProxyServer", &controllers.SshProxyServerController{})
	beego.Router("/sshKeys", &controllers.SshKeysController{})
	beego.Router("/sshKeys/:objectId/:keyname", &controllers.SshKeysController{})
	beego.Router("/sshRules", &controllers.SshRulesController{})
	beego.Router("/sshRules/:uid/:container", &controllers.SshRulesController{})
}
