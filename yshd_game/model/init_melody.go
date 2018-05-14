package model

import (
	//"github.com/olahol/melody"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/go-redis/redis"
	"github.com/yshd_game/common"
	"github.com/yshd_game/melody"
	"github.com/yshd_game/sensitive"
	"regexp"
	"strconv"
	"time"
)

type Response struct {
	MType       int    `json:"mtype"`
	Uid         int    `json:"uid"`
	Name        string `json:"name"`
	Msg         string `json:"msg"`
	Face        string `json:"face"`
	IsFocus     bool   `json:"is_focus"`
	Sex         int    `json:"sex"`
	Location    string `json:"location"`
	UserLevel   int    `json:"user_level"`
	AnchorLevel int    `json:"anchor_level"`
	IsSuperUser bool   `json:"is_super_user"`
	Guard       int    `json:"guard"`
}

type ResponseStar struct {
	MType int `json:"mtype"`
}

type ResponseErr struct {
	MType int `json:"mtype"`
	Err   int `json:"err"`
}

type ResponseGift struct {
	MType       int    `json:"mtype"`
	SendId      int    `json:"send_id"`
	SendName    string `json:"send_name"`
	SendImage   string `json:"send_Image"`
	RevId       int    `json:"rev_id"`
	RevName     string `json:"rev_name"`
	GiftId      int    `json:"gift_id"`
	GiftNum     int    `json:"gift_num"`
	SendLevel   int    `json:"send_level"`
	GiftDynamic int    `json:"gift_dynamic"`
	GiftName    string `json:"gift_name"`
}

type ResponseClose struct {
	MType int `json:"mtype"`
}

type ResponseSys struct {
	MType  int    `json:"mtype"`
	Notice string `json:"notice"`
}

type ResponseHints struct {
	MType int `json:"mtype"`
}

type ResponseInviteResult struct {
	UID      int    `json:"uid"`
	NickName string `json:"nick_name"`
	MType    int    `json:"mtype"`
	IsAgree  int    `json:"is_agree"`
}

type ResponseInvite struct {
	MType int    `json:"mtype"`
	User  int    `json:"user"`
	Other int    `json:"other"`
	Rid   string `json:"rid"`
}

type ResponseChatInfo struct {
	MType    int    `json:"mtype"`
	NickName string `json:"nikcname"`
	Uid      int    `json:"uid"`
}

type ResponseChatInfoV2 struct {
	MType     int    `json:"mtype"`
	NickName  string `json:"nikcname"`
	Uid       int    `json:"uid"`
	UserLevel int    `json:"user_level"`
	Guard     int    `json:"guard"`
	Super     bool   `json:"super"`
}
type ResponseTaskInfo struct {
	MType     int            `json:"mtype"`
	Tasks     []TaskInfoSend `json:"tasks"`
	Current   int            `json:"current"`
	Finish    int            `json:"finish"`
	TaskNum   int            `json:"task_num"`
	AllFinish int            `json:"all_finish"`
}

type OnConnectRes struct {
	MType int `json:"mtype"`
}

type ResponseNiuNiuRecord struct {
	MType  int            `json:"mtype"`
	Record []NiuNiuRecord `json:"record"`
}

type ResponseFocus struct {
	MType         int    `json:"mtype"`
	Uid           int    `json:"uid"`
	NickNameSelf  string `json:"nikcname_self"`
	Uid2          int    `json:"uid2"`
	NickNameOther string `json:"nikcname_other"`
}

type ResponseLetterUnread struct {
	MType int `json:"mtype"`
	Num   int `json:"num"`
}

type ResponseGuardOpen struct {
	MType     int    `json:"mtype"`
	Uid       int    `json:"uid"`
	NickName  string `json:"nick_name"`
	AnchorId  int    `json:"anchor_id"`
	AchorName string `json:"anchor_name"`
	First     int    `json:"first"`
}

var chat *melody.Melody

func MelodyInit() *melody.Melody {
	chat = melody.New()
	chat.HandleConnect(ConnectNewSession)
	chat.HandleMessage(ChoseRoom)
	chat.HandleError(ErrorSession)
	chat.HandleDisconnect(DisConnectNewSession)
	chat.HandleExpection(ExpectionSession)

	go func() {
		for v := range chat.PongChan {
			uid, _ := strconv.Atoi(v)
			u, err := GetCacheUser(uid)
			if err == redis.Nil {
			} else if err != nil {
			} else {
				SetCacheUser(uid, u)
			}
		}
	}()

	return chat
}

func GetChat() *melody.Melody {
	return chat
}

func ExpectionSession(session *melody.Session) {
	defer common.PrintPanicStack()
	req := session.Request
	uid := req.FormValue("uid")
	user, _ := GetUserByUidStr(uid)

	reg := regexp.MustCompile(`[0-9]+`)

	roomid := reg.FindAllString(session.Request.URL.Path, -1)
	if len(roomid) == 0 {
		common.Log.Errf("get room panic")
		//session.Close()
		return
	}

	DelUserSession(user.Uid)
}
func ConnectNewSession(session *melody.Session) {

	defer common.PrintPanicStack()
	req := session.Request
	channel := req.FormValue("channel")
	if channel == "web" {
		return
	}

	token := req.FormValue("token")
	if token == "" {

		session.Close()
		return
	}
	uid := req.FormValue("uid")
	user, _ := GetUserByUidStr(uid)
	if user.Token != token {

		session.Close()
		return
	}
	c_user := &UserInfo{}

	user.GetChatUser(c_user)
	reg := regexp.MustCompile(`[0-9]+`)

	roomid := reg.FindAllString(session.Request.URL.Path, -1)
	if len(roomid) == 0 {
		common.Log.Errf("get room panic")
		session.Close()
		return
	}

	roomid_ := roomid[0]

	c_room := GetChatRoom(roomid_)
	if c_room == nil {
		common.Log.Errf("room is not exist ")
		session.Close()
		return
	}
	ip := common.GetRemoteIp(req)
	ret := AddUserToChat(roomid_, c_user, ip)
	if !ret {

		common.Log.Errf("add user to chat uid=?", c_user.Chat.Uid)
		session.Close()
		return
	}

	AddUserSession(user.Uid, session, roomid_, token)

	var status int
	if c_room.room.Uid == user.Uid {
		status = common.USER_STATUE_LIVE
	} else {
		status = common.USER_STATUE_SEE
	}

	u, err := GetCacheUser(user.Uid)
	if err == redis.Nil {
		s := &CacheUser{
			Uid:    user.Uid,
			Status: status,
			RoomId: roomid_,
		}
		SetCacheUser(user.Uid, s)
	} else if err != nil {
		return
	} else {
		u.Status = status
		u.RoomId = roomid_
		SetCacheUser(user.Uid, u)

	}

	notice := GetSysNotice()
	if notice == "" {
		return
	}
	var data ResponseSys
	data.MType = common.MESSAGE_TYPE_SYS
	data.Notice = notice
	SendMsgToUser(user.Uid, data)

	var data2 OnConnectRes
	data2.MType = common.MESSAGE_TYPE_ON_CONNECT
	SendMsgToUser(user.Uid, data2)
}

func ErrorSession(session *melody.Session, err error) {
	defer common.PrintPanicStack()
	req := session.Request
	uid := req.FormValue("uid")

	reg := regexp.MustCompile(`[0-9]+`)

	roomid := reg.FindAllString(session.Request.URL.Path, -1)
	if len(roomid) == 0 {
		common.Log.Errf("get room panic")
		session.Close()
		return
	}

	roomid_ := roomid[0]

	common.Log.Errf("melody err is %s,uid=%s,rid=%s", err.Error(), uid, roomid_)
}

func ChoseRoom(s *melody.Session, msg []byte) {
	defer common.PrintPanicStack()
	req := s.Request
	token := req.FormValue("token")
	//godump.Dump(string([]byte(msg)))

	if token != "" {
		user, has := GetUserByToken(token)
		if !has {
			var res ResponseErr
			res.MType = common.MESSAGE_TYPE_ERR
			res.Err = common.ERR_TOEKN_EXPIRE
			if b, err := json.Marshal(res); err == nil {
				msg = b
			}
			chat.BroadcastFilter(msg, func(q *melody.Session) bool {
				return s == q
			})
			return
		}

		sess := GetUserSessByUid(user.Uid)
		if sess == nil {
			return
		}

		gag := GetGagByUidAndRoomID(user.Uid, sess.Roomid)
		if gag != nil {
			return
		}
		rinfo := GetChatRoom(sess.Roomid)
		if rinfo == nil {
			return
		}
		rinfo.LastSayTime = time.Now().Unix() + common.ROBOT_SAY_TIMER
		js, err := simplejson.NewJson(msg)
		if err != nil {
			common.Log.Errf("orm err is 1 %s", err.Error())
			return
		}
		mtype := js.Get("mtype")
		msgtype, err := mtype.Int()
		if err != nil {
			common.Log.Errf("orm err is 2 %s", err.Error())
			return
		}
		switch msgtype {
		case 1:
			var data Response
			clientmsg := js.Get("msg")

			str, err := clientmsg.String()
			if err != nil {
				common.Log.Errf("orm err is 3 %s", err.Error())
				return
			}

			data.Msg = sensitive.GetSensitiveWord(str)
			data.Uid = user.Uid
			data.Name = user.NickName
			data.MType = common.MESSAGE_TYPE_COMMOM
			data.Face = user.Image
			data.Sex = user.Sex
			data.Location = user.Location
			data.UserLevel = user.UserLevel
			data.AnchorLevel = user.AnchorLevel
			data.IsSuperUser = user.IsSuperUser()
			data.Guard = CheckGuard(user.Uid, rinfo.room.Uid)
			chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
				return s != q && q.Request.URL.Path == s.Request.URL.Path
			})
		case 3:
			var data ResponseStar
			data.MType = common.MESSAGE_TYPE_STAR
			/*
				if b, err := json.Marshal(data); err == nil {
					msg = b
				}
			*/
			chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
				return s != q && q.Request.URL.Path == s.Request.URL.Path
			})
		case common.MESSAGE_TYPE_USER_STATUE_JOIN:
			var data ResponseHints
			data.MType = common.MESSAGE_TYPE_USER_STATUE_JOIN
			/*
				if b, err := json.Marshal(data); err == nil {
					msg = b
				}
			*/
			chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
				return s != q && q.Request.URL.Path == s.Request.URL.Path
			})
		case common.MESSAGE_TYPE_USER_STATUE_LEVEL:
			var data ResponseHints
			data.MType = common.MESSAGE_TYPE_USER_STATUE_LEVEL
			/*
				if b, err := json.Marshal(data); err == nil {
					msg = b
				}
			*/
			chat.BroadcastFilterByJson(data, func(q *melody.Session) bool {
				return s != q && q.Request.URL.Path == s.Request.URL.Path
			})

		case common.MESSAGE_TYPE_MULTIPLE_IS_AGREE:
			var data ResponseInviteResult
			data.MType = common.MESSAGE_TYPE_MULTIPLE_IS_AGREE

			data.NickName = user.NickName
			data.UID = user.Uid

			anchorUID := js.Get("anchor_uid")
			anchorUID2, err := anchorUID.Int()
			if err != nil {
				common.Log.Errf("orm err is 2 %s", err.Error())
				return
			}

			isAgree := js.Get("is_agree")
			isAgree2, err := isAgree.Int()
			if err != nil {
				common.Log.Errf("orm err is 2 %s", err.Error())
				return
			}

			data.IsAgree = isAgree2

			SendMsgToUser(anchorUID2, data)

		case common.MESSAGE_TYPE_MULTIPLE_CANCEL:
			//向副主播发送取消聊天的消息
			var data ResponseHints
			data.MType = common.MESSAGE_TYPE_MULTIPLE_CANCEL

			deputyAnchorUID := js.Get("deputy_anchor_uid")
			deputyAnchorUID_, err := deputyAnchorUID.String()
			deputyAnchorUID2, _ := strconv.Atoi(deputyAnchorUID_)

			if err != nil {
				common.Log.Errf("orm err is 2 %s", err.Error())
				return
			}

			SendMsgToUser(deputyAnchorUID2, data)

		case common.MESSAGE_TYPE_SEND_GIFT:
			SendGitByWebSocketStyle(s, msg)

		case common.MESSAGE_TYPE_DIAMOND_SCORE_INFO:
			SendDiamondScoreByWebSocketStyle(s, msg, user.Uid)

		/////////////////////////////////////////////////
		case common.MESSAGE_TYPE_GAME_START:
			AnchorGameStartMsg(s, msg, user, rinfo)
		case common.MESSAGE_TYPE_GAME_STOP:

			userID := req.FormValue("uid")
			userID_, err := strconv.Atoi(userID)
			if err != nil {
				common.Log.Errf("orm err is 2 %s", err.Error())
				return
			}

			if userID_ != rinfo.room.Uid {
				var data ResponseHints
				data.MType = common.MESSAGE_TYPE_GAME_START_EOR
				SendMsgToUser(userID_, data)
				return
			} else if rinfo.GameType == common.GAME_TYPE_NIUNIU {
				t := GetNiuNiuByRoomid(rinfo.room.Rid)
				if t != nil {
					t.IsStarted = false
					var data ResponseHints
					data.MType = common.MESSAGE_TYPE_GAME_STOP
					return
				}
			} else if rinfo.GameType == common.GAME_TYPE_TEXAS {
				t := GetTexasByRoomid(rinfo.room.Rid)
				if t != nil {
					t.IsStarted = false
					var data ResponseHints
					data.MType = common.MESSAGE_TYPE_GAME_STOP
					return
				}
			} else if rinfo.GameType == common.GAME_TYPE_GOLDEN_FLOWER {
				t := GetGoldenFlowerByRoomid(rinfo.room.Rid)
				if t != nil {
					t.IsStarted = false
					var data ResponseHints
					data.MType = common.MESSAGE_TYPE_GAME_STOP
					return
				}
			}

		case common.MESSAGE_TYPE_GAME_STATE_USER_RAISING: //押分

			u, err := GetCacheUser(user.Uid)
			if err == redis.Nil {
				return
			} else if err != nil {
				return
			}

			if rinfo.GameType == common.GAME_TYPE_NIUNIU {

				UserRaiseMsgDispose(s, msg, u.RoomId, user)
			} else if rinfo.GameType == common.GAME_TYPE_TEXAS {
				TexasPokRaiseMsgDispose(s, msg, u.RoomId, user)
			} else if rinfo.GameType == common.GAME_TYPE_GOLDEN_FLOWER {
				GoldenFlowerUserRaiseMsgDispose(s, msg, u.RoomId, user)
			}

		case common.MESSAGE_TYPE_GAME_STATE:

			u, err := GetCacheUser(user.Uid)
			if err == redis.Nil {
				return
			} else if err != nil {
				return
			}

			if rinfo.GameType == common.GAME_TYPE_NIUNIU {
				NoticeNiuNiuStateToUser(s, user.Uid, u.RoomId)
			} else if rinfo.GameType == common.GAME_TYPE_TEXAS {
				NoticeTexasPokStateToUser(s, user.Uid, u.RoomId)
			} else if rinfo.GameType == common.GAME_TYPE_GOLDEN_FLOWER {
				NoticeGoldenFlowerStateToUser(s, user.Uid, u.RoomId)
			}

		case common.MESSAGE_TYPE_GAME_TOY_CATCH: //抓娃娃消息处理
			ToyCashMsgDispose(s, msg, user, rinfo)
		case common.MESSAGE_TYPE_GAME_TOY_CATCH_NEXT: //抓娃娃消息处理第二步【二次确认】
			ToyCashMsgDisposeNext(s, msg, user, rinfo)
		case common.MESSAGE_TYPE_GAME_TOY_NOT_CATCH: //娃娃没碰到时候【仅仅扣游戏币了事】
			ToyNotCashMsgDispose(s, msg, user, rinfo)

		/////////////////////////////////////////////
		case common.MESSAGE_TYPE_GAME_RECORD_REQ:

			roomId, err := js.Get("room_id").String()

			if err != nil {
				common.Log.Errf("orm err is 2 %s", err.Error())
				return
			}
			s := &ResponseNiuNiuRecord{}
			s.Record = make([]NiuNiuRecord, 0)
			s.Record = GetLastestRecordByRid(roomId)
			s.MType = common.MESSAGE_TYPE_GAME_RECORD_RES

			SendMsgToUser(user.Uid, s)

		default:
			return
		}
	} else {
		return
	}
}

func AnchorGameStartMsg(s *melody.Session, msg []byte, user *User, r *ChatRoomInfo) {
	req := s.Request

	userID := req.FormValue("uid")
	userID_, err := strconv.Atoi(userID)
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	if userID_ != r.room.Uid {
		var data ResponseHints
		data.MType = common.MESSAGE_TYPE_GAME_START_EOR
		SendMsgToUser(userID_, data)
		return
	} else if r.GameType == common.GAME_TYPE_NIUNIU {
		t := GetNiuNiuByRoomid(r.room.Rid)
		if t == nil {
			var data ResponseHints
			data.MType = common.MESSAGE_TYPE_CLOSE
			SendMsgToUser(userID_, data)
			return
		} else {
			if !t.IsStarted {
				t.IsStarted = true
				var data ResponseHints
				data.MType = common.MESSAGE_TYPE_GAME_START
				SendMsgToUser(r.room.Uid, data)

				go NewGameNiuNiu(s, r.room.Uid, r.room.Rid)
			} else {
				var data ResponseHints
				data.MType = common.MESSAGE_TYPE_GAME_START_EOR
				SendMsgToUser(r.room.Uid, data)
			}
		}
	} else if r.GameType == common.GAME_TYPE_TEXAS {
		t := GetTexasByRoomid(r.room.Rid)
		if t == nil {
			var data ResponseHints
			data.MType = common.MESSAGE_TYPE_CLOSE
			SendMsgToUser(userID_, data)
			return
		} else {
			if !t.IsStarted {
				t.IsStarted = true
				var data ResponseHints
				data.MType = common.MESSAGE_TYPE_GAME_START
				SendMsgToUser(r.room.Uid, data)

				go NewGameTexasPok(s, r.room.Uid, r.room.Rid)
			} else {
				var data ResponseHints
				data.MType = common.MESSAGE_TYPE_GAME_START_EOR
				SendMsgToUser(r.room.Uid, data)
			}
		}
	} else if r.GameType == common.GAME_TYPE_GOLDEN_FLOWER {
		t := GetGoldenFlowerByRoomid(r.room.Rid)
		if t == nil {
			var data ResponseHints
			data.MType = common.MESSAGE_TYPE_CLOSE
			SendMsgToUser(userID_, data)
			return
		} else {
			if !t.IsStarted {
				t.IsStarted = true
				var data ResponseHints
				data.MType = common.MESSAGE_TYPE_GAME_START
				SendMsgToUser(r.room.Uid, data)

				go NewGameGoldenFlower(s, r.room.Uid, r.room.Rid)
			} else {
				var data ResponseHints
				data.MType = common.MESSAGE_TYPE_GAME_START_EOR
				SendMsgToUser(r.room.Uid, data)
			}
		}
	}
}

type ResponseOfSendGitByWebSocket struct {
	MType    int    `json:"mtype"`
	ErrCode  int    `json:"ErrCode"`
	Diamond  int    `json:"diamond"`
	Score    int64  `json:"Score"`
	GiftName string `json:"gift_name"`
}

func SendGitByWebSocketStyle(s *melody.Session, msg []byte) {
	var data ResponseOfSendGitByWebSocket

	js, err := simplejson.NewJson(msg)

	req := s.Request
	token := req.FormValue("token")

	giftId := js.Get("giftid")
	giftId_, err := giftId.Int()

	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	uID := js.Get("uid")
	uID_, err := uID.Int()

	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	num := js.Get("num")
	num_, err := num.Int()
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	revid := js.Get("revid")
	revid_, err := revid.Int()
	if err != nil {
		common.Log.Errf("orm err is 2 %s", err.Error())
		return
	}

	data.MType = common.MESSAGE_TYPE_SEND_GIFT
	user, _ := GetUserByToken(token)

	g, ok := GetGiftById(giftId_)
	if ok == false {
		return
	}

	if g.Category == 1 || g.Category == 0 {
		data.ErrCode = user.SendGiftV2(giftId_, num_, revid_)
	} else if g.Category == 2 || g.Category == 3 {
		data.ErrCode = user.TipGiftV2(giftId_, num_, revid_)
	}

	data.GiftName = g.Name
	data.Diamond = user.Diamond
	data.Score = user.Score
	SendMsgToUser(uID_, data)
}

type ResponseDiamondScore struct {
	MType   int   `json:"mtype"`
	ErrCode int   `json:"ErrCode"`
	Diamond int   `json:"diamond"`
	Score   int   `json:"Score"`
	Coupons int   `json:"coupons"`
	Moon    int64 `json:"moon"`
}

//websocet方式获取用户钻石和游戏币余额
func SendDiamondScoreByWebSocketStyle(s *melody.Session, msg []byte, uID int) {
	var data ResponseDiamondScore

	data.MType = common.MESSAGE_TYPE_DIAMOND_SCORE_INFO
	user, _ := GetUserByUid(uID)

	data.Diamond = user.Diamond
	data.Score = int(user.Score)
	data.Coupons = user.Coupons
	data.Moon = user.Moon

	SendMsgToUser(uID, data)
}

//func StopNiuNiu(anchorID int, roomID string) {
//	room := GetChatRoom(roomID)
//	room.niuNiu.IsNiuNiu = false
//}

func DisConnectNewSession(session *melody.Session) {
	//godump.Dump("disconnect call back")
	defer common.PrintPanicStack()

	req := session.Request

	token := req.FormValue("token")

	if token == "" {
		return
	}
	uid := req.FormValue("uid")
	user, _ := GetUserByUidStr(uid)

	reg := regexp.MustCompile(`[0-9]+`)
	roomid := reg.FindAllString(session.Request.URL.Path, -1)
	if len(roomid) == 0 {
		return
	}

	roomid_ := roomid[0]

	//common.Log.Debugf("begin disconnect conn uid=%d,%s", uid, roomid_)
	DestoryChatInfo(user, roomid_)
}

func SysSay(rid string, msg []byte, uid int, token string) bool {
	url := fmt.Sprintf("/channel/%s/ws", rid)
	md := GetChat()
	md.BroadcastFilter(msg, func(q *melody.Session) bool {
		req := q.Request
		st := req.FormValue("token")

		if st != "" {
			return q.Request.URL.Path == url && st != token
		}
		return q.Request.URL.Path == url
	})
	return true
}

func SysSayToAll(rid string, msg []byte) bool {
	url := fmt.Sprintf("/channel/%s/ws", rid)
	md := GetChat()
	md.BroadcastFilter(msg, func(q *melody.Session) bool {
		return q.Request.URL.Path == url
	})
	return true
}

func SysNotice(rid string, res interface{}, uid int, token string) {
	if b, err := json.Marshal(res); err == nil {
		SysSay(rid, b, uid, token)
	}
}

func SysNoticeToAll(rid string, res interface{}) {
	if b, err := json.Marshal(res); err == nil {
		SysSayToAll(rid, b)
	}
}

func GiftSay(rid string, res *ResponseGift, uid int, token string) {
	if b, err := json.Marshal(res); err == nil {
		SysSay(rid, b, uid, token)
	}
}

func CloseSay(rid string, res *ResponseClose, uid int, token string) {
	if b, err := json.Marshal(res); err == nil {
		SysSay(rid, b, uid, token)
	}

}
func AdminSysToAll(res interface{}) {
	if b, err := json.Marshal(res); err == nil {
		GetChat().Broadcast(b)
	}
}

/*
func SysSayAll(rid string, msg []byte, uid int) {
	url := fmt.Sprintf("/channel/%s/ws", rid)
	md := GetChat()
	md.BroadcastFilter(msg, func(q *melody.Session) bool {
		return q.Request.URL.Path == url
	})
}
*/
/*
func MonnitorSay(rid string, res *ResponseSys) {
	chat := GetChatRoom(rid)
	if chat != nil {
		SendMsgToUser(chat.room.Uid, res)
	}
}
*/
//预先销毁房间1先删除用户2添加主播id到等待重连队列
func DestoryChatInfo(user *User, rid string) {
	if chat := GetChatRoom(rid); chat == nil {
		common.Log.Infof("already del room rid=?", rid)
		return
	}

	u, err := GetCacheUser(user.Uid)
	if err == redis.Nil {
		s := &CacheUser{
			Uid:    user.Uid,
			Status: common.USER_STATUE_LEAVE,
		}
		SetCacheUser(user.Uid, s)
	} else if err != nil {
		return
	} else {
		u.Status = common.USER_STATUE_LEAVE
		SetCacheUser(user.Uid, u)
	}
	DelUserFromChat(rid, user.Uid)

	DelUserSession(user.Uid)
}

//销毁房间设置房间状态
func CloseChat(uid int, rid string) {
	common.Log.Infof("call back now begin close room uid=?,rid=?,time=?", uid, rid, time.Now().Unix())
	//room, has := chat_room_manager[rid]
	room := GetChatRoom(rid)
	if room != nil {
		if room.Statue != common.ROOM_ONLIVE {
			common.Log.Debugf("now time ? begin close multiple statue 2 room=?,uid=?", time.Now().Unix(), room.room.Rid, uid)
			return
		}
		room.Statue = common.ROOM_FINISH

		if room.IsMultiple {
			common.Log.Infof("now time ? begin close multiple room=?,uid=?", time.Now().Unix(), room.room.Rid, uid)
			CheckAndCloseMultiple(uid, room.room.Rid)

			res := &ResponseClose{MType: common.MESSAGE_TYPE_CLOSE}
			user, _ := GetUserByUid(uid)
			SysNotice(rid, res, user.Uid, user.Token)
			user.LeaveRoom()
			room.Clear()

		} else {
			CloseRoom(room, uid)
			res := &ResponseClose{MType: common.MESSAGE_TYPE_CLOSE}
			user, _ := GetUserByUid(uid)
			SysNotice(rid, res, user.Uid, user.Token)

			user.LeaveRoom()
			common.Log.Infof("delete room nomarl uid=%d,rid=%s", uid, rid)

			if room.Statue != common.ROOM_PLAYBACK {
				DelChatRoom(rid)
			} else {
				common.Log.Infof("begin close play back room uid=%d,rid=%s", uid, rid)
			}

			//CloseMultipleRoom(uid)
			if room.GameType == common.GAME_TYPE_NIUNIU { // 如果客户端是要德州扑克删掉德州索引里的那玩意
				DelRoomInfoNiuNiu(rid)
			} else if room.GameType == common.GAME_TYPE_TEXAS { // 如果客户端是要德州扑克删掉德州索引里的那玩意
				DelRoomInfoTexas(rid)
			} else if room.GameType == common.GAME_TYPE_GOLDEN_FLOWER {
				DelRoomInfoGoldenFlower(rid)
			}

			UpdateRecommandFlag(uid, rid)

		}

	}

}

//直接删除房间不进入等待队列
func DirectCloseRoom(uid int, rid string) {
	u, err := GetCacheUser(uid)
	if err == redis.Nil {
		s := &CacheUser{
			Uid:    uid,
			Status: common.USER_STATUE_LEAVE,
		}
		SetCacheUser(uid, s)
	} else if err != nil {
		return
	} else {
		u.Status = common.USER_STATUE_LEAVE
		SetCacheUser(uid, u)
	}

	DelUserSession(uid)

	CloseChat(uid, rid)
}
