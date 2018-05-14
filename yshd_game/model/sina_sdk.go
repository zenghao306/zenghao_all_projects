package model

import (
	"encoding/json"
	"fmt"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"io/ioutil"
	"net/http"
)

var (
	sina_addr   = "https://api.weibo.com/2/users/show.json"
	sina_appkey = "410136203"
)

func InitSinaAppKey() {
	sina_appkey = common.Cfg.MustValue("sina", "sina_appkey")
}

type SinaUserRet struct {
	Request      int    `json:"request"`
	Error_code   string `json:"error_code"`
	Error        string `json:"error"`
	Gender       string `json:"gender"`
	Screen_name  string `json:"screen_name"`
	Province     string `json:"province"`
	City         string `json:"city"`
	Avatar_hd    string `json:"avatar_hd"`
	Avatar_large string `json:"avatar_large"`
	Location     string `json:"location"`
}

//2418231496 849981
func SinaSDKLogin(access_token, uid, channel_id, device string, registerFrom int) int {
	url := fmt.Sprintf("%s?access_token=%s&uid=%s", sina_addr, access_token, uid)
	resp, err := http.Get(url)
	if err != nil {
		common.Log.Errf("qq token check err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		common.Log.Errf("qq token check err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	var v SinaUserRet
	json.Unmarshal(body, &v)
	if v.Error_code != "" {
		common.Log.Errf("sina skd login err is %v", &v)
		return 1
	}
	//user, has := GetUserByAccountWithSinaPlatform(access_token)
	user, has := GetUserByAccountAndPlatfrom(access_token, common.PLATFORM_SINA)
	if has {
		ret := user.CheckAccountForbid()
		if ret != common.ERR_SUCCESS {
			return ret
		}
		return common.ERR_SUCCESS
	} else {
		sex := 1
		if v.Gender == "m" {
			sex = 1
		} else if v.Gender == "f" {
			sex = 0
		}
		retcode := CreateAccountBySina(access_token, common.Md5("yunshanghudong123456*&^"), v.Location, v.Avatar_large, v.Screen_name, common.PLATFORM_SINA, sex, registerFrom, 0, channel_id, device)
		if retcode != common.ERR_SUCCESS {
			return common.ERR_UNKNOWN
		} else {
			return common.ERR_SUCCESS
		}
	}
}
