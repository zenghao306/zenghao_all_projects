/*****************************************************************************
 *
 * Filename:
 * ---------
 *    texas_pok.go
 *
 * Project:
 * --------
 *   yshd_game
 *
 * Description:
 * ------------
 *   本文件的函数为德州核心代码，主要用于实现德州扑克业务逻辑和通信协议。
 *
 * Author:
 * sky.Zeng
 *
 * Date:
 * 2017-04-10
 *
 ****************************************************************************/
package model

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/yshd_game/common"
	"github.com/yshd_game/melody"
	"github.com/yshd_game/timer"
	"sort"
	"sync"
	"time"
)

const (
	HIG_CARD        = iota //0高牌
	ONE_PAIR        = 1    //1.一对
	TWO_PAIR        = 2    //2.两对
	THREE_OF_A_KIND = 3    //3.3条
	STRAIGHT        = 4    //4.顺子
	FLUSH           = 5    //5.同花
	FULL_HOUSE      = 6    //6葫芦
	FOUR_OF_A_KIND  = 7    //7.4条
	STRAIGHT_FLUSH  = 8    //8.同花顺
	ROYAL_FLUSH     = 9    //9.皇家同花顺
)

var TexasPork = [52]Card{
	// 0-12 方块
	{2, 0}, {3, 0}, {4, 0}, {5, 0}, {6, 0}, /*0,1,2,3,4*/
	{7, 0}, {8, 0}, {9, 0}, {10, 0}, {11, 0}, /*5,6,7,8,9*/
	{12, 0}, {13, 0}, {14, 0}, /*10,11,12*/

	// 13-25 梅花
	{2, 1}, {3, 1}, {4, 1}, {5, 1}, {6, 1}, /*13,14,15,16,17*/
	{7, 1}, {8, 1}, {9, 1}, {10, 1}, {11, 1}, /*18,19,20,21,22*/
	{12, 1}, {13, 1}, {14, 1}, /*23,24,25*/

	// 26-38 红桃
	{2, 2}, {3, 2}, {4, 2}, {5, 2}, {6, 2}, /*26,27,28,29,30*/
	{7, 2}, {8, 2}, {9, 2}, {10, 2}, {11, 2}, /*31,32,33,34,35*/
	{12, 2}, {13, 2}, {14, 2}, /*36,37,38*/

	// 39-51 黑桃
	{2, 3}, {3, 3}, {4, 3}, {5, 3}, {6, 3}, /*39,40,41,42,43*/
	{7, 3}, {8, 3}, {9, 3}, {10, 3}, {11, 3}, /*44,45,46,47,48*/
	{12, 3}, {13, 3}, {14, 3}, /*49,50,51*/
}

type TexasActorInfo struct {
	Cards [5]Card `json:"cards"` //5张牌
	Grade int     `json:"grade"` //2张牌+5张牌里取3张组合最大的德克萨斯等级[0-9]参见德克萨斯级别常量
}

type Struct5Numbers struct {
	N1 int
	N2 int
	N3 int
	N4 int
	N5 int
}

//0,1,2,3,4这5个数字里随机取3个的10种组合
var Texas5Numbers = [21]Struct5Numbers{
	{0, 1, 2, 3, 4},
	{0, 1, 2, 3, 5},
	{0, 1, 2, 4, 5},
	{0, 1, 3, 4, 5},
	{0, 2, 3, 4, 5},
	{1, 2, 3, 4, 5},
	{0, 1, 2, 3, 6},
	{0, 1, 2, 4, 6},
	{0, 1, 3, 4, 6},
	{0, 2, 3, 4, 6},
	{1, 2, 3, 4, 6},
	{0, 1, 2, 5, 6},
	{0, 1, 3, 5, 6},
	{0, 2, 3, 5, 6},
	{1, 2, 3, 5, 6},
	{0, 1, 4, 5, 6},
	{0, 2, 4, 5, 6},
	{1, 2, 4, 5, 6},
	{0, 3, 4, 5, 6},
	{1, 3, 4, 5, 6},
	{2, 3, 4, 5, 6},
}

type RoomInfo_Texas struct { // RoomInfo_Texas.房间里的德州扑克信息
	Rid              string
	Uid              int
	IsStarted        bool
	mutex_Texas_info sync.RWMutex
	GameState        int    //游戏状态【进行到哪一步了】
	CSStartTime      int64  //当前状态下开始时间戳
	GameID           string //当前局的游戏ID
	LScore           int
	MScore           int
	RScore           int
	TexasActors      []*TexasActorInfo //当前局下的牌
	AnyCards         []*Card           //当前局下的任意牌
	StaticCards      []*Card           //当前局下的万家手上牌序列
	LargeActors      int               //最大的玩家
}

var Texas_room_manager map[string]*RoomInfo_Texas
var mutex_Texas sync.RWMutex

//根据房间ID创建一个德州结构体对象并返回
func NewRoomInfoTexas() *RoomInfo_Texas {
	m := &RoomInfo_Texas{}
	return m
}

//德州管理者map[Texas_room_manager]初始化
func InitTexas() map[string]*RoomInfo_Texas {
	Texas_room_manager = make(map[string]*RoomInfo_Texas)
	return Texas_room_manager
}

//根据房间ID创建返回一个德州结构体对象
func GetTexasByRoomid(roomID string) *RoomInfo_Texas {
	v, ok := Texas_room_manager[roomID]
	if !ok {
		return nil
	}
	return v
}

//添加一个德州结构体r到德州管理者MAP里面[Texas_room_manager]
func AddRoomInfoTexas(r *RoomInfo_Texas) bool {
	_, ok := Texas_room_manager[r.Rid]
	if ok {
		return false
	}

	mutex_Texas.Lock()
	Texas_room_manager[r.Rid] = r
	mutex_Texas.Unlock()

	return true
}

//根据房间ID从德州管理者Texas_room_manager里面移除
func DelRoomInfoTexas(rid string) {
	mutex_Texas.Lock()
	delete(Texas_room_manager, rid)
	mutex_Texas.Unlock()
}

// 本函数主要用于德州扑克游戏的创建【入口】
// s 会话的Session对象，
// anchorID 调用者（主播）的UID
// roomID 房间ID
func NewGameTexasPok(s *melody.Session, anchorID int, roomID string) {
	var lTemp, mTemp, rTemp int
	var texas *RoomInfo_Texas

	defer common.PrintPanicStack()
	texas = GetTexasByRoomid(roomID)                                                  //根据房间ID获取德州扑克结构
	room, has := GetRoomById(roomID)                                                  //根据房间ID获取房间结构体，用作下面判断房间是否关闭等作用
	for texas != nil && texas.IsStarted && has && room.Statue == common.ROOM_ONLIVE { //当房间不为空且标识值为TRUE时候

		if !common.GameRunningSwitch { //游戏在后台关闭时候通知所有用户
			NoticeAllUserInRoomGameClose(s)
			break
		}

		texas.LScore = 0 //左边押分值设为0
		lTemp = 0
		texas.MScore = 0 //中间押分值设为0
		mTemp = 0
		texas.RScore = 0 //右边押分值设为0
		rTemp = 0
		gameStatTime := time.Now().Unix() //记录下游戏开始的时间戳

		actors, anyCards, cards := CreatTexasPokActors(3) //构造出三个玩家对象

		texas.LargeActors = getTexasMaxActors(actors, cards) //计算出最大的玩家[0-2]
		for _, m := range actors {
			for i := 0; i < 5; i++ {
				if m.Cards[i].Number == 14 {
					m.Cards[i].Number = 1
				}
			}
		}
		texas.TexasActors = actors //全局标识记录下

		for i := 0; i < 5; i++ {
			if anyCards[i].Number == 14 {
				anyCards[i].Number = 1
			}
		}
		texas.AnyCards = anyCards

		for i := 0; i < 6; i++ {
			if cards[i].Number == 14 {
				cards[i].Number = 1
			}
		}
		texas.StaticCards = cards

		gameID := fmt.Sprintf("%d_%d", anchorID, time.Now().Unix())
		texas.GameID = gameID //记录下游戏ID

		NoticeMsgToRoomWithGameID(s, common.MESSAGE_TYPE_GAME_GOING, gameID) //通知房间里所有用户游戏开始
		texas.GameState = common.GAME_GOING
		texas.CSStartTime = time.Now().Unix()

		f := timer.NewDispatcher(1)
		f.AfterFunc(4*time.Second, func() { //4秒后通知所有玩家
			texas.GameState = common.GAME_CAN_RAISE //标识出可以押分了
			texas.CSStartTime = time.Now().Unix()   //记录下押分状态的开始时间戳
			NoticeMsgBeginTexasRaise(s, roomID)     //通知所有用户德州任意牌中的明牌

		})
		(<-f.ChanTimer).Cb()

		flag := time.Now().Unix() + common.GAME_TEXAS_PORK_RAISE //记录下押注状态结束的时间戳，方便后续计算
		for time.Now().Unix() < flag {
			f2 := timer.NewDispatcher(1)
			f2.AfterFunc(1*time.Second, func() {
				//当有用户押分时候通知所有用户左中右三家押注情况的消息
				if lTemp != texas.LScore || mTemp != texas.MScore || rTemp != texas.RScore {
					NoticeLMRScoreMsgToRoom(s, common.MESSAGE_TYPE_GAME_RAISE_SCORE,
						texas.LScore, texas.MScore, texas.RScore)
				}
				lTemp = texas.LScore
				mTemp = texas.MScore
				rTemp = texas.RScore
			})
			(<-f2.ChanTimer).Cb()
		}

		//30秒后通知所有用户押分结束
		f3 := timer.NewDispatcher(1)
		f3.AfterFunc(1*time.Second, func() {
			texas.GameState = common.GAME_RAISE_END
			texas.CSStartTime = time.Now().Unix()
			NoticeMsgToRoomWithTime(s, common.MESSAGE_TYPE_GAME_RAISE_END, common.GAME_TEXAS_WAIT_TIME_RESULT)
		})
		(<-f3.ChanTimer).Cb()

		//押分结束后通知所有用户游戏结束并告知本局结果
		f4 := timer.NewDispatcher(1)
		f4.AfterFunc(common.GAME_TEXAS_WAIT_TIME_RESULT*time.Second, func() { //3秒后执行
			texas.GameState = common.GAME_RESULT
			texas.CSStartTime = time.Now().Unix()
			NoticeTexasResult(s, roomID, common.GAME_TEXAS_LOOK_RESULT)

			//游戏结束时候后台根据用户押分情况进行结算并告知对应的用户
			m := ResultBetReward(gameID, texas.LargeActors, common.GAME_BONUS_TIMES, roomID, gameStatTime, anchorID)
			ms := common.NewMapSorter(m)
			sort.Sort(ms)
			NoticeRaiselistMsgToAnchor(anchorID, common.MESSAGE_TYPE_GAME_WINNER_SORTER, ms)
		})
		(<-f4.ChanTimer).Cb()

		f5 := timer.NewDispatcher(1)
		f5.AfterFunc(common.GAME_TEXAS_LOOK_RESULT*time.Second, func() { //所有用户查看结果
		})
		(<-f5.ChanTimer).Cb()

		texas.GameState = 0 //房间里的游戏状态值设为默认，准备下一局

		room, has = GetRoomById(roomID)
	}
}

type ResponseTexasRaise struct { //牛牛押注
	MType         int  `json:"mtype"`
	OverTime      int  `json:"over_time"`
	AnyCard0      Card `json:"any_c0"`
	Player0_Score int  `json:"p0"`
	Player1_Score int  `json:"p1"`
	Player2_Score int  `json:"p2"`
	MyRaise0      int  `json:"m_r0"` //用户对第一个位置的押分
	MyRaise1      int  `json:"m_r1"` //用户对第二个位置的押分
	MyRaise2      int  `json:"m_r2"` //用户对第三个位置的押分
}

type ResponseTexasRaise2 struct { //牛牛押注
	MType         int  `json:"mtype"`
	OverTime      int  `json:"over_time"`
	BonusTimes    int  `json:"bonus_times"`
	AnyCard0      Card `json:"any_c0"`
	Player0_Score int  `json:"p0"`
	Player1_Score int  `json:"p1"`
	Player2_Score int  `json:"p2"`
	MyRaise0      int  `json:"m_r0"` //用户对第一个位置的押分
	MyRaise1      int  `json:"m_r1"` //用户对第二个位置的押分
	MyRaise2      int  `json:"m_r2"` //用户对第三个位置的押分
}

// 通知客户端游戏可开始押注
func NoticeMsgBeginTexasRaise(s *melody.Session, roomID string) {
	var data ResponseTexasRaise
	data.MType = common.MESSAGE_TYPE_GAME_CAN_RAISE
	data.OverTime = common.GAME_TEXAS_PORK_RAISE
	data.Player0_Score = 0
	data.Player1_Score = 0
	data.Player2_Score = 0
	data.MyRaise0 = 0
	data.MyRaise1 = 0
	data.MyRaise2 = 0
	t := GetTexasByRoomid(roomID)
	if t != nil {
		data.AnyCard0.Number = t.AnyCards[0].Number
		data.AnyCard0.Color = t.AnyCards[0].Color
	}

	chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
		return q.Request.URL.Path == s.Request.URL.Path
	})
}

type TexasPokReult struct {
	User1       TexasActorInfo `json:"user1"`
	User2       TexasActorInfo `json:"user2"`
	User3       TexasActorInfo `json:"user3"`
	AnyCards    []*Card        `json:"any_cards"` //当前局下的任意牌
	StaticCards []*Card        `json:"s_cards"`   //当前局下的固定牌
}

type ResponseTexasResult struct {
	MType    int           `json:"mtype"`
	OverTime int           `json:"over_time"`
	Puk      TexasPokReult `json:"puk"`
	Winner   int           `json:"winner"`
}

func NoticeTexasResult(s *melody.Session, roomID string, overTime int) {
	var data ResponseTexasResult
	data.MType = common.MESSAGE_TYPE_GAME_RESULT
	data.OverTime = overTime
	t := GetTexasByRoomid(roomID)
	if t != nil {
		data.Puk.User1.Cards = t.TexasActors[0].Cards
		data.Puk.User1.Grade = t.TexasActors[0].Grade
		data.Puk.User2.Cards = t.TexasActors[1].Cards
		data.Puk.User2.Grade = t.TexasActors[1].Grade
		data.Puk.User3.Cards = t.TexasActors[2].Cards
		data.Puk.User3.Grade = t.TexasActors[2].Grade
		data.Puk.AnyCards = t.AnyCards
		data.Puk.StaticCards = t.StaticCards
		data.Winner = t.LargeActors
	}

	chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
		return q.Request.URL.Path == s.Request.URL.Path
	})
}

//德州扑克押分消息处理
func TexasPokRaiseMsgDispose(s *melody.Session, msg []byte, roomID string, user *User) {
	js, err := simplejson.NewJson(msg)
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	room := GetTexasByRoomid(roomID)

	t := GetTexasByRoomid(roomID)
	if t == nil {
		return
	}

	score := js.Get("score") //获取客户端传递的分数值
	score_, err := score.Int()
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	if t.GameState != common.GAME_CAN_RAISE { //如果当前游戏已经不处于押分状态下
		var data ResponseRaiseResult
		data.MType = common.MESSAGE_TYPE_GAME_NOT_RAISE // 游戏不可押分的消息ID赋值给消息对象
		data.Player0_Score = room.LScore                //获取左边押分总分数
		data.Player1_Score = room.MScore                //获取中间押分总分数
		data.Player2_Score = room.RScore                //获取右边押分总分数
		//获取用户当前游戏局不同位置押分分值
		data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(room.GameID, user.Uid)
		data.Score = int(user.Score)

		SendMsgToUser(user.Uid, data) //发送当前游戏不处于押分状态的消息给客户端
		return
	} else if score_%10 != 0 || score_ > int(user.Score) || score_ <= 0 { //判断客户端传递的分值是否是一个有效值
		var data ResponseRaiseResult
		data.MType = common.MESSAGE_TYPE_GAME_RAISE_SCORE_ERROR
		data.Player0_Score = room.LScore //获取左边押分总分数
		data.Player1_Score = room.MScore //获取中间押分总分数
		data.Player2_Score = room.RScore //获取右边押分总分数
		//获取用户当前游戏局不同位置押分分值
		data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(room.GameID, user.Uid)
		data.Score = int(user.Score)

		SendMsgToUser(user.Uid, data) //不是有效值发送消息给客户，就不要正常押分了
		return
	}

	direction := js.Get("dir") //获取方向值
	direction_, err := direction.Int()
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	if direction_ == 0 { //如果是押的第一家
		error := AddGameBetByID(room.GameID, user.Uid, score_, 0) //存储押分数据到redis
		if error != nil {
			var data ResponseRaiseResult
			data.MType = common.MESSAGE_TYPE_GAME_USER_RAISE_EOR
			data.Player0_Score = room.LScore //获取左边押分总分数
			data.Player1_Score = room.MScore //获取中间押分总分数
			data.Player2_Score = room.RScore //获取右边押分总分数
			//获取用户当前游戏局不同位置押分分值
			data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(t.GameID, user.Uid)
			//记录下用户押分后当前的余额
			data.Score = int(user.Score)

			common.Log.Errf("orm err is 2 %s", err.Error())
			SendMsgToUser(user.Uid, data)
		} else {
			t.LScore += score_

			var data ResponseRaiseResult
			data.MType = common.MESSAGE_TYPE_GAME_STATE_USER_RAISING
			data.Player0_Score = room.LScore //获取左边押分总分数
			data.Player1_Score = room.MScore //获取中间押分总分数
			data.Player2_Score = room.RScore //获取右边押分总分数
			//获取用户当前游戏局不同位置押分分值
			data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(t.GameID, user.Uid)

			//user.subScore(int64(score_)) //数据库表里用户对应的分值减掉
			user.DelMoney(nil, common.MONEY_TYPE_SCORE, int64(score_), false)
			//记录下用户押分后当前的余额
			data.Score = int(user.Score)

			SendMsgToUser(user.Uid, data)
		}
	} else if direction_ == 1 {
		error := AddGameBetByID(room.GameID, user.Uid, score_, 1)
		if error != nil {
			var data ResponseRaiseResult
			data.MType = common.MESSAGE_TYPE_GAME_USER_RAISE_EOR
			data.Player0_Score = room.LScore //获取左边押分总分数
			data.Player1_Score = room.MScore //获取中间押分总分数
			data.Player2_Score = room.RScore //获取右边押分总分数
			//获取用户当前游戏局不同位置押分分值
			data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(t.GameID, user.Uid)

			data.Score = int(user.Score)
			common.Log.Errf("orm err is 2 %s", error.Error())
			SendMsgToUser(user.Uid, data)
		} else {
			t.MScore += score_

			var data ResponseRaiseResult
			data.MType = common.MESSAGE_TYPE_GAME_STATE_USER_RAISING
			data.Player0_Score = room.LScore //获取左边押分总分数
			data.Player1_Score = room.MScore //获取中间押分总分数
			data.Player2_Score = room.RScore //获取右边押分总分数
			//获取用户当前游戏局不同位置押分分值
			data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(t.GameID, user.Uid)

			//user.subScore(int64(score_))
			user.DelMoney(nil, common.MONEY_TYPE_SCORE, int64(score_), false)
			//记录下用户押分后当前的余额
			data.Score = int(user.Score)

			SendMsgToUser(user.Uid, data)
		}
	} else if direction_ == 2 {
		error := AddGameBetByID(room.GameID, user.Uid, score_, 2)
		if error != nil {
			var data ResponseRaiseResult
			data.MType = common.MESSAGE_TYPE_GAME_USER_RAISE_EOR
			data.Player0_Score = room.LScore //获取左边押分总分数
			data.Player1_Score = room.MScore //获取中间押分总分数
			data.Player2_Score = room.RScore //获取右边押分总分数
			//获取用户当前游戏局不同位置押分分值
			data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(t.GameID, user.Uid)

			data.Score = int(user.Score)
			common.Log.Errf("orm err is 2 %s", error.Error())

			SendMsgToUser(user.Uid, data)
		} else {
			t.RScore += score_

			var data ResponseRaiseResult
			data.MType = common.MESSAGE_TYPE_GAME_STATE_USER_RAISING
			data.Player0_Score = room.LScore //获取左边押分总分数
			data.Player1_Score = room.MScore //获取中间押分总分数
			data.Player2_Score = room.RScore //获取右边押分总分数

			//获取用户当前游戏局不同位置押分分值
			data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(t.GameID, user.Uid)

			//user.subScore(int64(score_))
			user.DelMoney(nil, common.MONEY_TYPE_SCORE, int64(score_), false)
			//记录下用户押分后当前的余额
			data.Score = int(user.Score)

			SendMsgToUser(user.Uid, data)
		}
	}
}

type ResponseTexasResult2 struct {
	MType      int           `json:"mtype"`
	OverTime   int           `json:"over_time"`
	BonusTimes int           `json:"bonus_times"`
	Puk        TexasPokReult `json:"puk"`
	Winner     int           `json:"winner"`
}

func NoticeTexasPokStateToUser(s *melody.Session, userID int, roomID string) {
	room := GetTexasByRoomid(roomID)
	if room == nil {
		return
	}

	switch room.GameState {
	case common.GAME_NOT_GOING: //游戏没开始状态时候
		var data ResponseNiuNiuState
		data.MType = common.MESSAGE_TYPE_GAME_NOT_GOING
		data.State = common.GAME_NOT_GOING

		//b, _ := json.Marshal(data)
		//godump.Dump(string([]byte(b)))
		SendMsgToUser(userID, data)
	case common.GAME_GOING: //游戏开始时候
		var data ResponseWithGameID
		data.MType = common.MESSAGE_TYPE_GAME_GOING
		data.BonusTimes = common.GAME_BONUS_TIMES
		data.GameID = room.GameID

		SendMsgToUser(userID, data)
	case common.GAME_CAN_RAISE: //游戏处于可押分状态下
		var data ResponseTexasRaise2
		data.MType = common.MESSAGE_TYPE_GAME_CAN_RAISE
		data.BonusTimes = common.GAME_BONUS_TIMES
		data.OverTime = common.GAME_NIUNIU_RAISE - int(time.Now().Unix()-room.CSStartTime)
		if data.OverTime < 0 {
			data.OverTime = 0
		}
		data.Player0_Score = 0
		data.Player1_Score = 0
		data.Player2_Score = 0
		data.MyRaise0 = 0
		data.MyRaise1 = 0
		data.MyRaise2 = 0
		data.AnyCard0.Number = room.AnyCards[0].Number
		data.AnyCard0.Color = room.AnyCards[0].Color

		data.Player0_Score = room.LScore
		data.Player1_Score = room.MScore
		data.Player2_Score = room.RScore
		data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(room.GameID, userID)

		SendMsgToUser(userID, data)
	case common.GAME_RAISE_END: //游戏押分结束时候，告知玩家扑克牌，并告知几家的押分值
		var data ResponseTexasRaise2
		data.MType = common.MESSAGE_TYPE_GAME_RAISE_END
		data.BonusTimes = common.GAME_BONUS_TIMES
		data.OverTime = common.GAME_NIUNIU_WAIT_TIME_RESULT - int(time.Now().Unix()-room.CSStartTime)
		if data.OverTime < 0 {
			data.OverTime = 0
		}
		data.Player0_Score = room.LScore
		data.Player1_Score = room.MScore
		data.Player2_Score = room.RScore
		data.MyRaise0 = 0
		data.MyRaise1 = 0
		data.MyRaise2 = 0
		data.AnyCard0.Number = room.AnyCards[0].Number
		data.AnyCard0.Color = room.AnyCards[0].Color
		//获取用户当前游戏局不同位置押分分值
		data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(room.GameID, userID)

		SendMsgToUser(userID, data)
	case common.GAME_RESULT: //游戏结束的时候，告知几家的牌，并告知客户端几个玩家扑克对应的牛牛等级
		var data ResponseTexasResult2
		data.MType = common.MESSAGE_TYPE_GAME_RESULT
		data.BonusTimes = common.GAME_BONUS_TIMES
		data.OverTime = common.GAME_TEXAS_LOOK_RESULT - int(time.Now().Unix()-room.CSStartTime)
		if data.OverTime < 0 {
			data.OverTime = 0
		}

		data.Puk.User1.Cards = room.TexasActors[0].Cards
		data.Puk.User1.Grade = room.TexasActors[0].Grade
		data.Puk.User2.Cards = room.TexasActors[1].Cards
		data.Puk.User2.Grade = room.TexasActors[1].Grade
		data.Puk.User3.Cards = room.TexasActors[2].Cards
		data.Puk.User3.Grade = room.TexasActors[2].Grade
		data.Puk.AnyCards = room.AnyCards
		data.Puk.StaticCards = room.StaticCards
		data.Winner = room.LargeActors

		SendMsgToUser(userID, data)
	}
}

//返回德州扑克最大组合者下标maxIndex
func getTexasMaxActors(actors []*TexasActorInfo, card []*Card) int {
	maxIndex := 0
	i := 0
	temp := &TexasActorInfo{}
	var card_m [2]Card
	var cardTemp [2]Card
	for _, m := range actors {
		if i == 0 {
			temp = m
			cardTemp[0].Number = card[0].Number
			cardTemp[0].Color = card[0].Color
			cardTemp[1].Number = card[1].Number
			cardTemp[1].Color = card[1].Color
		} else { //先比较等级

			card_m[0].Number = card[i*2].Number
			card_m[0].Color = card[i*2].Color
			card_m[1].Number = card[i*2+1].Number
			card_m[1].Color = card[i*2+1].Color

			flag := CompareTwoTexasActor(*m, *temp, card_m, cardTemp)
			//fmt.Printf("\n flag=%t", flag)
			if flag {
				temp = m
				cardTemp = card_m
				maxIndex = i
			}
		}

		i++
	}

	return maxIndex
}

func PrintTexasPok(c Card) { //打印牛的信息
	switch c.Color {
	case 0:
		fmt.Print("方块")
	case 1:
		fmt.Print("梅花")
	case 2:
		fmt.Print("红桃")
	case 3:
		fmt.Print("黑桃")
	}

	switch c.Number {
	case 2:
		fmt.Print("2")
	case 3:
		fmt.Print("3")
	case 4:
		fmt.Print("4")
	case 5:
		fmt.Print("5")
	case 6:
		fmt.Print("6")
	case 7:
		fmt.Print("7")
	case 8:
		fmt.Print("8")
	case 9:
		fmt.Print("9")
	case 10:
		fmt.Print("10")
	case 11:
		fmt.Print("J")
	case 12:
		fmt.Print("Q")
	case 13:
		fmt.Print("K")
	case 14:
		fmt.Print("A")
	}
	fmt.Print(" ")
}

//根据传入的游戏角色数量（几个角色在玩扑克）生成11个数，依次取两张赋给对应用户做任意牌，最后5张作为公共牌
//返回每个用户最大的组合
func CreatTexasPokActors(numberOfPlayer int) ([]*TexasActorInfo, []*Card, []*Card) {
	var i, grageTemp int

	Random := common.RandomRangeArr(0, 51, 2*numberOfPlayer+5)

	texasActors := make([]*TexasActorInfo, 0)
	anyCards := make([]*Card, 0)
	cards := make([]*Card, 0)

	for k := 0; k < numberOfPlayer; k++ { //
		var c [7]Card //临时数组，存储每个玩家前面两张牌和从5张任意牌里取3张去取等级
		temp := &TexasActorInfo{}

		//依次取出每个角色的7张牌
		c[0] = TexasPork[Random[k*2+5]]
		c[1] = TexasPork[Random[k*2+6]]
		c[2] = TexasPork[Random[0]]
		c[3] = TexasPork[Random[1]]
		c[4] = TexasPork[Random[2]]
		c[5] = TexasPork[Random[3]]
		c[6] = TexasPork[Random[4]]

		grageTemp = 0
		// 从7张牌里取出3张来
		temp.Cards[0] = c[Texas5Numbers[0].N1]
		temp.Cards[1] = c[Texas5Numbers[0].N2]
		temp.Cards[2] = c[Texas5Numbers[0].N3]
		temp.Cards[3] = c[Texas5Numbers[0].N4]
		temp.Cards[4] = c[Texas5Numbers[0].N5]

		//以下的操作是从每个用户的两张牌加上从前面5张牌组合里取5张依次排列组合并给出最大的等级
		for i = 0; i < 21; i++ {
			var cTemp [5]Card //临时数组，存储每个玩家前面两张牌和从5张任意牌里取3张去取等级
			cTemp[0] = c[Texas5Numbers[i].N1]
			cTemp[1] = c[Texas5Numbers[i].N2]
			cTemp[2] = c[Texas5Numbers[i].N3]
			cTemp[3] = c[Texas5Numbers[i].N4]
			cTemp[4] = c[Texas5Numbers[i].N5]

			if i == 0 {
				temp.Cards = cTemp
			}

			//获取临时组合牌的等级
			grageTemp_ := GetGradeOfTexasPok(cTemp)
			if grageTemp_ > grageTemp { //先比较等级
				grageTemp = grageTemp_
				temp.Cards = cTemp
			} else if grageTemp_ == grageTemp { //同一等级时候再比较同牌型大小，

				a := &TexasActorInfo{}
				a.Cards = cTemp
				a.Grade = grageTemp_

				temp.Grade = grageTemp

				//比较五张牌的大小
				if CompareTwoTexasActorWith5Cards(a, temp) {
					grageTemp = grageTemp_
					temp.Cards = cTemp
				}

			}

		}

		//存储等级
		temp.Grade = grageTemp

		texasActors = append(texasActors, temp) //追加到texasActors
	}

	for i = 0; i < 5; i++ {
		tempCard := &Card{}
		tempCard.Number = TexasPork[Random[i]].Number
		tempCard.Color = TexasPork[Random[i]].Color
		anyCards = append(anyCards, tempCard)
	}

	// 每个角色的自由牌存储
	for i = 0; i < 6; i++ {
		tempCard := &Card{}
		tempCard.Number = TexasPork[Random[i+5]].Number
		tempCard.Color = TexasPork[Random[i+5]].Color
		cards = append(cards, tempCard)
	}

	return texasActors, anyCards, cards //返回texasActors对象
}

// 根据传入的5张牌(参数c)获返回对应的德州等级
func GetGradeOfTexasPok(c [5]Card) int {
	var flag, i, j int
	var cTemp Card

	//排序
	for i = 0; i < 4; i++ {
		for j = i + 1; j < 5; j++ {
			if c[i].Number > c[j].Number {
				cTemp = c[i]
				c[i] = c[j]
				c[j] = cTemp
			}
		}
	}

	//记录重复的牌个数
	flag = 0
	for i = 0; i < 5; i++ {
		for j = i + 1; j < 5; j++ {
			if c[i].Number == c[j].Number {
				flag++
			}
		}
	}

	switch flag {
	case 6:
		//fmt.Println("四条")
		return FOUR_OF_A_KIND
	case 4:
		//fmt.Println("葫芦")
		return FULL_HOUSE
	case 3:
		//fmt.Println("三条")
		return THREE_OF_A_KIND
	case 2:
		//fmt.Println("两对")
		return TWO_PAIR
	case 1:
		//fmt.Println("一对")
		return ONE_PAIR
	case 0:
		//fmt.Println("高牌")
		isStraight := true
		for i = 0; i < 4; i++ {
			if c[i].Number+1 != c[i+1].Number {
				isStraight = false
			}
		}

		// 如果五张牌依次相差1的情况，表示是顺子哟
		if isStraight {
			if c[0].Color == c[1].Color && c[0].Color == c[2].Color && c[0].Color == c[3].Color &&
				c[0].Color == c[4].Color { //五张牌花色一样
				if c[0].Number == 10 {
					return ROYAL_FLUSH //皇家同花顺
				} else {
					return STRAIGHT_FLUSH //同花顺
				}
			} else {
				return STRAIGHT //普通的顺子
			}
		} else if c[0].Color == c[1].Color && c[0].Color == c[2].Color && c[0].Color == c[3].Color &&
			c[0].Color == c[4].Color {
			return FLUSH //同花顺
		}
		return HIG_CARD //高牌
	}
	return HIG_CARD //高牌
}

// 比较两家德州扑克玩家大小
// 如果a>b返回TRUE,否者返回false
func CompareTwoTexasActor(a, b TexasActorInfo, cardsA, cardsB [2]Card) bool {
	var i, j, k int
	var aMax, bMax Card

	if a.Grade > b.Grade {
		return true
	} else if a.Grade < b.Grade {
		return false
	}

	aMax = cardsA[0] //获取第一家手上牌最大的
	if cardsA[1].Number > cardsA[0].Number ||
		(cardsA[1].Number == cardsA[0].Number && cardsA[1].Color >= cardsA[0].Color) {
		aMax = cardsA[1]
	}

	bMax = cardsB[0] //获取第二家手上牌最大的
	if cardsB[1].Number > cardsB[0].Number ||
		(cardsB[1].Number == cardsB[0].Number && cardsB[1].Color >= cardsB[0].Color) {
		bMax = cardsB[1]
	}
	/*** 以下部分代码是两家扑克等级相同情况下对两家进行比较大小***/

	if a.Grade == HIG_CARD || a.Cards == b.Cards { //如果都为高牌情况,比较两家手上2牌最大的一张
		//比较两家最大的牌
		if aMax.Number > bMax.Number || ((aMax.Number == bMax.Number) && aMax.Color > bMax.Color) {
			return true
		}
		return false
	}

	// 对玩家a的5张牌从大到小排序。
	for i = 0; i < 4; i++ {
		for j = i + 1; j < 5; j++ {
			if a.Cards[i].Number < a.Cards[j].Number ||
				((a.Cards[i].Number == a.Cards[j].Number) && a.Cards[i].Color < a.Cards[j].Color) {
				temp := a.Cards[i]
				a.Cards[i] = a.Cards[j]
				a.Cards[j] = temp
			}
		}
	}

	// 对玩家b的5张牌从大到小进行排序。
	for i = 0; i < 4; i++ {
		for j = i + 1; j < 5; j++ {
			if b.Cards[i].Number < b.Cards[j].Number ||
				((b.Cards[i].Number == b.Cards[j].Number) && b.Cards[i].Color < b.Cards[j].Color) {
				temp := b.Cards[i]
				b.Cards[i] = b.Cards[j]
				b.Cards[j] = temp
			}
		}
	}

	// 如果是同花顺、顺子、同花，比较排序后最大的扑克就能比较出大小
	if a.Grade == STRAIGHT_FLUSH || a.Grade == STRAIGHT || a.Grade == FLUSH {
		//if a.Cards[0].Number > b.Cards[0].Number ||
		//	((a.Cards[0].Number == b.Cards[0].Number) && a.Cards[0].Color > b.Cards[0].Color) {
		//	fmt.Printf("\n aaaaaaaaaaa")
		//	return true
		//}
		for m := 0; m < 5; m++ {
			aTemp := a.Cards[m]
			bTemp := b.Cards[m]
			if aTemp.Number > bTemp.Number || ((aTemp.Number == bTemp.Number) && aTemp.Color > bTemp.Color) {
				return true
			} else if aTemp.Number < bTemp.Number {
				return false
			}
		}

		return false
	}

	index := 0 //记录有多少个唯一的数字挪动到最后面去了
	for i = 0; i < 5; i++ {
		flag := false
		for j = 0; j < 5; j++ {
			if a.Cards[index].Number == a.Cards[j].Number && index != j {
				flag = true
			}
		}
		if flag {
			index++
		} else {
			temp := a.Cards[index]
			for k = index + 1; k < 5; k++ {
				a.Cards[k-1] = a.Cards[k]
			}
			a.Cards[k-1] = temp
		}
	}

	index = 0 //记录有多少个唯一的数字挪动到最后面去了
	for i = 0; i < 5; i++ {
		flag := false
		for j = 0; j < 5; j++ {
			if b.Cards[index].Number == b.Cards[j].Number && index != j {
				flag = true
			}
		}
		if flag {
			index++
		} else {
			temp := b.Cards[index]
			for k = index + 1; k < 5; k++ {
				b.Cards[k-1] = b.Cards[k]
			}
			b.Cards[k-1] = temp
		}
	}

	if a.Grade == FULL_HOUSE { //葫芦要把三张一样的牌移到最前面
		if a.Cards[3].Number == a.Cards[2].Number {
			temp0 := a.Cards[0]
			temp1 := a.Cards[1]
			a.Cards[0] = a.Cards[2]
			a.Cards[1] = a.Cards[3]
			a.Cards[2] = a.Cards[4]
			a.Cards[3] = temp0
			a.Cards[4] = temp1
		}
	}
	if b.Grade == FULL_HOUSE { //葫芦要把三张一样的牌移到最前面
		if b.Cards[3].Number == b.Cards[2].Number {
			temp0 := b.Cards[0]
			temp1 := b.Cards[1]
			b.Cards[0] = b.Cards[2]
			b.Cards[1] = b.Cards[3]
			b.Cards[2] = b.Cards[4]
			b.Cards[3] = temp0
			b.Cards[4] = temp1
		}
	}

	// 如果两家都是葫芦那就根据葫芦来比较大小吧
	if a.Grade == FULL_HOUSE || a.Grade == THREE_OF_A_KIND {
		if a.Cards[0].Number == a.Cards[1].Number && a.Cards[0].Number == a.Cards[2].Number { //取出a的葫芦里最大的那张牌
			aMax = a.Cards[0]
		} else {
			aMax = a.Cards[3]
		}

		if b.Cards[0].Number == b.Cards[1].Number && b.Cards[0].Number == b.Cards[2].Number { //取出b的葫芦里最大的那张牌
			bMax = b.Cards[0]
		} else {
			bMax = b.Cards[3]
		}

		//比较两家最大的牌
		if aMax.Number > bMax.Number || ((aMax.Number == bMax.Number) && aMax.Color > bMax.Color) {
			return true
		} else if aMax.Number == bMax.Number && aMax.Color == bMax.Color {
			for m := 0; m < 5; m++ {
				aTemp := a.Cards[m]
				bTemp := b.Cards[m]
				if aTemp.Number > bTemp.Number || ((aTemp.Number == bTemp.Number) && aTemp.Color > bTemp.Color) {
					return true
				} else if aTemp.Number < bTemp.Number {
					return false
				}
			}
		}

		//比较两家最大的牌
		if aMax.Number > bMax.Number || ((aMax.Number == bMax.Number) && aMax.Color > bMax.Color) {
			return true
		}
		return false
	} else if a.Grade == ONE_PAIR || a.Grade == TWO_PAIR || a.Grade == HIG_CARD {
		for m := 0; m < 5; m++ {
			aTemp := a.Cards[m]
			bTemp := b.Cards[m]
			if aTemp.Number > bTemp.Number || ((aTemp.Number == bTemp.Number) && aTemp.Color > bTemp.Color) {
				return true
			} else if aTemp.Number < bTemp.Number {
				return false
			}
		}

		//如果上面还比较不出来大小就比较两家底牌最大的牌
		if aMax.Number > bMax.Number || ((aMax.Number == bMax.Number) && aMax.Color > bMax.Color) {
			return true
		}
		return false
	}

	/***以下是其它情况比较大小了，只需比较排序和分组后的第一张牌就可以比较出大小了***/

	aMaxTemp := a.Cards[0]
	bMaxTemp := b.Cards[0]
	//比较两家最大的牌
	if aMaxTemp.Number > bMaxTemp.Number || ((aMaxTemp.Number == bMaxTemp.Number) && aMaxTemp.Color > bMaxTemp.Color) {
		return true
	} else if aMaxTemp.Number < bMaxTemp.Number || ((aMaxTemp.Number == bMaxTemp.Number) && aMaxTemp.Color < bMaxTemp.Color) {
		return false
	}
	//第一张牌比较不出来大小那就比较每家底牌那两张里的大牌勒
	if aMax.Number > bMax.Number || ((aMax.Number == bMax.Number) && aMax.Color > bMax.Color) {
		return true
	}
	return false
}

// 比较两家德州扑克玩家大小
// 如果a>b返回TRUE,否者返回false
func CompareTwoTexasActorWith5Cards(a, b *TexasActorInfo) bool {
	var i, j, k int
	var aMax, bMax Card

	if a.Grade > b.Grade {
		return true
	} else if a.Grade < b.Grade {
		return false
	}

	// 对玩家a的5张牌从大到小排序。
	for i = 0; i < 4; i++ {
		for j = i + 1; j < 5; j++ {
			if a.Cards[i].Number < a.Cards[j].Number ||
				((a.Cards[i].Number == a.Cards[j].Number) && a.Cards[i].Color < a.Cards[j].Color) {
				temp := a.Cards[i]
				a.Cards[i] = a.Cards[j]
				a.Cards[j] = temp
			}
		}
	}

	// 对玩家b的5张牌从大到小进行排序。
	for i = 0; i < 4; i++ {
		for j = i + 1; j < 5; j++ {
			if b.Cards[i].Number < b.Cards[j].Number ||
				((b.Cards[i].Number == b.Cards[j].Number) && b.Cards[i].Color < b.Cards[j].Color) {
				temp := b.Cards[i]
				b.Cards[i] = b.Cards[j]
				b.Cards[j] = temp
			}
		}
	}

	// 如果是同花顺、顺子、同花，比较排序后最大的扑克就能比较出大小
	if a.Grade == STRAIGHT_FLUSH || a.Grade == STRAIGHT || a.Grade == FLUSH {
		if a.Cards[0].Number > b.Cards[0].Number ||
			((a.Cards[0].Number == b.Cards[0].Number) && a.Cards[0].Color > b.Cards[0].Color) {
			return true
		}
		return false
	}

	index := 0 //记录有多少个唯一的数字挪动到最后面去了
	for i = 0; i < 5; i++ {
		flag := false
		for j = 0; j < 5; j++ {
			if a.Cards[index].Number == a.Cards[j].Number && index != j {
				flag = true
			}
		}
		if flag {
			index++
		} else {
			temp := a.Cards[index]
			for k = index + 1; k < 5; k++ {
				a.Cards[k-1] = a.Cards[k]
			}
			a.Cards[k-1] = temp
		}
	}

	index = 0 //index记录有多少个唯一的数字挪动到最后面去了
	//如果有重复的牌往前挪
	for i = 0; i < 5; i++ {
		flag := false
		for j = 0; j < 5; j++ {
			if b.Cards[index].Number == b.Cards[j].Number && index != j {
				flag = true
			}
		}
		if flag {
			index++
		} else {
			temp := b.Cards[index]
			for k = index + 1; k < 5; k++ {
				b.Cards[k-1] = b.Cards[k]
			}
			b.Cards[k-1] = temp
		}
	}

	//葫芦要把三张一样的牌移到最前面
	if a.Grade == FULL_HOUSE {
		if a.Cards[3].Number == a.Cards[2].Number {
			temp0 := a.Cards[0]
			temp1 := a.Cards[1]
			a.Cards[0] = a.Cards[2]
			a.Cards[1] = a.Cards[3]
			a.Cards[2] = a.Cards[4]
			a.Cards[3] = temp0
			a.Cards[4] = temp1
		}
	}

	//葫芦要把三张一样的牌移到最前面
	if b.Grade == FULL_HOUSE {
		if b.Cards[3].Number == b.Cards[2].Number {
			temp0 := b.Cards[0]
			temp1 := b.Cards[1]
			b.Cards[0] = b.Cards[2]
			b.Cards[1] = b.Cards[3]
			b.Cards[2] = b.Cards[4]
			b.Cards[3] = temp0
			b.Cards[4] = temp1
		}
	}

	if a.Grade == FULL_HOUSE || a.Grade == THREE_OF_A_KIND { // 如果两家都是葫芦那就根据葫芦来比较大小吧
		if a.Cards[0].Number == a.Cards[1].Number && a.Cards[0].Number == a.Cards[2].Number { //取出a的葫芦里最大的那张牌
			aMax = a.Cards[0]
		} else {
			aMax = a.Cards[3]
		}

		if b.Cards[0].Number == b.Cards[1].Number && b.Cards[0].Number == b.Cards[2].Number { //取出b的葫芦里最大的那张牌
			bMax = b.Cards[0]
		} else {
			bMax = b.Cards[3]
		}

		//比较两家最大的牌
		if aMax.Number > bMax.Number || ((aMax.Number == bMax.Number) && aMax.Color > bMax.Color) {
			return true
		} else if aMax.Number == bMax.Number && aMax.Color == bMax.Color {
			for m := 0; m < 5; m++ {
				aTemp := a.Cards[m]
				bTemp := b.Cards[m]
				if aTemp.Number > bTemp.Number || ((aTemp.Number == bTemp.Number) && aTemp.Color > bTemp.Color) {
					return true
				} else if aTemp.Number < bTemp.Number {
					return false
				}
			}
		}
		return false
	} else if a.Grade == ONE_PAIR || a.Grade == TWO_PAIR || a.Grade == HIG_CARD {
		for m := 0; m < 5; m++ {
			aTemp := a.Cards[m]
			bTemp := b.Cards[m]
			if aTemp.Number > bTemp.Number || ((aTemp.Number == bTemp.Number) && aTemp.Color > bTemp.Color) {
				return true
			} else if aTemp.Number < bTemp.Number {
				return false
			}
		}
		return true
	}

	/***以下是其它情况比较大小了，只需比较排序和分组后的第一张牌就可以比较出大小了***/

	aMax = a.Cards[0]
	bMax = b.Cards[0]
	//比较两家最大的牌
	if aMax.Number > bMax.Number || ((aMax.Number == bMax.Number) && aMax.Color > bMax.Color) {
		return true
	}
	return false
}
