package controllers

import (
	//"agent_ssh/models"
	"bytes"
	"encoding/json"
	//"fmt"
	"github.com/astaxie/beego"
	"os"
	"os/exec"
	"ssh_proxy_manage/logs"
	"strings"
)

type SshAgentController struct {
	beego.Controller
}

func (this *SshAgentController) UpdateContainerRull() {
	container := this.GetString("container")
	if container == "" {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(`{"result":1,"message":"container missing"}`))
		this.StopRun()
	}

	keylist := this.GetString("keylist")
	if keylist == "" {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(`{"result":1,"message":"keylist missing"}`))
		this.StopRun()
	}

	logid := this.GetString("logid")
	if logid == "" {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(`{"result":1,"message":"logid missing"}`))
		this.StopRun()
	}

	token := this.GetString("token")
	if token == "" {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(`{"result":1,"message":"token missing"}`))
		this.StopRun()
	}

	uiip := this.GetString("uiip")
	if uiip == "" {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(`{"result":1,"message":"uiip missing"}`))
		this.StopRun()
	}

	uiport := this.GetString("uiport")
	if uiport == "" {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(`{"result":1,"message":"uiport missing"}`))
		this.StopRun()
	}

	if this.checkToken(token) == false {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(`{"result":1,"message":"token error"}`))
		this.StopRun()
	}

	keylistMap := make(map[string]string)
	err := json.Unmarshal([]byte(keylist), &keylistMap)
	if err != nil {
		logs.Error("UpdateContainerRull keylist error", keylistMap, err, logid)
		this.Ctx.Output.SetStatus(403)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"param error"}`))
		this.StopRun()
	}
	logs.Normal("UpdateContainerRull params: ", container, keylistMap, token, logid, uiip, uiport, logid)

	//增加用户
	this.system("/usr/sbin/useradd "+container, logid)
	authKeyStr := ""
	for _, authkey := range keylistMap {
		authKeyStr += `command="ssh bae@` + uiip + ` -p ` + uiport + `" ` + authkey + "\n"
	}

	logs.Normal("UpdateContainerRull authkey is ", authKeyStr, logid)

	//将内容写入authorized_key文件
	authkeyfile := "/home/ssh/authorized_key_" + container
	logs.Normal("UpdateContainerRull authkey file is ", authkeyfile, logid)
	fout, err := os.Create(authkeyfile)
	defer fout.Close()
	if err != nil {
		logs.Error("authkeyfile is :", authkeyfile, err, logid)
		this.Ctx.Output.SetStatus(500)
		this.Ctx.Output.Body([]byte(`{"result":1,"error":"error happened"}`))
		this.StopRun()
	}

	//生成id_rsa&id_rsa.pub
	this.system("/bin/echo -e 'y\n' | /usr/bin/ssh-keygen -t rsa -f /home/bae/.ssh/id_rsa -N '' -q > /dev/null", logid)

	/*
	*
	*  todo 将生成的id_rsa.pub 通过调度传到container为authorized_keys；
	*	同时将container中/home/bae/.ssh的文件设置为root，bae不可读
	 */

	fout.WriteString(authKeyStr)
	this.Ctx.Output.Body([]byte(`{"result":0}`))
	this.StopRun()
}

func (this *SshAgentController) DeleteContainerRull() {
	container := this.GetString("container")
	if container == "" {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(`{"result":1,"message":"container missing"}`))
		this.StopRun()
	}

	logid := this.GetString("logid")
	if logid == "" {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(`{"result":1,"message":"logid missing"}`))
		this.StopRun()
	}

	token := this.GetString("token")
	if token == "" {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(`{"result":1,"message":"token missing"}`))
		this.StopRun()
	}

	if this.checkToken(token) == false {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(`{"result":1,"message":"token error"}`))
		this.StopRun()
	}

	//删除用户
	this.system("/usr/sbin/userdel "+container, logid)

	//删除对应的authorized_key文件
	authkeyfile := "/home/ssh/authorized_key_" + container
	this.system("rm -f "+authkeyfile, logid)
	this.Ctx.Output.Body([]byte(`{"result":0}`))
	this.StopRun()
}

func (this *SshAgentController) GetRulesNumber() {
	token := this.GetString("token")
	if token == "" {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(`{"result":1,"message":"token missing"}`))
		this.StopRun()
	}

	if this.checkToken(token) == false {
		this.Ctx.Output.SetStatus(400)
		this.Ctx.Output.Body([]byte(`{"result":1,"message":"token error"}`))
		this.StopRun()
	}

	cmd := "cat /etc/passwd|wc -l"
	lnum, _ := this.system(cmd, "0")
	this.Ctx.Output.Body([]byte(`{"result":0,"rule_num":` + lnum + `}`))
	this.StopRun()
}

func (this *SshAgentController) checkToken(token string) bool {
	check_token := beego.AppConfig.String("token")
	return check_token == token
}

func (this *SshAgentController) system(s, logid string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logs.Error("exec command error happened:", err, "the command is :", s, logid)
		return "", err
	}
	output := strings.TrimSpace(out.String())
	logs.Normal("exec command result is: ", output, "process state is :", cmd.ProcessState, logid)
	//fmt.Printf("%s\n", out.String(), cmd.ProcessState)
	return output, err
}
