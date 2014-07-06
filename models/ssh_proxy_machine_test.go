package models

import (
	"fmt"
	"testing"
)

func Test_GetTableName(t *testing.T) {
	sshmodel := new(SshProxyManage)
	fmt.Println("get ssh_port_manage name:", sshmodel.TableName())
}

func noTest_Insert(t *testing.T) {
	sshProxyM := new(SshProxyManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1
	var sshProxyOb SshProxyMachine
	sshProxyOb.MachineHost = "liu"
	sshProxyOb.IsAvail = true
	sshProxyOb.UserNumber = 100
	err := sshProxyM.Insert(dbconn, sshProxyOb)
	fmt.Println("case 1:", sshProxyOb, err)
	if err != nil {
		t.Error(err)
	}

	//case 2
	var sshPortOb2 SshProxyMachine
	sshPortOb2.MachineHost = "yang"
	err = sshProxyM.Insert(dbconn, sshPortOb2)
	fmt.Println("case 2:", sshPortOb2, err)
	if err != nil {
		t.Error(err)
	}
}

func noTest_Query(t *testing.T) {
	sshProxyM := new(SshProxyManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1,table have this recode
	var sshProxyOb SshProxyMachine
	sshProxyOb, err := sshProxyM.Query(dbconn, "liu")
	fmt.Println("case 1:", sshProxyOb, err)
	if err != nil {
		t.Error(err)
	}

	//case 2,table do not have this recode
	var sshPortOb2 SshProxyMachine
	sshPortOb2, err = sshProxyM.Query(dbconn, "long")
	fmt.Println("case 2:", sshPortOb2, err)
	if err != nil {
		t.Error(err)
	}
}

func noTest_Delete(t *testing.T) {
	sshProxyM := new(SshProxyManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1,table have this recode
	err := sshProxyM.Delete(dbconn, "liu")
	fmt.Println("case 1:", err)
	if err != nil {
		t.Error(err)
	}

	//case 2,table do not have this recode
	err = sshProxyM.Delete(dbconn, "long")
	fmt.Println("case 2:", err)
	if err != nil {
		t.Error(err)
	}
}
