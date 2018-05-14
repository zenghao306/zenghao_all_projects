package common

const (
	TaskVersion = "1.2"
)

const (
	ERR_SUCCESS                    = iota // 成功 == 0
	ERR_UNKNOWN                           // 未知错误 == 1
	ERR_EXIST                             // 账号已经存在 == 2
	ERR_PWD                               // 密码错误 == 3
	ERR_ROOM_EXIST                        // 房间号不存在 == 4
	ERR_ALREAD_FINSIH                     // 直播已经结束 == 5
	ERR_FOCUS_EXIST                       //没有关注过这个人 == 6
	ERR_TOKEN_VALID                       //session不可用 == 7
	ERR_ACCOUNT_EXIST                     //账号不存在==8
	ERR_LOGIN_OUT                         //账号未登陆==9
	ERR_PARAM                             //参数错误==10
	ERR_GIFT_EXIST                        //礼品不存在==11
	ERR_DIAMOND                           //钻石不足==12
	ERR_VERIFT_MUTIPLY                    //多次请求验证码==13
	ERR_VERIFT_CODE                       //验证码错误==14
	ERR_REGISTER_TEL                      //手机号已被注册==15
	ERR_BIND_TEL                          //手机号没有绑定账户==16
	ERR_ALREADY_GAG                       //已经被禁言==17
	ERR_GAG_EXIST                         //不存在被禁言的记录==18
	ERR_GAG_INSIDE                        //主播不在房间内无法禁言==19
	ERR_CASH_TEL_BIND                     //提现手机已经绑定==20
	ERR_CASH_TEL_UNBIND                   //提现手机已经未绑定==21
	ERR_RICE                              //米粒不足==22
	ERR_SESSION_EXIST                     //私聊会话不存在==23
	ERR_TOEKN_EXPIRE                      //token过期==24
	ERR_BLACK_IN                          //已经在对方黑名单==25
	ERR_UPLOAD_EMPTY                      //没有上传图片==26
	ERR_CONFGI_ITEM                       //充值道具不存在==27
	ERR_GET_IMAGE_ERR                     //获取第三方图片失败==28
	ERR_WEIXIN_TOKEN_FAILED               //授权验证失败==29
	ERR_JSON_GET                          //json解码失败==30
	ERR_REPEAT_PAY                        //重复提交此次交易==31
	ERR_TOEKN_NULL                        //token空==32
	ERR_THIRD_PWD                         //密码不能修改=33
	ERR_TEL_DUPLICATE                     //手机号冲突=34
	ERR_CHAT_FINISH                       //房间状态不是在直播=35
	ERR_NOT_HOUSE_MANAGER                 //不是房主没有权限=36
	ERR_NOT_CHAT_EXIST                    //聊天室不存在=37
	ERR_ALREADY_IN_CHAT                   //已经在聊天室=38
	ERR_STAY_IN_SAME_CHAT                 //不在同一个聊天室=39
	ERR_LETTER_TO_SELF                    //不能发送私信给自己=40
	ERR_FOCUS_TO_SELF                     //不能关注给自己=41
	ERR_OTHER_LOGIN                       //你被挤掉了==42
	ERR_OVER_CASH                         //超过提现金额==43
	ERR_UPLOAD_NOTIFY                     //没货收到上传图片==44
	ERR_AUTH_REAL                         //没实名==45
	ERR_AUTH_VET                          //实名正在审核==46
	ERR_RECONNECT_ROOM                    //断线重连已经失效==47
	ERR_USER_SESS                         //用户已经存在链接=48
	ERR_USER_RECONNECT                    //用户正在断线重连=49
	ERR_USER_OFFLINE                      //用户不在线=50
	ERR_CASH_MIN                          //提现没有达到最低数量=51
	ERR_SEND_GIFT_ROBOT                   //测试账户类型该操作无效=52
	ERR_FORBID                            //该账户已经被封=53
	ERR_REPEAT_NICKNAME                   //昵称重复=54
	ERR_REPEAT_THIRD_NICKNAME             //第三方登陆昵称重复=55
	ERR_VERSION_UPDATE                    //需要更新=56
	ERR_SEND_GIFT_SELF                    //不能给自己送礼=57
	ERR_MONEY_LIMIT                       //金钱超过最大限制=58
	ERR_FORBID_POWER                      //临时封号不能操作=59
	ERR_FORBID_EVER                       //永久封号不能操作=60
	ERR_CONTAIN_SENSETIVE                 //包含敏感词=61
	ERR_CLOSE_ROOM                        //关闭的房间房主不一致=62
	ERR_EXIST_STREAM                      //没有直播流=63
	ERR_PLAYBACK_MAX                      //超过录播数量=64
	ERR_PLAYERBACK_REPEAT                 //重复设置推荐状态=65
	ERR_PLAYERBACK_OWNER                  //不能进入自己回播房间=66
	ERR_MULTIPLE_ROOM_NOT_FOUND           //找不到房间=67
	ERR_MULTIPLE_ROOM_IN_USE              //房间正在使用=68
	ERR_MULTIPLE_ROOM_PARM                //参数不可用=69
	ERR_MULTIPLE_ROOM_EXIST               //房间已存在=70
	ERR_MULTIPLE_ROOM_BUSY                //暂无可用房间=71
	ERR_MULTIPLE_ROOM_FIN                 //会议室关闭=72
	ERR_MULTIPLE_NO_POWER_LINK_MIC        //没有权限连麦=73
	ERR_MULTIPLE_HAS_MIC                  //已经申请过不能重复申请=74
	ERR_MULTIPLE_STATUE                   //七牛状态不对=75
	ERR_MULTIPLE_RID                      //会议室不存在=76
	ERR_MULTIPLE_OWNER                    //房主数据错误=77
	//ERR_MULTIPLE_STATUES                  //
	ERR_MONEY_MONEY_LESS          //积分太少=78
	ERR_MONEY_MONEY_MORE          //金钱限制=79
	ERR_TASK_NONE                 //没有找到任务数据=80
	ERR_TASK_STATUS               //任务未完成=81
	ERR_ONLINE_REWARD_EARLY       //在线奖励未达到领取时间=82
	ERR_DAILY_ONLINE_TIMES        //在线奖励领取次数超限制=83
	ERR_DAILY_TASK_CONFIG         //任务版本失效=84
	ERR_DAILY_BOUNTS_HAS_GET      //当日任务已经领取 = 85
	ERR_HAS_NO_DAILY_BUONTS       //当日没有任务可以领取 = 86
	ERR_BIND_USER                 //没有绑定用户=87
	ERR_WEIXIN_ACCOUNT_SIMPLE_BAN //微信非实名用户不可发放=88
	ERR_WEIXIN_NAME_MISMATCH      //真实姓名不一致=89
	ERR_WEIXIN_AMOUNT_LIMIT       //申请付款金额不在有效区间内=90
	ERR_WEIXIN_OPENID_ERROR       //提现的openid错误=91
	ERR_WEIXIN_TRADENO_ERROR      //提现的单号错误=92
	ERR_CASH_NOT_VALID            //不是有效的提现金额=93
	ERR_STOCK_NIL                 //兑换物品0库存=94
	ERR_MONEY_MOON_LESS           //月亮不足95
	ERR_BET_OVER_MAX              //压分超限制=96
	ERR_CONFIG_ADMIN              //没有配置绑定分成用户=97
	ERR_CONFIG_GROUP              //没有配置关联家族=98
	ERR_TOEKN_OVER_TIME           //token时间超过期限==99
	ERR_OPERATOR_EXPIRE           //操作时间过长请重新再来=100
	ERR_ROOM_READY                //已经准备就绪请等待=101
	ERR_PRE_LIVE_RID              //没有按照流程获取房间号=102
	ERR_PRE_LIVE_ADMIN            //家族长不能开播=103
	ERR_PLEASE_GOINGO_CHAT        //请进入直播间=104
	ERR_DUPLICATE_UV              //uv已经注册过=105
	ERR_TOEKN_DB                  //数据库繁忙==106
	ERR_CONFIG_SWITCH             //没有定义开关信息=107
	ERR_SELF_LOGIN                //你已经存在链接==108
	ERR_LETTER_OPT                //没有私信权限==109
	ERR_NO_DAY_SETTLE_POWER       //没有日结权限==110
	ERR_DAY_DURATION_NOT_ENOUGH   //当日在线时长不够==111
	ERR_FOCUS_ALREADY             //已经关注过==112
	ERR_WEIXIN_CASH_NOTENOUGH     //账户余额不足==113
	ERR_WEIXIN_CASH_SYSTEMERROR   //系统错误=114
	ERR_DB_DEL                    //数据库删除出错==115
	ERR_DB_ADD                    //数据库新增出错==116
	ERR_DB_UPDATE                 //数据库更新出错==117
	ERR_DB_FIND                   //数据库查询没有数据==118
	ERR_VERSION_MUST_UPDATE       //需要强制更新=119
	ERR_AUDIENCE_EXIST            //观众已经存在=120
	ERR_NO_POWER_TO_CLOSE_ROOM    //没有关闭房间的权限121
	ERR_TRADE_STATUS              //交易状态未完成=122
	ERR_TRADE_DATE                //没有交易记录=123
	ERR_SNS_TIMEOUT               //短信交易过期=124
	ERR_SNS_CORRECT               //短信验证失败=125
	ERR_GUARD_LIMIE_TIME          //守护开通时间过长=126
	ERR_GUARD_SUPER               //守护开通不能是超管=127
	ERR_SCORE_NOTENOUGH           //游戏币不足=128
	ERR_INNER_XML_ENCODE          = 300
)

var explain = [...]string{
	"成功",
	"未知错误",
	"账号已经存在",
	"密码错误",
	"房间号不存在",
	"直播已经结束",
	"没有关注过这个人",
	"参数错误",
	"礼品不存在",
	"钻石不足",
	"多次请求验证码",
	"验证码错误",
	"手机号已被注册",
}

func GetErrorInfo(err int) string {
	return explain[err]
}

//1qq2新浪3微信4 google 5 Facebook 6 twitter 7 Paypal
const (
	PLATFORM_SELF        = iota //自有平台
	PLATFORM_QQ                 //qq平台,有小写字母
	PLATFORM_SINA               //新浪平台
	PLATFORM_WEIXIN             //微信平台
	PLATFORM_GOOGLE_PLAY = 4    //Google play
	PLATFORM_FACE_BOOK   = 5    //FaceBook
	PLATFORM_TWITTER     = 6    //Twitter
	PLATFORM_PAYPAL      = 7    // Paypal
	PLATFORM_HUAWEI      = 8    //huawei
)

const (
	OS_TYPE_ANDROID = 1
	OS_TYPE_IOS     = 2
)
const (
	ROOM_LIST_PAGE_COUNT       = 10 //房间列表显示分页数目
	MSG_LIST_PAGER_COUNT       = 10 //个人私信分页数量
	FOCUS_LIST_PAGE_COUNT      = 10 //关注列表数量
	SEND_GIFT_PAGE_COUNT       = 10 //送礼排行榜显示分页数目
	CASH_RECORD_PAGE_COUNT     = 10 //提现列表分页数量
	BLACK_LIST_PAGE_COUNT      = 10 //黑名单分页数量
	LIVE_STATISTICS_PAGE_COUNT = 10 //直播统计分页
	MOON_ORDER_PAGE_COUNT      = 10 //月亮商城列表每页数量
)

const (
	VERIFT_CODE_MIN = 1000
	VERIFY_CODE_MAX = 9999
)

var (
	StaticPath             = "./image"
	RoomKey                = "chat"
	ListKey                = "ulist"
	InfoKey                = "info"
	MAX_ROBOT_ADD_ROOM_NUM = 40
	FIRST_CHARGE_DIAMON    = 200
)

const (
	MESSAGE_TYPE_COMMOM            = 1 //聊天室消息定义普通聊天
	MESSAGE_TYPE_GIFT              = 2
	MESSAGE_TYPE_STAR              = 3
	MESSAGE_TYPE_ERR               = 4
	MESSAGE_TYPE_CLOSE             = 5
	MESSAGE_TYPE_SYS               = 6
	MESSAGE_TYPE_ADMIN             = 7
	MESSAGE_TYPE_GAG               = 8
	MESSAGE_TYPE_UNGAG             = 9
	MESSAGE_TYPE_USER_STATUE_JOIN  = 10
	MESSAGE_TYPE_USER_STATUE_LEVEL = 11
	MESSAGE_TYPE_ADMIN_FORBID      = 14
	MESSAGE_TYPE_MULTIPLE_INVITE   = 15
	MESSAGE_TYPE_MULTIPLE_IS_AGREE = 16
	MESSAGE_TYPE_MULTIPLE_CANCEL   = 17
	MESSAGE_TYPE_CHAT_INFO         = 18
	MESSAGE_TYPE_SEND_GIFT         = 19
	MESSAGE_TYPE_ON_CONNECT        = 20
	//MESSAGE_TYPE_TIP_GIFT          = 20
	MESSAGE_TYPE_ADMIN_CLOSE = 21

	MESSAGE_TYPE_GAME_START              = 30
	MESSAGE_TYPE_GAME_STOP               = 31
	MESSAGE_TYPE_GAME_GOING              = 32 //游戏进行中
	MESSAGE_TYPE_GAME_CAN_RAISE          = 33
	MESSAGE_TYPE_GAME_STATE_USER_RAISING = 34 //用户押分
	MESSAGE_TYPE_GAME_RAISE_SCORE        = 35 //左中右押分
	MESSAGE_TYPE_GAME_RAISE_END          = 36 //押分结束
	MESSAGE_TYPE_GAME_RESULT             = 37 //本局结果
	MESSAGE_TYPE_GAME_STATE              = 38
	MESSAGE_TYPE_GAME_START_EOR          = 39 //牛牛创建游戏失败
	MESSAGE_TYPE_GAME_USER_RAISE_EOR     = 40 //用户押分失败
	MESSAGE_TYPE_GAME_NOT_GOING          = 41 //没在牛牛游戏中
	MESSAGE_TYPE_GAME_NOT_RAISE          = 42 //当前没在押分状态下，不可押分
	MESSAGE_TYPE_GAME_RAISE_SCORE_ERROR  = 43 //无效的押分金额，比如大于余额，或者对10取余不为0
	MESSAGE_TYPE_GAME_END_WIN_SCORE      = 46 //游戏结束发送的消息ID
	MESSAGE_TYPE_GAME_RECORD_REQ         = 47 //游戏记录查询
	MESSAGE_TYPE_GAME_RECORD_RES         = 48
	MESSAGE_TYPE_GAME_WINNER_SORTER      = 49 //游戏排行榜

	MESSAGE_TYPE_GAME_TOY_CATCH            = 51 //抓娃娃
	MESSAGE_TYPE_GAME_TOY_CATCH_NEXT       = 52 //抓娃娃第二步
	MESSAGE_TYPE_GAME_TOY_CATCH_NOTICE_ALL = 53 //通知所有用户某某人抓到娃娃了
	MESSAGE_TYPE_GAME_TOY_NOT_CATCH        = 54 //娃娃没碰到时候【仅仅扣游戏币了事】
	MESSAGE_TYPE_GAME_CLOSE                = 55 //游戏结束

	MESSAGE_TYPE_TASK_INFO          = 60 //任务状态更新通知
	MESSAGE_TYPE_DIAMOND_SCORE_INFO = 61 //钻石以及余额获取

	MESSAGE_TYPE_FOCUS         = 70 //关注系统通知
	MESSAGE_TYPE_LETTER_UNREAD = 71 //未读消息通知

	MESSAGE_OPEN_GUARD = 75 //开通守护
)
const (
	GAME_NIUNIU_RAISE            = 30 //牛牛押注环节时长（多少秒）
	GAME_TEXAS_PORK_RAISE        = 30 //德州扑克押注环节时长（多少秒）
	GAME_GF_RAISE                = 30 //砸金花押注环节时长（多少秒）
	GAME_NIUNIU_WAIT_TIME_RESULT = 3  //押注完后等待结果时间（多少秒）
	GAME_TEXAS_WAIT_TIME_RESULT  = 3  //德州扑克完后等待结果时间（多少秒）
	GAME_GF_WAIT_TIME_RESULT     = 3  //砸金花完后等待结果时间（多少秒）
	GAME_NIUNIU_LOOK_RESULT      = 10 //查看结果时间（多少秒）
	GAME_TEXAS_LOOK_RESULT       = 10 //
	GAME_GF_LOOK_RESULT          = 10 //
	GAME_BONUS_TIMES             = 3  //奖励倍数
)

const (
	GAME_NOT_GOING = 0
	GAME_GOING     = 1 // 1.牛牛游戏进行中
	GAME_CAN_RAISE = 2 // 1.用户可押分状态
	GAME_RAISE_END = 3 // 2.押分结束
	GAME_RESULT    = 4 // 3.揭晓结果[结束]
)

const (
	EXP_EVERY_NUM_WATCH        = 5          //每次观众达到观看时间增加的经验值
	EXP_LIMIT_NUM_WATCH        = 20         //每天观看获得经验上限
	FORBID_ACCOUNT_KEEP_TIME   = 1 * 3600   //账号被封时间多少秒之后
	CASH_MONEY_BASE_NUM        = 200        //提现金钱比例
	RECONNECT_TIMEOUT          = 20         //重连超时时间
	MAX_MONEY_LIMIT            = 2000000000 //最大金钱数量
	APPLE_TRADE_CONVERT        = 100
	MIN_CASH_RICE              = 1000 //最少提现米粒
	ROBOT_SAY_TIMER            = 180  //机器人发言间隔
	MAX_SAY_RECORD             = 512
	ONlINE_TIME_REWARD         = 360 //在线时长间隔时间
	ONLINE_TIME_SOCRE          = 15  //在线时长奖励游戏币
	TOKEN_EXPIRE_TIME          = 604800
	EVERY_DAY_LOWEST_CASH_TIME = 3600 //主播每日提现最低直播时长
	NICKNAMERESETMONEY         = 20000
	GUARD_KEEP_DAY             = 7
	GUARD_KEEP_TIME            = GUARD_KEEP_DAY * 24 * 3600

	GUARD_KEEP_TIME_MAX = 180 * 24 * 3600
)

const (
	ACCOUNT_TYPE_COMMON = iota
	ACCOUNT_TYPE_TEST   = 1
)

const (
	SIGNEL_ROOM  = iota
	MUTIPLE_ROOM = 1
)

const (
	TASK_STATUS_NONE    = iota
	TASK_STATUS_ACCECPT = 1 //已经接受任务
	TASK_STATUS_POST    = 2 //带提交任务
	TASK_STATUS_FINISH  = 3 //完成任务
)

//1钻石2星星3游戏币4月亮
const (
	MONEY_TYPE_NONE    = iota
	MONEY_TYPE_DIAMOND = 1
	MONEY_TYPE_RICE    = 2
	MONEY_TYPE_SCORE   = 3
	MONEY_TYPE_MOON    = 4
)
const (
	ACTION_TYPE_LOG_NONE         = iota
	ACTION_TYPE_LOG_LOGIN        = 1
	ACTION_TYPE_LOG_LOGINOUT     = 2
	ACTION_TYPE_LOG_SEND_GIFT    = 3
	ACTION_TYPE_LOG_REGISTER     = 4
	ACTION_TYPE_LOG_CREATE_ROOM  = 5
	ACTION_TYPE_LOG_LEAVE_ROOM   = 6
	ACTION_TYPE_LOG_FOCUS        = 7
	ACTION_TYPE_LOG_CANCEL_FOCUS = 8
	ACTION_TYPE_LOG_ADD_BLACK    = 9
	ACTION_TYPE_LOG_DEL_BLACK    = 10
	ACTION_TYPE_LOG_CLOSE_CHAT   = 11
	ACTION_TYPE_LOG_MAX          = 12
)

var Desc [ACTION_TYPE_LOG_MAX]string

func InitDesc() {

	//Desc = []string{"Red","Blue", "Green", "Yellow", "Pink"}
	Desc[ACTION_TYPE_LOG_NONE] = "无"
	Desc[ACTION_TYPE_LOG_LOGIN] = "登陆"
	Desc[ACTION_TYPE_LOG_LOGINOUT] = "登出"
	Desc[ACTION_TYPE_LOG_SEND_GIFT] = "发送礼物"
	Desc[ACTION_TYPE_LOG_REGISTER] = "注册"
	Desc[ACTION_TYPE_LOG_CREATE_ROOM] = "创建房间"
	Desc[ACTION_TYPE_LOG_LEAVE_ROOM] = "离开房间"
	Desc[ACTION_TYPE_LOG_FOCUS] = "关注"
	Desc[ACTION_TYPE_LOG_CANCEL_FOCUS] = "取消关注"
	Desc[ACTION_TYPE_LOG_ADD_BLACK] = "添加黑名单"
	Desc[ACTION_TYPE_LOG_DEL_BLACK] = "取消黑名单"
	Desc[ACTION_TYPE_LOG_CLOSE_CHAT] = "封禁房间"
	/*
		Desc=make
		Desc = append(Desc, "无")
		Desc = append(Desc, "登陆")
		Desc = append(Desc, "登出")
		Desc = append(Desc, "发送礼物")
		Desc = append(Desc, "注册")
		Desc = append(Desc, "创建房间")
		Desc = append(Desc, "离开房间")
		Desc = append(Desc, "关注")
		Desc = append(Desc, "取消关注")
		Desc = append(Desc, "添加黑名单")
		Desc = append(Desc, "取消黑名单")
	*/
}

/*
Desc:=[]string{
	"无",
	"登陆",
	"登出",
	"发送礼物",
	"注册",
	"创建房间",
	"离开房间",
	"关注",
	"取消关注",
	"添加黑名单",
	"取消黑名单"
}
*/
func GetDesc(optype int) string {
	return Desc[optype]
}

const (
	GAME_TYPE_DEFAULT       = iota // 0.默认
	GAME_TYPE_NIUNIU        = 1    // 1.牛牛
	GAME_TYPE_TEXAS         = 2    // 1.德州扑克
	GAME_TYPE_GOLDEN_FLOWER = 3
	GAME_TYPE_TOY_CATCH     = 4 //抓娃娃
	GAME_TYPE_MAX           = 4 //游戏最大个数
)

const (
	GIFT_CATEGORY_HOT = iota
	//GIFT_CATEGORY_COMMON           = 1
	GIFT_CATEGORY_EXTRAVAGANT      = 1
	GIFT_CATEGORY_COMMON_GAME      = 2
	GIFT_CATEGORY_EXTRAVAGANT_GAME = 3
)

//1用户2系统3管理员
const (
	LIVE_IDENTITY_USER  = 1
	LIVE_IDENTITY_SYS   = 2
	LIVE_IDENTITY_ADMIN = 3
)
