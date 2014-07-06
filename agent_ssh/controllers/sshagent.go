package controllers

import (
	//"agent_ssh/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"os/exec"
)

type SshAgentController struct {
	beego.Controller
}

func (this *SshAgentController) Post() {
	var command []string
	json.Unmarshal(this.Ctx.Input.RequestBody, &command)
	fmt.Println(command)
	for _, comm := range command {
		err := this.system(comm)
		if err != nil {
			this.Ctx.Output.Body([]byte(`{"result":1,"message":"` + err.Error() + `"}`))
			this.StopRun()
		}
	}
	this.Ctx.Output.Body([]byte(`{"result":0}`))
	this.StopRun()
}

func (this *SshAgentController) Get() {

}

func (this *SshAgentController) Put() {

}

func (this *SshAgentController) Delete() {

}

func (this *SshAgentController) system(s string) error {
	cmd := exec.Command("/bin/sh", "-c", s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("error happened:", err)
		return err
	}
	fmt.Printf("%s\n", out.String(), cmd.ProcessState)
	return err
}
