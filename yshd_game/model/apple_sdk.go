package model

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type InApp struct {
	Quantity       string `json:"quantity"`
	Product_id     string `json:"product_id"`
	Transaction_id string `json:"transaction_id"`
	Purchase_date  string `json:"purchase_date"`
	App_item_id    string `json:"app_item_id"`
}
type Receipt struct {
	Deitail []InApp `json:"in_app"`
}
type AppleRetRecord struct {
	Status      string  `json:"status"`
	Environment string  `json:"environment"`
	Ret         Receipt `json:"receipt"`
}

func AppleAuth(receipt string, uid int, sandbox int, channel_id string, device string, anchor_id int) int {
	a := make(map[string]string, 1)
	a["receipt-data"] = receipt
	b, _ := json.Marshal(a)

	body := strings.NewReader(string(b))

	var url string
	if sandbox == 1 {
		url = "https://sandbox.itunes.apple.com/verifyReceipt"
	} else {
		url = "https://buy.itunes.apple.com/verifyReceipt"
	}

	//发送unified order请求.
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		common.Log.Err("New Http Request发生错误，原因:%s", err.Error())
		return common.ERR_UNKNOWN
	}

	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req) //发送
	defer resp.Body.Close()     //一定要关闭resp.Body
	data, _ := ioutil.ReadAll(resp.Body)
	//	godump.Dump(string(data))

	js, err := simplejson.NewJson(data)
	if err != nil {
		common.Log.Errf("orm err is 1 %s", err.Error())
		return common.ERR_UNKNOWN
	}

	status, err := js.Get("status").Int()
	if err != nil {
		common.Log.Errf("json decode err is 2 %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if status == 21007 {
		return AppleAuth(receipt, uid, 1, channel_id, device, anchor_id)
	} else if status != 0 {
		return common.ERR_UNKNOWN
	}

	object := js.Get("receipt")
	in_app := object.Get("in_app").GetIndex(0)

	quantity, err := in_app.Get("quantity").String()
	if err != nil {
		common.Log.Errf("json decode err is  %s", err.Error())
		return common.ERR_JSON_GET
	}

	quantity_, err := strconv.Atoi(quantity)
	if err != nil {
		common.Log.Errf("strconv quantity err is  %s", err.Error())
		return common.ERR_JSON_GET
	}

	product_id, err := in_app.Get("product_id").String()
	if err != nil {
		common.Log.Errf("json decode err is  %s", err.Error())
		return common.ERR_JSON_GET
	}

	transaction_id, err := in_app.Get("transaction_id").String()
	if err != nil {
		common.Log.Errf("json decode err is  %s", err.Error())
		return common.ERR_JSON_GET
	}

	purchase_date_ori, err := in_app.Get("purchase_date").String()
	if err != nil {
		common.Log.Errf("json decode err is  %s", err.Error())
		return common.ERR_JSON_GET
	}

	purchase_date, err := in_app.Get("purchase_date_ms").String()
	if err != nil {
		common.Log.Errf("json decode err is  %s", err.Error())
		return common.ERR_JSON_GET
	}
	purchase_date_, err := strconv.Atoi(purchase_date)
	if err != nil {
		common.Log.Errf("strconv purchase_date err is  %s", err.Error())
		return common.ERR_JSON_GET
	}
	purchase_date_finish := purchase_date_ / 1000

	randnum := common.RadnomRange(100000, 999999)
	tradeno := fmt.Sprintf("%d%d", time.Now().Unix(), randnum)
	nowtime := time.Now()

	item, has := GetIOSItem(product_id)
	if has == false {
		return common.ERR_CONFGI_ITEM
	}

	has, err = orm.Where(" trade_type=2 and third_trade_id=?", transaction_id).Get(&Trade{})
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has {
		return common.ERR_REPEAT_PAY
	}

	active := GetChargeActive(product_id)

	session := orm.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		session.Rollback()
		return common.ERR_UNKNOWN
	}

	m := &Trade{
		TradeId:         tradeno,
		TimeStart:       nowtime,
		Money:           quantity_ * item.Money,
		Uid:             uid,
		Diamond:         quantity_ * item.Diamond,
		ThirdTradeId:    transaction_id,
		PurchaseTime:    purchase_date_finish,
		Status:          common.TRADE_HAVE_SUCCESS,
		TradeType:       common.TRADE_TYPE_APPLE,
		PurchaseDateOri: purchase_date_ori,
		ChannelId:       channel_id,
		Device:          device,
		AnchorId:        anchor_id,
	}
	aff_row, err := session.Insert(m)
	if err != nil || aff_row == 0 {
		if err != nil {
			common.Log.Errf("err is %s", err.Error())
		}
		session.Rollback()
		return common.ERR_UNKNOWN
	}
	/*
		m := &Trade{}
		has, err = orm.Where("trade_id=? and third_trade_id=?", tradeno, transaction_id).Get(m)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
	*/

	user, _ := GetUserByUid(uid)
	pay_num := quantity_ * item.Diamond
	ret := user.AddMoney(session, common.MONEY_TYPE_DIAMOND, int64(pay_num), true)
	if ret != common.ERR_SUCCESS {
		session.Rollback()
		return ret
	}
	exp := pay_num / 10

	aff, err := user.AddUserExp(session, exp, true)
	if err != nil || aff == 0 {
		session.Rollback()
		return common.ERR_UNKNOWN
	}
	ret = user.SetNewPay(session, pay_num)
	if ret != common.ERR_SUCCESS {
		session.Rollback()
		return ret
	}
	ntime := time.Now().Unix()
	if active != nil {
		if active.Status == 1 && active.BeginTime < ntime && active.FinishTime > ntime {
			ret := user.AddMoney(session, int32(active.MoneyType), active.ExtraNum, true)
			if ret != common.ERR_SUCCESS {
				session.Rollback()
				return ret
			}
		}
	}
	err = session.Commit()
	if err != nil {
		return common.ERR_UNKNOWN
	}

	return common.ERR_SUCCESS
}
