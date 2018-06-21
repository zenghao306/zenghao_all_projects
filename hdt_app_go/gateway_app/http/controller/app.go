package controller

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"github.com/kataras/iris"
	. "hdt_app_go/gateway_app/log"
	. "hdt_app_go/gateway_app/rpc"
	proto "hdt_app_go/protcol"
	"io/ioutil"
)

type AppListReq struct {
	Tel   string `json:"tel"`   //注册手机号
	Token string `json:"token"` //密码
	Index int    `json:"index"` //index[从0开始]
}

func AppList(ctx iris.Context) {
	data := map[string]interface{}{
		"errcode": proto.ERR_OK,
	}

	jsonStr, _ := ioutil.ReadAll(ctx.Request().Body)
	var req AppListReq
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		Log.Err(err.Error())
		data["errcode"] = proto.ERR_PARAM
		data["err_msg"] = err.Error()
		ctx.JSON(data)
		return
	}
	valid := validation.Validation{}
	valid.Required(req.Tel, "tel")
	valid.Required(req.Token, "token")
	if valid.HasErrors() || req.Index < 0 {
		for _, err1 := range valid.Errors {
			Log.Warningf("invalid args %s: %s", err1.Key, err1.Message)
		}
		data["errcode"] = proto.ERR_PARAM
		ctx.JSON(data)
		return
	}

	_, token := RpcClient.Register.GetUserToken(req.Tel)

	if token != req.Token {
		data["errcode"] = proto.ERR_EXPIRATION
		ctx.JSON(data)
		return
	}

	errCode, list := RpcClient.Register.AppList(req.Index)
	data["errcode"] = errCode
	if errCode == proto.ERR_OK {
		data["list"] = list
	}

	ctx.JSON(data)
	return
}

type AppDetailInfoReq struct {
	Tel   string `json:"tel"`    //注册手机号
	Token string `json:"token"`  //密码
	AppId int64  `json:"app_id"` //AppId
}

//AppDetailInfo
func AppDetailInfo(ctx iris.Context) {
	data := map[string]interface{}{
		"errcode": proto.ERR_OK,
	}

	jsonStr, _ := ioutil.ReadAll(ctx.Request().Body)
	var req AppDetailInfoReq
	err := json.Unmarshal([]byte(jsonStr), &req)
	if err != nil {
		Log.Err(err.Error())
		data["errcode"] = proto.ERR_PARAM
		data["err_msg"] = err.Error()
		ctx.JSON(data)
		return
	}
	valid := validation.Validation{}
	valid.Required(req.Tel, "tel")
	valid.Required(req.Token, "token")
	valid.Required(req.AppId, "app_id")
	if valid.HasErrors() {
		for _, err1 := range valid.Errors {
			Log.Warningf("invalid args %s: %s", err1.Key, err1.Message)
		}
		data["errcode"] = proto.ERR_PARAM
		ctx.JSON(data)
		return
	}

	_, token := RpcClient.Register.GetUserToken(req.Tel)

	if token != req.Token {
		data["errcode"] = proto.ERR_EXPIRATION
		ctx.JSON(data)
		return
	}

	errCode, rsp := RpcClient.Register.AppDetailInfo(req.Tel, req.AppId)
	data["errcode"] = errCode
	if errCode == proto.ERR_OK {
		data["user_app_hdt"] = rsp.UserAppHdt
		data["app_hdt_total"] = rsp.AppHdtTotal
		data["app_content"] = rsp.AppContent
		data["app_imgs"] = rsp.AppImg
		data["app_ios_address"] = rsp.IosAddress
		data["app_android_address"] = rsp.AndroidAddress
	}

	ctx.JSON(data)
	return
}
