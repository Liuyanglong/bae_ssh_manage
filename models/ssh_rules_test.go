package models

import (
	"fmt"
	"testing"
)

func Test_RuleGetTableName(t *testing.T) {
	sshmodel := new(SshRuleManage)
	fmt.Println("get ssh_port_manage name:", sshmodel.TableName())
}

func noTest_RuleInsert(t *testing.T) {
	sshRuleM := new(SshRuleManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1
	var sshRuleOb SshRule
	sshRuleOb.ContainerName = "liuyanglong"
	sshRuleOb.ProxyHost = "111111111:9999"
	sshRuleOb.PublicPort = 44
	sshRuleOb.UiIpPort = "22222222:8888"
	sshRuleOb.Uid = 898989
	err := sshRuleM.Insert(dbconn, sshRuleOb)
	fmt.Println("case 1:", sshRuleOb, err)
	if err != nil {
		t.Error(err)
	}

	//case 2
	var sshPortOb2 SshRule
	sshPortOb2.ContainerName = "yang_lizan"
	err = sshRuleM.Insert(dbconn, sshPortOb2)
	fmt.Println("case 2:", sshPortOb2, err)
	if err != nil {
		t.Error(err)
	}
}

func noTest_RuleQuery(t *testing.T) {
	sshRuleM := new(SshRuleManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1,table have this recode
	var sshRuleOb SshRule
	sshRuleOb, err := sshRuleM.Query(dbconn, "liuyanglong")
	fmt.Println("case 1:", sshRuleOb, err)
	if err != nil {
		t.Error(err)
	}

	//case 2,table do not have this recode
	var sshPortOb2 SshRule
	sshPortOb2, err = sshRuleM.Query(dbconn, "long")
	fmt.Println("case 2:", sshPortOb2, err)
	if err != nil {
		t.Error(err)
	}
}

func noTest_RuleDelete(t *testing.T) {
	sshRuleM := new(SshRuleManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1,table have this recode
	err := sshRuleM.Delete(dbconn, "liuyanglong")
	fmt.Println("case 1:", err)
	if err != nil {
		t.Error(err)
	}

	//case 2,table do not have this recode
	err = sshRuleM.Delete(dbconn, "longtest")
	fmt.Println("case 2:", err)
	if err != nil {
		t.Error(err)
	}
}

func Test_RuleUpdate(t *testing.T) {
	sshRuleM := new(SshRuleManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1
	var sshRuleOb1 SshRule
	sshRuleOb1.ContainerName = "yang_lizan"
	sshRuleOb1.ProxyHost = "test1test1"
	sshRuleOb1.PublicPort = 1234
	sshRuleOb1.UiIpPort = "11.11.111.1111:444"
	sshRuleOb1.Uid = 54321
	err := sshRuleM.Update(dbconn, sshRuleOb1)
	fmt.Println("case 1:", err)
	if err != nil {
		t.Error(err)
	}

	//case 2
	var sshRuleOb2 SshRule
	sshRuleOb2.ContainerName = "yang_lizan"
	err = sshRuleM.Update(dbconn, sshRuleOb2)
	fmt.Println("case 2:", err)
	if err != nil {
		t.Error(err)
	}

	//case 3
	var sshRuleOb3 SshRule
	sshRuleOb3.ContainerName = "yang_lizan_hahaha"
	err = sshRuleM.Update(dbconn, sshRuleOb3)
	fmt.Println("case 3:", err)
	if err != nil {
		t.Error(err)
	}
}
