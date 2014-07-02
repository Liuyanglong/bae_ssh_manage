package models

import (
	"fmt"
	"testing"
)

func Test_SshPortGetTableName(t *testing.T) {
	sshmodel := new(SshPortManage)
	fmt.Println("get ssh_port_manage name:", sshmodel.TableName())
}

func Test_GetDbConn(t *testing.T) {
	dbconn, err := InitDbConn()
	fmt.Println("get db conn:", dbconn, err)
	if err != nil {
		t.Error(err)
	}
}

func noTest_SshPortInsert(t *testing.T) {
	sshPortM := new(SshPortManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1
	var sshPortOb SshPort
	sshPortOb.ContainerNumber = 10
	sshPortOb.IsAvail = true
	sshPortOb.PublicPort = 100
	err := sshPortM.Insert(dbconn, sshPortOb)
	fmt.Println("case 1:", sshPortOb, err)
	if err != nil {
		t.Error(err)
	}

	//case 2
	var sshPortOb2 SshPort
	//sshPortOb.containerNumber = 1
	//sshPortOb.isAvail = true
	sshPortOb2.PublicPort = 101
	err = sshPortM.Insert(dbconn, sshPortOb2)
	fmt.Println("case 2:", sshPortOb2, err)
	if err != nil {
		t.Error(err)
	}
}

func noTest_SshPortQuery(t *testing.T) {
	sshPortM := new(SshPortManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1,table have this recode
	var sshPortOb SshPort
	sshPortOb, err := sshPortM.Query(dbconn, "100")
	fmt.Println("case 1:", sshPortOb, err)
	if err != nil {
		t.Error(err)
	}

	//case 2,table do not have this recode
	var sshPortOb2 SshPort
	sshPortOb2, err = sshPortM.Query(dbconn, "10000")
	fmt.Println("case 2:", sshPortOb2, err)
	if err != nil {
		t.Error(err)
	}
}

func noTest_SshPortDelete(t *testing.T) {
	sshPortM := new(SshPortManage)
	dbconn, dberr := InitDbConn()
	fmt.Println("dberror:", dberr)
	if dberr != nil {
		t.Error(dberr)
	}

	//case 1,table have this recode
	err := sshPortM.Delete(dbconn, "100")
	fmt.Println("case 1:", err)
	if err != nil {
		t.Error(err)
	}

	//case 2,table do not have this recode
	err = sshPortM.Delete(dbconn, "10000")
	fmt.Println("case 2:", err)
	if err != nil {
		t.Error(err)
	}
}
