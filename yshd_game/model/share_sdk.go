package model

import (
	"encoding/json"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	server_addr = "https://webapi.sms.mob.com/sms/verify"
	//ios_appkey  = "f286c3f2e400"
	ios_appkey = "18738cf7e2598"
	//android_appkey = "f3975b53da60"
	android_appkey = "18738cf7e2598"
	zone           = "86"
)

func InitAppKey() {
	android_appkey := common.Cfg.MustValue("share_sdk_appkey", "android_appkey")
	SetAndroidAppKey(android_appkey)
	ios_appkey := common.Cfg.MustValue("share_sdk_appkey", "ios_appkey")
	SetIOSAppKey(ios_appkey)
}

func SetIOSAppKey(key string) {
	if key == "" {
		common.Log.Panicf("ios app key nil")
	}
	ios_appkey = key
}
func SetAndroidAppKey(key string) {
	if key == "" {
		common.Log.Panicf("android app key nil")
	}
	android_appkey = key
}

type Message struct {
	Status int
}

// curl -d 'appkey=18738cf7e2598&phone=18588249532&zone=86&code=6827' 'https://webapi.sms.mob.com/sms/verify'
//curl -d 'appkey=1d7cd69c545fc&phone=13054066184&zone=86&code=1446' 'https://webapi.sms.mob.com/sms/verify'
func RequestSnsVerify(tel, code, ptype string) int {
	v := url.Values{}
	appkey := android_appkey
	if ptype == "0" {
		appkey = android_appkey
	} else if ptype == "1" {
		appkey = ios_appkey
	}

	v.Set("appkey", appkey)
	v.Set("phone", tel)
	v.Set("zone", zone)
	v.Set("code", code)
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := &http.Client{}
	reqest, _ := http.NewRequest("POST", server_addr, body)
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	response, err := client.Do(reqest)
	defer response.Body.Close()
	if err != nil {
		common.Log.Err("client error")
		return common.ERR_UNKNOWN
	}
	var m Message
	if response.StatusCode == 200 {

		body, _ := ioutil.ReadAll(response.Body)

		err = json.Unmarshal(body, &m)
		if err != nil {
			common.Log.Errf("error: %s", err.Error())
			return common.ERR_UNKNOWN
		}

		switch m.Status {
		case 200:
			return common.ERR_SUCCESS
		case 467:
			return common.ERR_VERIFT_MUTIPLY
		case 468:
			return common.ERR_VERIFT_CODE
		default:
			return common.ERR_UNKNOWN
		}

	}
	return common.ERR_UNKNOWN
}
