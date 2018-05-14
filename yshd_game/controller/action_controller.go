package controller

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
	"strconv"
	"strings"
	//"github.com/liudng/godump"
)

type PayConfigReq struct {
	Uid      int    `form:"uid" binding:"required"`
	Token    string `form:"token" binding:"required"`
	Category int    `form:"category"`
}

type RegisterReq struct {
	Account      string `form:"account" binding:"required"`
	Pwd          string `form:"pwd" binding:"required"`
	Code         string `form:"code" binding:"required"`
	Type         string `form:"type" binding:"required"`
	ChannelId    string `form:"channel_id"`
	Device       string `form:"device"`
	RegisterFrom int    `form:"register_from"`
	Version      string `form:"version"`
}

type QianXunReq struct {
	Mobile string `form:"mobile"`
}

//首页测试登陆状态页面
func Index(s sessions.Session, r render.Render, req *http.Request, resp http.ResponseWriter) {
	/*
		ret_value := make(map[string]interface{})
		if uid := s.Get("userID"); uid != nil {
			ret_value["ErrCode"] = x.GetToken()
			r.JSON(http.StatusOK, ret_value)

		} else {
			ret_value["ErrCode"] = common.ERR_LOGIN_OUT
			r.JSON(http.StatusOK, ret_value)
		}
	*/
}

//验证手机号
func VerifyTel(req *http.Request, r render.Render) {
	common.Log.Info("VerifyTel() called@@@@@@")
	ret_value := make(map[string]interface{})
	tel := req.FormValue("tel")
	_, has := model.GetUserByTel(tel)
	if has {
		ret_value["ErrCode"] = common.ERR_REGISTER_TEL
	} else {
		ret_value["ErrCode"] = common.ERR_SUCCESS
	}
	r.JSON(http.StatusOK, ret_value)

}

//curl -d 'account=18664328365&pwd=yunshanghudong123456&code=1222&type=0&channel_id=guanfang&version=1.1' 'http://192.168.1.12:3003/action/register'
//curl -d 'account=18664328365&pwd=yunshanghudong123456&code=4555&type=0&channel_id=guanfang' 'http://120.76.156.177:3003/action/register'
//curl -d 'account=13350377086&pwd=yunshanghudong123456&code=1234&type=0&channel_id=guanfang' 'http://120.76.156.177:3003/action/register'
//注册账号处理
func Register(req *http.Request, r render.Render, d RegisterReq) {
	common.Log.Info("Register() called@@@@@@")
	ret_value := make(map[string]interface{})
	/*
		ret_value["token"] = ""
		ret_value["nick_name"] = ""
		tel := req.FormValue("account")
		code := req.FormValue("code")
		pwd := req.FormValue("pwd")
		type_ := req.FormValue("type")
	*/
	//result := model.RequestSnsVerify(tel, code, type_)

	var result int

	if d.Version == "" {
		result = model.RequestSnsVerify(d.Account, d.Code, d.Type)
	} else {
		result = model.QianXunSnsVerify(d.Account, d.Code)
	}
	//result=common.ERR_SUCCESS
	/*
		result = model.RequestSnsVerify(d.Account, d.Code, d.Type)
	*/
	//godump.Dump(d)

	if result == common.ERR_SUCCESS {
		err_code := model.CreateAccountByTel(d.Pwd, d.Account, d.ChannelId, d.Device, d.RegisterFrom)
		ret_value["ErrCode"] = err_code
		if err_code == common.ERR_SUCCESS {
			user, has := model.GetUserByTel(d.Account)
			if has {
				user.GenNewNickName()
				token := common.GenUserToken(user.Uid)
				user.SetToken(token)
				//ret_value["token"] = token
				//ret_value["nick_name"] = user.NickName

				info := &model.LoginInfo{}
				user.GetLoginInfo(info)
				ret_value["user"] = info
			}
		}
	} else {
		ret_value["ErrCode"] = result
	}
	r.JSON(http.StatusOK, ret_value)
}

type LoginReq struct {
	account string
	pwd     string
}

type DailyRecordReq struct {
	Uid   int    `form:"uid" binding:"required"`
	Token string `form:"token" binding:"required"`
}

type DSReq struct {
	Uid   int    `form:"uid" binding:"required"`
	Token string `form:"token" binding:"required"`
}

type WeiXinCashReq struct {
	Uid   int    `form:"uid" binding:"required"`
	Token string `form:"token" binding:"required"`
	Money int    `form:"money" binding:"required"`
}

//登陆过程处理
//curl -d 'account=15311111601&pwd=21218cca77804d2ba1922c33e0151105' 'http://shangtv.cn:3003/action/login'
//curl -d 'account=15311111600&pwd=21218cca77804d2ba1922c33e0151105' 'http://192.168.1.12:3003/action/login'
//http://192.168.1.12:3000/action/login?account=15311111620&pwd=21218cca77804d2ba1922c33e0151105
func Login(req *http.Request, r render.Render, d LoginReq) {
	common.Log.Info("Login() called@@@@@@")
	req.ParseForm()
	account := req.FormValue("account")
	pwd := req.FormValue("pwd")

	req.ParseForm()
	platform := common.PLATFORM_SELF

	ret_value := make(map[string]interface{})
	err_code, user := model.Auth(account, pwd, platform)

	info := &model.LoginInfo{}
	if err_code == common.ERR_SUCCESS {
		token := common.GenUserToken(user.Uid)
		user.SetToken(token)
		user.GetLoginInfo(info)
		ret_value["user"] = info
		model.DailyRecordLog(user.Uid)

		go ChecOtherkLoginIn(user.Uid)
	}

	ret_value["ErrCode"] = err_code
	r.JSON(http.StatusOK, ret_value)
}

func DailyRecord(req *http.Request, r render.Render, d DailyRecordReq) {
	common.Log.Info("VerifyTel() called@@@@@@")
	req.ParseForm()
	uid := req.FormValue("uid")
	uid_, _ := strconv.Atoi(uid)
	req.ParseForm()
	ret_value := make(map[string]interface{})
	//platform := common.PLATFORM_SELF
	model.DailyRecordLog(uid_)
	ret_value["ErrCode"] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func LoginBonusList(req *http.Request, r render.Render, d DailyRecordReq) {
	req.ParseForm()
	uid := req.FormValue("uid")
	uid_, _ := strconv.Atoi(uid)
	req.ParseForm()
	ret_value := make(map[string]interface{})
	ret_value["ErrCode"] = common.ERR_SUCCESS
	ret_value["list"], ret_value["today_get"] = model.LoginBonusList(uid_)
	r.JSON(http.StatusOK, ret_value)
}

func GetLoginBonus(req *http.Request, r render.Render, d DailyRecordReq) {
	req.ParseForm()
	uid := req.FormValue("uid")
	uid_, _ := strconv.Atoi(uid)
	req.ParseForm()
	ret_value := make(map[string]interface{})
	ret_value["ErrCode"] = common.ERR_SUCCESS
	ret_value["status"] = model.GetDailyLoginBonus(uid_)
	ret_value["list"], ret_value["today_get"] = model.LoginBonusList(uid_)
	r.JSON(http.StatusOK, ret_value)
}

//修改密码
//curl -d 'account=18588249532&pwd=123456&code=6827&type=0' 'http://shangtv.cn:3000/action/modify_pwd'
func FindPwdByTel(req *http.Request, r render.Render) {
	common.Log.Info("FindPwdByTel() called@@@@@@")
	ret_value := make(map[string]interface{})
	tel := req.FormValue("account")
	pwd := req.FormValue("pwd")
	code := req.FormValue("code")
	type_ := req.FormValue("type")
	version_ := req.FormValue("version")
	//result := model.RequestSnsVerify(tel, code, type_)

	var result int

	if version_ == "" {
		result = model.RequestSnsVerify(tel, code, type_)
	} else {
		result = model.QianXunSnsVerify(tel, code)
	}

	if result == common.ERR_SUCCESS {
		ret_value["ErrCode"] = model.ModifyPwdByTel(tel, pwd)
	} else {
		ret_value["ErrCode"] = result
	}
	r.JSON(http.StatusOK, ret_value)
}

//curl  'http://192.168.1.12:3000/gift_config'
func GetGiftConfig(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag], ret_value["gift"] = model.GetGiftConfig()
	ret_value["tip"] = model.GetTipGift()
	r.JSON(http.StatusOK, ret_value)
}

//http://192.168.1.12:3000/pay/pay_config_android?token=fec1f362b900b1d0b3e4081bb047fabb&uid=1
//curl  'http://192.168.1.12:3003/pay/pay_config_android?token=f8774b3ff6de892115312d8c77bfa79f&uid=1'
func GetAndroidItemConfig(req *http.Request, r render.Render, d PayConfigReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag], ret_value["item"] = model.GetAndroidItemConfig()
	user, _ := model.GetUserByUid(d.Uid)
	ret_value["diamond"] = user.Diamond
	if user.NewPay == false {
		ret_value["ad"] = model.GetAd()
	}
	r.JSON(http.StatusOK, ret_value)
}

//curl  'http://t1.shangtv.cn:3003/pay/pay_config_ios?token=f8774b3ff6de892115312d8c77bfa79f&uid=1'
func GetIOSItemConfig(req *http.Request, r render.Render, d PayConfigReq) {
	ret_value := make(map[string]interface{})

	if d.Category == 0 {
		d.Category = 1
	}
	ret_value[ServerTag], ret_value["item"] = model.GetIOSItemConfigByChannel(d.Category)
	user, _ := model.GetUserByUid(d.Uid)
	ret_value["diamond"] = user.Diamond
	if user.NewPay == false {
		ret_value["ad"] = model.GetAd()
	}
	r.JSON(http.StatusOK, ret_value)
}

func GetQiNiuConfig(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["face"], ret_value["cover"] = model.GetBucket()
	r.JSON(http.StatusOK, ret_value)
}

func GetGoogleItemConfig(req *http.Request, r render.Render, d PayConfigReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag], ret_value["item"] = model.GetGoogleItemConfig()
	user, _ := model.GetUserByUid(d.Uid)
	ret_value["diamond"] = user.Diamond
	if user.NewPay == false {
		ret_value["ad"] = model.GetAd()
	}
	r.JSON(http.StatusOK, ret_value)
}

//shangtv.cn:3003/task_config?platform=1&version=0.0
func GetTaskConfig(req *http.Request, r render.Render, v TaskVersionReq) {
	ret_value := make(map[string]interface{})
	if strings.Compare(v.Version, common.TaskVersion) == 0 {
		ret_value[ServerTag] = common.ERR_SUCCESS

	} else {
		ret_value[ServerTag] = common.ERR_DAILY_TASK_CONFIG
		ret_value["task"] = model.GetTaskConfig()
		ret_value["version"] = common.TaskVersion
	}

	r.JSON(http.StatusOK, ret_value)
}

func GetScoreConfig(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["config"] = model.GetScoreConfig()
	r.JSON(http.StatusOK, ret_value)

}

func GetDiamondScoreController(req *http.Request, r render.Render, d DSReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = common.ERR_SUCCESS
	user, ret := model.GetUserByUid(d.Uid)
	if ret == common.ERR_SUCCESS {
		info := &model.LoginInfo{}
		user.GetLoginInfo(info)
		ret_value["coupons"] = info.Coupons
		ret_value["diamond"] = info.Diamond
		ret_value["score"] = info.Score
		ret_value["moon"] = info.Moon
	}

	r.JSON(http.StatusOK, ret_value)
}

func GetItemConfig(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["gift"] = model.GetItemConfig()

	r.JSON(http.StatusOK, ret_value)
}

func WeiXinCashController(req *http.Request, r render.Render, d WeiXinCashReq) {
	common.Log.Info("WeiXinCashController() called@@@@@@")

	ret_value := make(map[string]interface{})

	user, _ := model.GetUserByUid(d.Uid)
	if user.CheckAuthReal() == false {
		ret_value[ServerTag] = common.ERR_AUTH_REAL
	} else if user.HasDaySettle() == false {
		ret_value[ServerTag] = common.ERR_NO_DAY_SETTLE_POWER //返回错误码，没有日结权限
		//} else if user.GetDayDuration() <= model.EveryDayLowestCashTime {
		//	ret_value[ServerTag] = common.ERR_DAY_DURATION_NOT_ENOUGH //返回错误码，当日没有足够直播时长
	} else if d.Money < common.CASH_BASE_VALUE || d.Money > common.CASH_MAX_VALUE { //如果未能达到规定的最低提现金额(100元)或者超过最大金额5000
		ret_value[ServerTag] = common.ERR_CASH_NOT_VALID
	} else if user.CashQuota() >= d.Money {
		ret_value[ServerTag] = user.ExchangeRiceToWeiXin(d.Money)
	} else {
		ret_value[ServerTag] = common.ERR_OVER_CASH
	}

	r.JSON(http.StatusOK, ret_value)
}

//curl  -d "mobile=18664328365" 'http://192.168.1.12:3003/action/send_sns'
//curl  -d 'mobile=13243751583' 'http://120.76.156.177:3003/action/send_sns'
//curl  -d 'mobile=13350377086' 'http://192.168.1.12:3003/action/send_sns'
func QianXunSnsController(r render.Render, d QianXunReq) {
	ret_value := make(map[string]interface{})
	code := model.GenQianXunSnsCode(d.Mobile)
	ret_value[ServerTag] = model.RequestQianXunVerify(d.Mobile, code)
	r.JSON(http.StatusOK, ret_value)
}
