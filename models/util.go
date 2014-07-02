package models

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"ssh_proxy_manage/logs"
	"strconv"
)

/*
* 将 用户的公钥数据 更新到对应的 proxy host
* @param rulelist 
* @param keylist 用户公钥
*/
func UpdateRule(rulelist []SshRule, keylist map[string]string, logid int64) error {
	logs.Normal("rule list:", rulelist, "key list", keylist, "logid:", logid)
	token := beego.AppConfig.String("proxy_url_token")
	keylistStr, _ := json.Marshal(keylist)
	//useradd
	for _, sshRuleOb := range rulelist {
		container := sshRuleOb.ContainerName
		proxyHost := sshRuleOb.ProxyHost
		//this url todo
		proxyAddUrl := proxyHost + ":9090/updateContainer?container=" + container + "&keylist=" + string(keylistStr) + "&logid=" + strconv.FormatInt(logid, 10) + "&token=" + token
		logs.Normal("curl add container url:", proxyAddUrl, "logid:", logid)
		/**
		 *
		 * todo call proxyAddUrl
		 *
		 */
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
		/**
		 *
		 * todo call proxyAddUrl
		 *
		 */
	}
	return nil
}

//获取logid
func GetLogId(jsonstr []byte) (int64, error) {
	f := make(map[string]interface{})
	err := json.Unmarshal(jsonstr, &f)
	if err != nil {
		return 0, err
	}

	var logid int64
	//强制将logid统一成int64型
	f_logid := f["logid"]
	switch f_logid.(type) {
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
