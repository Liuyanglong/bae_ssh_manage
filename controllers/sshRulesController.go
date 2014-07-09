package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"ssh_proxy_manage/logs"
	"ssh_proxy_manage/models"
	"strconv"
)

type SshRulesController struct {
	beego.Controller
}

//curl -i -X POST localhost:9090/sshRules/9999 -d '{"name":"aiaiia","port":4242,"rule":"~~~~","uid":9999}'
func (this *SshRulesController) Post() {
	var sshRulesOb models.SshRule
	requestBody := string(this.Ctx.Input.RequestBody)

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &sshRulesOb)
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

	sshRulesM := new(models.SshRuleManage)
	err = sshRulesM.Insert(dbconn, sshRulesOb, logid)
	if err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("ssh rules insert error:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":` + err.Error() + `}`))
		this.StopRun()
	}

	if err = this.reloadRules(dbconn, sshRulesOb, logid); err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("reload rules error:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":` + err.Error() + `}`))
		this.StopRun()
	}
	dbconn.Exec("COMMIT")
	logs.Normal("COMMIT", "logid:", logid)
	logs.Normal("post OK!", "logid:", logid)
	this.Ctx.Output.Body([]byte(`{"result":0}`))
}

//curl -i -X GET localhost:9090/sshRules/9999
func (this *SshRulesController) GetByUid() {
	uid := this.Ctx.Input.Params[":uid"]
	if uid == "" {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"param error"}`))
		this.StopRun()
	}

	dbconn, err := models.InitDbConn(0)
	if err != nil {
		logs.Error("init db conn error:", err, 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"init db conn error"}`))
		this.StopRun()
	}
	defer dbconn.Close()

	sshRulesM := new(models.SshRuleManage)
	rulelist, err := sshRulesM.QueryByUid(dbconn, uid, 0)
	if err != nil {
		logs.Error("ssh rule query by uid error:", err, 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}
	this.Data["json"] = rulelist
	this.ServeJson()
}

//curl -i -X GET localhost:9090/sshRules/9999/aiaiia
func (this *SshRulesController) GetByContainer() {
	uid := this.Ctx.Input.Params[":uid"]
	containerName := this.Ctx.Input.Params[":container"]
	if uid == "" {
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"param error"}`))
		this.StopRun()
	}

	dbconn, err := models.InitDbConn(0)
	if err != nil {
		logs.Error("init db conn error:", err, 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"init db conn error"}`))
		this.StopRun()
	}
	defer dbconn.Close()

	sshRulesM := new(models.SshRuleManage)
	if containerName == "" {
		rulelist, err := sshRulesM.QueryByUid(dbconn, uid, 0)
		if err != nil {
			logs.Error("ssh rule query by uid error:", err, 0)
			this.Ctx.Output.SetStatus(500)
			this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
			this.StopRun()
		}
		this.Data["json"] = rulelist
		this.ServeJson()
	}

	sshRulesOb, err := sshRulesM.Query(dbconn, uid, containerName, 0)
	if err != nil {
		logs.Error("ssh rule query error:", err, 0)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}
	this.Data["json"] = sshRulesOb
	this.ServeJson()
}

//curl -i -X GET localhost:9090/sshRules/9999
func (this *SshRulesController) DeleteByUid() {
	uid := this.Ctx.Input.Params[":uid"]
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
	logs.Normal("delete param:", "uid:", uid, "logid:", logid)
	dbconn, err := models.InitDbConn(logid)
	if err != nil {
		logs.Error("init db conn err:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"init db conn error"}`))
		this.StopRun()
	}
	defer dbconn.Close()

	dbconn.Exec("START TRANSACTION")
	logs.Normal("start transaction", "logid:", logid)

	sshRulesM := new(models.SshRuleManage)
	delrulelist, _ := sshRulesM.QueryByUid(dbconn, uid, logid)
	err = sshRulesM.DeleteByUid(dbconn, uid, logid)
	if err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("delete by uid err:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}
	err = models.DeleteContainerUserFromProxy(delrulelist, logid)
	if err != nil {
		logs.Error("DeleteContainerUserFromProxy error:", err, "logid:", logid)
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", logid)
	}
	dbconn.Exec("COMMIT")
	logs.Normal("COMMIT", "logid:", logid)
	logs.Normal("delete OK!", "logid:", logid)
	this.Ctx.Output.Body([]byte(`{"result":0}`))
	this.StopRun()
}

//curl -i -X GET localhost:9090/sshRules/9999/aiaiia
func (this *SshRulesController) DeleteByContainer() {
	uid := this.Ctx.Input.Params[":uid"]
	containerName := this.Ctx.Input.Params[":container"]
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
	logs.Normal("delete param:", "uid:", uid, "containerName:", containerName, "logid:", logid)
	dbconn, err := models.InitDbConn(logid)
	if err != nil {
		logs.Error("init db conn err:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"init db conn error"}`))
		this.StopRun()
	}
	defer dbconn.Close()

	dbconn.Exec("START TRANSACTION")
	logs.Normal("start transaction", "logid:", logid)

	sshRulesM := new(models.SshRuleManage)
	delrule, _ := sshRulesM.Query(dbconn, uid, containerName, logid)
	err = sshRulesM.Delete(dbconn, uid, containerName, logid)
	if err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("delete err:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}

	err = models.DeleteContainerUserFromProxy([]models.SshRule{delrule}, logid)
	if err != nil {
		logs.Error("Delete single Container From Proxy error:", err, "logid:", logid)
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", logid)
	}

	dbconn.Exec("COMMIT")
	logs.Normal("COMMIT", "logid:", logid)

	logs.Normal("delete OK!", "logid:", logid)
	this.Ctx.Output.Body([]byte(`{"result":0}`))
	this.StopRun()
}

//curl -i -X PUT localhost:9090/sshRules/9999 -d '{"name":"aiaiia","port":4242,"rule":"~~~~","uid":9999}'
func (this *SshRulesController) Put() {
	requestBody := string(this.Ctx.Input.CopyBody())

	var sshRulesOb models.SshRule
	err := json.Unmarshal(this.Ctx.Input.CopyBody(), &sshRulesOb)
	fmt.Println(string(this.Ctx.Input.CopyBody()))
	if err != nil {
		fmt.Println(err)
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
		logs.Error("init db conn err:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"init db conn error"}`))
		this.StopRun()
	}
	defer dbconn.Close()

	dbconn.Exec("START TRANSACTION")
	logs.Normal("start transaction", "logid:", logid)

	sshRulesM := new(models.SshRuleManage)
	err = sshRulesM.Update(dbconn, sshRulesOb, logid)
	if err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("ssh rule update error:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}

	sshRuleMsg, err := sshRulesM.Query(dbconn, strconv.Itoa(sshRulesOb.Uid), sshRulesOb.ContainerName, logid)
	if err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("ssh rule query error:", err, logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"` + err.Error() + `"}`))
		this.StopRun()
	}

	if err = this.reloadRules(dbconn, sshRuleMsg, logid); err != nil {
		dbconn.Exec("ROLLBACK")
		logs.Normal("ROLLBACK", "logid:", logid)
		logs.Error("reload rules error:", err, "logid:", logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":` + err.Error() + `}`))
		this.StopRun()
	}

	dbconn.Exec("COMMIT")
	logs.Normal("COMMIT", "logid:", logid)
	logs.Normal("put OK", "logid:", logid)
	this.Ctx.Output.Body([]byte(`{"result":0}`))
	this.StopRun()
}

func (this *SshRulesController) reloadRules(conn *sql.DB, sshRulesOb models.SshRule, logid int64) error {
	sshKeyOb := new(models.SshKeysManage)
	keylist, err := sshKeyOb.GetAll(conn, strconv.Itoa(sshRulesOb.Uid), logid)
	if err != nil {
		logs.Error("ssh key get all error:", err, "logid:", logid)
		return err
	}
	if len(keylist) <= 0 {
		return nil
	}
	err = models.UpdateRule([]models.SshRule{sshRulesOb}, keylist, logid)
	if err != nil {
		logs.Error("update rule error:", err, "rule list:", sshRulesOb, "key list:", keylist, "logid:", logid)
	}
	return err
}
