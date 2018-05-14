package controller

import (
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
	//	"strconv"
)

type CommonReq struct {
	Uid   int    `form:"uid" binding:"required"`
	Token string `form:"token" binding:"required"`
}

type CommonWithIndexReq struct {
	Uid   int    `form:"uid" binding:"required"`
	Token string `form:"token" binding:"required"`
	Index int    `form:"index"`
}

type CommonWithOid2Req struct {
	Uid   int    `form:"uid" binding:"required"`
	Token string `form:"token" binding:"required"`
	Oid   int    `form:"oid" binding:"required"`
}

type ExchangeItemReq struct {
	Uid    int    `form:"uid" binding:"required"`
	Token  string `form:"token" binding:"required"`
	ItemId int    `form:"item_id" binding:"required"`
	Name   string `form:"name" binding:"required"`
	Tel    string `form:"tel" binding:"required"`
	Addr   string `form:"addr" binding:"required"`
}

//curl -d 'token=1f59dd2bce2c73d041556ae8f85f9341&uid=7&bank=工商&card_no=55555555555555&real_name=大大'  'http://192.168.1.12:3000/card/add_card'
func AddCardInfo(req *http.Request, r render.Render, w http.ResponseWriter) {
	ret_value := make(map[string]interface{})
	uid := req.FormValue("uid")
	user, _ := model.GetUserExtraByUidStr(uid)
	if user.CheckCashTel() == false {
		ret_value[ServerTag] = common.ERR_CASH_TEL_UNBIND
		r.JSON(http.StatusOK, ret_value)
		return
	}
	bank := req.FormValue("bank")

	card_no := req.FormValue("card_no")

	real_name := req.FormValue("real_name")
	user.CardNo = card_no
	user.RealName = real_name
	user.Bank = bank
	if _, err := user.UpdateByColS("card_no", "real_name", "bank"); err != nil {
		ret_value[ServerTag] = common.ERR_UNKNOWN

		r.JSON(http.StatusOK, ret_value)
	}
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func BindCashTel(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	uid := req.FormValue("uid")
	user, _ := model.GetUserExtraByUidStr(uid)

	/*
		if user.CheckCashTel() {
			ret_value[ServerTag] = common.ERR_CASH_TEL_BIND
		} else {
			cash_tel := req.FormValue("cash_tel")
			code := req.FormValue("code")
			type_ := req.FormValue("type")
			result := model.RequestSnsVerify(cash_tel, code, type_)
			if result == common.ERR_SUCCESS {
				user.SetCashTel(cash_tel)
			}
			ret_value[ServerTag] = result
		}
	*/
	cash_tel := req.FormValue("cash_tel")
	code := req.FormValue("code")
	type_ := req.FormValue("type")
	if user.CashTel != "" {
		if user.IsChangeTel {
			result := model.RequestSnsVerify(cash_tel, code, type_)
			if result == common.ERR_SUCCESS {
				ret_value[ServerTag] = user.SetCashTel(cash_tel)
				user.SetChangeFlag(false)
			} else {
				ret_value[ServerTag] = result
			}

		} else {
			ret_value[ServerTag] = common.ERR_CASH_TEL_BIND
		}
	} else {
		result := model.RequestSnsVerify(cash_tel, code, type_)
		if result == common.ERR_SUCCESS {
			ret_value[ServerTag] = user.SetCashTel(cash_tel)
		} else {
			ret_value[ServerTag] = result
		}

	}

	r.JSON(http.StatusOK, ret_value)
}

func CashRiceShow(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	uid := req.FormValue("uid")
	user, _ := model.GetUserByUidStr(uid)
	ret_value["rice"] = user.Coupons            //当前米粒(包含冻结的)
	ret_value["current_rice"] = user.CashRice() //当前可提的米粒（不包含冻结的）
	ret_value["money"] = user.CashQuota()
	ret_value["all_money"] = user.CashQuotaAll() //所有米粒换算成钱
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

type CashBankReq struct {
	Uid   int    `form:"uid" binding:"required"`
	Token string `form:"token" binding:"required"`
	Money int    `form:"money" binding:"required"`
	//Code  string `form:"code" binding:"required"`
	Ptype string `form:"ptype"`
}

//curl -d  'token=a32f019db21573c79e4fe5a61bd29392&uid=103&money=1&code=1&Ptype=1' 'http://shangtv.cn:3000/card/cash_bank'

func CashBank(r render.Render, d CashBankReq) {
	ret_value := make(map[string]interface{})
	//	user_extre_, _ := model.GetUserExtraByUid(d.Uid)

	//result := model.RequestSnsVerify(user_extre_.CashTel, d.Code, d.Ptype)
	//if result == common.ERR_SUCCESS {
	user, _ := model.GetUserByUid(d.Uid)
	if user.CashQuota() >= d.Money {
		//rice := user.CashExchange(d.Money)
		ret_value[ServerTag] = user.ExchangeRiceToBank(d.Money)
	} else {
		ret_value[ServerTag] = common.ERR_OVER_CASH
	}
	//} else {
	//	ret_value[ServerTag] = result
	//}

	r.JSON(http.StatusOK, ret_value)
}

func ChangeTel(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	uid := req.FormValue("uid")
	user, _ := model.GetUserExtraByUidStr(uid)
	code := req.FormValue("code")
	type_ := req.FormValue("type")
	result := model.RequestSnsVerify(user.CashTel, code, type_)
	if result == common.ERR_SUCCESS {
		user.SetChangeFlag(true)
	}
	ret_value[ServerTag] = result
	r.JSON(http.StatusOK, ret_value)
}

//curl 'http://192.168.1.12:3000/card/cash_record?token=ea8aaed458825bcf83dcb6ca42ec68cf&uid=5&index=0'
func UserCashRecord(req *http.Request, r render.Render, d CommonWithIndexReq) {

	ret_value := make(map[string]interface{})
	user, has := model.GetUserExtraByUid(d.Uid)
	if has {
		ret_value[ServerTag], ret_value["record"] = user.GetCashRecord(d.Index)
	} else {
		ret_value[ServerTag] = common.ERR_UNKNOWN
	}

	r.JSON(http.StatusOK, ret_value)
}

func CheckCashTel(req *http.Request, r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})
	user, has := model.GetUserExtraByUid(d.Uid)
	if !has {
		ret_value[ServerTag] = common.ERR_UNKNOWN
		r.JSON(http.StatusOK, ret_value)
		return
	}
	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["tel"] = user.CashTel
	ret_value["card"] = user.CardNo
	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'uid=7&token=1d29748e886274e242872d79b968c240&item_id=10&name=1&tel=123&addr=ss'  'http://shangtv.cn:3003/moon/exchange_item'
func ExchangeItemController(req *http.Request, r render.Render, d ExchangeItemReq) {
	ret_value := make(map[string]interface{})
	u, _ := model.GetUserByUid(d.Uid)
	ret_value[ServerTag] = u.ExchangeItem(d.ItemId, d.Name, d.Tel, d.Addr)
	r.JSON(http.StatusOK, ret_value)
}

func CashWeiXin(r render.Render, d CashBankReq) {
	ret_value := make(map[string]interface{})

	user, _ := model.GetUserByUid(d.Uid)
	if user.CashQuota() >= d.Money {
		ret_value[ServerTag] = user.ExchangeRiceToWeiXin(d.Money)
	} else {
		ret_value[ServerTag] = common.ERR_OVER_CASH
	}

	r.JSON(http.StatusOK, ret_value)
}

//curl 'http://192.168.1.12:3003/moon/order_list?uid=2&index=0'
func OrderListController(req *http.Request, r render.Render, d CommonWithIndexReq) {
	ret_value := make(map[string]interface{})
	u, _ := model.GetUserByUid(d.Uid)

	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["order_list"] = u.MoonOrderList(d.Index)
	r.JSON(http.StatusOK, ret_value)
}
