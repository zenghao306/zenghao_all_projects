package controller

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"hdt_app_go/common"
	. "hdt_app_go/gateway_app/log"
	. "hdt_app_go/gateway_app/model"
	. "hdt_app_go/gateway_app/rpc"
	proto "hdt_app_go/protcol"
	"io/ioutil"
	"net/http"
	"net/url"
	//"strconv"
	"strconv"
	//"github.com/astaxie/beego/context/param"
	"strings"
	"time"
)

const (
	QIANXUN_UID        = "319"
	QIANXUN_KEY        = "a227930be074d74da1a9c1eae28ef29f"
	QIANXUNSERVER_ADDR = "http://c.dzd.com/v4/sms/send.do"
)

type QianXun struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func QianXunSnsController(ctx iris.Context) {
	data := map[string]interface{}{
		"errcode": proto.ERR_OK,
	}

	body, _ := ioutil.ReadAll(ctx.Request().Body)
	jsonStr := string(body)

	var req QianXunReq
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil || len(req.Tel) < 11 { //手机号长度必须大于等于11位
		Log.Err(err.Error())
		data["errcode"] = proto.ERR_PARAM
		data["err_msg"] = err.Error()
		ctx.JSON(data)
		return
	}

	code := GenQianXunSnsCode(req.Tel)

	data["errcode"] = RequestQianXunVerify(req.Tel, code)

	data["errcode"] = proto.ERR_OK
	data["code"] = code

	ctx.JSON(data)
}

// 生成千讯验证码并存储验证码到后台redis数据库。
func GenQianXunSnsCode(tel string) string {
	code := common.RandnomRange64(1000, 9999)
	code_ := strconv.FormatInt(code, 10)
	v := &proto.QianxunReq{}
	v.Tel = tel
	v.Code = code_

	errCode := RpcClient.Register.AddQianXunCode(v)
	if errCode != proto.ERR_OK { //不行就再来一遍
		RpcClient.Register.AddQianXunCode(v)
	}

	return code_
}

func RequestQianXunVerify(tel string, code string) int {
	timeNow := time.Now().Unix()

	t := time.Unix(timeNow, 0).Format("20060102150405")

	v := url.Values{}
	v.Set("uid", QIANXUN_UID)
	v.Set("timestamp", t)
	v.Set("sign", GetQianXunSign(QIANXUN_UID, QIANXUN_KEY, t))
	v.Set("mobile", tel)
	v.Set("text", "【云尚互动APP】您的验证码是:"+code)

	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := &http.Client{}
	reqest, _ := http.NewRequest("POST", QIANXUNSERVER_ADDR, body)

	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded;  param=value ; charset=utf-8") //

	response, err := client.Do(reqest)
	defer response.Body.Close()
	if err != nil {
		Log.Err(err.Error())
		return proto.ERR_UNKNOWN
	}

	var m QianXun
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		//godump.Dump(string(body))
		err = json.Unmarshal(body, &m)
		if err != nil {
			Log.Err("error: %s,body %s", err.Error(), string(body))
			return proto.ERR_UNKNOWN
		}

		switch m.Code {
		case 0:
			return proto.ERR_OK
		case 3:
			return proto.ERR_VERIFT_CODE
		default:
			Log.Err("Qian Xun code=? ,msg=?,tel=? ", m.Code, m.Msg, tel)
			return proto.ERR_UNKNOWN
		}
	}
	return proto.ERR_UNKNOWN
}

func GetQianXunSign(uid, key string, time_format string) string {
	cal := fmt.Sprintf("%s%s%v", uid, key, time_format)
	return strings.ToLower(common.Md5(cal))
}
