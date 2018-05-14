package controller

import (
	"encoding/xml"
	//"github.com/liudng/godump"
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"io/ioutil"
	"net/http"
	"time"
)

type WinXinPayReq struct {
	Uid       int    `form:"uid" binding:"required"`
	Token     string `form:"token" binding:"required"`
	Itemid    string `form:"itemid"`
	ChannelId string `form:"channel_id"`
	Device    string `form:"device"`
	AnchorId  int    `form:"anchor_id"`
}

type ApplePayReq struct {
	Uid       int    `form:"uid" binding:"required"`
	Token     string `form:"token" binding:"required"`
	Receipt   string `form:"receipt" binding:"required" `
	ChannelId string `form:"channel_id"`
	Device    string `form:"device"`
	AnchorId  int    `form:"anchor_id"`
}

type GooglePayReq struct {
	Uid          int    `form:"uid" binding:"required"`
	Token        string `form:"token" binding:"required"`
	ProductId    string `form:"productid" binding:"required" `
	SigntureData string `form:"signtureData"`
	//Signture      string `form:"signture" binding:"required" `
	PurchaseToken string `form:"purchase_token" binding:"required" `
	OrderId       string `form:"order_id"`
	ChannelId     string `form:"channel_id"`
	Device        string `form:"device"`
}
type ScoreExchangeReq struct {
	Uid     int    `form:"uid" binding:"required"`
	Token   string `form:"token" binding:"required"`
	Daimond int    `form:"diamond" binding:"required" `
}

//curl -d 'uid=3&token=ffc210d5011be88c90aaf6e313276b8a&itemid=gold_600' 'http://t1.shangtv.cn:3003/pay/weixin_prepay'
//curl -d 'uid=3&token=ffc210d5011be88c90aaf6e313276b8a&itemid=gold_600' 'http://192.168.1.12:3003/pay/weixin_prepay'
func WinXinPrepay(req *http.Request, r render.Render, d WinXinPayReq) {
	ret_value := make(map[string]interface{})
	ip := common.GetRemoteIp(req)
	weixin := &model.RetWeiXinClient{}
	ret_value[ServerTag], weixin = model.GetWinXinPrepayId2(d.Itemid, d.Uid, ip, d.ChannelId, d.Device, d.AnchorId)
	ret_value["weixin"] = weixin
	r.JSON(http.StatusOK, ret_value)
}

func WinXinH5Prepay(req *http.Request, r render.Render, d WinXinPayReq) {
	ret_value := make(map[string]interface{})
	ip := common.GetRemoteIp(req)
	weixin := &model.RetWeiXinClient{}
	ret_value[ServerTag], weixin = model.GetWeiXinH5PrepayId(d.Itemid, d.Uid, ip, d.ChannelId, d.Device)
	ret_value["weixin"] = weixin
	r.JSON(http.StatusOK, ret_value)
}

func WinXinPayCallBack(req *http.Request, r render.Render) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		common.Log.Err("read weixin resp err, is:%s", err.Error())
		return
	}
	var s model.WinXinPayCallBackReq

	xml.Unmarshal(body, &s)
	//godump.Dump(s)
	if s.Return_code == "SUCCESS" {
		ret := model.ProgressPayStatueCallBack(s)
		if ret != common.ERR_SUCCESS {
			common.Log.Errf("%v", s)
		}
	}
	r.XML(200, "SUCCESS")
}

func ApplePay(req *http.Request, r render.Render, d ApplePayReq) {
	ret_value := make(map[string]interface{})

	//godump.Dump(d)
	ret := model.AppleAuth(d.Receipt, d.Uid, 0, d.ChannelId, d.Device, d.AnchorId)

	ret_value[ServerTag] = ret

	user, _ := model.GetUserByUid(d.Uid)
	if ret == common.ERR_SUCCESS {
		//user.SetNewPay()
	}
	ret_value["diamond"] = user.Diamond
	r.JSON(http.StatusOK, ret_value)
}

//{"orderId":"GPA.1371-9891-0458-21763","packageName":"com.yunshang.enabc","productId":"gold_600","purchaseTime":1484113992263,"purchaseState":0,"developerPayload":"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuSoKVQ","purchaseToken":"dikongkjdppmjlaaedbpljee.AO-J1OzCp5t_YiMSdrBVtIR3SgcJaO5FxLtZ9I3aTDVVDhP9F3f_sef6dXEOQI38CdcwuN7uYtoLYf9vyOuk461yilOeDw9G5VoE2OUGprSI8feWaHsWIOc"}
//http://localhost:3000/pay/google_pay?uid=1&token=fec1f362b900b1d0b3e4081bb047fabb&productid=gold_600&purchase_token=heoaapbghhidcgibhknmjdpo.AO-J1Ozlyj7p6BaaY2G7c831A0qD1O5yJtLxMiyaNXue7_339EOONdzAqU7wuwFfaLkhFBugz0GaRSzilojMzbon7dihr_J67tgGCWRFAxJqkiige6vTjUg&order_id=GPA.1321-9761-7773-13676

//{"orderId":"GPA.1321-9761-7773-13676","packageName":"com.yunshang.enabc","productId":"gold_600","purchaseTime":1484122489916,"purchaseState":0,"developerPayload":"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuSoKVQ","purchaseToken":"heoaapbghhidcgibhknmjdpo.AO-J1Ozlyj7p6BaaY2G7c831A0qD1O5yJtLxMiyaNXue7_339EOONdzAqU7wuwFfaLkhFBugz0GaRSzilojMzbon7dihr_J67tgGCWRFAxJqkiige6vTjUg"}
func GooglePay(req *http.Request, r render.Render, d GooglePayReq) {
	ret_value := make(map[string]interface{})

	common.Log.Debugf("google pay uid=? ,productid=?, purchase=? , order_id=? ,time=?", d.Uid, d.ProductId, d.PurchaseToken, d.OrderId, time.Now().Unix())
	ret_value[ServerTag] = model.GooglePay(d.Uid, d.ProductId, d.SigntureData, d.PurchaseToken, d.OrderId, d.ChannelId, d.Device)
	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'uid=6&token=5aba94937c32e3b280083f3aa7268896&diamond=100' 'shangtv.cn:3003/pay/score_exchange'
func ScoreExchangeController(req *http.Request, r render.Render, d ScoreExchangeReq) {
	ret_value := make(map[string]interface{})
	u, _ := model.GetUserByUid(d.Uid)

	ret_value[ServerTag] = u.ExchangeScore(d.Daimond)
	ret_value["diamond"] = u.Diamond
	ret_value["score"] = u.Score
	r.JSON(http.StatusOK, ret_value)
}
