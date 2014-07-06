package controllers

import (
	"database/sql"
	"encoding/json"
	"github.com/astaxie/beego"
	"ssh_proxy_manage/logs"
	"ssh_proxy_manage/models"
	"strconv"
)

type SshKeysController struct {
	beego.Controller
}

//curl localhost:9090/sshKeys -X POST -i -d '{"uid":4444444,"keyname":"testliuyanglong","publickey":"dasd"}''
func (this *SshKeysController) Post() {
	var sshKeysOb models.SshKeys
	requestBody := string(this.Ctx.Input.RequestBody)
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &sshKeysOb)
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
		logs.Error("init db conn error!", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"init db conn error"}`))
		this.StopRun()
	}
	defer dbconn.Close()

	dbconn.Exec("START TRANSACTION")
	logs.Normal("start transaction", "logid:", logid)
	sshKeysM := new(models.SshKeysManage)
	err = sshKeysM.Insert(dbconn, sshKeysOb, logid)
	if err != nil {
		logs.Error("sshKeysM.Insert error!", err, "logid:", logid)
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":` + err.Error() + `}`))
		this.StopRun()
	}

	//reload ssh public key
	err = this.reloadSshPublicKeys(dbconn, strconv.Itoa(sshKeysOb.Uid), logid)
	if err != nil {
		logs.Error("reloadSshPublicKeys error!", err, sshKeysOb, "logid:", logid)
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":` + err.Error() + `}`))
		this.StopRun()
	}
	dbconn.Exec("COMMIT")
	logs.Normal("COMMIT", "logid:", logid)
	logs.Normal("post OK!", requestBody)
	this.Ctx.Output.Body([]byte(`{"result":0}`))
}

//curl localhost:9090/sshKeys/123456
func (this *SshKeysController) Get() {
	uid := this.Ctx.Input.Params[":objectId"]
	if uid == "" {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"param error"}`))
		this.StopRun()
	}

	dbconn, err := models.InitDbConn(0)
	if err != nil {
		logs.Error("db init error!", err, 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"init db conn error"}`))
		this.StopRun()
	}
	defer dbconn.Close()

	sshKeysM := new(models.SshKeysManage)
	sshKeysOb, err := sshKeysM.GetAll(dbconn, uid, 0)
	if err != nil {
		logs.Error("sshkeysM getall error:", err, 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}
	this.Data["json"] = sshKeysOb
	this.ServeJson()
}

//curl -i -X GET 'localhost:9090/sshKeys/123456/testccc?logid=4223'
func (this *SshKeysController) GetMsgByKeyname() {
	uid := this.Ctx.Input.Params[":objectId"]
	keyname := this.Ctx.Input.Params[":keyname"]
	if uid == "" || keyname == "" {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"param error"}`))
		this.StopRun()
	}

	dbconn, err := models.InitDbConn(0)
	if err != nil {
		logs.Error("db init error!", err, 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"init db conn error"}`))
		this.StopRun()
	}
	defer dbconn.Close()

	sshKeysM := new(models.SshKeysManage)
	sshKeysOb, err := sshKeysM.Query(dbconn, uid, keyname, 0)
	if err != nil {
		logs.Error("sshkeysM getall error:", err, 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}
	this.Data["json"] = sshKeysOb
	this.ServeJson()
}

//curl -i -X DELETE 'localhost:9090/sshKeys/4444444?logid=4223'
func (this *SshKeysController) Delete() {
	uid := this.Ctx.Input.Params[":objectId"]
	keyname := this.Ctx.Input.Params[":keyname"]
	if uid == "" {
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
	logs.Normal("delete param:", "uid:", uid, "keyname:", keyname, "logid:", logid)
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
	sshKeysM := new(models.SshKeysManage)
	if keyname != "" {
		err = sshKeysM.DeleteByKey(dbconn, uid, keyname, logid)
		if err != nil {
			dbconn.Exec("ROLLBACK")
			logs.Normal("ROLLBACK", "logid:", logid)
			logs.Error("delete by key error:", err, "logid:", logid)
			this.Ctx.Output.SetStatus(500)
			this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
			this.StopRun()
		}
		dbconn.Exec("COMMIT")
		logs.Normal("COMMIT", "logid:", logid)
		logs.Normal("delete by key OK", "logid:", logid)
		this.Ctx.Output.Body([]byte(`{"result":0}`))
		this.StopRun()
	}

	err = sshKeysM.Delete(dbconn, uid, logid)
	if err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("delete by uid error:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}

	//reload ssh public key
	err = this.reloadSshPublicKeys(dbconn, uid, logid)
	if err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("reload ssh public key error:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":` + err.Error() + `}`))
		this.StopRun()
	}
	dbconn.Exec("COMMIT")
	logs.Normal("COMMIT", "logid:", logid)
	logs.Normal("delete over", "logid:", logid)
	this.Ctx.Output.Body([]byte(`{"result":0}`))
	this.StopRun()
}

//curl localhost:9090/sshKeys -X PUT -i -d '{"uid":5544,"keyname":"testliuyanglong","publickey":"d","logid":"444444"}'
func (this *SshKeysController) Put() {
	requestBody := string(this.Ctx.Input.CopyBody())
	var sshKeysOb models.SshKeys
	err := json.Unmarshal(this.Ctx.Input.CopyBody(), &sshKeysOb)
	if err != nil {
		//fmt.Println(err)
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"param error"}`))
		this.StopRun()
	}
	logid, err := models.GetLogId(this.Ctx.Input.CopyBody())
	if err != nil {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"logid param error"}`))
		this.StopRun()
	}
	logs.Normal("param is:", requestBody, "logid:", logid)
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
	sshKeysM := new(models.SshKeysManage)
	err = sshKeysM.Update(dbconn, sshKeysOb, logid)
	if err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("ssh keys update error:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}

	//reload ssh public key
	err = this.reloadSshPublicKeys(dbconn, strconv.Itoa(sshKeysOb.Uid), logid)
	if err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("reload ssh public key error:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":` + err.Error() + `}`))
		this.StopRun()
	}
	dbconn.Exec("COMMIT")
	logs.Normal("COMMIT", "logid:", logid)
	logs.Normal("ssh key put ok", "logid:", logid)
	this.Ctx.Output.Body([]byte(`{"result":0}`))
	this.StopRun()
}

func (this *SshKeysController) reloadSshPublicKeys(dbconn *sql.DB, uid string, logid int64) error {
	sshKeysM := new(models.SshKeysManage)
	pubSshKeyMap, errall := sshKeysM.GetAll(dbconn, uid, logid)
	logs.Normal("reloadSshPublicKeys get sshkey map:", pubSshKeyMap, "logid:", logid)
	if errall != nil {
		logs.Error("sshKeysM.GetAll error!", errall, " uid:", uid, "logid:", logid)
		return errall
	}
	//check ssh rules,if the rules is null,then only change ssh keys;
	// if the rules is not null,then change the proxy server
	sshRuleM := new(models.SshRuleManage)
	rulelist, err := sshRuleM.QueryByUid(dbconn, uid, logid)
	if err != nil {
		logs.Error("sshRuleM.QueryByUid error!", err, " uid:", uid, "logid:", logid)
		return err
	}
	logs.Normal("sshRuleM.QueryByUid get rule list:", rulelist, "logid:", logid)
	if len(rulelist) > 0 {
		err = models.UpdateRule(rulelist, pubSshKeyMap, logid)
	}
	return err
}
