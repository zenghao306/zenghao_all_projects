package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	GooglePublicKey = "123"
	PackageName     = "com.yunshang.enabc"
)

func DoCheck() {
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	a := b64.EncodeToString([]byte(GooglePublicKey))
	fmt.Printf("a=%s", a)
}

type GooglePlayRet struct {
	ConsumptionState   int    `json:"consumptionState	"`
	DeveloperPayload   string `json:"developerPayload"`
	Kind               string `json:"kind"`
	PurchaseState      int    `json:"purchaseState"`
	PurchaseTimeMillis uint64 `json:"purchaseTimeMillis"`
}

func GooglePay(uid int, productId, signData, purchase_token, order_id string, channel_id, device string) int {

	url := fmt.Sprintf("https://www.googleapis.com/androidpublisher/v2/applications/%s/purchases/products/%s/tokens/%s", PackageName, productId, purchase_token)
	resp, err := http.Get(url)
	if err != nil {
		common.Log.Errf("google check pay err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		common.Log.Errf("google read check pay err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	//godump.Dump(body)
	var v GooglePlayRet
	json.Unmarshal(body, &v)

	//godump.Dump(string(body))
	//godump.Dump(v)
	item, has := GetAndroidItem(productId)
	if !has {
		return common.ERR_CONFGI_ITEM
	}
	var statue int
	if v.PurchaseState == 0 {
		statue = common.TRADE_HAVE_SUCCESS
	} else if v.PurchaseState == 1 {
		statue = common.TRADE_CANCELED
	}
	randnum := common.RadnomRange(100000, 999999)
	tradeno := fmt.Sprintf("%d%d", time.Now().Unix(), randnum)

	session := orm.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		session.Rollback()
		return common.ERR_UNKNOWN
	}

	_, err = session.Insert(&Trade{
		TradeId:         tradeno,
		TimeStart:       time.Now(),
		Money:           item.Money,
		Uid:             uid,
		Diamond:         item.Diamond,
		ThirdTradeId:    order_id,
		PurchaseTime:    int(v.PurchaseTimeMillis / 1000),
		Status:          statue,
		TradeType:       common.TRADE_TYPE_GOOGLE,
		PurchaseDateOri: string(v.PurchaseTimeMillis),
		ChannelId:       channel_id,
		Device:          device,
	})
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	user, ret := GetUserByUid(uid)
	if ret == common.ERR_SUCCESS {
		ret := user.AddMoney(session, common.MONEY_TYPE_DIAMOND, int64(item.Diamond), true)
		if ret != common.ERR_SUCCESS {
			session.Rollback()
			return ret
		}
		aff, err := user.AddUserExp(session, item.Diamond, true)
		if err != nil || aff == 0 {
			session.Rollback()
			return common.ERR_UNKNOWN
		}
	}

	err = session.Commit()
	if err != nil {
		return common.ERR_UNKNOWN
	}

	return common.ERR_SUCCESS
}
