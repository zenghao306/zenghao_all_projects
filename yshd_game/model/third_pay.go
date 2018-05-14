package model

import (
	"fmt"
	"github.com/yshd_game/common"
	"regexp"
	"sort"
	"strconv"
	"time"
	//"github.com/smartwalle/alipay"
	"crypto"
	"encoding/base64"
	"github.com/smartwalle/alipay/encoding"
	"net/http"
	"net/url"
	"strings"
)

func ThirdPaySwitch(id string) bool {

	sql := fmt.Sprintf("SELECT * FROM php_pay_info WHERE channel_id = '%s' && pay_type = 2", id)

	rowArray, _ := orm.Query(sql)
	length := len(rowArray)
	if length >= 1 {
		return true
	}
	return false
}

var (
	nowpay_appid    = "149683301163651"
	nowpay_md5key   = "PWGd0Tn2wgQCGCZboJLx1KKpIWjZj04b" //聚合MD5密钥
	nowpayNotifyUrl = "http://shangtv.cn:3003/third/pay/notify"

	huawei_applicationID = "100018145"
	huawei_merchantId    = "890086000102020465"
	huawei_notifyUrl     = "http://shangtv.cn:3003/huawei/pay/notify"

	alipay_applicationID = "2017071207727765"
	alipay_notifyUrl     = "http://shangtv.cn:3003/alipay/pay/notify"
)

const (
	//NOWPAY_APPID     = "149683301163651"
	NOWPAY_ORDERNAME    = "17玩-钻石充值"
	HUAWEIPAY_ORDERNAME = "17玩-钻石充值"
	ALIPAY_ORDERNAME    = "17玩-钻石充值"
	//NOWPAY_MD5KEY    = "PWGd0Tn2wgQCGCZboJLx1KKpIWjZj04b" //聚合MD5密钥
	NOWPAY_TIMEOUT           = 300
	HUAWEIPAY_ORDE_NOT_EXIST = 3

	HUAWEI_PRIVATE_KEY = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCCpfncHoThYw9zmoYiR3RU9eCGFyTZbou6vkdKtB1kA9fE30ovOSBRE0xdg6F6fP5y/FN0NVNelFIwv/YkwYi5brTIY2bFmVQMWmWFSB4fPtJdBmolnFZtoniOQp2I2bLSE5RnlpS4NPIGME9AGsR6EDWvJwfcCc1QC32ieWZCZrim7IMN8K3cOE4u8pAsKZanQP2BmVsPkOpZJX/fsB9Mjus6tZoYiH3JbC0frlt0aSviMJkVxEFHPA7bnnVTdfM/LTizQ4VmmiamWHTvKtXlXgTMoRXrY66xan/UHN3/S6nif4SWPjv+oEziDTQh8qOG/SWANl+cbdgBd2YKo1TjAgMBAAECggEACmWqCj31hjhLcQBBn3W/SMmeeh1aZeFZxl1BMC1AT4bMw5KfhT2PGFSoVaLVlXlgCIeTHqLlxReZqN6F+KvcNSGdynq6oYwPt8Hz5VT1bLgjppqNlPuplyUAYhXkEpF8nSJIw6uknzo7bommrOvUagBjPVKmWfj/uViIwYbWv/7ufyvnt2NKv7lTXwoU79OpflkSuXve5/ndY1qi+IcQq6gP08v/LSpleCh1N6cxN2l53eZyUrH94poYBWxO96RyOVh0kpOWcUeLN8Fb5gcvBzsmJ7gHiOVFS8EQF+2HqR0F7yD+7IQvdJPkTvQ/TYHsekfYoJHHkTnFrLHM2RgJfQKBgQDKDxEcITGH8UJptVXWPIoEmhusOicvr37ksfcGe6CNO6Tokkt+9YS6PiN5PBw5fGctj9tPS5rSDvEB3xpUmfZfEols1KvEfpsfeq84q12cDxvUTUmNG/bAjJ0mIuzObN6d9AaxnadFxYGsgeyAXYiHA8ZDrvJ/91/ESPzYkmGZPQKBgQClhqIZ+HPxqqf9ozBl6G12zxmInnfRzbx/w7LKOZOS0OWFG1RvwBsfFXwf8WR9lG42A5YX42zGPfak2lZQF/Dr371/AG28fsKByxuxXIL7oGd5hF0w8+ehyqpcUZCt4wnFXXrbx/dXLTMZq7cPUmFTc9wOSsc90CJtbJIit1xInwKBgQCn1Hnutndwpej26oKgjupIxkRD+o1/4zHv/Q3kmZ0Skj74WkNQ8ddL5r8KPO5opTcCNiIALBktbvGqD0jMipGECF6TQdZmQI7SR9HwrQ30yOvhnEyCY37CEkmZWpr9HpqN8hn5P6ynnFSIF+Z7/LShCHaO02pi9fLak5FIxdVBNQKBgCwcceyL2pInkXc1Wvt9clZ8IGjZGkNOGZLO20GrEFnK06+iNeFGB7RD7H0yvj39JVW9pO0ezWkTMDyEdwYHK0wgpwZGXfgiq8EdGHcumTVnbMvka1VVWoAyWy4ZCn5ch8kR/WC+rHqN8BVCPpmVsZwkvxsL8IbXhaU1jsgzpyMvAoGBALEdxQTQfFuhd+KA9xqog+o9gXaVIW0iwOre+Ahyh16D3wRQV/qFdLET216xgzExut8dXtO+BVnguPvEGKkN7VWT4/rfRk9nYDwQwSdx+nDzNLo1oSCUV6db45G5oe8kzSuDYOUtIk86I1AEIp6d6ztY4ZMvg9Jo6MjQZI9+JsAL
-----END PRIVATE KEY-----`

	ALIPAY_PRIVATE_KEY = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCIKMUxm9+NoPbVEzaM8wdDBcORi9gsVLF0ULokwxtFWY4PCTZnOrcd0tqCM+yQB8XSV0kU8T5pOnRCl6vaIKOHREYpbCqaSTK4ff2OSiTJ6EmgsQAROc9iTOLhSjLciEYLkwZyXL5JtgFaCbXcQthKb8kkKRp5wCBNaiOMnJUb7ZvLWV0HrzY1HDd/L3Z2XxRYhWNIuYWCsja4nfqJ8RYOOx3NbqAGbSmz0vRLmo58D6VvBEGbU7hBHpPFqIaKGEZjqJsUxZDOF2kaKx7SwEGJpBvjZThP3AWS4A4aT9XhTixHz+re9n52BhR/zqVESSU3Gk9a9BoW+qkg6DasxOXpAgMBAAECggEARfeVZWmw1emKDXIjDQjxiVpT5d3Txuv6iEfXb36m69saKdXVE/TTFks8p72g6V5lJDJgRpe1N4OnLHUeBSfSgHbwCucfeUr1+mIbwluNTgfElgN+gluPmvbhe12Sh0qrm9UAchIAYoZZaXgl8LqUxKNu29sXVMsKjl1lSSNJaCDQlYWsAYzI+CaYno5wVDJn5WYWU6NF1pABVNjk87g4y9ZStFqH2suW1PyiQo34jAUaO0u2u8Qbq0j4J+kR1aawQW+/14WlunilmMFZg94P+q/7m/4BqG/5XW99DO6l7Wng7I0QMuQft/ng3jnV2y9sqsGJTaBFVU/Rlq1lxQWORQKBgQC7ljTigNAZaeKvceNRRW9Q2WFKUhI3bilRTS4tuN3Y7S4OPUhp/67SJNCO7IPgtTUQ9dSTf8Nwc9TSyZTFtm6y+OSvpilszXs8p2OLp2f6PeT8zzWMmcZU0QF9wzHKF7QBxA9JbyXUcW4M8SwxjaKtJ09ocBWYmzjbfsCMLoihrwKBgQC50RpnBuhGrealInW+BFzWtRrm0VFAPxrk0f0R+q9lJMRuREC5U6dGH6whfr1xSmr4kY63JmO3Z/A38YY+R1EzNbx/k/+k3zKzVg+g7yiKnDXH0ekZkv+YEltFvgnAMksVcoEkyxH3D0b3iiy3pikDSq5L8R0DqCgwuXYuB0pP5wKBgQCLqoOHnTbTpSW1UQtZ6GPAA4nPhxmvEaNLuDZIrprmt3kR+wjeexMTvXtW3rw141U2YoI6q+a85FEx/Ap7xp/XOz8xlHrFWpyGBW81fJgLFmhW3oRVQe0MG22L0HhSqqFIq0xidZHqAeZZVnt8DaNwXpNeBA3gSLnlmxMLjF4IWwKBgHNx25J48yS+dFbScw6MTVXEDSOslmtxCXdyk2VxNzmCv1u2ofPCamGh2eKxiGdzkcQ/Qsi9XCSdudw3/WyCCIvlbehhfenkFe7foDQfgjOj27H603TlJFFJzlUlPY+gb4+ypVPDqrSxVCkFOsUawc5evq1F7v3PorCq+SadtuivAoGAd7ISEuFRkpgIJEf1QzanS4HPtCq67rYH7QzTZJYnbfQ2FSont0LjMADnK4hdkLs2gaJeLBlWyqsGbDxE5UoDR1QAOD4RH7Sb5/K8ISNbis0rZ+AOp41LmmTRDl/MmW1J4Ol8xbti3+EmlovUwZbosRh3gbpzfaiN/6fJLPCUskA=
-----END PRIVATE KEY-----`

	ALIPAY_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwov2pIZv7apU9OaGG3NfyDTSVek622bIH+MvepDhrVlC8GUM55aulsw817KFyfH3Wc1jxcVShxErQmUBiUiOxs8VEMdJd1CYvsmLDu8mkODBChq2yanoSUJ0x03OuDjClLT6hp2MJZ8rq0ae18K9pLQvzPWmJSwToDRtbZw62YkwYHub61Pu/IHpGf/MkooVo+zPQenSt/DKUj9xxOmQa16Ify5u8rlqAe3JqMuaChc2aM7vQZvCwaHKSVor5chc6T1MBNKzG9AnQ4MiH/jwkn/5K/dCUvQ1C9IdXre9r37YPxc0RJYAtYd/PCjH7+L0Np7Z+d83htANudcu9tkEmwIDAQAB
-----END PUBLIC KEY-----`
)

func InitNowPayKey() {
	nowpay_appid = common.Cfg.MustValue("now_pay", "nowpay_appid")
	nowpay_md5key = common.Cfg.MustValue("now_pay", "nowpay_md5key")
	nowpayNotifyUrl = common.Cfg.MustValue("now_pay", "nowpay_notify_url")
	//fmt.Printf("\nnowpay_appid:%s",nowpay_appid)
	//fmt.Printf("\nowpay_md5key:%s",nowpay_md5key)
}

func InitHuaWeiPayKey() {
	huawei_applicationID = common.Cfg.MustValue("huawei_pay", "huawei_applicationID")
	huawei_merchantId = common.Cfg.MustValue("huawei_pay", "huawei_merchantId")
	huawei_notifyUrl = common.Cfg.MustValue("huawei_pay", "huawei_notify_url")
}

func InitApliPayKey() {
	alipay_applicationID = common.Cfg.MustValue("ali_pay", "alipay_applicationID")
	alipay_notifyUrl = common.Cfg.MustValue("ali_pay", "alipay_notifyUrl")
}

type Params map[string]string

func (p Params) SetString(k, s string) {
	p[k] = s
}

func (p Params) GetString(k string) string {
	s, _ := p[k]
	return s
}

func (p Params) SetInt64(k string, i int64) {
	p[k] = strconv.FormatInt(i, 10)
}

func (p Params) GetInt64(k string) int64 {
	i, _ := strconv.ParseInt(p.GetString(k), 10, 64)
	return i
}

// 三方支付[聚合]签名生成
// 第一步：对参与MD5签名的字段按字典升序排序后，分别取值后并排除值为空的字段键值对，最后组成key1=value1&key2=value2....keyn=valuen" 表单字符串"。
// 第二步：对MD5密钥进行加密得到"密钥MD5值"。
// 第三步：最后对第一步中得到的表单字符串&第二步得到的密钥MD5值做MD5签名
func NowPayCalcSign(mReq map[string]interface{}) (string, string) {
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	var signStr string
	for _, k := range sorted_keys {
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStr = signStr + k + "=" + value + "&"
		}
	}

	//fmt.Printf("\n第一次MD5:%s",common.Md5(nowpay_md5key))
	beforeSign := signStr
	signStr += common.Md5(nowpay_md5key)
	//fmt.Printf("\n第二次beforeSign:%s",beforeSign)
	//fmt.Printf("\n第二次MD5前:%s",signStr)

	return beforeSign, common.Md5(signStr)
}

func HuaWeiCalcSign(mReq map[string]interface{}) string {
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	var signStr string

	index := 0
	length := len(sorted_keys)

	for _, k := range sorted_keys {
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStr = signStr + k + "=" + value // + "&"
		}

		index++
		if index < length {
			signStr = signStr + "&"
		}
	}
	fmt.Printf("签名前:%s", signStr)

	mdStr := common.Rsa256Signal(signStr, HUAWEI_PRIVATE_KEY)

	reg := regexp.MustCompile("_")
	mdStr = reg.ReplaceAllString(mdStr, "/") //将_替换为/字符

	reg2 := regexp.MustCompile("-")
	mdStr = reg2.ReplaceAllString(mdStr, "+") //将_替换为/字符

	return mdStr
}

func AlipayCalcSign(mReq map[string]interface{}, key string) string {
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	var signStr string

	index := 0
	length := len(sorted_keys)

	for _, k := range sorted_keys {
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStr = signStr + k + "=" + value // + "&"
		}

		index++
		if index < length {
			signStr = signStr + "&"
		}
	}

	mdStr := common.Rsa256Signal(signStr, key)

	reg := regexp.MustCompile("_")
	mdStr = reg.ReplaceAllString(mdStr, "/") //将_替换为/字符

	reg2 := regexp.MustCompile("-")
	mdStr = reg2.ReplaceAllString(mdStr, "+") //将-替换为+字符

	return mdStr
}

func AlipayKeyVerify(req *http.Request) (ok bool, err error) {

	key := []byte(ALIPAY_PUBLIC_KEY)
	sign, err := base64.StdEncoding.DecodeString(req.PostForm.Get("sign"))
	if err != nil {
		return false, err
	}

	var keys = make([]string, 0, 0)
	for key, value := range req.PostForm {
		if key == "sign" || key == "sign_type" {
			continue
		}
		if len(value) > 0 {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(req.PostForm.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var s = strings.Join(pList, "&")

	err = encoding.VerifyPKCS1v15([]byte(s), sign, key, crypto.SHA256)
	if err != nil {
		return false, err
	}
	return true, nil
}

func NowPayPrepay(itemid, channelID string, tradeType, uid, anchorID int, ip string, device string) (int, string, string) {

	//tradeID := fmt.Sprintf("JH%d_%d", uid, time.Now().Unix())
	randnum := common.RadnomRange(100000, 999999)
	tradeID := fmt.Sprintf("%d%d", time.Now().Unix(), randnum)

	nowtime := time.Now()
	expiretime := nowtime.Add(NOWPAY_TIMEOUT * time.Second)

	item, has := GetAndroidItem(itemid)
	if has == false {
		return common.ERR_CONFGI_ITEM, "", ""
	}

	if tradeType != common.THIRD_PAY_UNIONPAY && tradeType != common.THIRD_PAY_ALIPAY && tradeType != common.THIRD_PAY_WEIXIN {
		return common.ERR_UNKNOWN, "", ""
	}

	_, err := orm.Insert(&Trade{
		TradeId:    tradeID,
		TimeStart:  nowtime,
		TimeExpire: expiretime,
		Uid:        uid,
		AnchorId:   anchorID,
		Money:      item.Money,
		Diamond:    item.Diamond,
		Status:     common.TRADE_PRE_CREATE,
		TradeType:  tradeType,
		ChannelId:  channelID,
		Device:     device,
	})

	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN, "", ""
	}

	if nowpayNotifyUrl == "" {
		nowpayNotifyUrl = "http://shangtv.cn:3003/third/pay/notify"
	}

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appId"] = nowpay_appid
	m["mhtOrderNo"] = tradeID
	m["mhtOrderName"] = NOWPAY_ORDERNAME
	m["mhtOrderType"] = "01"
	m["mhtCurrencyType"] = "156" //人民币
	m["mhtOrderAmt"] = item.Money
	m["mhtOrderDetail"] = item.Describe
	m["mhtOrderTimeOut"] = NOWPAY_TIMEOUT
	m["mhtOrderStartTime"] = nowtime.Format("20060102150405")
	m["notifyUrl"] = nowpayNotifyUrl //"http://shangtv.cn:3003/third/pay/notify"
	//m["notifyUrl"] = "http://shangtv.kmdns.net:3003/third/pay/notify"
	m["mhtCharset"] = nowtime.Format("UTF-8")
	m["payChannelType"] = strconv.Itoa(tradeType)
	//m["mhtSignType"] = "MD5"
	m["mhtReserved"] = itemid

	beforeSign, sign := NowPayCalcSign(m) //这个是计算JuHeCalcSign签名的函数上面已贴出

	return common.ERR_SUCCESS, beforeSign, sign
}

func HuaWeiPrepay(itemid, channelID string, uid, anchorID int, ip string, device string) (int, string, []map[string]string) {
	retMap := make([]map[string]string, 0)

	randnum := common.RadnomRange(100000, 999999)
	tradeID := fmt.Sprintf("%d%d", time.Now().Unix(), randnum)

	nowtime := time.Now()
	expiretime := nowtime.Add(NOWPAY_TIMEOUT * time.Second)

	item, has := GetAndroidItem(itemid)
	if has == false {
		return common.ERR_CONFGI_ITEM, "", retMap
	}

	_, err := orm.Insert(&Trade{
		TradeId:    tradeID,
		TimeStart:  nowtime,
		TimeExpire: expiretime,
		Uid:        uid,
		AnchorId:   anchorID,
		Money:      item.Money,
		Diamond:    item.Diamond,
		Status:     common.TRADE_PRE_CREATE,
		TradeType:  common.HUAWEI_PAY,
		ChannelId:  channelID,
		Device:     device,
	})

	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN, "", retMap
	}

	if huawei_notifyUrl == "" {
		huawei_notifyUrl = "http://shangtv.cn:3003/huawei/pay/notify"
	}

	//将货币单位转换为"6.00"这样的字符串
	MoneyStr := fmt.Sprintf("%d.%d%d", item.Money/100, item.Money%100/10, item.Money%10)

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["productName"] = HUAWEIPAY_ORDERNAME
	m["productDesc"] = item.Describe
	m["applicationID"] = huawei_applicationID
	m["requestId"] = tradeID
	m["amount"] = MoneyStr //item.Money
	m["merchantId"] = huawei_merchantId
	m["sdkChannel"] = 1
	m["urlver"] = "2"
	m["url"] = huawei_notifyUrl

	ss := make(map[string]string)
	ss["productName"] = HUAWEIPAY_ORDERNAME
	ss["requestId"] = tradeID
	ss["url"] = huawei_notifyUrl
	ss["extReserved"] = itemid
	retMap = append(retMap, ss)

	sign := HuaWeiCalcSign(m) //这个是计算HuaWeiCalcSign签名的函数上面已贴出

	return common.ERR_SUCCESS, sign, retMap //, tradeID, HUAWEIPAY_ORDERNAME
}

func HuaWeiNotify(requestId, orderId, tradeTime, result, extReserved string) int {
	trade := &Trade{}
	has, err := orm.Where("trade_id=? and status=?", requestId, common.TRADE_PRE_CREATE).Get(trade)
	if err != nil {
		common.Log.Err("get real auth error: , %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if !has {
		return HUAWEIPAY_ORDE_NOT_EXIST
	}
	if result == "0" { //华为支付成功
		trade.Status = common.TRADE_HAVE_SUCCESS
		trade.ThirdTradeId = orderId
		trade.PurchaseTime = int(time.Now().Unix())
		trade.PurchaseDateOri = tradeTime
		orm.Where("trade_id=?", requestId).Update(trade)

		//以下是对用户的砖石进行实际添加
		user, ret := GetUserByUid(trade.Uid)
		if ret != common.ERR_SUCCESS {
			return common.ERR_UNKNOWN
		}

		active := GetChargeActive(extReserved)
		now_time := time.Now().Unix()

		session := orm.NewSession()
		defer session.Close()
		err = session.Begin()
		if err != nil {
			common.Log.Errf("orm is error:  %s", err.Error())
			session.Rollback()
			return common.ERR_UNKNOWN
		}
		ret = user.AddMoney(session, common.MONEY_TYPE_DIAMOND, int64(trade.Diamond), true)
		if ret != common.ERR_SUCCESS {
			session.Rollback()
			return common.ERR_UNKNOWN
		}
		aff, err := user.AddUserExp(session, trade.Diamond/10, true) //加上经验值
		if err != nil || aff == 0 {
			session.Rollback()
			return common.ERR_UNKNOWN
		}
		ret = user.SetNewPay(session, trade.Diamond)
		if ret != common.ERR_SUCCESS {
			session.Rollback()
			return common.ERR_UNKNOWN
		}

		if active != nil {
			if active.Status == 1 && active.BeginTime < now_time && active.FinishTime > now_time {
				ret := user.AddMoney(session, int32(active.MoneyType), active.ExtraNum, true)
				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return common.ERR_UNKNOWN
				}
			}
		}

		err = session.Commit()
		if err != nil {
			return common.ERR_UNKNOWN
		}
	}
	return common.ERR_SUCCESS
}

func NowPayNotify(mhtOrderNo, nowPayOrderNo, TradeStatus, payConsumerId, extReserved string) {
	trade := &Trade{}
	has, err := orm.Where("trade_id=? and status=?", mhtOrderNo, common.TRADE_PRE_CREATE).Get(trade)
	if err != nil {
		common.Log.Err("get real auth error: , %s", err.Error())
		return
	}
	if !has {
		return
	}
	if TradeStatus == "A001" { //第三方支付成功
		trade.Status = common.TRADE_HAVE_SUCCESS
		trade.ThirdTradeId = nowPayOrderNo
		trade.PurchaseTime = int(time.Now().Unix())
		if payConsumerId != "" {
			trade.Extend_1 = payConsumerId
		}
		orm.Where("trade_id=?", mhtOrderNo).Update(trade)

		//以下是对用户的砖石进行实际添加
		user, ret := GetUserByUid(trade.Uid)
		if ret != common.ERR_SUCCESS {
			return
		}

		active := GetChargeActive(extReserved)
		now_time := time.Now().Unix()

		session := orm.NewSession()
		defer session.Close()
		err = session.Begin()
		if err != nil {
			common.Log.Errf("orm is error:  %s", err.Error())
			session.Rollback()
			return
		}
		ret = user.AddMoney(session, common.MONEY_TYPE_DIAMOND, int64(trade.Diamond), true)
		if ret != common.ERR_SUCCESS {
			session.Rollback()
			return
		}
		aff, err := user.AddUserExp(session, trade.Diamond/10, true) //加上经验值
		if err != nil || aff == 0 {
			session.Rollback()
			return
		}
		//user.SetNewPay(session,trade.Diamond)       //首充有礼
		ret = user.SetNewPay(session, trade.Diamond)
		if ret != common.ERR_SUCCESS {
			session.Rollback()
			return
		}

		if active != nil {
			if active.Status == 1 && active.BeginTime < now_time && active.FinishTime > now_time {
				ret := user.AddMoney(session, int32(active.MoneyType), active.ExtraNum, true)
				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return
				}
			}
		}

		err = session.Commit()
		if err != nil {
			common.Log.Err("NowPayNotify 事务 error: , %s", err.Error())
			common.Log.Err("NowPayNotify 事务 error: , uid=%d", user.Uid)
			return
		}

		return
	}
}

func AlipayPrepay(itemid, channelID string, uid, anchorID int, ip string, device string) (int, string, []map[string]string) {
	retMap := make([]map[string]string, 0)
	randnum := common.RadnomRange(100000, 999999)
	tradeID := fmt.Sprintf("%d%d", time.Now().Unix(), randnum)

	nowtime := time.Now()
	expiretime := nowtime.Add(NOWPAY_TIMEOUT * time.Second)

	item, has := GetAndroidItem(itemid)
	if has == false {
		return common.ERR_CONFGI_ITEM, "", retMap
	}

	_, err := orm.Insert(&Trade{
		TradeId:    tradeID,
		TimeStart:  nowtime,
		TimeExpire: expiretime,
		Uid:        uid,
		AnchorId:   anchorID,
		Money:      item.Money,
		Diamond:    item.Diamond,
		Status:     common.TRADE_PRE_CREATE,
		TradeType:  common.ALIPAY_PAY,
		ChannelId:  channelID,
		Device:     device,
	})

	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN, "", retMap
	}

	if alipay_notifyUrl == "" {
		alipay_notifyUrl = "http://t1.shangtv.cn:3003/alipay/pay/notify"
	}

	//将货币单位转换为"6.00"这样的字符串
	MoneyStr := fmt.Sprintf("%d.%d%d", item.Money/100, item.Money%100/10, item.Money%10)

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)

	eitemid := url.QueryEscape(itemid)

	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	m["app_id"] = alipay_applicationID
	bizContent := `{"timeout_express":"30m","passback_params":"` + eitemid + `","seller_id":"","product_code":"QUICK_MSECURITY_PAY","total_amount":"` + MoneyStr + `","subject":"` + ALIPAY_ORDERNAME + `","body":"17玩钻石充值","out_trade_no":"` + tradeID + `"}`
	m["biz_content"] = bizContent
	m["charset"] = "utf-8"
	m["format"] = "json"
	m["method"] = "alipay.trade.app.pay"
	m["notify_url"] = alipay_notifyUrl
	m["sign_type"] = "RSA2"
	//m["subject"] = ALIPAY_ORDERNAME
	m["timestamp"] = timeStamp
	m["version"] = "1.0"

	ss := make(map[string]string)
	ss["app_id"] = alipay_applicationID
	ss["biz_content"] = bizContent
	ss["charset"] = "utf-8"
	ss["format"] = "json"
	ss["method"] = "alipay.trade.app.pay"
	ss["notify_url"] = alipay_notifyUrl
	ss["sign_type"] = "RSA2"
	//ss["subject"] = ALIPAY_ORDERNAME
	ss["timestamp"] = timeStamp
	ss["version"] = "1.0"
	retMap = append(retMap, ss)

	sign := AlipayCalcSign(m, ALIPAY_PRIVATE_KEY) //这个是计算HuaWeiCalcSign签名的函数上面已贴出

	return common.ERR_SUCCESS, sign, retMap //, tradeID, HUAWEIPAY_ORDERNAME
}

/*
func AlipayPrepay(itemid, channelID string, uid, anchorID int, ip string, device string) (int, string, []map[string]string) {
	retMap := make([]map[string]string, 0)
	randnum := common.RadnomRange(100000, 999999)
	tradeID := fmt.Sprintf("%d%d", time.Now().Unix(), randnum)

	nowtime := time.Now()
	expiretime := nowtime.Add(NOWPAY_TIMEOUT * time.Second)

	item, has := GetAndroidItem(itemid)
	if has == false {
		return common.ERR_CONFGI_ITEM, "", retMap
	}

	_, err := orm.Insert(&Trade{
		TradeId:    tradeID,
		TimeStart:  nowtime,
		TimeExpire: expiretime,
		Uid:        uid,
		AnchorId:   anchorID,
		Money:      item.Money,
		Diamond:    item.Diamond,
		Status:     common.TRADE_PRE_CREATE,
		TradeType:  common.ALIPAY_PAY,
		ChannelId:  channelID,
		Device:     device,
	})

	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN, "", retMap
	}

	if alipay_notifyUrl == "" {
		alipay_notifyUrl = "http://t1.shangtv.cn:3003/alipay/pay/notify"
	}

	//将货币单位转换为"6.00"这样的字符串
	MoneyStr := fmt.Sprintf("%d.%d%d", item.Money/100, item.Money%100/10, item.Money%10)

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)

	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	m["app_id"] = alipay_applicationID
	bizContent := `{"timeout_express":"30m","seller_id":"","product_code":"QUICK_MSECURITY_PAY","total_amount":"` + MoneyStr + `","subject":"` + ALIPAY_ORDERNAME + `","body":"17玩钻石充值","out_trade_no":"` + tradeID + `"}`
	m["biz_content"] = bizContent
	m["charset"] = "utf-8"
	m["format"] = "json"
	m["method"] = "alipay.trade.app.pay"
	m["notify_url"] = alipay_notifyUrl
	m["sign_type"] = "RSA2"
	//m["subject"] = ALIPAY_ORDERNAME
	m["timestamp"] = timeStamp
	m["version"] = "1.0"

	ss := make(map[string]string)
	ss["app_id"] = alipay_applicationID
	ss["biz_content"] = bizContent
	ss["charset"] = "utf-8"
	ss["format"] = "json"
	ss["method"] = "alipay.trade.app.pay"
	ss["notify_url"] = alipay_notifyUrl
	ss["sign_type"] = "RSA2"
	//ss["subject"] = ALIPAY_ORDERNAME
	ss["timestamp"] = timeStamp
	ss["version"] = "1.0"
	retMap = append(retMap, ss)

	sign := AlipayCalcSign(m, ALIPAY_PRIVATE_KEY) //这个是计算HuaWeiCalcSign签名的函数上面已贴出

	return common.ERR_SUCCESS, sign, retMap //, tradeID, HUAWEIPAY_ORDERNAME
}
*/

func AlipayNotify(outTradeNo, tradeNo, notifyTime, tradeStatus, passBack string) int {
	trade := &Trade{}
	has, err := orm.Where("trade_id=? and status=?", outTradeNo, common.TRADE_PRE_CREATE).Get(trade)
	if err != nil {
		common.Log.Err("get real auth error: , %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if !has {
		return common.ERR_UNKNOWN
	}
	if tradeStatus == "TRADE_SUCCESS" { //支付宝付款成功
		trade.Status = common.TRADE_HAVE_SUCCESS
		trade.ThirdTradeId = tradeNo
		trade.PurchaseTime = int(time.Now().Unix())
		trade.PurchaseDateOri = notifyTime
		orm.Where("trade_id=?", outTradeNo).Update(trade)

		//以下是对用户的砖石进行实际添加
		user, ret := GetUserByUid(trade.Uid)
		if ret != common.ERR_SUCCESS {
			return common.ERR_UNKNOWN
		}
		active := GetChargeActive(passBack)
		now_time := time.Now().Unix()

		session := orm.NewSession()
		defer session.Close()
		err = session.Begin()
		if err != nil {
			common.Log.Errf("orm is error:  %s", err.Error())
			session.Rollback()
			return common.ERR_UNKNOWN
		}
		ret = user.AddMoney(session, common.MONEY_TYPE_DIAMOND, int64(trade.Diamond), true)
		if ret != common.ERR_SUCCESS {
			session.Rollback()
			return common.ERR_UNKNOWN
		}
		aff, err := user.AddUserExp(session, trade.Diamond/10, true) //加上经验值
		if err != nil || aff == 0 {
			session.Rollback()
			return common.ERR_UNKNOWN
		}
		ret = user.SetNewPay(session, trade.Diamond)
		if ret != common.ERR_SUCCESS {
			session.Rollback()
			return common.ERR_UNKNOWN
		}

		if active != nil {
			if active.Status == 1 && active.BeginTime < now_time && active.FinishTime > now_time {
				ret := user.AddMoney(session, int32(active.MoneyType), active.ExtraNum, true)
				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return common.ERR_UNKNOWN
				}
			}
		}

		err = session.Commit()
		if err != nil {
			return common.ERR_UNKNOWN
		}
	}
	return common.ERR_SUCCESS
}
