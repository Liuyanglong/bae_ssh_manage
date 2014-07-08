package models

import (
	"encoding/json"
	//"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"ssh_proxy_manage/logs"
	"strconv"
	"strings"
)

func UpdateRule(rulelist []SshRule, keylist map[string]string, logid int64) error {
	logs.Normal("update rule list:", rulelist, "key list", keylist, "logid:", logid)
	token := beego.AppConfig.String("proxy_url_token")
	keylistStr, _ := json.Marshal(keylist)
	//useradd
	for _, sshRuleOb := range rulelist {
		container := sshRuleOb.ContainerName
		proxyHost := sshRuleOb.ProxyHost
		uiIpPort := sshRuleOb.UiIpPort
		uiArr := strings.Split(uiIpPort, ":")
		if len(uiArr) != 2 {
			logs.Error("update rule uiIpPort Error!", uiIpPort, container, proxyHost)
			continue
		}
		//this url todo
		proxyAddUrl := "http://" + proxyHost + ":9090/updateContainer?container=" + container + "&uiip=" + uiArr[0] + "&uiport=" + uiArr[1] + "&keylist=" + string(keylistStr) + "&logid=" + strconv.FormatInt(logid, 10) + "&token=" + token
		logs.Normal("curl add container url:", proxyAddUrl, "logid:", logid)

		req := httplib.Get(proxyAddUrl)
		output := make(map[string]interface{})
		err := req.ToJson(&output)
		if err != nil {
			logs.Error("request from "+proxyAddUrl+" error:", err, logid)
			return err
		}

		if output["result"].(int) == 0 {
			logs.Normal(proxyAddUrl, "response ok!", logid)
			continue
		} else {
			logs.Error(proxyAddUrl+" error:", output["error"], logid)
		}

	}
	return nil
}

func DeleteContainerUserFromProxy(rulelist []SshRule, logid int64) error {
	logs.Normal("delete container from proxy:", rulelist, "logid:", logid)
	token := beego.AppConfig.String("proxy_url_token")
	//userdel
	for _, sshRuleOb := range rulelist {
		container := sshRuleOb.ContainerName
		proxyHost := sshRuleOb.ProxyHost
		//this url todo
		proxyDelUrl := proxyHost + ":9090/deleteContainer?container=" + container + "&logid=" + strconv.FormatInt(logid, 10) + "&token=" + token
		logs.Normal("curl delete container url:", proxyDelUrl, "logid:", logid)

		req := httplib.Get(proxyDelUrl)
		output := make(map[string]interface{})
		err := req.ToJson(&output)
		if err != nil {
			logs.Error("request from "+proxyDelUrl+" error:", err, logid)
			return err
		}

		if output["result"].(int) == 0 {
			logs.Normal(proxyDelUrl, "response ok!", logid)
			continue
		} else {
			logs.Error(proxyDelUrl+" error:", output["error"], logid)
		}
	}
	return nil
}

func GetLogId(jsonstr []byte) (int64, error) {
	f := make(map[string]interface{})
	err := json.Unmarshal(jsonstr, &f)
	if err != nil {
		return 0, err
	}

	var logid int64
	f_logid := f["logid"]
	switch f_logid.(type) {
	case string:
		logid, _ = strconv.ParseInt(f_logid.(string), 10, 64)
	case int:
		logid = int64(f["logid"].(int))
	case int64:
		logid, _ = f_logid.(int64)
	case float64:
		logid_f, _ := f_logid.(float64)
		logid = int64(logid_f)
	}
	return logid, err
}
