package controller

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"github.com/kataras/iris"
	//"hdt_app_go/common"
	//"hdt_app_go/common/hdtcodec"
	. "hdt_app_go/gateway_app/log"
	//. "hdt_app_go/gateway_app/model"
	. "hdt_app_go/gateway_app/rpc"
	proto "hdt_app_go/protcol"
	"io/ioutil"
)

type TelTokenReq struct {
	Tel   string `json:"tel"`   //注册手机号
	Token string `json:"token"` //Token
}

func GetUserRankingInfo(ctx iris.Context) {
	data := map[string]interface{}{
		"errcode": proto.ERR_OK,
	}

	//以下这段代码是对传入的参数进行校验
	//body, _ := ioutil.ReadAll(ctx.Request().Body)
	//jsonStr, _ := hdtcodec.HdtDecodeV0(string(body))
	jsonStr, _ := ioutil.ReadAll(ctx.Request().Body)
	var req TelTokenReq
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

	errCode, rsp := RpcClient.Register.GetUserRankingInfo(req.Tel)
	data["errcode"] = errCode
	if errCode == proto.ERR_OK {
		data["difficulty"] = rsp.DegreeOfDifficulty
		data["hdt_mining_last"] = rsp.HdtMiningLast
		data["hdt_mining_total"] = rsp.HdtMiningTotal
		data["mining_index"] = rsp.MiningIndex
	}

	ctx.JSON(data)
	return
}

func GetUseRankingHdtDig(ctx iris.Context) {
	data := map[string]interface{}{
		"errcode": proto.ERR_OK,
	}

	//以下这段代码是对传入的参数进行校验
	jsonStr, _ := ioutil.ReadAll(ctx.Request().Body)
	var req TelTokenReq
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

	errCode, rsp := RpcClient.Register.GetUseRankingHdtDig(req.Tel)
	data["errcode"] = errCode
	if errCode == proto.ERR_OK {
		data["ranking_hdt_dig"] = rsp.RankingOfHdtDig
	}

	ctx.JSON(data)
	return
}

func MinePoolInfo(ctx iris.Context) {
	data := map[string]interface{}{
		"errcode": proto.ERR_OK,
	}

	//以下这段代码是对传入的参数进行校验
	jsonStr, _ := ioutil.ReadAll(ctx.Request().Body)
	var req TelTokenReq
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

	errCode, rsp := RpcClient.Register.GetMinePoolInfo(req.Tel)
	data["errcode"] = errCode
	if errCode == proto.ERR_OK {
		data["errcode"] = rsp.ErrCode
		data["hdt_supply_limit"] = rsp.HdtSupplyLimit
		data["app_hdt_balance_total"] = rsp.AppHdtBalanceTotal
		data["hdt_total_supply"] = rsp.HdtTotalSupply
		data["degree_difficulty"] = rsp.DegreeOfDifficulty
	}

	ctx.JSON(data)
	return
}

func GetMinePoolTaskList(ctx iris.Context) {
	data := map[string]interface{}{
		"errcode": proto.ERR_OK,
	}

	//以下这段代码是对传入的参数进行校验
	jsonStr, _ := ioutil.ReadAll(ctx.Request().Body)
	var req TelTokenReq
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

	errCode, rsp := RpcClient.Register.GetMinePoolTaskList(req.Tel)
	data["errcode"] = errCode
	if errCode == proto.ERR_OK {
		data["list"] = rsp.MinePoolTasklist
	}

	ctx.JSON(data)
	return
}