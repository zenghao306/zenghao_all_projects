//砸金花相关的逻辑处理
package model

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/yshd_game/common"
	"github.com/yshd_game/melody"
	"github.com/yshd_game/timer"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	GF_LEAFLETS       = iota //0.单张
	GF_SUB            = 1    //1.对子
	GF_STRAIGHT       = 2    //2.顺子
	GF_SAME_FLOWER    = 3    //3.同花
	GF_STRAIGHT_FLUSH = 4    //4.同花顺
	GF_LEOPARD        = 5    //5.豹子
)

type GFActorInfo struct {
	Cards  [3]Card //3张牌
	MIndex int     //最大的一张牌的下标
	Grade  int     //3张牌对应的金花等级[0-5]参见金花级别常量
}

var gf_room_manager map[string]*GFGameInfo
var mutex_GoldenFlower sync.RWMutex

//本局游戏信息
type GFGameInfo struct {
	Rid           string //房间ID
	Uid           int    //用户ID
	IsStarted     bool   //是否已开始
	mutex_gf_info sync.RWMutex
	GameState     int            //游戏状态【进行到哪一步了】
	CSStartTime   int64          //当前状态下开始时间戳
	GameID        string         //当前局的游戏ID
	GFActors      []*GFActorInfo //当前局下的牌
	LScore        int            //最左边押分数据
	MScore        int            //中间押分数据
	RScore        int            //右边押分数据
	LargeActors   int            //哪家最大
}

//新创建一牛牛游戏
func NewRoomInfoGoldenFlower() *GFGameInfo {
	m := &GFGameInfo{}
	return m
}

//初始化牛牛Map
func InitGoldenFlower() map[string]*GFGameInfo {
	gf_room_manager = make(map[string]*GFGameInfo)
	return gf_room_manager
}

// 根据房间号（roomid）获取扎金花游戏对象（GFGameInfo）
// roomid 直播的房间ID
func GetGoldenFlowerByRoomid(roomid string) *GFGameInfo {
	v, ok := gf_room_manager[roomid]
	if !ok {
		return nil
	}
	return v
}

//添加牛牛游戏对象到房间Map中
// r牛牛游戏对象
func AddRoomInfoGoldenFlower(r *GFGameInfo) bool {
	_, ok := gf_room_manager[r.Rid]
	if ok {
		return false
	}

	mutex_GoldenFlower.Lock()
	gf_room_manager[r.Rid] = r
	mutex_GoldenFlower.Unlock()

	return true
}

// 根据房间ID从牛牛Map里删掉牛牛对象
// rid房间ID
func DelRoomInfoGoldenFlower(rid string) {

	mutex_GoldenFlower.Lock()
	delete(gf_room_manager, rid)
	mutex_GoldenFlower.Unlock()
}

func CreatGoldenFlower(numberOfPlayer int) []*GFActorInfo {
	Random := common.RandomRangeArr(0, 51, 3*numberOfPlayer)

	actors := make([]*GFActorInfo, 0)

	for k := 1; k <= numberOfPlayer; k++ {
		temp := &GFActorInfo{}

		//temp.MIndex =
		//构造每个角色的三张牌
		temp.Cards[0] = TexasPork[Random[(k-1)*3]]
		temp.Cards[1] = TexasPork[Random[(k-1)*3+1]]
		temp.Cards[2] = TexasPork[Random[(k-1)*3+2]]
		temp.Grade, temp.MIndex = GetGradeOfJH(temp.Cards)
		actors = append(actors, temp) //追加到actors
	}

	return actors
}

// 根据传入的三张牌获取对应金花等级（牌型），返回对应金花等级以及最大牌的下标
func GetGradeOfJH(c [3]Card) (int, int) {
	var i, j, mIndex int

	//mIndex主要用来记录三张里最大牌的下标
	mIndex = 0
	if c[1].Number > c[0].Number || (c[1].Number == c[0].Number && c[1].Color > c[0].Color) {
		mIndex = 1
	}
	if c[2].Number > c[mIndex].Number || (c[2].Number == c[mIndex].Number && c[2].Color > c[mIndex].Color) {
		mIndex = 2
	}

	// 对玩家b的5张牌从大到小进行排序。
	for i = 0; i < 2; i++ {
		for j = i + 1; j < 3; j++ {
			if c[i].Number < c[j].Number ||
				((c[i].Number == c[j].Number) && c[i].Color < c[j].Color) {
				temp := c[i]
				c[i] = c[j]
				c[j] = temp
			}
		}
	}

	if c[0].Number == c[1].Number || c[0].Number == c[2].Number || c[1].Number == c[2].Number { //有相等的情况
		if c[0].Number == c[1].Number && c[0].Number == c[2].Number {
			return GF_LEOPARD, mIndex //豹子
		} else {
			return GF_SUB, mIndex //对子
		}
	}

	if c[0].Number == c[1].Number+1 && c[1].Number == c[2].Number+1 { // 顺子
		if c[0].Color == c[1].Color && c[0].Color == c[2].Color { //同花顺
			return GF_STRAIGHT_FLUSH, mIndex
		} else {
			return GF_STRAIGHT, mIndex
		}
	}

	if c[0].Color == c[1].Color && c[0].Color == c[2].Color { //同花
		return GF_SAME_FLOWER, mIndex
	}

	return GF_LEAFLETS, mIndex
}

func getGFMaxActors(actors []*GFActorInfo) int {
	maxIndex := 0 //最大者下标
	i := 0
	temp := &GFActorInfo{} //临时存储GFActorInfo对象

	//遍历actors寻找最大者
	for _, m := range actors {
		if i == 0 { //先取第一个值赋值给temp
			temp = m
		} else if m.Grade > temp.Grade { //先比较等级
			temp = m
			maxIndex = i
		} else if m.Grade == temp.Grade {
			flag := CompareTwoGoldenFlowerActor(*m, *temp)
			if flag {
				temp = m
				maxIndex = i
			}
		}

		i++
	}

	return maxIndex
}

//两家金花扑克牌一样时候比较大小
func CompareTwoGoldenFlowerActor(a, b GFActorInfo) bool {
	var i, j, k int

	if a.Grade > b.Grade {
		return true
	} else if a.Grade < b.Grade {
		return false
	}

	/*** 以下部分代码是两家扑克等级相同情况下对两家进行比较大小***/
	// 对玩家a的3张牌从大到小排序。
	for i = 0; i < 2; i++ {
		for j = i + 1; j < 3; j++ {
			if a.Cards[i].Number < a.Cards[j].Number ||
				((a.Cards[i].Number == a.Cards[j].Number) && a.Cards[i].Color < a.Cards[j].Color) {
				temp := a.Cards[i]
				a.Cards[i] = a.Cards[j]
				a.Cards[j] = temp
			}
		}
	}

	// 对玩家b的3张牌从大到小进行排序。
	for i = 0; i < 2; i++ {
		for j = i + 1; j < 3; j++ {
			if b.Cards[i].Number < b.Cards[j].Number ||
				((b.Cards[i].Number == b.Cards[j].Number) && b.Cards[i].Color < b.Cards[j].Color) {
				temp := b.Cards[i]
				b.Cards[i] = b.Cards[j]
				b.Cards[j] = temp
			}
		}
	}

	index := 0 //记录有多少个唯一的数字挪动到最后面去了
	for i = 0; i < 3; i++ {
		flag := false
		for j = 0; j < 3; j++ {
			if a.Cards[index].Number == a.Cards[j].Number && index != j {
				flag = true
			}
		}
		if flag {
			index++
		} else {
			temp := a.Cards[index]
			for k = index + 1; k < 3; k++ {
				a.Cards[k-1] = a.Cards[k]
			}
			a.Cards[k-1] = temp
		}
	}

	index = 0 //记录有多少个唯一的数字挪动到最后面去了
	for i = 0; i < 3; i++ {
		flag := false
		for j = 0; j < 3; j++ {
			if b.Cards[index].Number == b.Cards[j].Number && index != j {
				flag = true
			}
		}
		if flag {
			index++
		} else {
			temp := b.Cards[index]
			for k = index + 1; k < 3; k++ {
				b.Cards[k-1] = b.Cards[k]
			}
			b.Cards[k-1] = temp
		}
	}
	//fmt.Printf("\n 挪动后a: (%d,%d,%d) ",a.Cards[0].Number,a.Cards[1].Number,a.Cards[2].Number)
	//fmt.Printf("\n 挪动后b: (%d,%d,%d) ",b.Cards[0].Number,b.Cards[1].Number,b.Cards[2].Number)

	if a.Grade == GF_LEAFLETS || a.Grade == GF_SUB {
		// 单牌情况下先从大到小依次比较大小
		for m := 0; m < 3; m++ {
			aTemp := a.Cards[m]
			bTemp := b.Cards[m]
			if aTemp.Number > bTemp.Number {
				return true
			} else if aTemp.Number < bTemp.Number {
				return false
			}
		}

		//依次比较大小还比较不出来那就比较第一张【排序后最大的那张】花色大小
		if a.Cards[0].Color > b.Cards[0].Color {
			return true
		} else {
			return false
		}
	}

	// 依次比较出大小
	for m := 0; m < 3; m++ {
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

/////////////////////////////////////////////////////////
// 本函数主要用于砸金花游戏的创建【入口函数】
// s：长链接的Session对象，在游戏过程中用来发送消息
// anchorID：主播ID
// roomID：房间ID
func NewGameGoldenFlower(s *melody.Session, anchorID int, roomID string) {
	//左中右三家分值的临时存储变量
	var lTemp, mTemp, rTemp int
	//砸金花游戏对象，用来保存每局金花游戏相关的数据
	var roomGF *GFGameInfo
	//fmt.Printf("\nNewGameGoldenFlower@@@")

	defer common.PrintPanicStack()

	//根据房间ID获取牛牛游戏对象
	roomGF = GetGoldenFlowerByRoomid(roomID)

	// 房间对象
	room, has := GetRoomById(roomID)

	//当房间不为空且标识值为TRUE时候
	for roomGF != nil && roomGF.IsStarted && has && room.Statue == common.ROOM_ONLIVE {
		//fmt.Printf("\n roomGF != nil=%t,roomGF.IsStarted=%t,has=%t,room.Statue=%d", roomGF != nil, roomGF.IsStarted, has, room.Statue)

		if !common.GameRunningSwitch { //游戏在后台关闭时候通知所有用户
			NoticeAllUserInRoomGameClose(s)
			break
		}

		//左边押分值设为0
		roomGF.LScore = 0
		lTemp = 0
		//中间押分值设为0
		roomGF.MScore = 0
		mTemp = 0
		//右边押分值设为0
		roomGF.RScore = 0
		rTemp = 0

		gameStatTime := time.Now().Unix() //记录下游戏开始的时间戳

		//构造出三个玩家对象
		actors := CreatGoldenFlower(3)
		//计算出最大的玩家[0-2]
		roomGF.LargeActors = getGFMaxActors(actors)

		for _, m := range actors {
			for i := 0; i < 3; i++ {
				if m.Cards[i].Number == 14 {
					m.Cards[i].Number = 1
				}
			}
		}

		//全局标识记录下
		roomGF.GFActors = actors

		//fmt.Printf("\n")
		//for _, m := range roomGF.GFActors {
		//	for i := 0; i < 3; i++ {
		//		fmt.Printf(" (%d,%d),", m.Cards[i].Number, m.Cards[i].Color)
		//	}
		//	fmt.Printf("\n")
		//}

		//fmt.Printf("\n 235 235 235")
		//构造游戏ID
		gameID := fmt.Sprintf("%d_%d", anchorID, time.Now().Unix())
		roomGF.GameID = gameID //记录下游戏ID

		//通知房间里所有用户游戏开始
		NoticeMsgToRoomWithGameID(s, common.MESSAGE_TYPE_GAME_GOING, gameID)
		//存储游戏状态
		roomGF.GameState = common.GAME_GOING
		//记录下本状态下开始的时间戳，以免让中途有用户进入时候好计算倒计时
		roomGF.CSStartTime = time.Now().Unix()

		//fmt.Printf("\n 247")

		f := timer.NewDispatcher(1)
		f.AfterFunc(4*time.Second, func() { //4秒后通知所有玩家

			//fmt.Printf("\n f f f f 248")

			//标识出可以押分了
			roomGF.GameState = common.GAME_CAN_RAISE
			//记录下押分状态的开始时间戳
			roomGF.CSStartTime = time.Now().Unix()
			//通知所有用户可以押注了
			NoticeBeginGFRaise(s)
		})
		(<-f.ChanTimer).Cb()

		flag := time.Now().Unix() + common.GAME_GF_RAISE //记录下押注状态结束的时间戳，方便后续计算
		for time.Now().Unix() < flag {
			f2 := timer.NewDispatcher(1)
			f2.AfterFunc(1*time.Second, func() {
				//fmt.Printf(" f2 ")
				//当有用户押分时候通知所有用户左中右三家押注情况的消息
				if lTemp != roomGF.LScore || mTemp != roomGF.MScore || rTemp != roomGF.RScore {
					//fmt.Printf("\n f2 f2 f2 f2 265")
					NoticeLMRScoreMsgToRoom(s, common.MESSAGE_TYPE_GAME_RAISE_SCORE,
						roomGF.LScore, roomGF.MScore, roomGF.RScore)
				}

				//记录啊下左中右三家当前押分数据，下次再跟当前记录的押分数据比较判断是否有变化。
				lTemp = roomGF.LScore
				mTemp = roomGF.MScore
				rTemp = roomGF.RScore
			})
			(<-f2.ChanTimer).Cb()
		}

		//30秒后通知所有用户押分结束
		f3 := timer.NewDispatcher(1)
		f3.AfterFunc(1*time.Second, func() {
			//fmt.Printf("\n f3 f3 f3 f3 281")
			roomGF.GameState = common.GAME_RAISE_END //保存游戏状态
			roomGF.CSStartTime = time.Now().Unix()   //记录下本状态下开始的时间戳，以免让中途有用户进入时候好计算倒计时
			NoticeMsgToRoomWithTime(s, common.MESSAGE_TYPE_GAME_RAISE_END, common.GAME_GF_WAIT_TIME_RESULT)
		})
		(<-f3.ChanTimer).Cb()

		//押分结束后通知所有用户游戏结束并告知本局结果
		f4 := timer.NewDispatcher(1)
		f4.AfterFunc(common.GAME_GF_WAIT_TIME_RESULT*time.Second, func() { //10秒后执行
			//fmt.Printf("\n f f f f 291")
			roomGF.GameState = common.GAME_RESULT  //保存游戏状态
			roomGF.CSStartTime = time.Now().Unix() //记录下本状态下开始的时间戳，以免让中途有用户进入时候好计算倒计时

			NoticeGoldenFlowerResult(s, actors, roomGF.LargeActors, common.GAME_GF_LOOK_RESULT) //通知所有用户游戏结果

			//游戏结束时候后台根据用户押分情况进行结算并告知对应的用户
			m := ResultBetReward(gameID, roomGF.LargeActors, common.GAME_BONUS_TIMES, roomID, gameStatTime, anchorID)
			ms := common.NewMapSorter(m)

			sort.Sort(ms)
			NoticeRaiselistMsgToAnchor(anchorID, common.MESSAGE_TYPE_GAME_WINNER_SORTER, ms)

		})
		(<-f4.ChanTimer).Cb()

		f5 := timer.NewDispatcher(1)
		f5.AfterFunc(common.GAME_GF_LOOK_RESULT*time.Second, func() { //所有用户查看结果
			//	fmt.Printf("\n f5 ")
		})
		(<-f5.ChanTimer).Cb()

		roomGF.GameState = 0 //房间里的游戏状态值设为默认，准备下一局

		room, has = GetRoomById(roomID)

		//fmt.Printf("\n roomGF != nil=%t,roomGF.IsStarted=%t,has=%t,room.Statue=%d", roomGF != nil, roomGF.IsStarted, has, room.Statue)
	}

}

type ResponseGFRaise struct { //牛牛押注
	MType         int `json:"mtype"`       //消息ID
	BonusTimes    int `json:"bonus_times"` //奖励倍数
	OverTime      int `json:"over_time"`   //倒计时
	Player0_Score int `json:"p0"`          //最左边所有用户的押分
	Player1_Score int `json:"p1"`          //第二个位置所有用户押分
	Player2_Score int `json:"p2"`          //第三个位置所有用户押分
	MyRaise0      int `json:"m_r0"`        //用户对第一个位置的押分
	MyRaise1      int `json:"m_r1"`        //用户对第二个位置的押分
	MyRaise2      int `json:"m_r2"`        //用户对第三个位置的押分
}

//通知客户端牛牛游戏可开始押注
func NoticeBeginGFRaise(s *melody.Session) {
	var data ResponseGFRaise
	data.MType = common.MESSAGE_TYPE_GAME_CAN_RAISE //消息ID
	data.BonusTimes = common.GAME_BONUS_TIMES       //游戏奖励倍数
	data.OverTime = common.GAME_GF_RAISE            //倒计时

	chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
		return q.Request.URL.Path == s.Request.URL.Path
	})
}

type GFCardInfo struct {
	Cards [3]Card `json:"cards"` //3张牌
	Grade int     `json:"grade"` //级别常量
}

type PukGF struct {
	User1 GFCardInfo `json:"user1"` //用户1的牌
	User2 GFCardInfo `json:"user2"` //用户2的牌
	User3 GFCardInfo `json:"user3"` //用户3的牌
}

type ResponseGFRaiseResult struct { //牛牛押注
	MType         int   `json:"mtype"`       //消息ID
	BonusTimes    int   `json:"bonus_times"` //奖励倍数
	OverTime      int   `json:"over_time"`   //倒计时
	Puk           PukGF `json:"puk"`         //扑克信息
	Player0_Score int   `json:"p0"`          //最左边所有用户的押分
	Player1_Score int   `json:"p1"`          //第二个位置所有用户押分
	Player2_Score int   `json:"p2"`          //第三个位置所有用户押分
	Winner        int   `json:"winner"`
}

//根据传入的GFActorInfo对象和本阶段剩余时间[overTime]通知客户端牛牛游戏结果
func NoticeGoldenFlowerResult(s *melody.Session, actors []*GFActorInfo, largeNumber, overTime int) {
	//largeNumber := getGFMaxActors(actors)
	var data ResponseGFRaiseResult
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

type ResponseGFState struct {
	MType      int `json:"mtype"`
	BonusTimes int `json:"bonus_times"`
	State      int `json:"state"`
}

// 用户在游戏中途进入房间时候根据游戏状态值发送对应消息给用户
// userID被通知者ID
// roomID所在直播房间ID
func NoticeGoldenFlowerStateToUser(s *melody.Session, userID int, roomID string) {
	//获取砸金花游戏对象
	room := GetGoldenFlowerByRoomid(roomID)
	if room == nil {
		return
	}

	switch room.GameState { //判断房间对象的状态
	case common.GAME_NOT_GOING: //游戏没开始状态时候
		var data ResponseGFState
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
		var data ResponseGFRaise
		data.MType = common.MESSAGE_TYPE_GAME_CAN_RAISE
		data.BonusTimes = common.GAME_BONUS_TIMES
		data.OverTime = common.GAME_GF_RAISE - int(time.Now().Unix()-room.CSStartTime)
		if data.OverTime < 0 {
			data.OverTime = 0
		}
		data.Player0_Score = room.LScore
		data.Player1_Score = room.MScore
		data.Player2_Score = room.RScore
		data.MyRaise0, data.MyRaise1, data.MyRaise2 = GetUserRaiseInfo(room.GameID, userID)

		SendMsgToUser(userID, data)
	case common.GAME_RAISE_END: //游戏押分结束时候，告知玩家扑克牌，并告知几家的押分值
		var data ResponseGFRaiseResult
		data.BonusTimes = common.GAME_BONUS_TIMES
		data.MType = common.MESSAGE_TYPE_GAME_RAISE_END
		data.OverTime = common.GAME_GF_WAIT_TIME_RESULT - int(time.Now().Unix()-room.CSStartTime)
		if data.OverTime < 0 {
			data.OverTime = 0
		}

		data.Player0_Score = room.LScore
		data.Player1_Score = room.MScore
		data.Player2_Score = room.RScore

		SendMsgToUser(userID, data)
	case common.GAME_RESULT: //游戏结束的时候，告知几家的牌，并告知客户端几个玩家扑克对应的牛牛等级
		var data ResponseGFRaiseResult
		data.MType = common.MESSAGE_TYPE_GAME_RESULT
		data.BonusTimes = common.GAME_BONUS_TIMES
		//本状态的倒计时
		data.OverTime = common.GAME_GF_LOOK_RESULT - int(time.Now().Unix()-room.CSStartTime)
		if data.OverTime < 0 {
			data.OverTime = 0
		}
		data.Puk.User1.Cards = room.GFActors[0].Cards
		data.Puk.User1.Grade = room.GFActors[0].Grade
		data.Puk.User2.Cards = room.GFActors[1].Cards
		data.Puk.User2.Grade = room.GFActors[1].Grade
		data.Puk.User3.Cards = room.GFActors[2].Cards
		data.Puk.User3.Grade = room.GFActors[2].Grade
		data.Player0_Score = room.LScore
		data.Player1_Score = room.MScore
		data.Player2_Score = room.RScore
		data.Winner = room.LargeActors

		SendMsgToUser(userID, data)
	}
}

type ResponseMyGFRaiseResult struct {
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
func GoldenFlowerUserRaiseMsgDispose(s *melody.Session, msg []byte, roomID string, user *User) {
	//req := s.Request
	js, err := simplejson.NewJson(msg)
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	room := GetGoldenFlowerByRoomid(roomID)
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
		var data ResponseMyGFRaiseResult
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
		var data ResponseMyGFRaiseResult
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
			var data ResponseMyGFRaiseResult
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
			var data ResponseMyGFRaiseResult
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
			var data ResponseMyGFRaiseResult
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

			var data ResponseMyGFRaiseResult
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
			var data ResponseMyGFRaiseResult
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

			var data ResponseMyGFRaiseResult
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
