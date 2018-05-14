package controller

import (
	"fmt"
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"io/ioutil"
	"net/http"
	//"encoding/json"
	//"github.com/codinl/go-logger"
	//"github.com/bitly/go-simplejson"
	//"github.com/liudng/godump"
	"strings"
)

type ThirdPaySwitchReq struct {
	ChannelID string `form:"channel_id" binding:"required"`
	Uid       int    `form:"uid" binding:"required"`
	Token     string `form:"token" binding:"required"`
}

type ThirdPayReq struct {
	Uid       int    `form:"uid" binding:"required"`
	Token     string `form:"token" binding:"required"`
	TradeType int    `form:"trade_type" binding:"required"`
	Itemid    string `form:"itemid" binding:"required"`
	AnchorID  int    `form:"anchorid"`
	ChannelId string `form:"channel_id"`
	Device    string `form:"device"`
}

type HuaWeiPayReq struct {
	Uid       int    `form:"uid" binding:"required"`
	Token     string `form:"token" binding:"required"`
	Itemid    string `form:"itemid" binding:"required"`
	AnchorID  int    `form:"anchorid"`
	ChannelId string `form:"channel_id"`
	Device    string `form:"device"`
}

type AliPayReq struct {
	Uid       int    `form:"uid" binding:"required"`
	Token     string `form:"token" binding:"required"`
	Itemid    string `form:"itemid" binding:"required"`
	AnchorID  int    `form:"anchorid"`
	ChannelId string `form:"channel_id"`
	Device    string `form:"device"`
}

type NotifyReq struct {
	//FunCode           string `form:"funcode" binding:"required"`
	//AppId             string `form:"appId" binding:"required"`
	MhtOrderNo string `form:"mhtOrderNo"`
	//MhtOrderName      string `form:"mhtOrderName" binding:"required"`
	//MhtOrderType      string `form:"mhtOrderType" binding:"required"`
	//MhtCurrencyType   string `form:"mhtCurrencyType" binding:"required"`
	//MhtOrderAmt       int    `form:"mhtOrderAmt" binding:"required"`
	//MhtOrderTimeOut   int    `form:"mhtOrderTimeOut"`
	//MhtOrderStartTime string `form:"mhtOrderStartTime" binding:"required"`
	//MhtCharset        string `form:"mhtCharset" binding:"required"`
	NowPayOrderNo string `form:"nowPayOrderNo"`
	//DeviceType        string `form:"deviceType" binding:"required"`
	//PayChannelType    string `form:"payChannelType"`
	TradeStatus   string `form:"tradeStatus"`
	PayConsumerId string `form:"payConsumerId"`
	//MhtReserved       string `form:"mhtReserved"`
	//SignType          string `form:"signType" binding:"required"`
	//signature         string `form:"signature" binding:"required"`
}

func ThirdPaySwitch(req *http.Request, r render.Render, d ThirdPaySwitchReq) {
	ret_value := make(map[string]interface{})

	ret_value["ErrCode"] = common.ERR_SUCCESS
	ret_value["thirdpay_switch"] = model.ThirdPaySwitch(d.ChannelID)
	r.JSON(http.StatusOK, ret_value)
}

func NowPayPrepay(req *http.Request, r render.Render, d ThirdPayReq) {
	common.Log.Info("NowPayPrepay() called@@@@@@")
	ret_value := make(map[string]interface{})
	ip := common.GetRemoteIp(req)

	ret_value[ServerTag], ret_value["before_sign"], ret_value["sign"] = model.NowPayPrepay(d.Itemid, d.ChannelId, d.TradeType, d.Uid, d.AnchorID, ip, d.Device)

	r.JSON(http.StatusOK, ret_value)
}

func HuaWeiPrepay(req *http.Request, r render.Render, d HuaWeiPayReq) {
	ret_value := make(map[string]interface{})
	ip := common.GetRemoteIp(req)

	ret_value[ServerTag], ret_value["sign"], ret_value["m_list"] = model.HuaWeiPrepay(d.Itemid, d.ChannelId, d.Uid, d.AnchorID, ip, d.Device)

	r.JSON(http.StatusOK, ret_value)
}

func AliPayPrepay(req *http.Request, r render.Render, d AliPayReq) {
	ret_value := make(map[string]interface{})
	ip := common.GetRemoteIp(req)

	ret_value[ServerTag], ret_value["sign"], ret_value["m_list"] = model.AlipayPrepay(d.Itemid, d.ChannelId, d.Uid, d.AnchorID, ip, d.Device)

	r.JSON(http.StatusOK, ret_value)
}

type ThirdPartyReq struct {
	ThirdUID   string `json:"third_uid"` //第三方uid
	From       int    `json:"from"`      //来自哪个平台：1:QQ,2:微信，3:SINA
	HeadImgUrl string `json:"head_img_url,omitempty"`
	UserName   string `json:"username,omitempty"`
}

type NotifyReq2 struct {
	MhtOrderNo    string `json:"mhtOrderNo"`
	NowPayOrderNo string `json:"nowPayOrderNo"`
	TradeStatus   string `json:"tradeStatus"`
	PayConsumerId string `json:"payConsumerId"`
}

func split(s rune) bool {
	if s == '&' {
		return true
	}
	return false
}

func split2(s rune) bool {
	if s == '=' {
		return true
	}
	return false
}

func URLDecode(str string) (string, string, string, string, string) {
	var (
		mhtOrderNo    string
		nowPayOrderNo string
		tradeStatus   string
		payConsumerId string
		mhtReserved   string
	)
	strAray := strings.FieldsFunc(str, split) //	[widuu hello word]根据n字符分割

	retMap := make([]map[string]string, 0)
	for _, row := range strAray {
		ss := make(map[string]string)
		rowStr := strings.FieldsFunc(row, split2)
		ss[rowStr[0]] = rowStr[1]
		if rowStr[0] == "mhtOrderNo" {
			mhtOrderNo = rowStr[1]
		} else if rowStr[0] == "nowPayOrderNo" {
			nowPayOrderNo = rowStr[1]
		} else if rowStr[0] == "tradeStatus" {
			tradeStatus = rowStr[1]
		} else if rowStr[0] == "payConsumerId" {
			payConsumerId = rowStr[1]
		} else if rowStr[0] == "mhtReserved" {
			mhtReserved = rowStr[1]
		}

		retMap = append(retMap, ss)
	}
	return mhtOrderNo, nowPayOrderNo, tradeStatus, payConsumerId, mhtReserved
}

func NowPayNotify(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})

	body, _ := ioutil.ReadAll(req.Body)

	bodyStr := fmt.Sprintf("%s", body)
	mhtOrderNo, nowPayOrderNo, tradeStatus, payConsumerId, ext := URLDecode(bodyStr)

	model.NowPayNotify(mhtOrderNo, nowPayOrderNo, tradeStatus, payConsumerId, ext)
	ret_value["success"] = "Y"

	r.JSON(http.StatusOK, ret_value)
}

type HuaWeiNotifyReq struct {
	RequestId string `form:"requestId" `
	OrderId   string `form:"orderId"`
	TradeTime string `form:"tradeTime"`
	Result    string `form:"result"`

	ExtReserved string `form:"extReserved"`
}

func HuaWeiNotify(req *http.Request, r render.Render, d HuaWeiNotifyReq) {

	ret_value := make(map[string]interface{})
	//godump.Dump(d)
	ret_value["result"] = model.HuaWeiNotify(d.RequestId, d.OrderId, d.TradeTime, d.Result, d.ExtReserved)

	r.JSON(http.StatusOK, ret_value)
}

type AlipayNotifyReq struct {
	OutTradeNo  string `form:"out_trade_no" `
	TradeNo     string `form:"trade_no" `
	TradeStatus string `form:"trade_status" `
	NotifyTime  string `form:"notify_time" `
	AppId       string `form:"app_id" `
	Body        string `form:"body" `
	BuyerId     string `form:"buyer_id" `
	Charset     string `form:"charset" `
	GmtClose    string `form:"gmt_close" `
	GmtPayment  string `form:"gmt_payment" `
	Notifyid    string `form:"notify_id" `
	NotifyType  string `form:"notify_type" `
	RefundFee   string `form:"refund_fee" `
	SellerId    string `form:"seller_id" `
	Subject     string `form:"subject" `
	TotalAmount string `form:"total_amount" `
	Version     string `form:"version" `
	Sign        string `form:"sign" `
	PassBack    string `form:"passback_params"`
}

func AlipayNotify(req *http.Request, r render.Render, d AlipayNotifyReq) {
	ok, _ := model.AlipayKeyVerify(req)

	if ok {
		flag := model.AlipayNotify(d.OutTradeNo, d.TradeNo, d.NotifyTime, d.TradeStatus, d.PassBack)

		if flag == common.ERR_SUCCESS {
			ret_value := "success"
			r.JSON(http.StatusOK, ret_value)
		} else {
			ret_value := "fail"
			r.JSON(http.StatusOK, ret_value)
		}
	} else {
		ret_value := "fail"
		r.JSON(http.StatusOK, ret_value)
	}
}
