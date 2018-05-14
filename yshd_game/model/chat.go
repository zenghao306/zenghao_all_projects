package model

import (
	"github.com/yshd_game/common"
	"github.com/yshd_game/melody"
	"github.com/yshd_game/timer"
	//"github.com/yshd/melody"
	"encoding/json"
	"fmt"
	"github.com/yshd_game/wrap"
	"net/http"
	"sync"
	"time"
	//"github.com/liudng/godump"
)

type ChatUser struct {
	Uid         int
	Image       string
	NickName    string
	Location    string
	UserLevel   int
	AnchorLevel int
	Sex         int
	Signature   string
	IsSuperUser bool
	Guard       int
}

type UserInfo struct {
	Chat    ChatUser
	IsRobit bool

	//Session *melody.Session
}

//type NiuNiuGameInfo struct {
//	IsNiuNiu    bool            //是否在牛牛游戏中
//	NiuNiuState int             //牛牛状态
//	CSStartTime int64           //当前状态下开始时间戳
//	GameID      string          //当前局的游戏ID
//	niuActors   []*NiuActorInfo //当前局下的牌
//	LScore      int
//	MScore      int
//	RScore      int
//}

type ChatRoomInfo struct {
	//	UList        list.List
	//	UserMap      map[int]bool

	UserInfoMap map[int]*UserInfo
	//mutex_room_info sync.RWMutex

	room         RoomInfo
	count        int
	rice         int
	moon         int
	Statue       int
	Save         int
	VRobotNumber int //虚拟机器人数额
	//Session         *melody.Session //房主的sesion
	LastSayTime int64

	SayRecord  RobotSayInfo
	RoomType   int
	IsMultiple bool
	GameType   int
	//niuNiu     NiuNiuGameInfo
	MsgProc chan *wrap.Call
	die     chan bool
}

type SessUser struct {
	Uid    int
	Token  string
	Sess   *melody.Session
	Roomid string
}

//用户的ID session索引管理器
var mutex_session_guard sync.RWMutex
var user_session_manager map[int]*SessUser

func NewChatRoomInfo() *ChatRoomInfo {
	m := &ChatRoomInfo{}
	m.UserInfoMap = make(map[int]*UserInfo, 0)
	m.SayRecord.Init(lenghtSayRobot)
	m.MsgProc = make(chan *wrap.Call, 50)
	m.die = make(chan bool, 1)
	go func(c *ChatRoomInfo) {
	loop2:
		for {
			select {
			case msg, ok := <-c.MsgProc:
				if ok {
					SendGift(msg)
					msg.DoneV2()
				} else {
					common.Log.Infof("close2 null go rid=%s", c.room.Rid)
					break loop2
				}
			case <-c.die:
				common.Log.Infof("close go rid=%s", c.room.Rid)
				return
			}
		}
	}(m)

	return m
}

func SendGift(call *wrap.Call) {
	r := call.Render
	//d := call.Request.(map[string]interface{})
	var ok bool
	/*
		var uid_ int
		uid,ok:=d["uid"]
		if !ok {
			uid_=0
		}else{
			uid_=uid.(int)
		}
	*/

	uid := call.Uid
	token_ := call.Token
	gift_id_ := call.GiftID
	num_ := call.Num
	rev_id_ := call.RevId
	/*
		token,_:=d["token"]
		token_:=token.(string)
		gift_id,_:=d["gift_id"]
		gift_id_:=gift_id.(int)
		num,_:=d["num"]
		num_:=num.(int)
		rev_id,_:=d["rev_id"]
		rev_id_:=rev_id.(int)
	*/

	ret_value := make(map[string]interface{})

	var user *User

	var ret int
	if uid == 0 {
		user, ok = GetUserByToken(token_)
		if ok == false {
			ret_value["ErrCode"] = common.TOKEN_EXPIRE_TIME
			r.JSON(http.StatusOK, ret_value)
			return
		}
	} else {
		user, ret = GetUserByUid(uid)
		if ret != common.ERR_SUCCESS {
			ret_value["ErrCode"] = common.ERR_ACCOUNT_EXIST
			r.JSON(http.StatusOK, ret_value)

			return
		}
	}
	gift, ok := GetGiftById(gift_id_)
	if !ok {
		ret_value["ErrCode"] = common.ERR_GIFT_EXIST
		r.JSON(http.StatusOK, ret_value)

		return
	}

	if gift.Category == common.GIFT_CATEGORY_EXTRAVAGANT || gift.Category == common.GIFT_CATEGORY_HOT {
		ret = user.SendGiftV2(gift_id_, num_, rev_id_)
	} else {
		ret = user.TipGiftV2(gift_id_, num_, rev_id_)
	}

	//ret := user.SendGift(d.GiftId, d.Num, d.Revid)
	ret_value["ErrCode"] = ret

	user_new, _ := GetUserByToken(token_)
	ret_value["diamond"] = user_new.Diamond
	ret_value["score"] = user_new.Score

	r.JSON(http.StatusOK, ret_value)
}

func (self *ChatRoomInfo) GetChatInfo() *RoomInfo {
	return &self.room
}

func (self *ChatRoomInfo) Clear() {
	/*
		for e := self.UList.Front(); e != nil; e = e.Next() {
			self.UList.Remove(e)
		}
	*/
	self.count = 0
	self.LastSayTime = 0
	self.rice = 0
	self.Save = 0
	self.moon = 0
	self.SayRecord.Reset()
}
func (self *ChatRoomInfo) GetCount() int {
	//self.mutex_room_info.RLock()
	//defer self.mutex_room_info.RUnlock()
	return len(self.UserInfoMap)
	//return self.UList.Len()
}

func (self *ChatRoomInfo) GetChatRealUserCount() int {
	//self.mutex_room_info.RLock()
	//defer self.mutex_room_info.RUnlock()
	num := 0
	for _, u := range self.UserInfoMap {
		if u.IsRobit == false {
			num++
		}
	}
	return num
}

func (self *ChatRoomInfo) GetRice() int {
	return self.rice
}

func (self *ChatRoomInfo) GetVRobotCount() int {
	return self.VRobotNumber
}

func (self *ChatRoomInfo) AddRice(num int) {
	self.rice += num
}

func (self *ChatRoomInfo) AddMoon(num int) {
	self.moon += num
}

func (self *ChatRoomInfo) GetMoon() int {
	return self.moon
}

func (self *ChatRoomInfo) UpdateAudience(uid int) {
	var score int
	c, ok := self.UserInfoMap[uid]
	if ok {
		c.Chat.Guard = CheckGuard(uid, self.room.Uid)
		if c.Chat.Guard == 1 {
			score = 10000000000
		}
		self.UserInfoMap[uid] = c

		IncrAudience(self.room.Rid, uid, score)
	}
}

//新增观众使用redis保存列表
func (self *ChatRoomInfo) AddAudience(u *UserInfo, score int) int {
	//self.mutex_room_info.Lock()
	//defer self.mutex_room_info.Unlock()

	_, ok := self.UserInfoMap[u.Chat.Uid]
	if ok {
		return common.ERR_AUDIENCE_EXIST
	}

	u.Chat.Guard = CheckGuard(u.Chat.Uid, self.room.Uid)
	if u.Chat.Guard == 1 {
		score += 10000000000
	}
	self.UserInfoMap[u.Chat.Uid] = u
	err := AddAudience(self.room.Rid, u.Chat.Uid, score)
	if err != nil {
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS
}

//删除观众从redis中优先删除观众
func (self *ChatRoomInfo) DelAudience(uid int) int {
	//	self.mutex_room_info.Lock()
	//	defer self.mutex_room_info.Unlock()
	err := DelAudience(self.room.Rid, uid)
	if err != nil {
		return common.ERR_UNKNOWN
	}
	delete(self.UserInfoMap, uid)
	return common.ERR_SUCCESS
}

//根据redis中排序，并且取出本地数据返回
func (self *ChatRoomInfo) GetAudienceList(index int) (users []*ChatUser, ret int) {

	//	self.mutex_room_info.RLock()
	//	defer self.mutex_room_info.RUnlock()
	us, err := GetAudience(self.room.Rid, int64(index*10), int64((index+1)*10))
	if err != nil {
		ret = common.ERR_UNKNOWN
		return
	}

	users = make([]*ChatUser, 0)
	for _, v := range us {
		u, ok := self.UserInfoMap[v]
		if ok {
			users = append(users, &u.Chat)
		}
	}
	ret = common.ERR_SUCCESS
	return
}

func (self *ChatRoomInfo) CheckAudienceByUid(uid int) bool {
	//	self.mutex_room_info.RLock()
	//	defer self.mutex_room_info.RUnlock()
	_, ok := self.UserInfoMap[uid]
	return ok
}

func GetDump(rid string) int {
	res, err := orm.Query("select sum(money_num) as dump from go_niu_niu_record where room_id=?", rid)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return 0
	}

	if len(res) == 0 {
		return 0
	}

	dump, ok := res[0]["dump"]
	if ok {
		return common.BytesToInt(dump)
	}
	return 0
}

func (self *ChatRoomInfo) GetRobotSay() (string, bool) {
	id, ret := self.SayRecord.GetRandomSayInChat()
	if ret == false {
		return "", false
	}
	say, ret := GetRobotSayByID(id)
	if ret == false {
		return "", false
	}
	return say, true
}

func InitChatSession() map[int]*SessUser {
	user_session_manager = make(map[int]*SessUser)
	return user_session_manager
}

/*
func (v *SessUser) SendMsg(data interface{})  {
	if v.Sess != nil {
		v.Sess.CloseWithMsgAndJson(data)
	}
}

func (v *SessUser) CloseSesion() {
	if v.Sess != nil {
		var data ResponseErr
		data.Err = common.ERR_OTHER_LOGIN
		data.MType = common.MESSAGE_TYPE_ERR
		//SendMsgToUserWithClose(v.Uid, data)
		//v.Sess.Close()
		v.Sess.CloseWithMsgAndJson(data)
	}
}
*/

func GetUserSessCount() int {
	mutex_session_guard.RLock()
	defer mutex_session_guard.RUnlock()
	return len(user_session_manager)
}

func GetUserSessByUid(uid int) *SessUser {
	mutex_session_guard.RLock()
	defer mutex_session_guard.RUnlock()
	v, ok := user_session_manager[uid]
	if !ok {
		return nil
	}
	return v
}
func AddUserSession(uid int, sess *melody.Session, roomid string, token string) {
	sess_user_ := GetUserSessByUid(uid)
	if sess_user_ != nil {
		common.Log.Errf("add user session err uid is %d", uid)
		if sess_user_.Roomid != roomid {
			sess_user_.Sess.Close()
			user, _ := GetUserByUid(uid)
			DestoryChatInfo(user, sess_user_.Roomid)
		}
	}
	sess_user := &SessUser{}
	sess_user.Uid = uid
	sess_user.Sess = sess
	sess_user.Roomid = roomid
	sess_user.Token = token
	mutex_session_guard.Lock()
	defer mutex_session_guard.Unlock()
	user_session_manager[uid] = sess_user

}

func DelUserSession(uid int) bool {
	mutex_session_guard.Lock()
	defer mutex_session_guard.Unlock()
	_, ok := user_session_manager[uid]
	if !ok {
		common.Log.Debugf("del user_session_manager uid=%d", uid)
		return false
	}

	delete(user_session_manager, uid)
	return true
}

//房间ID 的索引管理器
var chat_room_manager map[string]*ChatRoomInfo
var mutex_chat_guardv2 sync.RWMutex

func InitChat() map[string]*ChatRoomInfo {
	chat_room_manager = make(map[string]*ChatRoomInfo)
	return chat_room_manager
}
func GetChatRoom(roomid string) *ChatRoomInfo {
	mutex_chat_guardv2.RLock()
	defer mutex_chat_guardv2.RUnlock()
	v, ok := chat_room_manager[roomid]
	if !ok {
		return nil
	}
	return v
}

func AddChatRoom(chat *ChatRoomInfo) bool {
	mutex_chat_guardv2.Lock()
	defer mutex_chat_guardv2.Unlock()
	_, ok := chat_room_manager[chat.room.Rid]
	if ok {
		return false
	}
	chat_room_manager[chat.room.Rid] = chat

	return true
}

func CheckChatRobotTimer() {
	curtime := time.Now().Unix()
	mutex_chat_guardv2.Lock()
	defer mutex_chat_guardv2.Unlock()
	for k, v := range chat_room_manager {
		vlen := len(v.UserInfoMap)
		if vlen <= 2 {
			continue
		}

		if v.LastSayTime <= curtime {
			r := common.RadnomRange(0, vlen-1)
			i := 0
			var robot int

			for _, e := range v.UserInfoMap {
				if i <= r {
					i++
					continue
				} else {
					i++
				}
				if e.IsRobit == true {
					robot = e.Chat.Uid
					break
				}
			}
			/*
				for e := v.UList.Front(); e != nil; e = e.Next() {
					if i <= r {
						i++
						continue
					} else {
						i++
					}

					s := e.Value.(*UserInfo)
					if s.IsRobit == true {
						robot = s.Chat.Uid
						break
					}
				}
			*/
			user, ret_has := GetUserByUid(robot)
			if ret_has != common.ERR_SUCCESS {
				break
			}
			var data Response

			say, ret := v.GetRobotSay()

			if ret == false {
				break
			}
			data.Msg = say
			//data.Msg = GetRandomSay()

			data.Uid = user.Uid
			data.Name = user.NickName
			data.MType = common.MESSAGE_TYPE_COMMOM
			data.Face = user.Image
			data.Sex = user.Sex
			data.Location = user.Location
			data.UserLevel = user.UserLevel
			data.AnchorLevel = user.AnchorLevel
			data.IsSuperUser = user.IsSuperUser()
			SysNoticeToAll(k, &data)

			v.LastSayTime = curtime + common.ROBOT_SAY_TIMER
		}
	}
}

func DelChatRoom(rid string) {
	common.Log.Debugf("DelChatRoom rid=%s,time=%d", rid, time.Now().Unix())

	room := GetChatRoom(rid)
	if room != nil {
		close(room.MsgProc)
		//room.die<-true
	}

	ReportAnchorDate("delete", rid, 0, 0, "", time.Now().Unix())
	mutex_chat_guardv2.Lock()
	defer mutex_chat_guardv2.Unlock()

	delete(chat_room_manager, rid)
	DelAudienceKey(rid)
}

func AddUserToChat(rid string, user *UserInfo, ip string) bool {
	room := GetChatRoom(rid)

	user_, _ := GetUserByUid(user.Chat.Uid)
	if user.Chat.Uid == room.room.Uid {
		if flag := AnchorMgr.DelReconnectMap(user.Chat.Uid); flag {
			return true
		}
		if room.IsMultiple == false {
			db_room, has := GetRoomById(rid)
			if !has {
				return false
			}
			db_room.Statue = common.ROOM_ONLIVE

			_, err := orm.Where("room_id=?", rid).MustCols("statue").Update(db_room)
			if err != nil {
				common.Log.Errf("mysql error is %s", err.Error())
				return false
			}
			ReportAnchorDate("add", db_room.RoomId, db_room.OwnerId, db_room.Weight, db_room.RoomName, db_room.CreateTime.Unix())
			room.LastSayTime = time.Now().Unix() + common.ROBOT_SAY_TIMER

			user_.JoinRoom(rid, true, room.RoomType, ip)

			room.Statue = common.ROOM_ONLIVE
		} else if room.IsMultiple == true {
			has, db_room := GetMultipleRoomByRid(rid)
			if !has {
				return false
			}
			db_room.Statue = common.MULTIPLE_ROOM_BUSY

			_, err := orm.Where("room_id=?", rid).MustCols("statue").Update(db_room)
			if err != nil {
				common.Log.Errf("mysql error is %s", err.Error())
				return false
			}
			user_.JoinRoom(rid, true, room.RoomType, ip)
			room.Statue = common.ROOM_ONLIVE
		}

		//AddRobot(common.MAX_ROBOT_ADD_ROOM_NUM, &room.UList)

		// go func() {
		// 	d := timer.NewDispatcher(1)
		// 	t := common.RadnomRange(10, 30)
		// 	d.AfterFunc(time.Duration(t)*time.Second, func() {
		// 		AddRobotFirst(&room.UList, 6, 18)
		// 	})
		// 	(<-d.ChanTimer).Cb()
		// 	//	g := common.RadnomRange(60, 120)

		// 	f := timer.NewDispatcher(1)
		// 	if room.UList.Len() >= 21 {
		// 		return
		// 	}
		// 	g := common.RadnomRange(30, 90)
		// 	f.AfterFunc(time.Duration(g)*time.Second, func() {

		// 		AddRobotFirst(&room.UList, 21-room.UList.Len(), 40-room.UList.Len())
		// 	})
		// 	(<-f.ChanTimer).Cb()

		// }()
		go func() {
			defer common.PrintPanicStack()
			d := timer.NewDispatcher(1)
			t := common.RadnomRange(0, 5)
			d.AfterFunc(time.Duration(t)*time.Second, func() {
				//vNumber := AddRobotOfNumber(&room.UList, 2)
				vNumber := AddRobotOfNumber(room.room.Rid, 2)
				room.VRobotNumber += vNumber
			})
			(<-d.ChanTimer).Cb()
			//	g := common.RadnomRange(60, 120)

			f := timer.NewDispatcher(1)
			g := common.RadnomRange(6, 10)
			f.AfterFunc(time.Duration(g)*time.Second, func() {
				//vNumber := AddRobotOfNumber(&room.UList, 3)
				vNumber := AddRobotOfNumber(room.room.Rid, 3)
				room.VRobotNumber += vNumber //虚拟机器人数量
			})
			(<-f.ChanTimer).Cb()

			//11-20秒增加4个
			f2 := timer.NewDispatcher(1)
			g = common.RadnomRange(11, 20)
			f2.AfterFunc(time.Duration(g)*time.Second, func() {
				//vNumber := AddRobotOfNumber(&room.UList, 4)
				vNumber := AddRobotOfNumber(room.room.Rid, 4)
				room.VRobotNumber += vNumber //虚拟机器人数量
			})
			(<-f2.ChanTimer).Cb()

			//21-30秒增加5个
			f3 := timer.NewDispatcher(1)
			g = common.RadnomRange(21, 30)
			f3.AfterFunc(time.Duration(g)*time.Second, func() {
				//vNumber := AddRobotOfNumber(&room.UList, 5)
				vNumber := AddRobotOfNumber(room.room.Rid, 5)
				room.VRobotNumber += vNumber //虚拟机器人数量
			})
			(<-f3.ChanTimer).Cb()

			//31-60秒增加5个
			f4 := timer.NewDispatcher(1)
			g = common.RadnomRange(31, 60)
			f4.AfterFunc(time.Duration(g)*time.Second, func() {
				//vNumber := AddRobotOfNumber(&room.UList, 5)
				vNumber := AddRobotOfNumber(room.room.Rid, 5)
				room.VRobotNumber += vNumber //虚拟机器人数量
			})
			(<-f4.ChanTimer).Cb()

			//61-120秒增加5个
			f5 := timer.NewDispatcher(1)
			g = common.RadnomRange(61, 120)
			f5.AfterFunc(time.Duration(g)*time.Second, func() {
				//vNumber := AddRobotOfNumber(&room.UList, 5)
				vNumber := AddRobotOfNumber(room.room.Rid, 5)
				room.VRobotNumber += vNumber //虚拟机器人数量
			})
			(<-f5.ChanTimer).Cb()

			//121-200秒增加5个
			f6 := timer.NewDispatcher(1)
			g = common.RadnomRange(121, 200)
			f6.AfterFunc(time.Duration(g)*time.Second, func() {
				//vNumber := AddRobotOfNumber(&room.UList, 5)
				vNumber := AddRobotOfNumber(room.room.Rid, 5)
				room.VRobotNumber += vNumber //虚拟机器人数量
			})
			(<-f6.ChanTimer).Cb()

			//201-300秒增加5个
			f7 := timer.NewDispatcher(1)
			g = common.RadnomRange(201, 300)
			f7.AfterFunc(time.Duration(g)*time.Second, func() {
				//vNumber := AddRobotOfNumber(&room.UList, 5)
				vNumber := AddRobotOfNumber(room.room.Rid, 5)
				room.VRobotNumber += vNumber //虚拟机器人数量
			})
			(<-f7.ChanTimer).Cb()

			for room.GetCount()+room.VRobotNumber < 1083 {
				f8 := timer.NewDispatcher(1)
				g = common.RadnomRange(28, 32)
				num := common.RadnomRange(13, 27)
				f8.AfterFunc(time.Duration(g)*time.Second, func() {
					//vNumber := AddRobotOfNumber(&room.UList, num)
					vNumber := AddRobotOfNumber(room.room.Rid, num)
					room.VRobotNumber += vNumber //虚拟机器人数量
				})
				(<-f8.ChanTimer).Cb()
			}

		}()

		return true
	}

	cuser := &UserInfo{}
	user_.GetChatUser(cuser)

	sess := GetUserSessByUid(user_.Uid)
	if sess != nil {
		return false
	}
	var score int
	if user_.Robot {
		score = -1
	} else {
		score = HistoryCountGiftNum(user_.Uid, room.room.Uid)
	}

	mutex_chat_guardv2.Lock()
	room.AddAudience(cuser, score)
	mutex_chat_guardv2.Unlock()

	//直播间每进入一个用户增加8个机器人
	if user.Chat.Uid != room.room.Uid && room.GetCount()+room.VRobotNumber < 1083 {
		//vNumber := AddRobotOfNumber(&room.UList, 8)
		vNumber := AddRobotOfNumber(room.room.Rid, 8)
		room.VRobotNumber += vNumber //虚拟机器人数量
	}

	var data ResponseChatInfoV2
	data.MType = common.MESSAGE_TYPE_CHAT_INFO
	data.NickName = user_.NickName
	data.Uid = user_.Uid
	data.UserLevel = user_.UserLevel

	data.Guard = CheckGuard(user_.Uid, room.room.Uid)
	data.Super = user_.IsSuperUser()
	SendMsgToRoom(rid, data)

	user_.JoinRoom(rid, false, room.RoomType, ip)

	return true
}

func DelUserFromChat(rid string, uid int) {
	room := GetChatRoom(rid)
	if room == nil {
		common.Log.Infof("user close ws uid is %d,rid is %s", uid, rid)
		return
	}
	if room.room.Uid == uid {
		AnchorMgr.AddAnchor(uid, rid)
		return
	}

	mutex_chat_guardv2.Lock()
	defer mutex_chat_guardv2.Unlock()

	if ok := room.CheckAudienceByUid(uid); ok {
		room.DelAudience(uid)
		//每退出一个用户，减去3个机器人
		if room.VRobotNumber >= 3 {
			room.VRobotNumber -= 3
			//return
		} else if room.VRobotNumber > 0 {
			room.VRobotNumber = 0
			RemoveRobotByNumber(room, 3-room.VRobotNumber)
		} else {
			RemoveRobotByNumber(room, 3)
		}

		user2, ret := GetUserByUid(uid)
		if ret == common.ERR_SUCCESS {
			user2.LeaveRoom()
		}
	}

}

func GetChatBaseInfo(rid string, index int) []*ChatUser {
	defer common.PrintPanicStack()
	chat := GetChatRoom(rid)
	if chat != nil {
		uArray := make([]*ChatUser, 0)
		mutex_chat_guardv2.RLock()
		defer mutex_chat_guardv2.RUnlock()
		uArray, ret := chat.GetAudienceList(index)
		if ret == common.ERR_SUCCESS {
			return uArray
		}
	}
	return nil
}

func SendMsgToUser(uid int, msg interface{}) bool {
	user := GetUserSessByUid(uid)
	if user != nil {
		if b, err := json.Marshal(msg); err == nil {
			GetChat().SendToSelf(b, user.Sess)
			return true
		}
	}
	return false
}

func SendMsgToRoom(rid string, msg interface{}) {
	path := fmt.Sprintf("/channel/%s/ws", rid)
	chat.BroadcastFilterByJson(msg, func(q *melody.Session) bool {
		return q.Request.URL.Path == path
	})
}
