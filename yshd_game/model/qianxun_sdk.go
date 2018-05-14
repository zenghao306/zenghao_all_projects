package model

import (
	"encoding/json"
	"fmt"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	QianXunUid        string
	QianXunKey        string
	QianXunServerAddr string
)

type QianXun struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func init() {
	QianXunUid = "319"
	QianXunKey = "a227930be074d74da1a9c1eae28ef29f"
	QianXunServerAddr = "http://c.dzd.com/v4/sms/send.do"
	//QianXunServerAddr="http://c.dzd.com/v4/sms/get_balance.do"
}
func GetQianXunSign(uid, key string, time_format string) string {
	//godump.Dump(time_format)
	cal := fmt.Sprintf("%s%s%v", uid, key, time_format)
	//godump.Dump(cal)
	return strings.ToLower(common.Md5(cal))
}

func GenQianXunSnsCode(tel string) string {
	code := common.RandnomRange64(1000, 9999)
	code_ := strconv.FormatInt(code, 10)
	AddQianXunCode(tel, code_)
	return code_
}

func QianXunSnsVerify(tel string, sns_text string) int {
	result, code := GetQianXunCode(tel)
	if result != common.ERR_SUCCESS {
		return result
	}

	if code == sns_text {
		DelQianXunCode(tel)
		return common.ERR_SUCCESS
	}
	return common.ERR_SNS_CORRECT
}

//20170727777517
func RequestQianXunVerify(tel string, code string) int {
	timeNow := time.Now().Unix()

	t := time.Unix(timeNow, 0).Format("20060102150405")
	//godump.Dump(GetQianXunSign(QianXunUid,QianXunKey,t))
	/*
		v := url.Values{}
		v.Set("uid",QianXunUid)
		v.Set("timestamp",t)
		v.Set("sign",GetQianXunSign(QianXunUid,QianXunKey,t))
	*/

	v := url.Values{}
	v.Set("uid", QianXunUid)
	v.Set("timestamp", t)
	v.Set("sign", GetQianXunSign(QianXunUid, QianXunKey, t))
	v.Set("mobile", tel)
	v.Set("text", "【云尚互动APP】您的验证码是:"+code)

	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := &http.Client{}
	reqest, _ := http.NewRequest("POST", QianXunServerAddr, body)

	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded;  param=value ; charset=utf-8") //

	response, err := client.Do(reqest)
	defer response.Body.Close()
	if err != nil {
		common.Log.Err("client error")
		return common.ERR_UNKNOWN
	}

	var m QianXun
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		//godump.Dump(string(body))
		err = json.Unmarshal(body, &m)
		if err != nil {
			common.Log.Errf("error: %s,body %s", err.Error(), string(body))
			return common.ERR_UNKNOWN
		}

		//godump.Dump(m)
		switch m.Code {
		case 0:
			return common.ERR_SUCCESS
		case 3:
			return common.ERR_VERIFT_CODE
		default:
			common.Log.Errf("Qian Xun code=? ,msg=?,tel=? ", m.Code, m.Msg, tel)
			return common.ERR_UNKNOWN
		}
	}
	return common.ERR_UNKNOWN
}
