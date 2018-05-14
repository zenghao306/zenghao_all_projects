package model

import (
	"github.com/bitly/go-simplejson"
	"github.com/yshd_game/common"
	"github.com/yshd_game/confdata"
	"github.com/yshd_game/melody"
	"time"
)

type ToyInfo struct {
	Id           int     `xorm:"int(11) not null pk autoincr"` //id
	Image1       string  `xorm:"varchar(11) "`                 //名称1
	Price1       int     `xorm:"int(11) not null "`            //第二类押注额对应的娃娃价格1
	Percent1     float64 `xorm:"float(11) not null"`
	PlatformPum1 float64 `xorm:"float(11) "`        //平台抽成比例1
	Image2       string  `xorm:"varchar(11) "`      //名称2
	Price2       int     `xorm:"int(11) not null "` //第二类押注额对应的娃娃价格2
	Percent2     float64 `xorm:"float(11) not null"`
	PlatformPum2 float64 `xorm:"float(11) "`        //平台抽成比例2
	Image3       string  `xorm:"varchar(11) "`      //名称2
	Price3       int     `xorm:"int(11) not null "` //第二类押注额对应的娃娃价格3
	Percent3     float64 `xorm:"float(11) "`
	PlatformPum3 float64 `xorm:"float(11) "` //平台抽成比例
}

type ToyInfoCache struct {
	Id     int    `json:"id`     //id
	Image1 string `json:"image1` //名称1
	Image2 string `json:"image2` //名称2
	Image3 string `json:"image3` //名称3
	Price1 int    `json:"price1` //第二类押注额对应的娃娃价格1
	Price2 int    `json:"price2` //第二类押注额对应的娃娃价格2
	Price3 int    `json:"price3` //第二类押注额对应的娃娃价格3
}

type ToyRecord struct {
	Id         int    `xorm:"int(20) not null pk autoincr"` //id
	RoomId     string `xorm:"varchar(64) not null"`         //房间ID
	OwnerId    int    `xorm:"varchar(11) not null"`         //用户ID
	RaiseScore int    `xorm:"int(11) "`                     //押注分值(10,100,1000)
	ToyId      int    `xorm:"int(11) not null"`             //娃娃ID
	ToyPrice   int    `xorm:"int(11) "`                     //娃娃价格
	Status     int    `xorm:"int(11) "`                     //娃娃价格
	Time       int64  `xorm:"int(20) "`                     //操作时间戳
}

var toy_info_map map[int]ToyInfo
var toy_info_cache []ToyInfoCache

func LoadToyInfo() map[int]ToyInfo {
	toy_info_map = make(map[int]ToyInfo)
	toy_info_cache = make([]ToyInfoCache, 0)

	err := orm.Where("status=1 ").Find(&toy_info_map)
	if err != nil {
		common.Log.Panic("orm err is %s", err.Error())
	}

	for _, v := range toy_info_map {
		var a ToyInfoCache
		a.Id = v.Id
		a.Image1 = v.Image1
		a.Image2 = v.Image2
		a.Image3 = v.Image3
		a.Price1 = v.Price1
		a.Price2 = v.Price2
		a.Price3 = v.Price3

		toy_info_cache = append(toy_info_cache, a)
	}

	return toy_info_map
}

func GetToyInfoConfig() (int, []ToyInfoCache) {
	return common.ERR_SUCCESS, toy_info_cache
}

func GetToyInfoById(id int) (ToyInfo, bool) {
	toy, exist := toy_info_map[id]
	return toy, exist
}

type ToyCashMsgRespose struct {
	MType      int   `json:"mtype"`
	ServerTime int64 `json:"server_time"`
	CashSucc   bool  `json:"cash_succ"` //根据概率算出的是否抓到了
	Score      int   `json:"Score"`     //余额
}

type ToyCashMsgResult struct {
	MType int `json:"mtype"`
	Score int `json:"Score"` //余额
}

func ToyCashMsgDispose(s *melody.Session, msg []byte, user *User, r *ChatRoomInfo) {
	var cashSucc bool
	//req := s.Request
	js, err := simplejson.NewJson(msg)
	if err != nil {
		common.Log.Errf("orm err is 1 %s", err.Error())
		return
	}

	if r.GameType != common.GAME_TYPE_TOY_CATCH { //如果不是抓娃娃房间，表明是客户端处理错误，不予理会，返回
		return
	}

	if !common.GameRunningSwitch { //游戏在后台关闭时候通知所有用户
		NoticeAllUserInRoomGameClose(s)
		return
	}

	timeOpration := time.Now().Unix() //记录下操作的时间戳

	toyID := js.Get("toy_id")
	toyID_, err := toyID.Int()
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	toy, exist := GetToyInfoById(toyID_)
	if !exist {
		return
	}

	score := js.Get("score")
	score_, err := score.Int()
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	if score_%10 != 0 || score_ > int(user.Score) || score_ <= 0 { //判断客户端传递的分值是否是一个有效值
		var data ToyCashMsgResult
		data.MType = common.MESSAGE_TYPE_GAME_RAISE_SCORE_ERROR
		data.Score = int(user.Score)

		SendMsgToUser(user.Uid, data) //不是有效值发送消息给客户，就不要正常押分了
		return
	}

	toyRecord := &ToyRecord{
		RoomId:  r.room.Rid,
		OwnerId: user.Uid,
		ToyId:   toyID_,
		Status:  0,
		Time:    timeOpration,
	}

	switch score_ {
	case 10:
		toyRecord.RaiseScore = 10
		toyRecord.ToyPrice = toy.Price1
		cashSucc = GetRaiseSuccByPercent(toy.Percent1)
	case 100:
		toyRecord.RaiseScore = 100
		toyRecord.ToyPrice = toy.Price2
		cashSucc = GetRaiseSuccByPercent(toy.Percent2)
	case 1000:
		toyRecord.RaiseScore = 1000
		toyRecord.ToyPrice = toy.Price3
		cashSucc = GetRaiseSuccByPercent(toy.Percent3)
	}

	if cashSucc {
		toyRecord.Status = 1
	} else {
		toyRecord.Status = 0
	}

	_, err = orm.InsertOne(toyRecord)
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return
	}

	user.DelMoney(nil, common.MONEY_TYPE_SCORE, int64(toyRecord.RaiseScore), false)

	var data ToyCashMsgRespose
	data.MType = common.MESSAGE_TYPE_GAME_TOY_CATCH
	data.ServerTime = timeOpration
	data.CashSucc = cashSucc
	data.Score = int(user.Score)

	SendMsgToUser(user.Uid, data)
}

//抓娃娃消息处理第二步
func ToyCashMsgDisposeNext(s *melody.Session, msg []byte, user *User, r *ChatRoomInfo) {
	js, err := simplejson.NewJson(msg)
	if err != nil {
		common.Log.Errf("orm err is 1 %s", err.Error())
		return
	}

	serverTime := js.Get("server_time")
	serverTime_, err := serverTime.Int()
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	toyID := js.Get("toy_id")
	toyID_, err := toyID.Int()
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	toyR := &ToyRecord{}
	has, err := orm.Where("owner_id=? and toy_id=? and time=? and status = 1", user.Uid, toyID_, serverTime_).Get(toyR)
	if err != nil {
		common.Log.Err("get real auth error: , %s", err.Error())
		return
	}
	if has { //确认对应的记录存在【加这判断是为了防止后台被劫持】
		// 获取平台抽成比
		platformPum := 0.1 //定义并初始化平台抽成比变量
		toyTemp, exist := GetToyInfoById(toyID_)
		if !exist {
			platformPum = 0.1
		} else {
			switch toyR.ToyPrice {
			case 10:
				platformPum = toyTemp.PlatformPum1
			case 100:
				platformPum = toyTemp.PlatformPum2
			case 1000:
				platformPum = toyTemp.PlatformPum3
			default:
				platformPum = 0.1
			}

		}
		if platformPum < 0.1 || platformPum > 0.9 { //如果不是有效范围内的值则给一个默认值0.1
			platformPum = 0.1
		}

		gainScore := float32(toyR.ToyPrice) - float32(toyR.RaiseScore)*float32(platformPum) //计算用户赢取分数值

		user.AddMoney(nil, common.MONEY_TYPE_SCORE, int64(gainScore), false)

		sess := GetUserSessByUid(user.Uid)
		if sess != nil && sess.Roomid == r.room.Rid {
			//通知用户自己实际赢取金额和用户货币余额
			UserCashResult(user.Uid, int(user.Score), toyID_, toyR.ToyPrice, int(gainScore))
		}

		toyR.Status = 2
		orm.Where("id=?", toyR.Id).Update(toyR)

		isGuard := CheckGuard(user.Uid, r.room.Uid)

		//通知所有用户某某人赢取游戏币
		NoticeToAllUserToyCatch(s, user, toyID_, toyR.ToyPrice, int(gainScore), isGuard)

		//[记录用户每局实际赢取的游戏币] added by zenghao 2017-09-14
		UserWinScoreRecord("", user.Uid, int(gainScore))
		//added end

		TriggerTask(confdata.TargetType_win, user.Uid, 1)

	}
}

type ResponseAllUsersCatchResult struct { //所有用户赢取金额
	MType int `json:"mtype"`  //消息ID
	ToyID int `json:"toy_id"` //娃娃ID
	//ToyName       string `json:"toy_name"`  //娃娃名
	ToyPrice      int    `json:"toy_price"` //娃娃价格
	Winner        int    `json:"winner"`    //赢家ID
	UserLevel     int    `json:"user_level"`
	NickName      string `json:"nick_name"`
	WinnerIsGuard int    `json:"winner_is_guard"` //赢家是否守护
	WinCoins      int    `json:"win_coins"`       //赢取实际金额
}

func NoticeToAllUserToyCatch(s *melody.Session, user *User, toyID, toyPrice, winScore, isGuard int) {
	var data ResponseAllUsersCatchResult
	data.MType = common.MESSAGE_TYPE_GAME_TOY_CATCH_NOTICE_ALL
	data.ToyID = toyID
	//data.ToyName = toyName
	data.ToyPrice = toyPrice
	data.Winner = user.Uid
	data.UserLevel = user.UserLevel
	data.NickName = user.NickName
	data.WinnerIsGuard = isGuard
	data.WinCoins = winScore

	//b, _ := json.Marshal(data)
	//godump.Dump(string([]byte(b)))
	chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
		return q.Request.URL.Path == s.Request.URL.Path
	})
}

type ResponseUserCatchResult struct {
	MType int `json:"mtype"`
	ToyID int `json:"toy_id"` //娃娃ID
	//ToyName  string `json:"toy_name"`  //娃娃名
	ToyPrice int `json:"toy_price"` //娃娃价格
	Score    int `json:"Score"`     //余额
	WinCoins int `json:"win_coins"` //赢取金额
}

// 每局结束后发送用户押分输赢结果【还剩多少金币，赢了多少金币】
// uID:用户ID
// score：用户当前游戏币
// winScores:当局赢下的游戏币
func UserCashResult(uID, score, toyID, toyPrice, winScore int) {
	var data ResponseUserCatchResult
	data.MType = common.MESSAGE_TYPE_GAME_END_WIN_SCORE
	data.Score = score
	data.WinCoins = winScore
	data.ToyPrice = toyPrice
	//data.ToyName = toyName
	data.ToyID = toyID

	SendMsgToUser(uID, data)
}

func ToyNotCashMsgDispose(s *melody.Session, msg []byte, user *User, r *ChatRoomInfo) {
	js, err := simplejson.NewJson(msg)
	if err != nil {
		common.Log.Errf("orm err is 1 %s", err.Error())
		return
	}

	if r.GameType != common.GAME_TYPE_TOY_CATCH { //如果不是抓娃娃房间，表明是客户端处理错误，不予理会，返回
		return
	}

	score := js.Get("score")
	score_, err := score.Int()
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	if score_%10 != 0 || score_ > int(user.Score) || score_ <= 0 { //判断客户端传递的分值是否是一个有效值
		var data ToyCashMsgResult
		data.MType = common.MESSAGE_TYPE_GAME_RAISE_SCORE_ERROR
		data.Score = int(user.Score)

		SendMsgToUser(user.Uid, data) //不是有效值发送消息给客户，就不要扣了
		return
	} else {
		user.DelMoney(nil, common.MONEY_TYPE_SCORE, int64(score_), false)

		var data ToyCashMsgResult
		data.MType = common.MESSAGE_TYPE_GAME_TOY_NOT_CATCH
		data.Score = int(user.Score)

		SendMsgToUser(user.Uid, data)
	}
}
