package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"ssh_proxy_manage/logs"
	"ssh_proxy_manage/models"
)

type SshProxyServerController struct {
	beego.Controller
}

func (this *SshProxyServerController) Post() {
	var sshProxyOb models.SshProxyMachine
	requestBody := this.Ctx.Input.RequestBody

	err := json.Unmarshal(requestBody, &sshProxyOb)
	if err != nil {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"param error"}`))
		this.StopRun()
	}
	logid, err := models.GetLogId(this.Ctx.Input.RequestBody)
	if err != nil {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"logid param error"}`))
		this.StopRun()
	}
	logs.Normal("post data:", string(requestBody), "logid:", logid)
	dbconn, err := models.InitDbConn(logid)
	if err != nil {
		logs.Error("init db conn error:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"init db conn error"}`))
		this.StopRun()
	}
	defer dbconn.Close()

	dbconn.Exec("START TRANSACTION")
	logs.Normal("start transaction", "logid:", logid)

	sshProxyM := new(models.SshProxyManage)
	err = sshProxyM.Insert(dbconn, sshProxyOb, logid)
	if err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("ssh proxy insert error:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":` + err.Error() + `}`))
		this.StopRun()
	}
	dbconn.Exec("COMMIT")
	logs.Normal("COMMIT", "logid:", logid)
	logs.Normal("post ok!", "logid:", logid)
	this.Ctx.Output.Body([]byte(`{"result":0}`))
}

func (this *SshProxyServerController) Get() {
	host := this.Ctx.Input.Params[":objectId"]
	if host == "" {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"param error"}`))
		this.StopRun()
	}

	dbconn, err := models.InitDbConn(0)
	if err != nil {
		logs.Error("init db conn error:", err, "logid:", 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"init db conn error"}`))
		this.StopRun()
	}
	defer dbconn.Close()

	sshProxyM := new(models.SshProxyManage)
	sshProxyOb, err := sshProxyM.Query(dbconn, host, 0)
	if err != nil {
		logs.Error("ssh proxy query error:", err, "logid:", 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}
	logs.Normal("ssh post ok!", sshProxyOb, "logid:", 0)
	this.Data["json"] = sshProxyOb
	this.ServeJson()
}

func (this *SshProxyServerController) Delete() {
	host := this.Ctx.Input.Params[":objectId"]
	if host == "" {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"param error"}`))
		this.StopRun()
	}
	logid, _ := this.GetInt("logid")
	if logid == 0 {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"logid param error"}`))
		this.StopRun()
	}
	logs.Normal("delete param:", "host:", host, "logid:", logid)
	dbconn, err := models.InitDbConn(logid)
	if err != nil {
		logs.Error("init db conn error:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"init db conn error"}`))
		this.StopRun()
	}
	defer dbconn.Close()

	dbconn.Exec("START TRANSACTION")
	logs.Normal("start transaction", "logid:", logid)

	sshProxyM := new(models.SshProxyManage)
	err = sshProxyM.Delete(dbconn, host, logid)
	if err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("ssh proxy delete err:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}
	dbconn.Exec("COMMIT")
	logs.Normal("COMMIT", "logid:", logid)
	logs.Normal("ssh proxy delete OK!", "logid:", logid)
	this.Ctx.Output.Body([]byte(`{"result":0}`))
	this.StopRun()
}