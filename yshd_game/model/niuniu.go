/*****************************************************************************
 *
 * Filename:
 * ---------
 *    niuniu.c
 *
 * Project:
 * --------
 *   game_niuniu
 *
 * Description:
 * ------------
 *   本文件的函数为牛牛核心代码，主要用于实现牛牛业务逻辑以及牛牛通信协议。
 *
 * Author:
 * sky Zeng
 *
 * Date:
 * 2017-04-09
 *
 ****************************************************************************/
package model

import (
	//"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/yshd_game/common"
	"github.com/yshd_game/melody"
	"github.com/yshd_game/timer"
	"sort"
	"strings"
	//"sync"
	"time"
	//"github.com/liudng/godump"
	"sync"
)

const (
	NIU_ZERO  = iota //0没牛
	NIU_ONE          //1.牛一
	NIU_TWO          //2.牛二
	NIU_THREE        //3.牛三
	NIU_FOUR         //4.牛四
	NIU_FIVE         //5.牛五
	NIU_SIX          //6.牛六
	NIU_SEVEN        //7.牛七
	NIU_EIGHT        //8.牛八
	NIU_NINE         //9.牛九
	NIU_NIU          //10.牛牛
	YIN_NIU          //11.银牛
	JIN_NIU          //12.金牛
	BOMB_A           //13.炸弹A
	BOMB_2           //14.炸弹2
	BOMB_3           //15.炸弹3
	BOMB_4           //16.炸弹4
	BOMB_5           //17.炸弹5
	BOMB_6           //18.炸弹6
	BOMB_7           //19.炸弹7
	BOMB_8           //20.炸弹8
	BOMB_9           //21.炸弹9
	BOMB_10          //22.炸弹10
	BOMB_J           //23.炸弹J
	BOMB_Q           //24.炸弹Q
	BOMB_K           //25.炸弹K
)

const ( //黑桃＞红桃＞草花＞方块
	COLOR_BLOCK  = iota //0-方块
	COLOR_CLUB          //1-草花
	COLOR_HEARDS        // 2-红桃
	COLOR_SPADE         //3-黑桃
)

//牛牛扑克的初始数据结构
var Pork = [52]Card{
	{1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}, {6, 0}, {7, 0}, {8, 0}, {9, 0}, {10, 0}, {11, 0}, {12, 0}, {13, 0},
	{1, 1}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {6, 1}, {7, 1}, {8, 1}, {9, 1}, {10, 1}, {11, 1}, {12, 1}, {13, 1},
	{1, 2}, {2, 2}, {3, 2}, {4, 2}, {5, 2}, {6, 2}, {7, 2}, {8, 2}, {9, 2}, {10, 2}, {11, 2}, {12, 2}, {13, 2},
	{1, 3}, {2, 3}, {3, 3}, {4, 3}, {5, 3}, {6, 3}, {7, 3}, {8, 3}, {9, 3}, {10, 3}, {11, 3}, {12, 3}, {13, 3}}

//每一张牌
type Card struct {
	Number int `json:"n"` //点数[1-13]
	Color  int `json:"c"` //花色
}

//左中右每一家的相关信息
type NiuActorInfo struct {
	Cards   [5]Card //5张牌
	MIndex  int     //最大的一张牌的下标
	Grade   int     //5张牌对应的牛牛等级[0-10]参见牛牛级别常量
	Cardall int
}

//本局游戏信息
type NiuNiuGameInfo struct {
	Rid               string //房间ID
	Uid               int    //用户ID
	IsStarted         bool   //是否已开始
	mutex_NiuNiu_info sync.RWMutex
	GameState         int             //游戏状态【进行到哪一步了】
	CSStartTime       int64           //当前状态下开始时间戳
	GameID            string          //当前局的游戏ID
	niuActors         []*NiuActorInfo //当前局下的牌
	LScore            int             //最左边押分数据
	MScore            int             //中间押分数据
	RScore            int             //右边押分数据
	LargeActors       int             //哪家最大
}

var NiuNiu_room_manager map[string]*NiuNiuGameInfo
var mutex_NiuNiu sync.RWMutex

//新创建一牛牛游戏
func NewRoomInfoNiuNiu() *NiuNiuGameInfo {
	m := &NiuNiuGameInfo{}
	return m
}

//初始化牛牛Map
func InitNiuNiu() map[string]*NiuNiuGameInfo {
	NiuNiu_room_manager = make(map[string]*NiuNiuGameInfo)
	return NiuNiu_room_manager
}

// 根据房间号（roomid）获取牛牛游戏对象（NiuNiuGameInfo）
// roomid 直播的房间ID
func GetNiuNiuByRoomid(roomid string) *NiuNiuGameInfo {
	v, ok := NiuNiu_room_manager[roomid]
	if !ok {
		return nil
	}
	return v
}

//添加牛牛游戏对象到房间Map中
// r牛牛游戏对象
func AddRoomInfoNiuNiu(r *NiuNiuGameInfo) bool {
	_, ok := NiuNiu_room_manager[r.Rid]
	if ok {
		return false
	}

	mutex_NiuNiu.Lock()
	NiuNiu_room_manager[r.Rid] = r
	mutex_NiuNiu.Unlock()

	return true
}

// 根据房间ID从牛牛Map里删掉牛牛对象
// rid房间ID
func DelRoomInfoNiuNiu(rid string) {

	mutex_NiuNiu.Lock()
	delete(NiuNiu_room_manager, rid)
	mutex_NiuNiu.Unlock()
}

// 本函数主要用于牛牛游戏的创建【入口函数】
// s：长链接的Session对象，在游戏过程中用来发送消息
// anchorID：主播ID
// roomID：房间ID
func NewGameNiuNiu(s *melody.Session, anchorID int, roomID string) {
	//左中右三家分值的临时存储变量
	var lTemp, mTemp, rTemp int
	//牛牛游戏对象，用来保存每局牛牛游戏相关的数据
	var roomNiuNiu *NiuNiuGameInfo

	defer common.PrintPanicStack()

	//根据房间ID获取牛牛游戏对象
	roomNiuNiu = GetNiuNiuByRoomid(roomID)

	// 房间对象
	room, has := GetRoomById(roomID)

	//当房间不为空且标识值为TRUE时候
	for roomNiuNiu != nil && roomNiuNiu.IsStarted && has && room.Statue == common.ROOM_ONLIVE {

		if !common.GameRunningSwitch { //游戏在后台关闭时候通知所有用户
			NoticeAllUserInRoomGameClose(s)
			break
		}

		//左边押分值设为0
		roomNiuNiu.LScore = 0
		lTemp = 0
		//中间押分值设为0
		roomNiuNiu.MScore = 0
		mTemp = 0
		//右边押分值设为0
		roomNiuNiu.RScore = 0
		rTemp = 0

		gameStatTime := time.Now().Unix() //记录下游戏开始的时间戳

		randonArry := common.RandomRangeArr(5000, 10000, 1)
		randon := randonArry[0]
		if randon < 5000 || randon > 10000 {
			randon = 8500
		}

		//构造出三个玩家对象
		actors := CreateNiuActor(3)
		//计算出最大的玩家[0-2]
		roomNiuNiu.LargeActors = getMaxActors(actors)
		//全局标识记录下
		roomNiuNiu.niuActors = actors

		//fmt.Printf("\n")
		//PrintfNiuActor(actors[0])
		//PrintfNiuActor(actors[1])
		//PrintfNiuActor(actors[2])
		//fmt.Printf("刚开始 LargeActors=%d,randon=%d\n",roomNiuNiu.LargeActors,randon)
		//构造游戏ID
		gameID := fmt.Sprintf("%d_%d", anchorID, time.Now().Unix())
		roomNiuNiu.GameID = gameID //记录下游戏ID

		//通知房间里所有用户游戏开始
		NoticeMsgToRoomWithGameID(s, common.MESSAGE_TYPE_GAME_GOING, gameID)
		//存储游戏状态
		roomNiuNiu.GameState = common.GAME_GOING
		//记录下本状态下开始的时间戳，以免让中途有用户进入时候好计算倒计时
		roomNiuNiu.CSStartTime = time.Now().Unix()

		f := timer.NewDispatcher(1)
		f.AfterFunc(4*time.Second, func() { //4秒后通知所有玩家
			//标识出可以押分了
			roomNiuNiu.GameState = common.GAME_CAN_RAISE
			//记录下押分状态的开始时间戳
			roomNiuNiu.CSStartTime = time.Now().Unix()
			//通知所有用户分给三个玩家的第一张牌
			NoticeBeginNiuNiuRaise(s, actors[0].Cards[0], actors[1].Cards[0], actors[2].Cards[0])
		})
		(<-f.ChanTimer).Cb()

		flag := time.Now().Unix() + common.GAME_NIUNIU_RAISE //记录下押注状态结束的时间戳，方便后续计算
		for time.Now().Unix() < flag {
			f2 := timer.NewDispatcher(1)
			f2.AfterFunc(1*time.Second, func() {
				//当有用户押分时候通知所有用户左中右三家押注情况的消息
				if lTemp != roomNiuNiu.LScore || mTemp != roomNiuNiu.MScore || rTemp != roomNiuNiu.RScore {
					NoticeLMRScoreMsgToRoom(s, common.MESSAGE_TYPE_GAME_RAISE_SCORE,
						roomNiuNiu.LScore, roomNiuNiu.MScore, roomNiuNiu.RScore)
				}

				//记录啊下左中右三家当前押分数据，下次再跟当前记录的押分数据比较判断是否有变化。
				lTemp = roomNiuNiu.LScore
				mTemp = roomNiuNiu.MScore
				rTemp = roomNiuNiu.RScore
			})
			(<-f2.ChanTimer).Cb()
		}

		//30秒后通知所有用户押分结束
		f3 := timer.NewDispatcher(1)
		f3.AfterFunc(1*time.Second, func() {
			roomNiuNiu.GameState = common.GAME_RAISE_END //保存游戏状态
			roomNiuNiu.CSStartTime = time.Now().Unix()   //记录下本状态下开始的时间戳，以免让中途有用户进入时候好计算倒计时
			NoticeMsgToRoomWithTime(s, common.MESSAGE_TYPE_GAME_RAISE_END, common.GAME_NIUNIU_WAIT_TIME_RESULT)
		})
		(<-f3.ChanTimer).Cb()

		//押分结束后通知所有用户游戏结束并告知本局结果
		f4 := timer.NewDispatcher(1)
		f4.AfterFunc(common.GAME_NIUNIU_WAIT_TIME_RESULT*time.Second, func() { //10秒后执行
			roomNiuNiu.GameState = common.GAME_RESULT  //保存游戏状态
			roomNiuNiu.CSStartTime = time.Now().Unix() //记录下本状态下开始的时间戳，以免让中途有用户进入时候好计算倒计时

			RaiseMinIndex := GetMinIndexOfThree(lTemp, mTemp, rTemp)
			if lTemp+mTemp+rTemp >= randon && RaiseMinIndex != roomNiuNiu.LargeActors && GetRaiseSuccByPercent(NiuNiuPem) {
				for i := 0; i < 10; i++ {
					//for _, v := range actors {
					//	ReCreateNiuActor(v)
					//}
					ReCreateNiuActor(actors[0])
					ReCreateNiuActor(actors[1])
					ReCreateNiuActor(actors[2])
					//计算出最大的玩家[0-2]
					roomNiuNiu.LargeActors = getMaxActors(actors)
					if roomNiuNiu.LargeActors == RaiseMinIndex {
						break
					}
				}
			}
			//PrintfNiuActor(actors[0])
			//PrintfNiuActor(actors[1])
			//PrintfNiuActor(actors[2])
			//fmt.Printf("后来 LargeActors=%d\n",roomNiuNiu.LargeActors)
			//fmt.Printf("------------------------------------------------------------\n")

			NoticeNiuNiuResult(s, actors, common.GAME_NIUNIU_LOOK_RESULT) //通知所有用户游戏结果

			//游戏结束时候后台根据用户押分情况进行结算并告知对应的用户
			m := ResultBetReward(gameID, roomNiuNiu.LargeActors, common.GAME_BONUS_TIMES, roomID, gameStatTime, anchorID)
			ms := common.NewMapSorter(m)

			sort.Sort(ms)
			NoticeRaiselistMsgToAnchor(anchorID, common.MESSAGE_TYPE_GAME_WINNER_SORTER, ms)

		})
		(<-f4.ChanTimer).Cb()

		f5 := timer.NewDispatcher(1)
		f5.AfterFunc(common.GAME_NIUNIU_LOOK_RESULT*time.Second, func() { //所有用户查看结果
		})
		(<-f5.ChanTimer).Cb()

		roomNiuNiu.GameState = 0 //房间里的游戏状态值设为默认，准备下一局

		room, has = GetRoomById(roomID)
	}

}

type ResponseWithGameID struct {
	MType      int    `json:"mtype"`       // 消息ID
	GameID     string `json:"game_id"`     // 游戏ID
	BonusTimes int    `json:"bonus_times"` // 奖励倍数
}

//根据传入的msgID[游戏ID]和消息ID[msgID]将游戏ID和奖励倍数告知客户端
func NoticeMsgToRoomWithGameID(s *melody.Session, msgID int, gameID string) {
	var data ResponseWithGameID
	data.MType = msgID
	data.GameID = gameID
	data.BonusTimes = common.GAME_BONUS_TIMES

	chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
		return q.Request.URL.Path == s.Request.URL.Path
	})
}

type ResponseWithOverTime struct { //牛牛押注
	MType    int `json:"mtype"`     // 消息ID
	OverTime int `json:"over_time"` // 奖励倍数
}

//根据本一阶段剩余时间overTime[用来做倒计时的]通知所有玩家游戏即将开始
func NoticeMsgToRoomWithTime(s *melody.Session, msgID, overTime int) {
	var data ResponseWithOverTime
	data.MType = msgID
	data.OverTime = overTime

	chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
		return q.Request.URL.Path == s.Request.URL.Path
	})
}

type ResponseLMRScore struct {
	MType         int `json:"mtype"` // 消息ID
	Player0_Score int `json:"p0"`    //最左边所有用户的押分
	Player1_Score int `json:"p1"`    //第二个位置所有用户押分
	Player2_Score int `json:"p2"`    //第三个位置所有用户押分
}

//通知所有玩家游戏左中右三家的分数
func NoticeLMRScoreMsgToRoom(s *melody.Session, msgID, lScore, mScore, rScore int) {
	var data ResponseLMRScore
	data.MType = msgID
	data.Player0_Score = lScore //左边押的分数
	data.Player1_Score = mScore //中间押的分数
	data.Player2_Score = rScore //右边押的分数

	chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
		return q.Request.URL.Path == s.Request.URL.Path
	})
}

type Puk1 struct {
	User1 Card `json:"user1"` //用户1的牌
	User2 Card `json:"user2"` //用户2的牌
	User3 Card `json:"user3"` //用户3的牌
}

type ResponseNiuNiuRaise struct { //牛牛押注
	MType         int  `json:"mtype"`       //消息ID
	BonusTimes    int  `json:"bonus_times"` //奖励倍数
	OverTime      int  `json:"over_time"`   //倒计时
	Puk           Puk1 `json:"puk"`         //扑克信息
	Player0_Score int  `json:"p0"`          //最左边所有用户的押分
	Player1_Score int  `json:"p1"`          //第二个位置所有用户押分
	Player2_Score int  `json:"p2"`          //第三个位置所有用户押分
	MyRaise0      int  `json:"m_r0"`        //用户对第一个位置的押分
	MyRaise1      int  `json:"m_r1"`        //用户对第二个位置的押分
	MyRaise2      int  `json:"m_r2"`        //用户对第三个位置的押分
}

//通知客户端牛牛游戏可开始押注
func NoticeBeginNiuNiuRaise(s *melody.Session, c1, c2, c3 Card) {
	var data ResponseNiuNiuRaise
	data.MType = common.MESSAGE_TYPE_GAME_CAN_RAISE //消息ID
	data.BonusTimes = common.GAME_BONUS_TIMES       //奖励倍数
	data.OverTime = common.GAME_NIUNIU_RAISE        //倒计时
	data.Puk.User1 = c1                             //第一个用户对应的牌
	data.Puk.User2 = c2                             //第二个用户对应的牌
	data.Puk.User3 = c3                             //第三个用户对应的牌

	chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
		return q.Request.URL.Path == s.Request.URL.Path
	})
}

type UserNiuNiuPukReult struct {
	Cards [5]Card `json:"cards"` //牌
	Grade int     `json:"grade"` //牛牛等级
}

type PukReult struct { //游戏结果
	User1 UserNiuNiuPukReult `json:"user1"` //用户1
	User2 UserNiuNiuPukReult `json:"user2"` //用户2
	User3 UserNiuNiuPukReult `json:"user3"` //用户3
}
type ResponseNiuNiuResult struct { //牛牛押注
	MType      int      `json:"mtype"`       //消息ID
	BonusTimes int      `json:"bonus_times"` //奖励倍数
	OverTime   int      `json:"over_time"`   //倒计时
	Puk        PukReult `json:"puk"`         //扑克
	Winner     int      `json:"winner"`      //赢家
}

type ResponseNiuNiuRaiseResult struct { //牛牛押注
	MType         int  `json:"mtype"`       //消息ID
	BonusTimes    int  `json:"bonus_times"` //奖励倍数
	OverTime      int  `json:"over_time"`   //倒计时
	Puk           Puk1 `json:"puk"`         //扑克信息
	Player0_Score int  `json:"p0"`          //最左边所有用户的押分
	Player1_Score int  `json:"p1"`          //第二个位置所有用户押分
	Player2_Score int  `json:"p2"`          //第三个位置所有用户押分
}

//根据传入的NiuActorInfo对象和本阶段剩余时间[overTime]通知客户端牛牛游戏结果
func NoticeNiuNiuResult(s *melody.Session, actors []*NiuActorInfo, overTime int) {
	largeNumber := getMaxActors(actors)
	var data ResponseNiuNiuResult
	data.MType = common.MESSAGE_TYPE_GAME_RESULT
	data.BonusTimes = common.GAME_BONUS_TIMES
	data.OverTime = overTime
	data.Puk.User1.Cards = actors[0].Cards
	data.Puk.User1.Grade = actors[0].Grade
	data.Puk.User2.Cards = actors[1].Cards
	data.Puk.User2.Grade = actors[1].Grade
	data.Puk.User3.Cards = actors[2].Cards
	data.Puk.User3.Grade = actors[2].Grade
	data.Winner = largeNumber

	chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
		return q.Request.URL.Path == s.Request.URL.Path
	})
}

type ResponseNiuNiuState struct {
	MType      int `json:"mtype"`
	BonusTimes int `json:"bonus_times"`
	State      int `json:"state"`
}

// 用户在游戏中途进入房间时候根据游戏状态值发送对应消息给用户
// userID被通知者ID
// roomID所在直播房间ID
func NoticeNiuNiuStateToUser(s *melody.Session, userID int, roomID string) {
	//获取牛牛游戏对象
	room := GetNiuNiuByRoomid(roomID)
	if room == nil {
		return
	}

	switch room.GameState { //判断房间对象的状态
	case common.GAME_NOT_GOING: //游戏没开始状态时候
		var data ResponseNiuNiuState
		data.MType = common.MESSAGE_TYPE_GAME_NOT_GOING
		data.State = common.GAME_NOT_GOING
		data.BonusTimes = common.GAME_BONUS_TIMES

		SendMsgToUser(userID, data)

	case common.GAME_GOING: //游戏开始时候
		var data ResponseWithGameID
		data.MType = common.MESSAGE_TYPE_GAME_GOING
		data.BonusTimes = common.GAME_BONUS_TIMES
		data.GameID = room.GameID

		SendMsgToUser(userID, data)
	case common.GAME_CAN_RAISE: //游戏处于可押分状态下
		var data ResponseNiuNiuRaise
		data.MType = common.MESSAGE_TYPE_GAME_CAN_RAISE
		data.BonusTimes = common.GAME_BONUS_TIMES
		data.OverTime = common.GAME_NIUNIU_RAISE - int(time.Now().Unix()-room.CSStartTime)
		if data.OverTime < 0 {
			data.OverTime = 0
		}
		data.Puk.User1 = room.niuActors[0].Cards[0]
		data.Puk.User2 = room.niuActors[1].Cards[0]
		data.Puk.User3 = room.niuActors[2].Cards[0]
		data.Player0_Score = room.LScore
		data.Player1_Score = room.MScore
		data.Player2_Score = room.RScore
		data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(room.GameID, userID)

		SendMsgToUser(userID, data)
	case common.GAME_RAISE_END: //游戏押分结束时候，告知玩家扑克牌，并告知几家的押分值
		var data ResponseNiuNiuRaiseResult
		data.BonusTimes = common.GAME_BONUS_TIMES
		data.MType = common.MESSAGE_TYPE_GAME_RAISE_END
		data.OverTime = common.GAME_NIUNIU_WAIT_TIME_RESULT - int(time.Now().Unix()-room.CSStartTime)
		if data.OverTime < 0 {
			data.OverTime = 0
		}
		data.Puk.User1 = room.niuActors[0].Cards[0]
		data.Puk.User2 = room.niuActors[1].Cards[0]
		data.Puk.User3 = room.niuActors[2].Cards[0]
		data.Player0_Score = room.LScore
		data.Player1_Score = room.MScore
		data.Player2_Score = room.RScore

		SendMsgToUser(userID, data)
	case common.GAME_RESULT: //游戏结束的时候，告知几家的牌，并告知客户端几个玩家扑克对应的牛牛等级
		//largeNumber := getMaxActors(room.niuActors)
		var data ResponseNiuNiuResult
		data.MType = common.MESSAGE_TYPE_GAME_RESULT
		data.BonusTimes = common.GAME_BONUS_TIMES
		//本状态的倒计时
		data.OverTime = common.GAME_NIUNIU_LOOK_RESULT - int(time.Now().Unix()-room.CSStartTime)
		if data.OverTime < 0 {
			data.OverTime = 0
		}
		data.Puk.User1.Cards = room.niuActors[0].Cards //第一个玩家的扑克
		data.Puk.User1.Grade = room.niuActors[0].Grade //第一个玩家的扑克的牛牛等级
		data.Puk.User2.Cards = room.niuActors[1].Cards //第二个玩家的扑克
		data.Puk.User2.Grade = room.niuActors[1].Grade //第二个玩家的扑克的牛牛等级
		data.Puk.User3.Cards = room.niuActors[2].Cards //第三个玩家的扑克
		data.Puk.User3.Grade = room.niuActors[2].Grade //第三个玩家的扑克的牛牛等级
		data.Winner = room.LargeActors

		SendMsgToUser(userID, data)
	}
}

type ResponseRaiseResult struct {
	MType         int `json:"mtype"`
	Score         int `json:"Score"` //余额
	Player0_Score int `json:"p0"`
	Player1_Score int `json:"p1"`
	Player2_Score int `json:"p2"`
	MyRaise0      int `json:"m_r0"` //用户对第一个位置的押分
	MyRaise1      int `json:"m_r1"` //用户对第二个位置的押分
	MyRaise2      int `json:"m_r2"` //用户对第三个位置的押分
	ErrCode       int `json:"ErrCode"`
}

// 用户押分消息处理
// roomID房间ID
// user：User对象
func UserRaiseMsgDispose(s *melody.Session, msg []byte, roomID string, user *User) {
	//req := s.Request
	js, err := simplejson.NewJson(msg)
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	room := GetNiuNiuByRoomid(roomID)
	if room == nil {
		return
	}

	score := js.Get("score") //获取客户端传递的分数值
	score_, err := score.Int()
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	if room.GameState != common.GAME_CAN_RAISE { //如果当前游戏已经不处于押分状态下
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
			rets := strings.EqualFold(error.Error(), ERR_REDIS_STR)
			if rets {
				data.ErrCode = common.ERR_BET_OVER_MAX
			}
			data.MType = common.MESSAGE_TYPE_GAME_USER_RAISE_EOR
			data.Player0_Score = room.LScore //获取左边押分总分数
			data.Player1_Score = room.MScore //获取中间押分总分数
			data.Player2_Score = room.RScore //获取右边押分总分数
			//获取用户当前游戏局不同位置押分分值s

			data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(room.GameID, user.Uid)
			//记录下用户押分后当前的余额
			data.Score = int(user.Score)
			common.Log.Errf("orm err is 2 %s", error.Error())
			SendMsgToUser(user.Uid, data)
		} else {
			room.LScore += score_
			var data ResponseRaiseResult
			data.MType = common.MESSAGE_TYPE_GAME_STATE_USER_RAISING
			data.Player0_Score = room.LScore //获取左边押分总分数
			data.Player1_Score = room.MScore //获取中间押分总分数
			data.Player2_Score = room.RScore //获取右边押分总分数
			//获取用户当前游戏局不同位置押分分值
			data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(room.GameID, user.Uid)

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
			if error.Error() == ERR_REDIS_STR {
				data.ErrCode = common.ERR_BET_OVER_MAX
			}
			data.MType = common.MESSAGE_TYPE_GAME_USER_RAISE_EOR
			data.Player0_Score = room.LScore //获取左边押分总分数
			data.Player1_Score = room.MScore //获取中间押分总分数
			data.Player2_Score = room.RScore //获取右边押分总分数
			//获取用户当前游戏局不同位置押分分值
			data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(room.GameID, user.Uid)
			common.Log.Errf("orm err is 2 %s", error.Error())
			data.Score = int(user.Score)
			SendMsgToUser(user.Uid, data)
		} else {
			room.MScore += score_

			var data ResponseRaiseResult
			data.MType = common.MESSAGE_TYPE_GAME_STATE_USER_RAISING
			data.Player0_Score = room.LScore //获取左边押分总分数
			data.Player1_Score = room.MScore //获取中间押分总分数
			data.Player2_Score = room.RScore //获取右边押分总分数
			//获取用户当前游戏局不同位置押分分值
			data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(room.GameID, user.Uid)

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
			if error.Error() == ERR_REDIS_STR {
				data.ErrCode = common.ERR_BET_OVER_MAX
			}
			data.MType = common.MESSAGE_TYPE_GAME_USER_RAISE_EOR
			data.Player0_Score = room.LScore //获取左边押分总分数
			data.Player1_Score = room.MScore //获取中间押分总分数
			data.Player2_Score = room.RScore //获取右边押分总分数
			//获取用户当前游戏局不同位置押分分值
			data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(room.GameID, user.Uid)
			common.Log.Errf("orm err is 2 %s", error.Error())
			data.Score = int(user.Score)
			SendMsgToUser(user.Uid, data)
		} else {
			room.RScore += score_

			var data ResponseRaiseResult
			data.MType = common.MESSAGE_TYPE_GAME_STATE_USER_RAISING
			data.Player0_Score = room.LScore //获取左边押分总分数
			data.Player1_Score = room.MScore //获取中间押分总分数
			data.Player2_Score = room.RScore //获取右边押分总分数

			//获取用户当前游戏局不同位置押分分值
			data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(room.GameID, user.Uid)

			//user.subScore(int64(score_))
			user.DelMoney(nil, common.MONEY_TYPE_SCORE, int64(score_), false)
			//记录下用户押分后当前的余额
			data.Score = int(user.Score)
			SendMsgToUser(user.Uid, data)
		}
	}
}

type ResponseUserWinnerScore struct {
	MType    int `json:"mtype"`
	Score    int `json:"Score"`     //余额
	WinCoins int `json:"win_coins"` //赢取余额
}

// 每局结束后发送用户押分输赢结果【还剩多少金币，赢了多少金币】
// uID:用户ID
// score：用户当前游戏币
// winScores:当局赢下的游戏币
func UserWinnerScore(uID, score, winScores int) {
	var data ResponseUserWinnerScore
	data.MType = common.MESSAGE_TYPE_GAME_END_WIN_SCORE
	data.Score = score
	data.WinCoins = winScores

	SendMsgToUser(uID, data)
}

type GameRankingRecord struct {
	Uid      int    `json:"uid"`
	WinCoins int    `json:"win_coins"`
	NickName string `json:"nick_name"`
	Image    string `json:"image"`
}

type ResponseWinnerSorter struct {
	MType int                  `json:"mtype"`
	List  []*GameRankingRecord `json:"list"` //赢家列表
}

// 通知主播赢家排行榜
// anchorID：主播ID
// msgID：消息ID
// ms：排序后的Map对象
func NoticeRaiselistMsgToAnchor(anchorID, msgID int, ms common.MapSorter) {
	var data ResponseWinnerSorter
	data.MType = msgID

	gameRanking := make([]*GameRankingRecord, 0)

	for _, item := range ms {
		user, _ := GetUserByUid(item.Uid)
		if user != nil && item.WinCoins > 0 {
			m := &GameRankingRecord{}
			m.Uid = item.Uid
			m.WinCoins = item.WinCoins
			m.Image = user.Image
			m.NickName = user.NickName
			gameRanking = append(gameRanking, m)
		}
	}
	data.List = gameRanking

	SendMsgToUser(anchorID, data)

}

// 根据传入的参数值构造扑克玩家对象[NiuActorInfo对象]
// numberOfPlayer 有几个玩家
// NiuActorInfo 扑克玩家对象
func CreateNiuActor(numberOfPlayer int) []*NiuActorInfo {
	var i int //Cardall计算5张牌总值，cow计算牛几。

	Random := common.RandomRangeArr(0, 51, 5*numberOfPlayer)

	niuActors := make([]*NiuActorInfo, 0)

	for k := 1; k <= numberOfPlayer; k++ {
		temp := &NiuActorInfo{}
		//默认第一张牌为最大的一张牌的点数
		temp.MIndex = (Random[(k-1)*5]) / 13
		temp.Cardall = 0

		//根据给出的随机数[0-51]给出牌的点数和花色
		for i = 0; i < 5; i++ {
			temp.Cards[i] = Pork[Random[(k-1)*5+i]] //5张牌

			//存储当前用户最大的一张牌的下标
			if temp.Cards[i].Number > temp.Cards[temp.MIndex].Number {
				temp.MIndex = i
			} else if temp.Cards[i].Number == temp.Cards[temp.MIndex].Number &&
				temp.Cards[i].Color > temp.Cards[temp.MIndex].Color {
				temp.MIndex = i
			}

			temp.Cardall += Pork[Random[(k-1)*5+i]].Number //所有数值求和，方便后续计算
		}

		hasBomb, grade1 := HasBomb(temp.Cards) //是否有炸弹
		if hasBomb {
			temp.Grade = grade1 //有炸弹，记录下炸弹的值
		} else {
			hasJinOrYinNiu, grade2 := HasJinOrYinNiu(temp.Cards) //是否有金银牛
			if hasJinOrYinNiu {
				temp.Grade = grade2 //有金银牛记录下金银牛的值
			} else {
				temp.Grade, _ = GetNumberOfNiu(temp.Cards) //几牛啊
			}
		}

		niuActors = append(niuActors, temp) //追加到niuActors
	}
	return niuActors //返回niuActors对象
}

func PrintfNiuActor(actor *NiuActorInfo) {
	fmt.Printf("(")
	for _, ac := range actor.Cards {
		if ac.Number >= 1 && ac.Number <= 10 {
			fmt.Printf("%d,", ac.Number)
		} else if ac.Number == 11 {
			fmt.Printf("%s,", "J")
		} else if ac.Number == 12 {
			fmt.Printf("%s,", "Q")
		} else if ac.Number == 13 {
			fmt.Printf("%s,", "K")
		}
	}
	fmt.Printf(") ")
}

// 根据传入的NiuActorInfo对象重新构造每个对象的后面四张牌
func ReCreateNiuActor(actor *NiuActorInfo) { //改变actor后面四张牌
	Random := common.RandomRangeArr(0, 51, 4)

	actor.MIndex = 0
	actor.Cardall = actor.Cards[0].Number
	for i := 1; i < 5; i++ { //根据给出的随机数[0-51]给出牌的点数和花色
		actor.Cards[i] = Pork[Random[i-1]] //后面四张牌

		//存储当前用户最大的一张牌的下标
		if actor.Cards[i].Number > actor.Cards[actor.MIndex].Number {
			actor.MIndex = i
		} else if actor.Cards[i].Number == actor.Cards[actor.MIndex].Number &&
			actor.Cards[i].Color > actor.Cards[actor.MIndex].Color {
			actor.MIndex = i
		}

		actor.Cardall += actor.Cards[i].Number
	}

	hasBomb, grade1 := HasBomb(actor.Cards) //是否有炸弹
	if hasBomb {
		actor.Grade = grade1 //记录下炸弹等级
	} else {
		hasJinOrYinNiu, grade2 := HasJinOrYinNiu(actor.Cards) // 是否有金银牛
		if hasJinOrYinNiu {
			actor.Grade = grade2
		} else {
			actor.Grade, _ = GetNumberOfNiu(actor.Cards)
		}
	}
}

// 给一组扑克获取有几个牛（[1-13]间的5个随机数），返回几牛和最大的牌
func GetNumberOfNiu(c [5]Card) (int, int) {
	var i, j, cow, numberTotal int //Cardall计算5张牌总值，cow计算牛几。
	var n int = 0                  //存储10、J、Q、K张数。
	var c_Temp Card

	numberTotal = 0
	for i = 0; i < 5; i++ { //统计本组数字里大于等于10的个数。
		if c[i].Number >= 10 {
			n++
		}
		numberTotal += c[i].Number
	}

	//对5张牌从大到小排序。
	for i = 0; i < 4; i++ {
		for j = i + 1; j < 5; j++ {
			if c[i].Number < c[j].Number {
				c_Temp = c[i]
				c[i] = c[j]
				c[j] = c_Temp
			}
		}
	}

	switch n {
	case 0: //5张牌中没有一张10、J、Q、K。
		for i = 0; i < 4; i++ {
			for j = 1; j < 5; j++ {
				if (numberTotal-c[i].Number-c[j].Number)%10 == 0 && i != j {
					cow = (c[i].Number + c[j].Number) % 10
					return GetNiuNumberByN(cow), c[0].Number
				}
			}
		}

		return GetNiuNumberByN(-1), c[0].Number
	case 1: //5张牌中有一张10、J、Q、K。
		//剩下四张牌里2张加起来为10的倍数
		numberTotal = numberTotal - c[0].Number
		for k := 1; k < 4; k++ {
			for i = k + 1; i < 5; i++ {
				if (numberTotal-c[k].Number-c[i].Number)%10 == 0 {
					cow = numberTotal % 10
					return GetNiuNumberByN(cow), c[0].Number
				}
			}
		}

		//剩下四张牌里3张加起来为10的倍数
		if (c[1].Number+c[2].Number+c[3].Number)%10 == 0 {
			return GetNiuNumberByN(c[4].Number % 10), c[0].Number
		} else if (c[1].Number+c[2].Number+c[4].Number)%10 == 0 {
			return GetNiuNumberByN(c[3].Number % 10), c[0].Number
		} else if (c[1].Number+c[3].Number+c[4].Number)%10 == 0 {
			return GetNiuNumberByN(c[2].Number % 10), c[0].Number
		} else if (c[2].Number+c[3].Number+c[4].Number)%10 == 0 {
			return GetNiuNumberByN(c[1].Number % 10), c[0].Number
		}

		return GetNiuNumberByN(-1), c[0].Number
	case 2: //5张牌中有两张10、J、Q、K。
		if (c[2].Number+c[3].Number+c[4].Number)%10 == 0 {
			return GetNiuNumberByN(0), c[0].Number
		} else if c[2].Number+c[3].Number == 10 {
			return GetNiuNumberByN(c[4].Number % 10), c[0].Number
		} else if c[2].Number+c[4].Number == 10 {
			return GetNiuNumberByN(c[3].Number % 10), c[0].Number
		} else if c[3].Number+c[4].Number == 10 {
			return GetNiuNumberByN(c[2].Number % 10), c[0].Number
		}

		return GetNiuNumberByN(-1), c[0].Number
	case 3: //5张牌中有三张10、J、Q、K。
		for i = 0; i < n; i++ { //总值减去10、J、Q、K的牌。
			numberTotal -= c[i].Number
		}
		cow = numberTotal % 10
		return GetNiuNumberByN(cow), c[0].Number
	case 4: //5张牌中有四张10、J、Q、K。
		for i = 0; i < n; i++ { //总值减去10、J、Q、K的牌。
			numberTotal -= c[i].Number
		}
		cow = numberTotal % 10
		return GetNiuNumberByN(cow), c[0].Number
	case 5: //5张牌中五张都是10、J、Q、K。
		for i = 0; i < n; i++ { //总值减去10、J、Q、K的牌。
			numberTotal -= c[i].Number
		}
		cow = numberTotal % 10
		return GetNiuNumberByN(cow), c[0].Number
	}
	return NIU_ZERO, c[0].Number
}

// 根据传入的值（n）得到牛牛等级
func GetNiuNumberByN(n int) int {
	switch n {
	case 1:
		return NIU_ONE
	case 2:
		return NIU_TWO
	case 3:
		return NIU_THREE
	case 4:
		return NIU_FOUR
	case 5:
		return NIU_FIVE
	case 6:
		return NIU_SIX
	case 7:
		return NIU_SEVEN
	case 8:
		return NIU_EIGHT
	case 9:
		return NIU_NINE
	case 0:
		return NIU_NIU
	case -1:
		return NIU_ZERO
	}
	return NIU_ZERO
}

// 根据5张扑克判断是否有炸弹
func HasBomb(c [5]Card) (bool, int) {
	n1 := c[0].Number
	n2 := c[1].Number
	number1 := 1
	number2 := 1

	//统计本组数字里大于等于10的个数。
	for i := 2; i < 5; i++ {
		if n1 == n2 {
			if c[i].Number != n2 {
				n2 = c[i].Number
			}
			number1 += 1
			number2 = 1
		} else if n1 != n2 {
			if c[i].Number != n1 && c[i].Number != n2 {
				return false, NIU_ZERO
			} else if c[i].Number == n1 {
				number1 += 1
			} else {
				number2 += 1
			}
		}
	}

	//返回是否有炸弹，并且返回炸弹的值
	if number1 > number2 && number1 >= 4 {
		return true, n1 + (BOMB_A - YIN_NIU)
	} else if number2 > number1 && number2 >= 4 {
		return true, n2 + (BOMB_A - YIN_NIU)
	} else {
		return false, NIU_ZERO
	}
}

// 根据传入的5张牌判断是否有金牛或者银牛，返回的第一个参数为true时表示有金银牛，
// 第二个参数区别金银牛,如果没有金银牛第二个参数返回0
func HasJinOrYinNiu(c [5]Card) (bool, int) {
	var a, i, j int
	var n int = 0 //存储10、J、Q、K张数。

	//统计本组数字里大于等于10的个数。
	for i = 0; i < 5; i++ {
		if c[i].Number >= 10 {
			n++
		}
	}

	//没有金银牛，返回
	if n < 5 {
		return false, 0
	}

	//对5张牌从大到小排序。
	for i = 0; i < 4; i++ {
		for j = i + 1; j < 5; j++ {
			if c[i].Number < c[j].Number {
				a = c[i].Number
				c[i].Number = c[j].Number
				c[j].Number = a
			}
		}
	}

	//第5张牌也大于10为金牛，否者为银牛
	if c[4].Number > 10 {
		return true, JIN_NIU
	} else {
		return true, YIN_NIU
	}
}

// 给一组NiuActorInfo数组获取最大的下标
// actors为要检索的NiuActorInfo对象数组
func getMaxActors(actors []*NiuActorInfo) int {
	maxIndex := 0 //最大者下标
	i := 0
	temp := &NiuActorInfo{} //临时存储NiuActorInfo对象

	//遍历actors寻找最大者
	for _, m := range actors {
		if i == 0 { //先取第一个值赋值给temp
			temp = m
		} else if m.Grade > temp.Grade { //先比较等级
			temp = m
			maxIndex = i
		} else if m.Grade == temp.Grade &&
			m.Cards[m.MIndex].Number > temp.Cards[temp.MIndex].Number { //等级相同点数大些的时候
			temp = m
			maxIndex = i
		} else if m.Grade == temp.Grade &&
			m.Cards[m.MIndex].Number == temp.Cards[temp.MIndex].Number &&
			m.Cards[m.MIndex].Color > temp.Cards[temp.MIndex].Color { // 等级和点数相同再比较花色
			temp = m
			maxIndex = i
		}

		i++
	}

	return maxIndex
}

// 通知房间所有用户游戏关闭
func NoticeAllUserInRoomGameClose(s *melody.Session) {
	var data ResponseHints
	data.MType = common.MESSAGE_TYPE_GAME_CLOSE

	chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
		return q.Request.URL.Path == s.Request.URL.Path
	})
}
