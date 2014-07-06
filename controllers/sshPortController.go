package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"ssh_proxy_manage/logs"
	"ssh_proxy_manage/models"
)

type SshPortController struct {
	beego.Controller
}

//curl -i -X POST 'localhost:9090/sshPort' -d '{"port":2300,"host":"testinfo2","number":10,"avail":"true"}'
func (this *SshPortController) Post() {
	var sshPortOb models.SshPort
	requestBody := string(this.Ctx.Input.RequestBody)
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &sshPortOb)
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
	logs.Normal("post data:", requestBody, "logid:", logid)

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

	sshPortM := new(models.SshPortManage)
	err = sshPortM.Insert(dbconn, sshPortOb, logid)
	if err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("ssh port insert error:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":` + err.Error() + `}`))
		this.StopRun()
	}
	dbconn.Exec("COMMIT")
	logs.Normal("COMMIT", "logid:", logid)
	logs.Normal("ssh post OK!", "logid:", logid)
	this.Ctx.Output.Body([]byte(`{"result":0}`))
}

//curl -i -X GET 'localhost:9090/sshPort[/2200]'
func (this *SshPortController) Get() {
	portNum := this.Ctx.Input.Params[":objectId"]
	dbconn, err := models.InitDbConn(0)
	if err != nil {
		logs.Error("db init error:", err, 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"init db conn error"}`))
		this.StopRun()
	}
	defer dbconn.Close()

	sshPortM := new(models.SshPortManage)
	//当没有传入具体的port值时，输出全部sshport对象
	if portNum == "" {
		sshPortObs, _ := sshPortM.QueryAll(dbconn, 0)
		this.Data["json"] = sshPortObs
		this.ServeJson()
		this.StopRun()
	}

	//输出portNum对应的sshPort信息
	sshPortOb, err := sshPortM.Query(dbconn, portNum, 0)
	if err != nil {
		logs.Error("ssh port query error:", err, 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}
	this.Data["json"] = sshPortOb
	this.ServeJson()
}

//curl -i -X DELETE 'localhost:9090/sshPort/2200?logid=24232323'
func (this *SshPortController) Delete() {
	portNum := this.Ctx.Input.Params[":objectId"]
	if portNum == "" {
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

	sshPortM := new(models.SshPortManage)
	err = sshPortM.Delete(dbconn, portNum, logid)
	if err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("ssh port delete err:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}
	dbconn.Exec("COMMIT")
	logs.Normal("COMMIT", "logid:", logid)
	logs.Normal("ssh port delete OK", "logid:", logid)
	this.Ctx.Output.Body([]byte(`{"result":0}`))
	this.StopRun()
}

//curl localhost:9090/sshPort/getAvailPort
func (this *SshPortController) GetAvailPort() {
	sshport, err := this.GetBestUsePort(0)
	if err != nil {
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1}`))
		this.StopRun()
	}
	this.Data["json"] = sshport
	this.ServeJson()
}

func (this *SshPortController) UpdateNumber(port, num string, isplus bool, logid int64) error {
	dbconn, err := models.InitDbConn(logid)
	if err != nil {
		logs.Error("init db conn err:", err, "logid:", logid)
		return err
	}
	defer dbconn.Close()

	sshPortM := new(models.SshPortManage)
	err = sshPortM.Update(dbconn, port, num, isplus, logid)
	if err != nil {
		logs.Error("ssh port update error:", err, "logid:", logid)
		return err
	}
	logs.Normal("ssh updatenumber OK", logid)
	return nil
}

func (this *SshPortController) GetBestUsePort(logid int64) (models.SshPort, error) {
	dbconn, err := models.InitDbConn(logid)
	if err != nil {
		logs.Error("init db conn err:", err, "logid:", logid)
		return models.SshPort{}, err
	}
	defer dbconn.Close()

	sshPortM := new(models.SshPortManage)
	portOb, err := sshPortM.GetBestUsePort(dbconn, logid)
	if err != nil {
		logs.Error("ssh port get best userport error:", err, "logid:", logid)
	}
	logs.Normal("get best use port ok!", portOb, "logid:", logid)
	return portOb, err
}
