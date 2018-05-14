package model

import (
	"encoding/json"
	"fmt"
	"github.com/yshd_game/common"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var (
	//qq_addr           = "https://openmobile.qq.com/user/get_simple_userinfo"
	qq_addr       = "https://graph.qq.com/user/get_user_info"
	qq_union_addr = "https://graph.qq.com/oauth2.0/me"
	//qq_ios_appkey     = "1105172675"
	qq_ios_appkey = "1105711217"
	//qq_android_appkey = "1105265368"
	qq_android_appkey = "1105711217"
	qq_h5_appkey      = "1105711217"
)

type QQUserRet struct {
	Ret            int    `json:"ret"`
	Msg            string `json:"msg"`
	Sex            int    `json:"sex"`
	Nickname       string `json:"nickname"`
	Figureurl_qq_2 string `json:"figureurl_qq_2"`
	Gender         string `json:"gender"`
}

type QQUniobRet struct {
	ClientId         string `json:"client_id"`
	OpenId           string `json:"openid"`
	UnionbId         string `json:"unionid"`
	Error            int    `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func InitQQKey() {
	qq_ios_appkey = common.Cfg.MustValue("tencent", "qq_appkey")
	qq_android_appkey = qq_ios_appkey
	qq_h5_appkey = common.Cfg.MustValue("tencent", "qq_h5_appkey")
}

func GetQQUnionId(access_token string) (unionid string, opendid string, ret int) {
	url_union := fmt.Sprintf("%s?access_token=%s&unionid=1", qq_union_addr, access_token)
	resp, err := http.Get(url_union)
	if err != nil {
		common.Log.Errf("qq token check err is %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		common.Log.Errf("qq token check err is %s", err.Error())

		ret = common.ERR_UNKNOWN
		return
	}

	re, _ := regexp.Compile("callback\\(")
	src2 := re.ReplaceAllString(string(body), "")
	re2, _ := regexp.Compile("\\)\\;")
	src3 := re2.ReplaceAllString(src2, "")

	var v QQUniobRet
	json.Unmarshal([]byte(src3), &v)

	if v.Error != 0 {
		common.Log.Errf("qq req unionid err is %d,%s", v.Error, v.ErrorDescription)
		ret = common.ERR_UNKNOWN
		return
	}

	unionid = v.UnionbId
	opendid = v.OpenId
	ret = common.ERR_SUCCESS
	return
}

func QQLogin(access_token string, os int, channel_id, device string, registerFrom int) (unionid string, ret int, new_register int) {
	//func QQLogin(access_token, openid string, os int, channel_id, device string) int {
	appkey := qq_ios_appkey
	if os == 1 {
		appkey = qq_ios_appkey
	} else if os == 0 {
		appkey = qq_android_appkey
	} else if os == 2 {
		appkey = qq_h5_appkey
	}

	unionid, opend_id, ret := GetQQUnionId(access_token)
	if ret != common.ERR_SUCCESS {
		return
	}
	/*
		url := fmt.Sprintf("%s?access_token=%s&openid=%s&oauth_consumer_key=%s", qq_addr, access_token, opend_id, appkey)

		resp, err := http.Get(url)
		if err != nil {
			common.Log.Errf("qq token check err is %s", err.Error())
			 ret=common.ERR_UNKNOWN
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			common.Log.Errf("qq token check err is %s", err.Error())
			 ret=common.ERR_UNKNOWN
			return
		}

		var v QQUserRet
		json.Unmarshal(body, &v)
	*/

	////if v.Ret == 0 {

	user, has := GetUserByAccountAndPlatfrom(unionid, common.PLATFORM_QQ)
	if has {
		ret = user.CheckAccountForbid()

		if ret != common.ERR_SUCCESS {
			return
		}
		InsertLog(common.ACTION_TYPE_LOG_LOGIN, user.Uid, "")
		ret = common.ERR_SUCCESS
		return
	} else {

		new_register = 1
		url := fmt.Sprintf("%s?access_token=%s&openid=%s&oauth_consumer_key=%s", qq_addr, access_token, opend_id, appkey)

		resp, err := http.Get(url)
		if err != nil {
			common.Log.Errf("qq token check err is %s", err.Error())
			ret = common.ERR_UNKNOWN
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			common.Log.Errf("qq token check err is %s", err.Error())
			ret = common.ERR_UNKNOWN
			return
		}

		var v QQUserRet
		err = json.Unmarshal(body, &v)
		if err != nil {
			common.Log.Errf("qq pull info is %s", err.Error())
			ret = common.ERR_UNKNOWN
			return
		}

		v.Nickname = strings.Trim(v.Nickname, " ")
		if v.Ret == 0 {
			sex := 1
			if v.Gender == "男" {
				sex = 1
			} else if v.Gender == "女" {
				sex = 0
			}
			retcode := CreateAccountByQQ(unionid, common.Md5("yunshanghudong123456@#$"), "地球", v.Figureurl_qq_2, v.Nickname, common.PLATFORM_QQ, sex, registerFrom, 0, channel_id, device, opend_id)
			if retcode != common.ERR_SUCCESS {
				ret = common.ERR_UNKNOWN
				return
			} else {
				ret = common.ERR_SUCCESS
				return
			}
		}

	}

	return
}
