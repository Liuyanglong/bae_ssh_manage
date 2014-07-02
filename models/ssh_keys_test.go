package models

import (
	"fmt"
	"testing"
)

func Test_KeysGetTableName(t *testing.T) {
	sshmodel := new(SshKeysManage)
	fmt.Println("get ssh_port_manage name:", sshmodel.TableName())
}

func noTest_KeysInsert(t *testing.T) {
	sshKeysM := new(SshKeysManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1
	var sshKeysOb SshKeys
	sshKeysOb.Uid = 33333
	sshKeysOb.KeyName = "test keys"
	sshKeysOb.PublicKey = "sa;ljdfosudoasufoasudfoisajdflsjlsjlsjlsjsljlkajsdlkfj"
	err := sshKeysM.Insert(dbconn, sshKeysOb)
	fmt.Println("case 1:", sshKeysOb, err)
	if err != nil {
		t.Error(err)
	}

	//case 2
	var sshPortOb2 SshKeys
	sshPortOb2.Uid = 9999
	err = sshKeysM.Insert(dbconn, sshPortOb2)
	fmt.Println("case 2:", sshPortOb2, err)
	if err != nil {
		t.Error(err)
	}
}

func noTest_KeysQuery(t *testing.T) {
	sshKeysM := new(SshKeysManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1,table have this recode
	var sshKeysOb SshKeys
	sshKeysOb, err := sshKeysM.Query(dbconn, "33333")
	fmt.Println("case 1:", sshKeysOb, err)
	if err != nil {
		t.Error(err)
	}

	//case 2,table do not have this recode
	var sshPortOb2 SshKeys
	sshPortOb2, err = sshKeysM.Query(dbconn, "555555")
	fmt.Println("case 2:", sshPortOb2, err)
	if err != nil {
		t.Error(err)
	}
}

func Test_KeysDelete(t *testing.T) {
	sshKeysM := new(SshKeysManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1,table have this recode
	err := sshKeysM.Delete(dbconn, "33333")
	fmt.Println("case 1:", err)
	if err != nil {
		t.Error(err)
	}

	//case 2,table do not have this recode
	err = sshKeysM.Delete(dbconn, "888888888")
	fmt.Println("case 2:", err)
	if err != nil {
		t.Error(err)
	}
}

func Test_RuleUpdate(t *testing.T) {
	sshKeysM := new(SshKeysManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1
	var sshRuleOb1 SshKeys
	sshRuleOb1.Uid = 33333
	sshRuleOb1.KeyName = "test1test1"
	sshRuleOb1.PublicKey = "11.11.111.1111:444"
	err := sshKeysM.Update(dbconn, sshRuleOb1)
	fmt.Println("case 1:", err)
	if err != nil {
		t.Error(err)
	}

	//case 2
	var sshRuleOb2 SshKeys
	sshRuleOb2.Uid = 9999
	err = sshKeysM.Update(dbconn, sshRuleOb2)
	fmt.Println("case 2:", err)
	if err != nil {
		t.Error(err)
	}

	//case 3
	var sshRuleOb3 SshKeys
	sshRuleOb3.Uid = 9999
	sshRuleOb3.KeyName = "qqqqqq"
	sshRuleOb3.PublicKey = "ttttttttttttttttttttt"
	err = sshKeysM.Update(dbconn, sshRuleOb3)
	fmt.Println("case 3:", err)
	if err != nil {
		t.Error(err)
	}
}
