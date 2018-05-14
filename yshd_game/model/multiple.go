package model

import (
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"github.com/yshd_game/timer"
	//"strconv"
	"fmt"
	"sync"
	"time"
	//"os/user"
	"github.com/liudng/godump"
)

type MultipleRoomList struct {
	/*
		Id         int64
		OwnerId    int    `xorm:"not null "`
		RoomName   string `xorm:"varchar(255) UNIQUE(ROOM_NAME) "  `
		CreateTime int64
		Count      int
		Status     int //会议室状态
		CloseTime  int64
		Cover      string
		LockTime   int64
	*/

	RoomId          string    `xorm:"varchar(128)  pk not null"` //房间ID
	RoomName        string    `xorm:"varchar(255)"  `            //房间名字
	OwnerId         int       `xorm:"not null "`                 //主播ID
	CreateTime      time.Time //创建时间
	FinishTime      time.Time //结束时间
	Location        string    `xorm:"varchar(128)"  ` //定位
	Cover           string    `xorm:"varchar(255)"  ` //封面图片
	Statue          int       //房间状态
	LiveUrl         string    `xorm:"varchar(255)"  ` //直播流
	MobileUrl       string    `xorm:"varchar(255)"  ` //移动流
	Rice            int       //收到的米粒
	Count           int       //人数
	Weight          int       `xorm:"not null default(0)"` //排序权重
	LockTime        int64
	MutipleRecordId int64
	Moon            int
}

var MultipleRoom chan *MultipleRoomList

//var MutipleTokenMgr map[string]*MultipleRoomList
var MutipleLock *sync.Mutex

func InitMultipleRoom() {
	MultipleRoom = make(chan *MultipleRoomList)
	/*
		MutipleTokenMgr = make(map[string]*MultipleRoomList)
		err := orm.Find(&MutipleTokenMgr)

		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
	*/
	MutipleLock = new(sync.Mutex)
}

func AddMultipleRoomList(owner int, roomid string) int {
	m := &MultipleRoomList{}
	m.CreateTime = time.Now()
	m.OwnerId = owner
	m.RoomId = roomid
	m.Statue = common.MULTIPLE_ROOM_PRE
	_, err := orm.Insert(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS
}

func GetFreeMultipleRoom(uid int) (int, *MultipleRoomList) {
	m := &MultipleRoomList{}
	has, err := orm.Where("owner_id=?", uid).Limit(1, 0).Get(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN, m
	}

	if has {
		ComfirmMultipleRoomStatue(m.RoomId)
		return common.ERR_MULTIPLE_HAS_MIC, m
	}

	MutipleLock.Lock()
	defer MutipleLock.Unlock()
	m = &MultipleRoomList{}
	has, err = orm.Where("statue =? or statue=? ", common.MULTIPLE_ROOM_PRE, common.MULTIPLE_ROOM_FIN).OrderBy("create_time").Limit(1, 0).Get(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN, m
	}

	if has {
		m.OwnerId = uid
		m.LockTime = time.Now().Unix() + 60
		m.Statue = common.MULTIPLE_ROOM_LOCK
		_, err := orm.Where("room_id=?", m.RoomId).Update(m)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN, m
		}
		//MultipleRoom <- m

		go func(rid string) {
			d := timer.NewDispatcher(1)
			d.AfterFunc(5*time.Second, func() {
				godump.Dump(rid)
				CheckRoomStatue(rid)
			})

			(<-d.ChanTimer).Cb()
		}(m.RoomId)

		return common.ERR_SUCCESS, m
	}
	return common.ERR_MULTIPLE_ROOM_BUSY, m
}

func GetMultipleRoomByRid(rid string) (bool, *MultipleRoomList) {
	m := &MultipleRoomList{}
	has, err := orm.Where("room_id=?", rid).Get(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return false, m
	}
	return has, m
}

func GetMultipleRoomByUid(uid int) (bool, *MultipleRoomList) {
	m := &MultipleRoomList{}
	has, err := orm.Where("uid=?", uid).Get(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return false, m
	}
	return has, m
}

func CheckRoomStatue(rid string) int {
	has, n := GetMultipleRoomByRid(rid)
	if has == false {
		return common.ERR_UNKNOWN
	}

	if n.Statue == common.MULTIPLE_ROOM_LOCK {
		n.Statue = common.MULTIPLE_ROOM_FIN
		_, err := orm.Where("room_id=?", rid).Update(n)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
	}
	return common.ERR_SUCCESS
}

func ComfirmMultipleRoomStatue(rid string) int {
	has, n := GetMultipleRoomByRid(rid)
	if has == false {
		return common.ERR_UNKNOWN
	}

	ret, m := ReqGetMutipleRoom(rid)
	//godump.Dump(m)
	if ret == common.ERR_SUCCESS {
		if n.Statue != m.RoomStatus {
			//godump.Dump(n)
			n.Statue = m.RoomStatus
			//godump.Dump(m)

			if m.RoomStatus == common.MULTIPLE_ROOM_FIN {
				n.FinishTime = time.Now()
				n.OwnerId = 0
				n.Cover = ""

			}
			//godump.Dump(n)
			_, err := orm.Where("room_id=?", n.RoomId).MustCols("statue", "owner_id").Update(n)
			if err != nil {
				common.Log.Errf("orm err is  %s", err.Error())
				return common.ERR_UNKNOWN
			}

		} else {
			if m.RoomStatus == common.MULTIPLE_ROOM_FIN && n.OwnerId != 0 {
				n.OwnerId = 0
				n.Cover = ""
				_, err := orm.Where("room_id=?", n.RoomId).MustCols("owner_id").Update(n)
				if err != nil {
					common.Log.Errf("orm err is  %s", err.Error())
					return common.ERR_UNKNOWN
				}
			}
		}
	} else if ret == common.ERR_MULTIPLE_ROOM_FIN {
		_, err := orm.Where("room_id=?", n.RoomId).Delete(n)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
	}
	return ret
}

func CheckLockMultiple() {
	InitMultipleRoom()
	for v := range MultipleRoom {
		d := timer.NewDispatcher(1)
		d.AfterFunc(5*time.Second, func() {
			//ComfirmMultipleRoomStatue(v.RoomId)
			CheckRoomStatue(v.RoomId)
			/*
				ret, m := ReqGetMutipleRoom(v.RoomName)
				if ret == common.ERR_SUCCESS {
					if v.Status != m.RoomStatus {
						//godump.Dump(v)
						if m.RoomStatus == 1 {
							ComfirmMultipleRoomStatue(v.RoomName, common.MULTIPLE_ROOM_BUSY)
						} else if m.RoomStatus == 2 {
							ComfirmMultipleRoomStatue(v.RoomName, common.MULTIPLE_ROOM_FIN)
						} else if m.RoomStatus == 0 {
							ComfirmMultipleRoomStatue(v.RoomName, common.MULTIPLE_ROOM_PRE)
						}
					}
				} else if ret == common.ERR_MULTIPLE_ROOM_FIN {
					ComfirmMultipleRoomStatue(v.RoomName, common.MULTIPLE_ROOM_FIN)
				}
			*/
		})
		(<-d.ChanTimer).Cb()
	}
}

func RefreshMultiple() {
	m := make([]MultipleRoomList, 0)
	err := orm.Find(&m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return
	}
	for _, v := range m {

		ret := ComfirmMultipleRoomStatue(v.RoomId)
		if ret == common.ERR_SUCCESS {
			c := NewChatRoomInfo()
			o := c.GetChatInfo()
			o.Rid = v.RoomId
			o.Uid = v.OwnerId
			c.Statue = common.ROOM_ONLIVE
			c.IsMultiple = true
			AddChatRoom(c)
		}
		//MultipleRoom <- &v
	}
}

func CheckAndCloseMultiple(uid int, rid string) int {
	MutipleLock.Lock()
	defer MutipleLock.Unlock()
	m := &MultipleRoomList{}

	chat := GetChatRoom(rid)
	if chat == nil {
		return common.ERR_NOT_CHAT_EXIST
	}

	has, err := orm.Where("room_id=?", rid).Get(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has == false {
		return common.ERR_PARAM
	}

	//if chat.Save == 1 {
	//	godump.Dump(SaveMutipleM3u8File(rid, uid, m.MutipleRecordId))
	//}

	if m.Statue != common.MULTIPLE_ROOM_FIN {
		m.Statue = common.MULTIPLE_ROOM_FIN
		m.OwnerId = 0
		m.FinishTime = time.Now()
		m.Rice = chat.GetRice()
		m.Count = chat.GetCount()
		m.Moon = chat.GetMoon()
		_, err := orm.Where("room_id=?", rid).MustCols("owner_id").Update(m)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		CloseMutipleRecord(m.MutipleRecordId, m.Rice, m.Count)
	}

	return common.ERR_SUCCESS
	//return ComfirmMultipleRoomStatue(rid)
}

func CloseMultiple(uid int, rid string) int {
	AnchorMgr.SetCloseRoom(uid)

	return ComfirmMultipleRoomStatue(rid)
}

func ReadyMutipleChat(uid int, rname, location, rid string, save, GameType int) int {
	//return 0
	m := &MultipleRoomList{}
	has, err := orm.Where("room_id=?", rid).Get(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	if has {
		room, has := chat_room_manager[rid]
		if has {
			room.room.Uid = uid
			room.room.Rid = rid
			room.Statue = common.ROOM_ONLIVE
			room.Save = save
			room.GameType = GameType
			user, ret := GetUserByUid(uid)
			if ret == common.ERR_SUCCESS {
				m.Cover = user.Image
				m.Location = location
				m.RoomName = rname
				m.CreateTime = time.Now()
				m.OwnerId = uid
				m.MutipleRecordId = AddMutipleRecord(rid, rname, location, user.Image, m.LiveUrl, m.MobileUrl, uid)
				_, err = orm.Where("room_id=?", rid).Update(m)
				if err != nil {
					common.Log.Errf("orm err is %s", err.Error())
					return common.ERR_UNKNOWN
				}

				if GameType == common.GAME_TYPE_NIUNIU { // 如果客户端是要牛牛添加相应的结构并放到索引里去吧
					r := NewRoomInfoNiuNiu()
					r.Rid = rid
					r.Uid = uid
					AddRoomInfoNiuNiu(r)
				} else if GameType == common.GAME_TYPE_TEXAS { // 如果客户端是要德州扑克添加相应的结构并放到索引里去吧
					r := NewRoomInfoTexas()
					r.Rid = rid
					r.Uid = uid
					AddRoomInfoTexas(r)
				} else if GameType == common.GAME_TYPE_GOLDEN_FLOWER { // 如果客户端是要砸金花添加相应的结构并放到索引里去吧
					r := NewRoomInfoGoldenFlower()
					r.Rid = rid
					r.Uid = uid
					AddRoomInfoGoldenFlower(r)
				}

				return common.ERR_SUCCESS
			}
		}
	}
	return common.ERR_MULTIPLE_RID
}

func ComfirmAddMutipleChat(uid int, rname, location, rid string) int {

	//godump.Dump("begin set room 3 rid true")
	common.Log.Debugf("begin set room 3 rid=? ,uid=?", rid, uid)
	m := &MultipleRoomList{}
	has, err := orm.Where("room_id=?", rid).Get(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has {
		room, has := chat_room_manager[rid]
		if has {
			//godump.Dump("begin set room rid true")
			common.Log.Debugf("begin set room 1 rid=? ,uid=?", rid, uid)
			room.room.Uid = uid
			room.room.Rid = rid

			room.Statue = common.ROOM_ONLIVE
		} else {
			//godump.Dump("begin set room rid false")
			common.Log.Debugf("begin set room 2 rid=? ,uid=?", rid, uid)
			c := NewChatRoomInfo()
			o := c.GetChatInfo()
			o.Rid = rid
			o.Uid = uid
			c.Statue = common.ROOM_ONLIVE
			c.IsMultiple = true
			AddChatRoom(c)
		}

		//ret := ComfirmMultipleRoomStatue(rid)

		ret, s := ReqGetMutipleRoom(rid)

		common.Log.Debugf("comfirm room statue  set room 2 rid=? ,uid=?,ret=?,status=?", rid, uid, ret, s.RoomStatus)
		if ret == common.ERR_SUCCESS {
			if s.RoomStatus == common.MULTIPLE_ROOM_BUSY {
				user, ret := GetUserByUid(uid)
				if ret == common.ERR_SUCCESS {
					m.Cover = user.Image
					m.Location = location
					m.RoomName = rname
					m.CreateTime = time.Now()

				}
				_, err = orm.Where("room_id=?", rid).Update(m)
				if err != nil {
					common.Log.Errf("orm err is %s", err.Error())
					return common.ERR_UNKNOWN
				}
				return common.ERR_SUCCESS
			} else {
				return common.ERR_MULTIPLE_STATUE
			}
		} else {
			return common.ERR_UNKNOWN
		}

		if ret != common.ERR_SUCCESS {
			return ret
		}

	}
	return common.ERR_MULTIPLE_RID
}

func GenMultiplePull(uid, line int, rid string) (int, string) {
	//res:=common.RandnomRange64(1,100)
	pullId := fmt.Sprintf("%d_%d_%s", uid, time.Now().Unix(), rid)
	push, pull := GenPullAddr(uid, line, pullId)

	//godump.Dump(pull)
	has, m := GetMultipleRoomByRid(rid)

	s := fmt.Sprintf("%s_%d", rid, time.Now().Unix())
	if has {
		m.LiveUrl = pull
		mobile := ""
		mobilepath := common.Cfg.MustValue("video", "mobile_addr1")
		if line == 1 {
			mobile = fmt.Sprintf("http://%s/%s.m3u8", mobilepath, s)
		} else if line == 2 {
			mobile = fmt.Sprintf("http://%s/%s/playlist.m3u8", mobilepath, s)
		}

		m.MobileUrl = mobile
		//return 0,pull
		_, err := orm.Where("room_id=?", rid).Update(m)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN, ""
		}
		return common.ERR_SUCCESS, push
	}
	return common.ERR_PARAM, ""
}

func CloseMultipleRoom(uid int) {
	ok, room := GetMultipleRoomByUid(uid)
	if ok {
		room.Statue = common.ROOM_FINISH
		room.OwnerId = 0
		_, err := orm.Where("room_id=?", room.RoomId).Update(room)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return
		}
	}
}
