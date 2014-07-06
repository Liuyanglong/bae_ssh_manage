package routers

import (
	"github.com/astaxie/beego"
	"ssh_proxy_manage/controllers"
)

func init() {
	beego.Router("/sshPort/getAvailPort", &controllers.SshPortController{}, "get:GetAvailPort")
	beego.RESTRouter("/sshPort", &controllers.SshPortController{})
	beego.RESTRouter("/sshProxyServer", &controllers.SshProxyServerController{})
	beego.Router("/sshKeys/:objectId/:keyname", &controllers.SshKeysController{}, "get:GetMsgByKeyname")
	beego.RESTRouter("/sshKeys", &controllers.SshKeysController{})
	beego.Router("/sshRules/:uid/:container", &controllers.SshRulesController{}, "get:GetByContainer;delete:DeleteByContainer")
	beego.Router("/sshRules/:uid", &controllers.SshRulesController{}, "get:GetByUid;delete:DeleteByUid;put:Put;post:Post")
}
