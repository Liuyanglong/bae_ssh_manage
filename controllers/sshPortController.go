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

func (this *SshPortController) Post() {
	var sshPortOb models.SshPort
	
	//获取post的body体中的数据
	requestBody := string(this.Ctx.Input.RequestBody)
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &sshPortOb)
	if err != nil {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"param error"}`))
		this.StopRun()
	}
	//获取参数中的logid数值
	logid, err := models.GetLogId(this.Ctx.Input.RequestBody)
	if err != nil {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"logid param error"}`))
		this.StopRun()
	}
	logs.Normal("post data:", requestBody, "logid:", logid)

	//初始化db句柄
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

	//将数据插入数据库
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

func (this *SshPortController) Get() {
	portNum := this.Ctx.Input.Params[":objectId"]
	if portNum == "" {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"param error"}`))
		this.StopRun()
	}
	dbconn, err := models.InitDbConn(0)
	if err != nil {
		logs.Error("db init error:", err, 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"init db conn error"}`))
		this.StopRun()
	}
	defer dbconn.Close()

	//获取数据库中的对应的数据
	sshPortM := new(models.SshPortManage)
	sshPortOb, err := sshPortM.Query(dbconn, portNum, 0)
	if err != nil {
		logs.Error("ssh port query error:", err, 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}
	//返回json数据
	this.Data["json"] = sshPortOb
	this.ServeJson()
}

func (this *SshPortController) Delete() {
	portNum := this.Ctx.Input.Params[":objectId"]
	if portNum == "" {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"param error"}`))
		this.StopRun()
	}
	
	//从参数中获取logid
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

/*
**更新这个port中对应映射的container的数量
* @param port 对外的ssh port
* @param num  变化的增量值
* @param isplus  当为true时，num为增量；为false时，num为减量
*/
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

//获取当前最佳可使用的对外port
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
