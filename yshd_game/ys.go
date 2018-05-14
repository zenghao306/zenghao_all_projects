package main

import (
	"github.com/go-martini/martini"
	//"github.com/go-xorm/xorm"
	//"github.com/liudng/godump"
	//"github.com/apiguy/go-hmacauth"
	//"github.com/martini-contrib/csrf"
	//github.com/beatrichartz/martini-sockets
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/yshd_game/common"
	"github.com/yshd_game/controller"
	"github.com/yshd_game/model"
	"github.com/yshd_game/sensitive"
	"net/http"
	"os"
	"os/signal"
	//"regexp"
	"strconv"
	//"time"
	//"github.com/liudng/godump"
	"github.com/liudng/godump"
	"github.com/yshd_game/confdata"
	"syscall"
	"time"
)

var IsServerRestart bool

func CheckValidatePath(r *http.Request, w http.ResponseWriter, c martini.Context, q render.Render) {
	uid := r.FormValue("uid")
	if uid == "" {
		checkValidate(r, w, c, q)
	} else {
		checkValidateAndUid(r, w, c, q)
	}
}
func checkValidate(r *http.Request, w http.ResponseWriter, c martini.Context, q render.Render) {
	/*
		cookies := r.Cookies()
		for _, c := range cookies {
			if c.Name == "my_session" {
				fmt.Println(c.Value)
			}
		}

			csrf.Validate(r, w, x)
			if s.Get("userID") == nil {
				ret_value := make(map[string]interface{})
				ret_value["ErrCode"] = common.ERR_TOKEN_VALID
				//x.Redirect("/login", 401)
				q.JSON(http.StatusOK, ret_value)
			}
	*/
	token := r.FormValue("token")
	ret_value := make(map[string]interface{})
	if token == "" {
		ret_value["ErrCode"] = common.ERR_TOEKN_NULL
		q.JSON(http.StatusOK, ret_value)
		return
	}
	user, has := model.GetUserByToken(token)
	if has == false {
		ret_value["ErrCode"] = common.ERR_TOEKN_EXPIRE
		q.JSON(http.StatusOK, ret_value)
		return
	}

	ret := user.CheckAccountForbid()
	if ret != common.ERR_SUCCESS {
		ret_value[controller.ServerTag] = ret
		q.JSON(http.StatusOK, ret_value)
		return
	}

	//c.MapTo(user, (*model.User)(nil))
}

func checkValidateAndUid(r *http.Request, w http.ResponseWriter, c martini.Context, q render.Render) {
	token := r.FormValue("token")
	uid := r.FormValue("uid")
	uid_, _ := strconv.Atoi(uid)
	ret_value := make(map[string]interface{})
	user, ret := model.GetUserByUid(uid_)
	nowtime := time.Now().Unix()
	if ret == common.ERR_SUCCESS {
		if user.ExpireTime < nowtime {
			ret_value := make(map[string]interface{})
			ret_value["ErrCode"] = common.ERR_TOEKN_OVER_TIME
			q.JSON(http.StatusOK, ret_value)
			return
		}
		if user.Token != token {
			ret_value := make(map[string]interface{})
			godump.Dump(user.Token)
			ret_value["ErrCode"] = common.ERR_TOEKN_EXPIRE
			q.JSON(http.StatusOK, ret_value)
			//http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		ret := user.CheckAccountForbid()
		if ret != common.ERR_SUCCESS {
			ret_value[controller.ServerTag] = ret
			q.JSON(http.StatusOK, ret_value)
			return
		}

		user.SetExpireTime()
	} else {
		ret_value := make(map[string]interface{})
		ret_value["ErrCode"] = ret
		q.JSON(http.StatusOK, ret_value)
		return
	}
}

func Init() {
	common.SetConfig()
	//创建配置的存放图片目录
	var path string
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = "\\"
	} else {
		path = "/"
	}
	dir, _ := os.Getwd()

	root := common.Cfg.MustValue("path", "root_path")
	//当前的目录
	cover := common.Cfg.MustValue("path", "cover")

	face := common.Cfg.MustValue("path", "face")

	gift := common.Cfg.MustValue("path", "gift")

	rpath := dir + path + root + cover
	lpath := dir + path + root + face
	gpath := dir + path + root + gift
	common.Logpath = dir + path + root + "/log/"

	common.SetLog()

	if !common.IsDirExists(rpath) {
		err := os.MkdirAll(rpath, os.ModePerm)
		if err != nil {
			common.Log.Panicf("sys config path %s, err %s", rpath, err.Error())
		}
	}

	if !common.IsDirExists(lpath) {
		err := os.MkdirAll(lpath, os.ModePerm)
		if err != nil {
			common.Log.Panicf("sys config path %s, err %s", lpath, err.Error())
		}
	}

	if !common.IsDirExists(gpath) {
		err := os.MkdirAll(gpath, os.ModePerm)
		if err != nil {
			common.Log.Panicf("sys config path %s, err %s", gpath, err.Error())
		}
	}

	if !common.IsDirExists(common.Logpath) {
		err := os.MkdirAll(common.Logpath, os.ModePerm)
		if err != nil {
			common.Log.Panicf("sys config path %s, err %s", common.Logpath, err.Error())
		}
	}
	//加载appkey
	model.InitAppKey()
	controller.InitServerRestartTimeSet()
	model.InitWeiXinKey()
	model.InitNowPayKey()
	model.InitHuaWeiPayKey()
	model.InitApliPayKey()
	model.InitQiNiuKey()
	model.InitSinaAppKey()
	model.InitQQKey()
	model.InitWeiXinKey()
	model.InitNowPayKey()
	model.InitHuaWeiPayKey()
	model.InitApliPayKey()
	sensitive.InitkeyWord()
	model.InitSayWord()
	model.LoadAllRobotNickName()
	common.Log.Debug("server initializing...")
}

func newMartini() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.Map(model.SetEngine())
	//m.Map(model.RedisInit())
	//
	common.InitDesc()
	m.Map(model.MelodyInit())
	m.Map(model.LoadGift())
	m.Map(model.LoadAndroidPay())
	m.Map(model.LoadIOSPay())
	m.Map(model.LoadComsumer())
	m.Map(model.LoadToyInfo())
	m.Map(model.LoadConfigUserExp())
	m.Map(model.LoadAnchorExp())
	m.Map(model.LoadScoreExchange())
	m.Map(model.LoadConfigItem())

	model.SystemVariableInitOrReset() //系统变量初始化操作
	model.InitChat()
	model.InitNiuNiu()
	model.InitTexas()
	model.InitGoldenFlower()
	model.InitRobot()
	model.InitChatSession()

	model.InitSysNotice()
	model.InitAdMgr()
	model.InitAchorRoom()
	model.InitCahce()
	//model.InitVersionData()

	model.InitRedis()
	confdata.InitConfData()
	model.LoadConfigTask()

	model.InitConsistData() //崩溃时候补偿时间

	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("my_session", store))
	//m.Use(martini.Static(common.StaticPath))
	m.Use(martini.Static("wangye"))
	m.Use(martini.Static("config"))
	//m.Use(martini.Static("./"))
	/*
		m.Use(csrf.Generate(&csrf.Options{
			Secret:     "token123",
			SessionKey: "userID",
			SetCookie:  true,
			// Custom error response.
			ErrorFunc: func(w http.ResponseWriter) {
				http.Error(w, "CSRF token validation failed", http.StatusBadRequest)
			},
		}))
	*/
	m.Use(render.Renderer())
	m.Action(r.Handle)
	return &martini.ClassicMartini{m, r}
}

func main() {
	defer common.PrintPanicStack()
	Init()

	m := newMartini()
	//按照模块分组当前是房间操作功能
	//shangtv.cm:8082/admin/switch_auth
	m.Group("/admin", func(r martini.Router) {
		//自用测试
		r.Get("/gen_name", controller.GenNameAdmin)
		r.Get("/refresh_room_status", controller.RefreshRoomStatus)
		r.Get("/get_room", binding.Bind(controller.RoomReq{}), controller.GetRoomInfoAdmin)
		r.Get("/get_sess", binding.Bind(controller.CommonReq{}), controller.GetUserInfoAdmin)
		r.Get("/close_sess", binding.Bind(controller.CommonReq{}), controller.CloseUserInfoAdmin)
		r.Get("/save_m8u3", binding.Bind(controller.GetStreamReq{}), controller.SaveM8u3)
		r.Get("/fresh_play_back", controller.FreshPlayBack)
		r.Any("/hidden_recommand_play_back", binding.Bind(controller.GetStreamReq{}), controller.MarkRecommandPlayBack)
		r.Get("/TEST", binding.Bind(controller.SwitchReq{}), controller.TEST)
		r.Any("/session/close", binding.Bind(controller.UIDReq{}), controller.AdminCloseSession)
		r.Get("/SELF", binding.Bind(controller.SelfReq{}), controller.SelfTest)
		r.Get("/SELF2", controller.SelfTest2)

		r.Get("/SELF3", controller.SelfTest3)
		r.Get("/audience_count", binding.Bind(controller.SelfReq{}), controller.AudienceCount)

		//加载所有配置
		r.Get("/reload_config", controller.ReloadAllConfig)
		r.Get("/reload/system/variable", controller.SystemVariableInitOrReset)

		//刷新系统公告
		r.Get("/flush_notice", controller.FlushNotice)
		//实名开关
		r.Get("/switch_auth", binding.Bind(controller.SwitchReq{}), controller.AuthReal)
		//获取当前实名开关状态
		r.Get("/switch_auth_status", controller.AuthRealStatus)
		//测试账户权限开关
		r.Get("/switch_account_auth", binding.Bind(controller.SwitchReq{}), controller.AuthAccountType)
		//设置提现开关
		r.Get("/set_switch_cash_bank", binding.Bind(controller.SwitchReq{}), controller.SetSwitchCashBank)
		r.Get("/set_switch_online_shopping", binding.Bind(controller.SwitchReq{}), controller.SetSwitchOnlineShopping)
		r.Get("/set_switch_game_running", binding.Bind(controller.SwitchReq{}), controller.SetGameRunningSwitch)
		//关闭直播间
		r.Any("/close_chat", binding.Bind(controller.ChatCloseReq{}), controller.CloseChat)
		//r.Any("/close_chat", binding.Bind(controller.ChatCloseReq{}), controller.CloseChat)
		//r.Get("/add_diamond", binding.Bind(controller.AddDiamondReq{}), controller.TestAddDiamond)
		//发送系统公告
		r.Get("/send_notice", binding.Bind(controller.NoticeToAllReq{}), controller.SendNotice)
		r.Get("/send_notice", binding.Bind(controller.NoticeToAllReq{}), controller.SendNotice)
		//监控直播间
		r.Get("/monitor", controller.MonitorRoom)
		//封号
		r.Post("/forbid_user_power", binding.Bind(controller.ForbidUserReq{}), controller.ForbidAccount)

		//删除会议室
		r.Delete("/delete_multiple", binding.Bind(controller.DelMultipleReq{}), controller.DelMultopleRoom)

		r.Get("/weixin/pay", binding.Bind(controller.WeiXinPayTadeNoReq{}), controller.WeiXinPay)
		r.Get("/weixin/reject_pay", binding.Bind(controller.WeiXinPayTadeNoReq{}), controller.RejectWeiXinPay)
		r.Get("/moon/weixin/pay", binding.Bind(controller.WeiXinPayTadeNoReq{}), controller.MoonWeiXinPay)
		r.Any("/send_letter", binding.Bind(controller.WriteLetterUserReq{}), controller.AdminWriterLetterController)
		r.Any("/user/raise_info", binding.Bind(controller.UserGameRaiseReq{}), controller.UserGameRaise)
		r.Get("/reset_redis", controller.ResetKey)

	})

	m.Group("/switch", func(r martini.Router) {
		r.Get("/get_switch_cash_bank", controller.GetSwitchCashBank)
		r.Get("/get_switch_online_shopping", controller.GetSwitchOnlineShopping)
		r.Get("/get_switch_by_channel", binding.Bind(controller.ChannelSwitchReq{}), controller.ChannelSwitchController)
		r.Get("/get_switch_game_running", controller.GetSwitchGameRunningSwitch)
	})

	m.Group("/room", func(r martini.Router) {
		//直播预创建
		r.Post("/live_create", binding.Bind(controller.PreLiveReq{}), controller.LiveCreate)
		//创建房间
		r.Post("/create_room", binding.Bind(controller.CreateRoomReq{}), controller.CreatRoom)
		//关闭房间
		r.Post("/close_room", binding.Bind(controller.CloseRoomReq{}), controller.CloseRoom)
		//获取拉流地址
		r.Get("/pull_addr", binding.Bind(controller.PreLiveReq{}), controller.PullAddr)
		//回播列表
		r.Get("/play_list", binding.Bind(controller.CommonReq{}), controller.PlayBackList)
		//删除回播记录
		r.Post("/del_play", binding.Bind(controller.CloseRoomReq{}), controller.DelPlayBack)
		//推荐回播房间
		r.Any("/recommand", binding.Bind(controller.RecommandPlayBackReq{}), controller.RecommandPlayBack)
		//取消推荐房间
		r.Any("/cancel_recommand", binding.Bind(controller.CommonReq{}), controller.CancelRecommand)
		//检查回播数量
		r.Get("/check_play_count", binding.Bind(controller.CommonReq{}), controller.CheckPlayBack)
		//检查实名认证
		r.Get("/check_auth", binding.Bind(controller.CommonReq{}), controller.CheckAuthRealController)
	}, CheckValidatePath)

	m.Group("/multiple", func(r martini.Router) {
		//创建房间
		r.Get("/create_room", binding.Bind(controller.CreateMultipleReq{}), controller.CreateMultiple)
		//获取房间
		r.Get("/get_room", binding.Bind(controller.GetMultipleReq{}), controller.GetMutipleRoom)
		//获取七牛直播token
		r.Get("/gen_token", binding.Bind(controller.GenMultpleTokenReq{}), controller.GenQiNiuTokenController)
		//邀请别人加入直播间
		r.Any("/invite", binding.Bind(controller.InviteReq{}), controller.InviteController)
		//退出会议室
		r.Any("/exit", binding.Bind(controller.CloseMultpleReq{}), controller.CloseMultiple)
		//确认状态
		r.Post("/comfirm", binding.Bind(controller.AddInfoMultpleReq{}), controller.ComfirmMultiple)
		//生产直播流
		r.Post("/gen_live_addr", binding.Bind(controller.MultipleLiveReq{}), controller.GenMultipleLive)
		//准备开播
		r.Any("/ready_room", binding.Bind(controller.ReadyRoomReq{}), controller.ReadyMutipleController)
	}, CheckValidatePath)

	m.Group("/action", func(r martini.Router) {
		//注册
		r.Post("/register", binding.Bind(controller.RegisterReq{}), controller.Register)
		//登录
		r.Any("/login", binding.Bind(controller.LoginReq{}), controller.Login)
		//改密码
		r.Post("/modify_pwd", controller.FindPwdByTel)

		//验证电话
		r.Post("/verify_tel", controller.VerifyTel)

		//短信验证
		r.Post("/send_sns", binding.Bind(controller.QianXunReq{}), controller.QianXunSnsController)
	})

	m.Group("/action", func(r martini.Router) {
		r.Any("/daily_record", binding.Bind(controller.DailyRecordReq{}), controller.DailyRecord)
		r.Any("/login/bonus/list", binding.Bind(controller.DailyRecordReq{}), controller.LoginBonusList)
		r.Any("/login/bonus/get", binding.Bind(controller.DailyRecordReq{}), controller.GetLoginBonus)
		r.Any("/get/diamond_score", binding.Bind(controller.DSReq{}), controller.GetDiamondScoreController)
		r.Any("/super/user/close_chat", binding.Bind(controller.ChatCloseReq{}), controller.SuperUserCloseChat)
	}, CheckValidatePath)

	m.Group("/user", func(r martini.Router) {
		//修改角色信息
		r.Post("/modify_char", binding.Bind(controller.ModifyCharReq{}), controller.ModifyChar)
		//修改密码
		r.Post("/modify_pwd", binding.Bind(controller.ModifyPwdReq{}), controller.ModifyPwd)
		//修改昵称
		r.Post("/modify_info_nick", binding.Bind(controller.ModifyInfoNickReq{}), controller.ModifyInfoNick)
		//修改昵称2
		r.Any("/modify_nickname_money", binding.Bind(controller.ModifyInfoNickReq{}), controller.ModifyNickNameWithScore)
		//获取修改昵称信息
		r.Any("/get_reseted_info", binding.Bind(controller.GetResetNickReq{}), controller.GetResetedInfo)
		//修改性别
		r.Post("/modify_info_sex", binding.Bind(controller.ModifyInfoSexReq{}), controller.ModifyInfoSex)
		//修改定位
		r.Post("/modify_info_location", binding.Bind(controller.ModifyInfoLocationReq{}), controller.ModifyInfoLocation)
		//修改签名
		r.Post("/modify_info_sign", binding.Bind(controller.ModifyInfoSignReq{}), controller.ModifyInfoSign)
		//设置推送
		r.Post("/set_push", binding.Bind(controller.SetPushReq{}), controller.SetPush)
		//重新获取基础信息
		r.Get("/refresh_info", binding.Bind(controller.CommonReq{}), controller.RefreshUserInfo)
		//r.Any("/upload_image", binding.Bind(controller.UploadFaceReq{}), controller.UploadFace)
		//实名认证
		r.Post("/auth_real", binding.Bind(controller.AddRealNameReq{}), controller.AuthRealName)
		//r.Get("live_statistics")
	}, CheckValidatePath)

	m.Group("/letter", func(r martini.Router) {
		//打开私信
		r.Get("/letter_open", controller.OpenLetter)
		//写私信
		r.Post("/write_letter", controller.WriteLetter)
		//私信列表
		r.Get("/letter_list", controller.LetterList)
		//私信详情
		r.Get("/letter_show", controller.ShowLetter)
		//删除会话
		r.Post("/del_session", controller.DelSession)
		//获取未读数量
		r.Get("/letter_unread", binding.Bind(controller.CommonReq{}), controller.UnreadLetter)
	}, checkValidateAndUid)

	m.Group("/focus", func(r martini.Router) {
		//关注别人
		r.Post("/focus_other", binding.Bind(controller.FocusOtherReq{}), controller.FocusOther)
		//取消关注
		r.Post("/cancel_focus", binding.Bind(controller.CancleFocusReq{}), controller.CancleFocus)
		//关注列表
		r.Get("/focus_list", binding.Bind(controller.LiveFocusReq{}), controller.GetFocusList)
		//粉丝列表
		r.Get("/fans_list", binding.Bind(controller.FansListReq{}), controller.GetFansList)
		//查看关注信息
		r.Get("/search_focus", binding.Bind(controller.FocusInfoReq{}), controller.GetFocusInfo)
		//关注的人在直播列表
		r.Get("/focus_live_list", binding.Bind(controller.FocusLiveListReq{}), controller.GeFocustLiveList)
	}, CheckValidatePath)

	m.Group("/search", func(r martini.Router) {
		//搜索界面推荐列表，返回用户没关注的主播
		r.Get("/recommend/live/list", binding.Bind(controller.FocusLiveListReq{}), controller.RecommendLiveList)
		//搜索用户
		r.Get("/user/list", binding.Bind(controller.UserSearchReq{}), controller.UserSearch)
	}, CheckValidatePath)

	m.Group("/gift", func(r martini.Router) {
		//送礼物
		r.Any("/send_gift", binding.Bind(controller.SendGiftReq{}), controller.WrapSendGift)
		//r.Post("/send_gift", binding.Bind(controller.SendGiftReq{}), controller.WrapSendGift)
		//r.Get("/send_rank", binding.Bind(controller.AddBlackReq{}), controller.SendGiftRank)
		//送礼排行
		r.Get("/send_rank", binding.Bind(controller.CommonWithOidReq{}), controller.SendGiftRank)

		r.Get("/send_rank_all", binding.Bind(controller.CommonReq{}), controller.GetSendGiftRankAll)
		r.Get("/send_rank_week", binding.Bind(controller.CommonReq{}), controller.GetSendGiftRankWeek)

		r.Get("/send_rank_month", binding.Bind(controller.CommonReq{}), controller.GetSendGiftRankMonth)
		//获得礼物排行
		r.Get("/gain_rank", binding.Bind(controller.CommonWithOidReq{}), controller.GainGiftRank)
		//r.Get("/gain_rank", binding.Bind(controller.AddBlackReq{}), controller.GainGiftRank)

		m.Get("/coupons/rank/list", binding.Bind(controller.CouponsRankReq{}), controller.CouponsRankList)

		m.Get("/coupons/rank/list_all", binding.Bind(controller.CommonReq{}), controller.CouponsRankAllList)
		m.Get("/coupons/rank/list_week", binding.Bind(controller.CommonReq{}), controller.CouponsRankWeekList)
		m.Get("/coupons/rank/list_month", binding.Bind(controller.CommonReq{}), controller.CouponsRankMonthList)

		m.Get("/moon/rank/list", binding.Bind(controller.CouponsRankReq{}), controller.MoonRankList)
		//月亮赠送排行榜
		r.Get("/send_rank_moon", binding.Bind(controller.CommonWithOidReq{}), controller.SendGameGiftRank)
		//月亮获得排行
		r.Get("/gain_rank_moon", binding.Bind(controller.CommonWithOidReq{}), controller.GainMoonGiftRank)
		//m.Post("/tip_gift", binding.Bind(controller.SendGiftReq{}), controller.TipToAnchor)
	}, CheckValidatePath)

	m.Group("/game/winner/rank", func(r martini.Router) {
		r.Get("/week", binding.Bind(controller.GameRankReq{}), controller.GetGameWinScoreRankWeek)
		r.Get("/month", binding.Bind(controller.GameRankReq{}), controller.GetGameWinScoreRankMonth)
		r.Get("/all", binding.Bind(controller.GameRankReq{}), controller.GetGameWinScoreRankAll)
	}, CheckValidatePath)
	/*
		m.Post("/gift/send_gift", checkValidate, binding.Bind(controller.SendGiftReq{}), controller.SendGift)
		m.Get("/gift/send_rank", controller.SendGiftRank)
		m.Get("/gift/gain_rank", controller.GainGiftRank)
	*/
	m.Group("/card", func(r martini.Router) {
		//添加银行卡
		r.Post("/add_card", controller.AddCardInfo)
		//绑定提现手机
		r.Post("/cash_tel", controller.BindCashTel)
		//展示可提现金额
		r.Get("/show_cash", controller.CashRiceShow)
		//提现到银行卡
		r.Post("/cash_bank", binding.Bind(controller.CashBankReq{}), controller.CashBank)
		//重置提现手机
		r.Post("/reset_tel", controller.ChangeTel)
		//提现记录
		r.Get("/cash_record", binding.Bind(controller.CommonWithIndexReq{}), controller.UserCashRecord)
		//检查绑定电话重复性
		r.Get("/check_tel", binding.Bind(controller.CommonReq{}), controller.CheckCashTel)
		r.Any("/weixin/cash", binding.Bind(controller.WeiXinCashReq{}), controller.WeiXinCashController)
		//r.Post("/exchange_item", binding.Bind(controller.ExchangeItemReq{}), controller.ExchangeItemController)
	}, checkValidateAndUid)

	m.Group("/black", func(r martini.Router) {
		//添加黑名单
		r.Post("/add_black", binding.Bind(controller.AddBlackReq{}), controller.AddBlackFunc)
		//删除黑名单
		r.Post("/del_black", binding.Bind(controller.DelBlackReq{}), controller.DelBlackFunc)
		//黑名单列表
		r.Get("/black_list", binding.Bind(controller.BlackListReq{}), controller.BlackListFunc)
	}, checkValidateAndUid)

	m.Group("/pay", func(r martini.Router) {
		//微信预创建
		r.Post("/weixin_prepay", binding.Bind(controller.WinXinPayReq{}), controller.WinXinPrepay)
		//r.Post("/h5_prepay", binding.Bind(controller.WinXinPayReq{}), controller.WinXinH5Prepay)

		//苹果支付
		r.Post("/apple_pay", binding.Bind(controller.ApplePayReq{}), controller.ApplePay)
		//谷歌支付
		r.Any("/google_pay", binding.Bind(controller.GooglePayReq{}), controller.GooglePay)
		//安卓支付配置
		r.Get("/pay_config_android", binding.Bind(controller.PayConfigReq{}), controller.GetAndroidItemConfig)
		//ios支付配置
		r.Get("/pay_config_ios", binding.Bind(controller.PayConfigReq{}), controller.GetIOSItemConfig)
		//google支付配置
		r.Get("/pay_config_google", binding.Bind(controller.PayConfigReq{}), controller.GetGoogleItemConfig)
		//游戏币兑换
		r.Post("/score_exchange", binding.Bind(controller.ScoreExchangeReq{}), controller.ScoreExchangeController)
	}, checkValidateAndUid)

	m.Group("/third/pay", func(r martini.Router) {
		r.Any("/channel/switch", binding.Bind(controller.ThirdPaySwitchReq{}), controller.ThirdPaySwitch)
		r.Any("/prepay", binding.Bind(controller.ThirdPayReq{}), controller.NowPayPrepay)
	}, checkValidateAndUid)
	m.Any("/third/pay/notify", controller.NowPayNotify)

	m.Group("/huawei/pay", func(r martini.Router) {
		r.Any("/prepay", binding.Bind(controller.HuaWeiPayReq{}), controller.HuaWeiPrepay)
	}, checkValidateAndUid)
	m.Post("/huawei/pay/notify", binding.Bind(controller.HuaWeiNotifyReq{}), controller.HuaWeiNotify)

	m.Group("/alipay/pay", func(r martini.Router) {
		r.Any("/prepay", binding.Bind(controller.AliPayReq{}), controller.AliPayPrepay)
	}, checkValidateAndUid)
	m.Any("/alipay/pay/notify", binding.Bind(controller.AlipayNotifyReq{}), controller.AlipayNotify)

	m.Group("/third_login", func(r martini.Router) {
		//微信登录
		r.Post("/weixin_login", binding.Bind(controller.WinXinLoginReq{}), controller.WeiXinLogin)
		//qq登录
		r.Post("/qq_login", binding.Bind(controller.QQLoginReq{}), controller.QQLogin)
		//r.Post("qq_h5_login",binding.Bind(controller.QQLoginReq{}), controller.QQLogin)
		//新浪登录
		r.Post("/sina_login", binding.Bind(controller.SinaLoginReq{}), controller.SinaLogin)
		//facebook登录
		r.Post("/facebook_login", binding.Bind(controller.FaceBookLoginReq{}), controller.FaceBookLogin)
		r.Get("/facebook_login", binding.Bind(controller.FaceBookLoginReq{}), controller.FaceBookLoginGet)
		//twitter登录
		r.Post("/twitter_login", binding.Bind(controller.TwitterLoginReq{}), controller.TwitterLogin)

		//华为登录
		r.Post("/huawei_login", binding.Bind(controller.HuaWeiLoginReq{}), controller.HuaWeiLogin)
	})

	m.Group("/channel", func(r martini.Router) {
		r.Any("/version/config", binding.Bind(controller.ChannelVersionConfigReq{}), controller.ChannelVersionConfig)
	}, checkValidateAndUid)

	m.Group("/goods", func(r martini.Router) {
		r.Any("/sale/list", binding.Bind(controller.UserSaleGoodsListReq{}), controller.UserSaleGoodsListFunc)
		r.Any("/show/list", binding.Bind(controller.GoodsShowReq{}), controller.ShowListFunc)
		r.Any("/click", binding.Bind(controller.SelectShowClickReq{}), controller.SelectShowClick)
	}, checkValidateAndUid)

	m.Group("/task", func(r martini.Router) {

		//r.Post("/refresh_task", binding.Bind(controller.CommonReq{}), controller.RefreshTaskController)
		//任务列表
		r.Get("/list_task", binding.Bind(controller.CommonReq{}), controller.DailyTaskListController)
		//提交任务
		r.Post("/post_task", binding.Bind(controller.PostTaskListReq{}), controller.PostTaskController)
	}, checkValidateAndUid)

	m.Group("/upfile", func(r martini.Router) {
		//生产七牛token
		m.Get("/gen_7niu_token", binding.Bind(controller.UpLoadTokenReq{}), controller.GetSevenNiuToken)
		//m.Post("/confirm_7niu_file", binding.Bind(controller.UpLoadFileReq{}), controller.Set7NiuFileName)
	}, checkValidateAndUid)

	m.Group("/daily", func(r martini.Router) {
		//在线奖励剩余时长
		r.Get("/online_time", binding.Bind(controller.CommonReq{}), controller.OnlineTimeController)
		//领取在线奖励
		r.Any("/online_reward", binding.Bind(controller.CommonReq{}), controller.PostOnlineRewardController)
	}, checkValidateAndUid)

	m.Group("/live_statistics", func(r martini.Router) {
		//直播统计信息
		r.Get("/live_info", binding.Bind(controller.LiveStatisticsReq{}), controller.LiveStatisticsController)
	}, checkValidateAndUid)

	//直播统计索引
	m.Get("/live_statistics/month_index", controller.MonthStatisticsController)

	m.Group("/moon", func(r martini.Router) {
		//r.Get("/item_config", controller.GetItemConfig)
		//兑换实物
		r.Post("/exchange_item", binding.Bind(controller.ExchangeItemReq{}), controller.ExchangeItemController)
		//订单列表
		r.Get("/order_list", binding.Bind(controller.CommonWithIndexReq{}), controller.OrderListController)
	}, checkValidateAndUid)

	m.Group("/guard", func(r martini.Router) {
		r.Post("/open", binding.Bind(controller.OpenGuard{}), controller.OpenGuardController)
		r.Get("/anchor_info", binding.Bind(controller.OpenGuard{}), controller.ListAnchorGuardController)
		r.Get("/self_info", binding.Bind(controller.OpenGuard{}), controller.ListSelfGuardController)
	}, checkValidateAndUid)

	//月亮商城配置
	m.Get("/moon/item_config", controller.GetItemConfig)

	m.Get("/dami", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "wangye/damilive.html")
	})
	m.Get("/yingke", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "wangye/player.html")

	})
	m.Get("/channel/:num", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "chan.html")
	})

	m.Get("/channel/:name/ws", controller.Join)
	/*
		//ws协议处理聊天室逻辑
		m.Get("/channel/:name/ws", func(w http.ResponseWriter, r *http.Request, q render.Render) {
			//godump.Dump("one people in room ")
			//ip := common.GetRemoteIp(r)
			//godump.Dump(ip)

			isWebChannel := false
			//ret_value := make(map[string]interface{})
			reg := regexp.MustCompile(`[0-9]+`)
			roomid := reg.FindAllString(r.URL.Path, -1)

			token := r.FormValue("token")
			//"channel" "web"
			uid := r.FormValue("uid")

			channel := r.FormValue("channel")

			if channel == "web" {
				isWebChannel = true
			}

			if len(roomid) == 0 {
				common.Log.Errf("get room panic")
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			roomid_ := roomid[0]
			c_room := model.GetChatRoom(roomid_)
			if c_room == nil {
				common.Log.Errf("room null panic")
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			if token != "" && !isWebChannel {
				user, ret := model.GetUserByUidStr(uid)
				if ret != common.ERR_SUCCESS {
					return
				} else {
					if token != user.Token {
						return
					}

					if forbid := user.CheckPowerAccount(); forbid == true {
						return
					}

					if sess := model.GetUserSessByUid(user.Uid); sess != nil {
						common.Log.Errf("sess is not nill 已存在@,uid= %d", user.Uid)
						controller.ChecOtherkLoginIn(user.Uid)
						return
					}

					if c_room.GetChatInfo().Uid == user.Uid {

						if c_room.Statue == common.ROOM_PLAYBACK {
							return
						}
						if c_room.Statue == common.ROOM_PRE || c_room.Statue == common.ROOM_READY {
							controller.Join(w, r)
							return
						}
					}
				}
			} else if isWebChannel {
				controller.Join(w, r)
				return
			} else {
				return
			}

			if c_room.Statue == common.ROOM_ONLIVE || c_room.Statue == common.ROOM_PLAYBACK || c_room.Statue == common.ROOM_READY {
				controller.Join(w, r)
			} else {
				common.Log.Errf("room is close")
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
		})
	*/
	//检查聊天室是否存在
	//shangtv.cn:3003/check_channel/2466149577842992759?uid=100060&token=284e342a9a32211b9a052becfd555de6

	/*
		m.Get("/check_channel/:num", func(w http.ResponseWriter, r *http.Request, q render.Render) {
			ret_value := make(map[string]interface{})
			reg := regexp.MustCompile(`[0-9]+`)
			roomid := reg.FindAllString(r.URL.Path, -1)
			if len(roomid) == 0 {
				common.Log.Errf("get room panic")
				ret_value[controller.ServerTag] = common.ERR_PARAM
				q.JSON(http.StatusOK, ret_value)
				return
			}
			roomid_ := roomid[0]
			c_room := model.GetChatRoom(roomid_)
			if c_room == nil {
				ret_value[controller.ServerTag] = common.ERR_ROOM_EXIST
				q.JSON(http.StatusOK, ret_value)
				return
			}

			if c_room.Statue == common.ROOM_PRE || c_room.Statue == common.ROOM_ONLIVE || c_room.Statue == common.ROOM_PLAYBACK || c_room.Statue == common.ROOM_READY {
				token := r.FormValue("token")
				if token != "" {
					user, has := model.GetUserByToken(token)
					if !has {
						ret_value[controller.ServerTag] = common.ERR_TOEKN_EXPIRE
						q.JSON(http.StatusOK, ret_value)
						return
					}

					if forbid := user.CheckPowerAccount(); forbid == true {
						ret_value[controller.ServerTag] = common.ERR_FORBID
						q.JSON(http.StatusOK, ret_value)
						return
					}

					if sess := model.GetUserSessByUid(user.Uid); sess != nil {
						ret_value[controller.ServerTag] = common.ERR_USER_SESS
						q.JSON(http.StatusOK, ret_value)
						return
					}
					if c_room.GetChatInfo().Uid == user.Uid {
						if c_room.Statue == common.ROOM_PLAYBACK {
							ret_value[controller.ServerTag] = common.ERR_PLAYERBACK_OWNER
							q.JSON(http.StatusOK, ret_value)
							return
						}
					}
					ret_value[controller.ServerTag] = common.ERR_SUCCESS
					q.JSON(http.StatusOK, ret_value)
					return
				}
			}
			ret_value[controller.ServerTag] = common.ERR_CHAT_FINISH
			q.JSON(http.StatusOK, ret_value)

			return
		})

	*/
	m.Get("/check_channel/:num", controller.CheckChannelController)

	m.Group("/chat", func(r martini.Router) {
		//禁言
		//r.Post("/gag", controller.GagUser)
		r.Any("/gag", binding.Bind(controller.GagReq{}), controller.GagUser)

		//取消禁言
		r.Post("/cancel_gag", controller.CancelGagUser)

		r.Get("/gag_status", binding.Bind(controller.CommonReqWithRid{}), controller.GagStatusUser)
	}, CheckValidatePath)

	//m.Any("/upload_image", controller.Upload)

	//查看个人信息
	m.Get("/look_info", controller.LookInfo)
	//查看个人所有信息
	m.Get("/look_info_all", checkValidateAndUid, binding.Bind(controller.CommonWithOid2Req{}), controller.LookInfoAll)
	//
	m.Get("/list_room", controller.ListRoom)

	m.Get("/list/room/real/user/count", controller.ListRoomRealUserCount)

	m.Get("/chat_info", binding.Bind(controller.ChatInfoReq{}), controller.GetChatInfo)
	m.Get("/audience_info", binding.Bind(controller.AudienceInfoReq{}), controller.GetAudienceInfo)

	//礼物列表
	m.Get("/gift_config", controller.GetGiftConfig)
	//m.Get("/item_config", controller.GetItemConfig)
	m.Get("/qi_niu_config", controller.GetQiNiuConfig)
	//任务列表
	m.Get("/task_config", binding.Bind(controller.TaskVersionReq{}), controller.GetTaskConfig)

	//积分兑换配置
	m.Get("/score_config", controller.GetScoreConfig)
	//土豪排行
	m.Get("/all_rank", controller.AllGiftRank)

	//德玛测试专用
	m.Any("/upload_html", controller.UploadHtml)
	//获取房间信息
	m.Get("/room_info", binding.Bind(controller.CommonReqOnlyUid{}), controller.GetRoomInfo)
	//m.Get("/multiple_room_info", binding.Bind(controller.CommonReqOnlyUid{}), controller.GetMultipleRoomInfo)
	m.Get("/play_back_room_info", binding.Bind(controller.CommonReqOnlyRid{}), controller.GetRoomInfoByRid)

	//上传敏感词
	m.Post("/upload_sensitive", controller.UploadSensitive)
	//上传机器人发言
	m.Post("/upload_robot", controller.UploadRobot)
	//微信支付回调
	m.Any("/weixin_pay_notify", controller.WinXinPayCallBack)

	//七牛回调处理图片
	m.Any("/sever_niu_notify", controller.QiNiuNotify)

	//热玩列表
	m.Get("/recommand_list", binding.Bind(controller.RecommandListReq{}), controller.RecommandListRoom)
	//m.Get("/recommand_list_v2", binding.Bind(controller.RecommandListReq{}), controller.RecommandListRoomV2)

	//m.Get("/ping", binding.Bind(controller.CommonReqOnlyUid{}), controller.PingReconnect)

	//举报
	m.Post("/report", binding.Bind(controller.ReportUserReq{}), controller.ReportUser)

	m.Get("/toy/config", controller.GetToyConfig)

	//版本信息
	m.Get("/version", binding.Bind(controller.VersionReq{}), controller.GetVersionInfo)
	m.Get("/cash_sys_say", binding.Bind(controller.LangeReq{}), controller.CashSayController)
	m.Get("/present/bank/list", controller.GetPresentBankList)

	m.Get("/admin_list", controller.GetAdminList)

	m.Get("/app/banner/list", controller.BannerList)

	m.Post("/register_device", binding.Bind(controller.ChannelReq{}), controller.RegisterDevice)
	//uv统计
	m.Any("/uv", binding.Bind(controller.UvReq{}), controller.AddUv)

	m.Any("/server_version", controller.ServerController)

	common.Log.Info("server is started...")

	go m.RunOnAddr(common.Cfg.MustValue("base", "add"))
	common.Log.Info("server is finishd...")

	//启动定时任务
	go model.TimerTask()
	go model.TimerTaskGameRaiseBk()
	go model.TimerTaskCouponsMonthRecord()
	//机器人自动发言
	go model.TimerTaskRobotSay()
	go model.TimerTaskMulitple()
	go model.TimerTaskNiuNiu()
	model.RefreshMultiple()
	//启动重连等待
	go model.AnchorMgr.ReconnectRoom()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	sig := <-c
	common.Log.Infof("server is finishd sig is %v", sig)

	var sys model.ResponseSys
	sys.MType = common.MESSAGE_TYPE_SYS
	sys.Notice = "服务器5秒后重启重!!!!"
	model.AdminSysToAll(sys)
	model.AdminSysToAll(sys)
	model.AdminSysToAll(sys)

	time.Sleep(5 * time.Second)
	model.GetChat().Close()
	model.AnchorMgr.Close()

	time.Sleep(2 * time.Second)
	common.Log.Info("server is finishd step2")
	os.Exit(0)
}

/*
package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/liudng/godump"
	"github.com/martini-contrib/encoder"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	"strconv"
)

type Some struct {
	Login    string `json:"login"`
	Password string `json:"password,omitempty" xml:",omitempty"`
	Url      string `json:"url"`
}

func (this Some) Filter() interface{} {
	this.Password = ""
	this.Url = "http://some-origin/" + this.Login
	return this
}

func init() {
	log.SetFlags(log.Lshortfile)
}
func Say(name, gender string, age int) {
	fmt.Printf("My name is %s, gender is %s, age is %d!\n", name, gender, age)
}
func main() {
	m := martini.New()
	route := martini.NewRouter()

	m.Use(func(c martini.Context, w http.ResponseWriter, r *http.Request) {
		// Use indentations. &pretty=1
		pretty, _ := strconv.ParseBool(r.FormValue("pretty"))
		// Use null instead of empty object for json &null=1
		null, _ := strconv.ParseBool(r.FormValue("null"))
		// Some content negotiation
		switch r.Header.Get("Content-Type") {
		case "application/xml":
			c.MapTo(encoder.XmlEncoder{PrettyPrint: pretty}, (*encoder.Encoder)(nil))
			w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		default:
			c.MapTo(encoder.JsonEncoder{PrettyPrint: pretty, PrintNull: null}, (*encoder.Encoder)(nil))
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}
	})

		route.Get("/user/:uid", func(enc encoder.Encoder, parms martini.Params) int {

			age, _ := strconv.ParseInt(parms["uid"], 10, 32)
			Say("陈一回", "男", int(age))
			return http.StatusOK
		})



	m.Get("/user", func(enc encoder.Encoder, parms martini.Paramsr, r *http.Request) (int, []byte) {

		fmt.Printf("param is  =  %s ", parms["uid"])
		godump.Dump(r.FormValue("uid"))
		//age, _ := strconv.ParseInt(parms["uid"], 10, 32)
		//Say("陈一回", "男", int(age))
		result := Some{"user1", "passwordhash", "/user1"}
		return http.StatusOK, encoder.Must(enc.Encode(result))
	})


	//route.URLFor("/user", "id")

	route.Get("/users", func(enc encoder.Encoder) (int, []byte) {
		result := []Some{
			Some{"user1", "somehash", "/user1"},
			Some{"user2", "somehash", "/user2"},
		}

		return http.StatusOK, encoder.Must(enc.Encode(result))
	})

	m.Action(route.Handle)

	//url := route.URLFor("user", "uid", "5")
	//fmt.Println(url)

	log.Println("Waiting for connections...")

	if err := http.ListenAndServe(":8000", m); err != nil {
		log.Fatal(err)
	}
}
*/
