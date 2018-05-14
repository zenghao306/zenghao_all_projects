package model

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"io/ioutil"
	"net/http"
	//	"net/url"
	"sort"
	//	"strconv"
	"github.com/go-with/wxpay"
	"strings"
	"time"

	"bufio"
	"io"
	"os"
)

type Trade struct {
	TradeId    string    `xorm:"varchar(20) pk not null "` //交易号
	TimeStart  time.Time //交易创建时间
	TimeExpire time.Time //交易过期时间
	Uid        int       `xorm:"int(20) not null "` //用户ID
	FeeType    string    `xorm:"varchar(8)`
	Money      int       `xorm:"Float not null "` //充值金额
	//Money         float32 `xorm:"Float not null "`   //充值金额
	Diamond int `xorm:"int(11) not null "` //获得钻石
	Status  int `xorm:"int(11) not null "` // 交易状态
	//WeiXinTradeNo string `xorm:"varchar(20)  "`     // 微信交易号
	//AppleTradeNo  string `xorm:"varchar(20)  "`     //苹果交易号
	//PurchaseData  string    `xorm:"varchar(128) not null "` // 支付完成时间
	TradeType int `xorm:"int(11) not null "`
	//P               time.Time //测试时间
	PurchaseDateOri string `xorm:"varchar(128) not null "` // 支付完成时间
	//PurchaseData    int64  `xorm:"bigint(20) not null "`
	PurchaseTime int `xorm:"int(11) not null "`
	//TotalFee     float32 `xorm:"Float not null "`

	//TradeNoCommon string `xorm:"varchar(64)  not null "`

	ChannelId    string `xorm:"varchar(64)  not null "`
	Device       string `xorm:"varchar(64)  not null "`
	AnchorId     int    `xorm:"  not null "`
	ThirdTradeId string `xorm:"varchar(40)  not null "`
	Extend_1     string `xorm:"varchar(100)  not null "`
}

var (
	weixin_addr = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	//weixin_ios_appkey     = "wxd9ed146926363774"
	//weixin_android_appkey = "wxd9ed146926363774"
	//
	weixin_ios_appkey     = "wxd69dc17cc2dcdb4d"
	weixin_android_appkey = "wxd69dc17cc2dcdb4d"
	//mfc_id                = "1325714901"
	mfc_id     = "1466981502"
	fee        = "CNY"
	trade_type = "APP"
	secret     = "yunshangyunshangyunshangyunshang"
	//secret     = "702388cde96549467aea6cb2f24dabdb"
	notify_url = "http://shangtv.kmdns.net:12346/weixin_pay_notify"
)

type UnifyOrderReq struct {
	Appid            string `xml:"appid"`
	Body             string `xml:"body"`
	Mch_id           string `xml:"mch_id"`
	Nonce_str        string `xml:"nonce_str"`
	Notify_url       string `xml:"notify_url"`
	Trade_type       string `xml:"trade_type"`
	Spbill_create_ip string `xml:"spbill_create_ip"`
	Total_fee        int    `xml:"total_fee"`
	Out_trade_no     string `xml:"out_trade_no"`
	Sign             string `xml:"sign"`
	Time_start       string `xml:"time_start"`
	Time_expire      string `xml:"time_expire"`
	Attach           string `xml:"attach"`
}

type ClientOrderReq struct {
	PartnerId string
	prepayId  string
	Package   string
	Nonce_str string
	TimeStamp string
}

type PreResult struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
	Appid       string `xml:"appid"`
	Mch_id      string `xml:"mch_id"`
	Nonce_str   string `xml:"nonce_str"`
	Sign        string `xml:"sign"`
	Result_code string `xml:"result_code"`
	Trade_type  string `xml:"trade_type"`
	Prepay_id   string `xml:"prepay_id"`
}

type WinXinPayCallBackReq struct {
	Return_code    string `xml:"return_code"`
	Return_msg     string `xml:"return_msg"`
	Appid          string `xml:"appid"`
	Mch_id         string `xml:"nch_id"`
	Nonce_str      string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Result_code    string `xml:"result_code"`
	Openid         string `xml:"openid"`
	Trade_type     string `xml:"trade_type"`
	Bank_type      string `xml:"bank_type"`
	Total_fee      int    `xml:"total_fee"`
	Cash_fee       int    `xml:"cash_fee"`
	Transaction_id string `xml:"transaction_id"`
	Out_trade_no   string `xml:"out_trade_no"`
	Time_end       string `xml:"time_end"`
	Fee_type       string `xml:"fee_type	"`
	Attach         string `xml:"attach"`
}

type RetWeiXinClient struct {
	Noncestr  string
	Prepayid  string
	Timestamp string
	Sign      string
}

func InitWeiXinKey() {
	weixin_notify := common.Cfg.MustValue("host", "weixin_notify")
	notify_url = fmt.Sprintf("http://%s/weixin_pay_notify", weixin_notify)
	//godump.Dump(notify_url)
	weixin_ios_appkey = common.Cfg.MustValue("tencent", "weixin_appkey")
	weixin_android_appkey = weixin_ios_appkey
	secret = common.Cfg.MustValue("tencent", "secret")
}

func wxpayCalcSign(mReq map[string]interface{}, key string) (sign string) {
	//fmt.Println("微信支付签名计算, API KEY:", key)
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		//fmt.Printf("k=%v, v=%v\n", k, mReq[k])
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	//STEP3, 在键值对的最后加上key=API_KEY
	if key != "" {
		signStrings = signStrings + "key=" + key
	}
	//godump.Dump(signStrings)
	//STEP4, 进行MD5签名并且将所有字符转为大写.
	/*
		md5Ctx := common.md5.New()
		md5Ctx.Write([]byte(signStrings))
		cipherStr := md5Ctx.Sum(nil)

	*/

	upperSign := strings.ToUpper(common.Md5(signStrings))
	return upperSign
}

func GetWeiXinH5PrepayId(itemid string, uid int, ip string, channel_id string, device string) (int, *RetWeiXinClient) {

	var ret_client RetWeiXinClient
	randnum := common.RadnomRange(100000, 999999)
	tradeno := fmt.Sprintf("%d%d", time.Now().Unix(), randnum)
	nowtime := time.Now()
	ret_client.Timestamp = fmt.Sprintf("%d", nowtime.Unix())
	expiretime := nowtime.Add(300 * time.Second)

	item, has := GetAndroidItem(itemid)
	if has == false {
		return common.ERR_CONFGI_ITEM, &ret_client
	}

	aff_row, err := orm.Insert(&Trade{
		TradeId:    tradeno,
		TimeStart:  nowtime,
		TimeExpire: expiretime,
		Money:      item.Money,
		Uid:        uid,
		Diamond:    item.Diamond,
		Status:     common.TRADE_PRE_CREATE,
		TradeType:  common.TRADE_TYPE_WEIXIN_H5,
		ChannelId:  channel_id,
		Device:     device,
	})

	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN, &ret_client
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD, &ret_client
	}

	attach := fmt.Sprintf("%s", itemid)
	nonc := common.GenWeiXinRandom()
	ret_client.Noncestr = nonc
	var yourReq UnifyOrderReq
	yourReq.Appid = "wxd0929116c03429bb" //微信开放平台我们创建出来的app的app id
	yourReq.Body = item.Describe
	yourReq.Mch_id = "1480014982"
	yourReq.Nonce_str = nonc
	yourReq.Notify_url = notify_url
	yourReq.Trade_type = trade_type
	yourReq.Spbill_create_ip = ip
	yourReq.Total_fee = item.Money
	yourReq.Out_trade_no = tradeno
	yourReq.Time_start = nowtime.Format("20060102150405")
	yourReq.Time_expire = expiretime.Format("20060102150405")
	yourReq.Attach = attach

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = yourReq.Appid
	m["body"] = yourReq.Body
	m["mch_id"] = yourReq.Mch_id
	m["notify_url"] = yourReq.Notify_url
	m["trade_type"] = yourReq.Trade_type
	m["spbill_create_ip"] = yourReq.Spbill_create_ip
	m["total_fee"] = yourReq.Total_fee
	m["out_trade_no"] = yourReq.Out_trade_no
	m["nonce_str"] = yourReq.Nonce_str
	m["time_start"] = yourReq.Time_start
	m["time_expire"] = yourReq.Time_expire
	m["attach"] = yourReq.Attach
	yourReq.Sign = wxpayCalcSign(m, secret) //这个是计算wxpay签名的函数上面已贴出
	bytes_req, err := xml.Marshal(yourReq)
	if err != nil {
		common.Log.Err("以xml形式编码发送错误, 原因:%s", err.Error())
		return common.ERR_INNER_XML_ENCODE, &ret_client
	}

	str_req := string(bytes_req)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	str_req = strings.Replace(str_req, "XUnifyOrderReq", "xml", -1)
	bytes_req = []byte(str_req)

	//发送unified order请求.
	req, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/pay/unifiedorder", bytes.NewReader(bytes_req))
	if err != nil {
		common.Log.Err("New Http Request发生错误，原因:%s", err.Error())
		return common.ERR_UNKNOWN, &ret_client

	}
	req.Header.Set("Accept", "application/xml")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	c := http.Client{}

	resp, _err := c.Do(req)

	if _err != nil {
		common.Log.Err("请求微信支付统一下单接口发送错误, 原因:%s", _err.Error())
		return common.ERR_UNKNOWN, &ret_client
	}

	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		common.Log.Err("read weixin resp err, is:%s", err2.Error())
		return common.ERR_UNKNOWN, &ret_client
	}

	var s PreResult
	xml.Unmarshal(body, &s)
	//godump.Dump(s)
	if s.Return_code == "SUCCESS" && s.Result_code == "SUCCESS" {
		ret_client.Prepayid = s.Prepay_id
		var t map[string]interface{}
		t = make(map[string]interface{}, 0)
		t["appid"] = weixin_ios_appkey
		t["partnerid"] = mfc_id
		t["prepayid"] = s.Prepay_id
		t["package"] = "Sign=WXPay"
		t["noncestr"] = ret_client.Noncestr
		t["timestamp"] = ret_client.Timestamp
		ret_client.Sign = wxpayCalcSign(t, secret)

		return common.ERR_SUCCESS, &ret_client
	}

	common.Log.Err("weixin return  err Return_code is:%s  ,Result_msg is %s", s.Return_code, s.Return_msg)

	return common.ERR_UNKNOWN, &ret_client

}

func GetWinXinPrepayId2(itemid string, uid int, ip string, channel_id string, device string, anchor_id int) (int, *RetWeiXinClient) {
	//ret_client := &RetWeiXinClient{}
	var ret_client RetWeiXinClient
	randnum := common.RadnomRange(100000, 999999)
	tradeno := fmt.Sprintf("%d%d", time.Now().Unix(), randnum)
	nowtime := time.Now()
	ret_client.Timestamp = fmt.Sprintf("%d", nowtime.Unix())
	expiretime := nowtime.Add(300 * time.Second)

	item, has := GetAndroidItem(itemid)
	if has == false {
		return common.ERR_CONFGI_ITEM, &ret_client
	}

	//u, _ := GetUserByUid(uid)
	aff_row, err := orm.Insert(&Trade{
		TradeId:    tradeno,
		TimeStart:  nowtime,
		TimeExpire: expiretime,
		Money:      item.Money,
		Uid:        uid,
		Diamond:    item.Diamond,
		Status:     common.TRADE_PRE_CREATE,
		TradeType:  common.TRADE_TYPE_WEIXIN,
		ChannelId:  channel_id,
		Device:     device,
		AnchorId:   anchor_id,
	})

	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN, &ret_client
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD, &ret_client
	}

	attach := fmt.Sprintf("%s", itemid)
	nonc := common.GenWeiXinRandom()
	ret_client.Noncestr = nonc
	var yourReq UnifyOrderReq
	yourReq.Appid = weixin_ios_appkey //微信开放平台我们创建出来的app的app id
	yourReq.Body = item.Describe
	yourReq.Mch_id = mfc_id
	yourReq.Nonce_str = nonc
	yourReq.Notify_url = notify_url
	yourReq.Trade_type = trade_type
	yourReq.Spbill_create_ip = ip
	yourReq.Total_fee = item.Money
	yourReq.Out_trade_no = tradeno
	yourReq.Time_start = nowtime.Format("20060102150405")
	yourReq.Time_expire = expiretime.Format("20060102150405")
	yourReq.Attach = attach

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = yourReq.Appid
	m["body"] = yourReq.Body
	m["mch_id"] = yourReq.Mch_id
	m["notify_url"] = yourReq.Notify_url
	m["trade_type"] = yourReq.Trade_type
	m["spbill_create_ip"] = yourReq.Spbill_create_ip
	m["total_fee"] = yourReq.Total_fee
	m["out_trade_no"] = yourReq.Out_trade_no
	m["nonce_str"] = yourReq.Nonce_str
	m["time_start"] = yourReq.Time_start
	m["time_expire"] = yourReq.Time_expire
	m["attach"] = attach
	yourReq.Sign = wxpayCalcSign(m, secret) //这个是计算wxpay签名的函数上面已贴出

	bytes_req, err := xml.Marshal(yourReq)
	if err != nil {
		common.Log.Err("以xml形式编码发送错误, 原因:%s", err.Error())
		return common.ERR_INNER_XML_ENCODE, &ret_client
	}

	str_req := string(bytes_req)
	//wxpay的unifiedorder接口需要http body中xmldoc的根节点是<xml></xml>这种，所以这里需要replace一下
	str_req = strings.Replace(str_req, "XUnifyOrderReq", "xml", -1)
	bytes_req = []byte(str_req)

	//发送unified order请求.
	req, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/pay/unifiedorder", bytes.NewReader(bytes_req))
	if err != nil {
		common.Log.Err("New Http Request发生错误，原因:%s", err.Error())
		return common.ERR_UNKNOWN, &ret_client

	}
	req.Header.Set("Accept", "application/xml")
	//这里的http header的设置是必须设置的.
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	c := http.Client{}

	resp, _err := c.Do(req)

	if _err != nil {
		common.Log.Err("请求微信支付统一下单接口发送错误, 原因:%s", _err.Error())
		return common.ERR_UNKNOWN, &ret_client
	}

	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		common.Log.Err("read weixin resp err, is:%s", err2.Error())
		return common.ERR_UNKNOWN, &ret_client
	}

	var s PreResult
	xml.Unmarshal(body, &s)
	//godump.Dump(s)
	if s.Return_code == "SUCCESS" && s.Result_code == "SUCCESS" {
		ret_client.Prepayid = s.Prepay_id
		var t map[string]interface{}
		t = make(map[string]interface{}, 0)
		t["appid"] = weixin_ios_appkey
		t["partnerid"] = mfc_id
		t["prepayid"] = s.Prepay_id
		t["package"] = "Sign=WXPay"
		t["noncestr"] = ret_client.Noncestr
		t["timestamp"] = ret_client.Timestamp
		ret_client.Sign = wxpayCalcSign(t, secret)

		return common.ERR_SUCCESS, &ret_client
	}

	common.Log.Err("weixin return  err Return_code is:%s  ,Result_msg is %s", s.Return_code, s.Return_msg)

	return common.ERR_UNKNOWN, &ret_client
}

func ProgressPayStatueCallBack(d WinXinPayCallBackReq) int {
	//godump.Dump("dddwin xin call back")
	//godump.Dump(d)
	if d.Result_code != "SUCCESS" {
		common.Log.Err("weixin pay err is %s,%s,%s", d.Return_code, d.Return_msg, d.Out_trade_no)
		return common.ERR_UNKNOWN
	}

	trade := &Trade{}
	has, err := orm.Where("trade_id=?", d.Out_trade_no).Get(trade)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	active := GetChargeActive(d.Attach)
	now_time := time.Now().Unix()
	if has {
		if trade.Status == common.TRADE_PRE_CREATE {
			user, ret := GetUserByUid(trade.Uid)
			if ret != common.ERR_SUCCESS {
				return common.ERR_UNKNOWN
			}

			_, err := orm.Exec("call trade_purchase_data_ori_format(?,?,?,?)", d.Out_trade_no, d.Time_end, common.TRADE_HAVE_SUCCESS, d.Fee_type)
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())
				return common.ERR_UNKNOWN
			}

			m := &Trade{}
			_, err = orm.Where("trade_id=?", d.Out_trade_no).Get(m)
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())
				return common.ERR_UNKNOWN
			}

			session := orm.NewSession()
			defer session.Close()
			err = session.Begin()
			if err != nil {
				common.Log.Errf("orm is error:  %s", err.Error())
				session.Rollback()
				return common.ERR_UNKNOWN
			}
			if m.Status == common.TRADE_HAVE_SUCCESS {

				ret := user.AddMoney(session, common.MONEY_TYPE_DIAMOND, int64(trade.Diamond), true)
				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return common.ERR_UNKNOWN
				}

				aff, err := user.AddUserExp(session, trade.Diamond/10, true)
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

				m.ThirdTradeId = d.Transaction_id
				aff, err = session.Where("trade_id=?", d.Out_trade_no).Update(m)
				if err != nil || aff == 0 {
					if err != nil {
						common.Log.Errf("orm err is %s", err.Error())
					}
					common.Log.Errf("aff is %d", aff)
					return common.ERR_UNKNOWN
				}

			} else {
				common.Log.Errf("weixin notify is err check db  trade_id=?", d.Transaction_id)
				return common.ERR_TRADE_STATUS
			}
			err = session.Commit()
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())
				return common.ERR_UNKNOWN
			}
			return common.ERR_SUCCESS
		}
	}
	return common.ERR_TRADE_DATE

}

type LoginResult struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type WeiXinUser struct {
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Language   string   `json:"language"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

func CheckWeixinToekn(access_token, openid string) int {
	url := fmt.Sprintf("http://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s", access_token, openid)
	resp, err := http.Get(url)
	if err != nil {
		common.Log.Errf("weixin token check err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		common.Log.Errf("weixin token check err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	var v LoginResult
	json.Unmarshal(body, &v)
	if v.Errcode != 0 {
		common.Log.Errf("weixin check err code is %d,msg is %s", v.Errcode, v.Errmsg)
	}
	return v.Errcode

}

func ReadLine2(filePth string) error {
	//godump.Dump(filePth)
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer f.Close()
	bfRd := bufio.NewReader(f)
	for {
		line_, err := bfRd.ReadBytes('\n')
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				return nil
			}
			return err
		}
		//line := string(line_)
		hookfn2(line_) //放在错误处理前面，即使发生错误，也会处理已经读取到的数据。
	}
	return nil
}

func hookfn2(line []byte) {
	//line = strings.Replace(line, "\n", "", -1)
	//line = strings.Replace(line, "\r", "", -1)
	var v WeiXinUser
	err := json.Unmarshal(line, &v)
	if err != nil {
		common.Log.Err(string(line))
		common.Log.Err(err.Error())
		return
	}
	_, has := GetUserByAccountAndPlatfrom(v.Unionid, common.PLATFORM_WEIXIN)
	if !has {
		//	common.Log.Debugf("weixin login res info is %s", string(li))
		sex := 1 //男
		if v.Sex == 1 {
			sex = 1
		} else {
			sex = 0
		}
		v.Nickname = strings.Trim(v.Nickname, " ")
		retcode := CreateAccountByWeiXin(v.Unionid, v.Openid, common.Md5("yunshanghudong123456@#$"), v.City, v.Headimgurl, v.Nickname, common.PLATFORM_WEIXIN, sex, 1, 0, "", "")
		if retcode != common.ERR_SUCCESS {
			return
		}

		//godump.Dump(retcode)
	}

	//godump.Dump("err3")
}

func GetWeiXinUserinfo(access_token, openid, channel_id, devie string, registerFrom int) (ret int, union string, new_register int) {
	url := fmt.Sprintf("http://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", access_token, openid)
	resp, err := http.Get(url)
	if err != nil {
		common.Log.Errf("weixin get user info err is %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		common.Log.Errf("weixin get user info is %s", err.Error())
		ret = common.ERR_UNKNOWN
		return

	}
	var v WeiXinUser
	json.Unmarshal(body, &v)

	//user, has := GetUserByAccount(v.Openid)
	user, has := GetUserByAccountAndPlatfrom(v.Unionid, common.PLATFORM_WEIXIN)
	if !has {
		new_register = 1
		//common.Log.Debugf("weixin login res info is %s", string(body))
		sex := 1 //男
		if v.Sex == 1 {
			sex = 1
		} else {
			sex = 0
		}
		v.Nickname = strings.Trim(v.Nickname, " ")

		retcode := CreateAccountByWeiXin(v.Unionid, openid, common.Md5("yunshanghudong123456@#$"), v.City, v.Headimgurl, v.Nickname, common.PLATFORM_WEIXIN, sex, registerFrom, 0, channel_id, devie)
		if retcode != common.ERR_SUCCESS {
			ret = common.ERR_UNKNOWN
			return
		}
	} else {

		if user.OpenId != openid && openid != "" {
			user.OpenId = openid
			user.UpdateByColS("open_id")
		}

		ret = user.CheckAccountForbid()
		if ret != common.ERR_SUCCESS {
			return
		}
		InsertLog(common.ACTION_TYPE_LOG_LOGIN, user.Uid, "")
		union = v.Unionid
		ret = common.ERR_SUCCESS
		return
	}

	union = v.Unionid
	ret = common.ERR_SUCCESS
	return
}

func GetThirdPartyUserInfo(openID, city, headImg, nickName, channel_id, device string, sex, platform, registerFrom int) int {
	var retcode int
	user, has := GetUserByAccountAndPlatfrom(openID, platform)

	if !has {
		if platform == common.PLATFORM_FACE_BOOK || platform == common.PLATFORM_TWITTER {
			retcode = CreateAccountByFacebookOrTwitter(openID, common.Md5("yunshanghudong123456@#$"), city, headImg, nickName, channel_id, device, sex, platform, 0)
		} else {
			retcode = CreateAccountByThirdParty(openID, common.Md5("yunshanghudong123456@#$"), city, headImg, nickName, channel_id, device, sex, platform, registerFrom, 0)
		}

		if retcode != common.ERR_SUCCESS {
			//common.Log.Errf("GetThirdPartyUserInfo() 435 retcode is %d", retcode)
			return common.ERR_UNKNOWN
		}
	} else {
		ret := user.CheckAccountForbid()
		if ret != common.ERR_SUCCESS {
			return ret
		}
		InsertLog(common.ACTION_TYPE_LOG_LOGIN, user.Uid, "")
		return common.ERR_SUCCESS
	}
	return common.ERR_SUCCESS
}

const (
	appId  = "wxd69dc17cc2dcdb4d"               // 微信公众平台应用ID
	mchId  = "1466981502"                       // 微信支付商户平台商户号
	apiKey = "yunshangyunshangyunshangyunshang" // 微信支付商户平台API密钥

	// 微信支付商户平台证书路径
	certFile   = "cert/apiclient_cert.pem"
	keyFile    = "cert/apiclient_key.pem"
	rootcaFile = "cert/rootca.pem"
)

func Weixinpay(openID, relName, desc, tradeNO string, amount int64) (int, string) {
	fmt.Printf("\nWeixinpay() \n")
	c := wxpay.NewClient(appId, mchId, apiKey)

	// 附着商户证书
	err := c.WithCert(certFile, keyFile, rootcaFile)
	if err != nil {
		common.Log.Errf(err.Error())
	}

	params := make(wxpay.Params)
	// 查询企业付款接口请求参数
	params.SetString("mch_appid", c.AppId)
	params.SetString("mchid", c.MchId)
	params.SetString("nonce_str", "5K8264ILTKCH16CQ2502SI8ZNMTM67VS") // 随机字符串
	params.SetString("partner_trade_no", tradeNO)                     // 商户订单号

	params.SetString("openid", openID)
	params.SetString("check_name", "FORCE_CHECK")
	params.SetString("re_user_name", relName)
	params.SetString("desc", desc)
	params.SetInt64("amount", amount)
	localIP := common.GetLocalIpStr()
	if localIP == "" {
		common.Log.Errf("get localIP err")
		localIP = "120.76.96.73"
	}
	params.SetString("spbill_create_ip", localIP)
	params.SetString("sign", c.Sign(params)) // 签名

	// 查询企业付款接口请求URL
	//url := "https://api.mch.weixin.qq.com/mmpaymkttransfers/gettransferinfo"
	url := "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
	//url := "https://api.mch.weixin.qq.com/secapi/pay/refund"

	// 发送查询企业付款请求
	ret, err := c.Post(url, params, true)
	if err != nil {
		common.Log.Errf(err.Error())
	}

	succStr := ret.GetString("return_code")
	if succStr == "SUCCESS" {
		errCodeDes := ret.GetString("err_code_des")
		resultCode := ret.GetString("result_code")
		errCode := ret.GetString("err_code")
		if resultCode == "SUCCESS" {
			return common.ERR_SUCCESS, ""
		} else if errCode == "V2_ACCOUNT_SIMPLE_BAN" {
			return common.ERR_WEIXIN_ACCOUNT_SIMPLE_BAN, errCodeDes
		} else if errCode == "NAME_MISMATCH" {
			return common.ERR_WEIXIN_NAME_MISMATCH, errCodeDes
		} else if errCode == "AMOUNT_LIMIT" {
			return common.ERR_WEIXIN_AMOUNT_LIMIT, errCodeDes
		} else if errCode == "OPENID_ERROR" {
			return common.ERR_WEIXIN_OPENID_ERROR, errCodeDes
		} else if errCode == "CASH_NOTENOUGH" {
			return common.ERR_WEIXIN_CASH_NOTENOUGH, errCodeDes
		} else if errCode == "CASH_SYSTEMERROR" {
			return common.ERR_WEIXIN_CASH_SYSTEMERROR, errCodeDes
		} else {
			return common.ERR_UNKNOWN, errCodeDes
		}
	} else {
		return common.ERR_UNKNOWN, ""
	}

	return common.ERR_UNKNOWN, ""
}
