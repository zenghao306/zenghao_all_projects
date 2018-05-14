package model

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/yshd_game/common"
	"strconv"
	"time"
)

type Focus struct {
	Id         int64
	User1      int `xorm:"int(11) not null UNIQUE(FOCUSE_USER)"` //用户ID1 ID1>ID2
	User2      int `xorm:"int(11) not null UNIQUE(FOCUSE_USER)"` //用户ID2
	OneFocus   int `xorm:"int(11) default(0)"`                   //user1关注user2
	TwoFocus   int `xorm:"int(11) default(0)"`                   //user2关注user1
	Push       bool
	FocusTime1 time.Time
	FocusTime2 time.Time
}

type OutUserInfo struct {
	Uid         int    `json:"owner_id"`
	Image       string `json:"image"`
	Signature   string `json:"signature"`
	Sex         int    `json:"sex"`
	Userlevel   int    `json:"user_level"`
	AnchorLevel int    `json:"anchor_level"`
	NickName    string `json:"nick_name"`
	Location    string `json:"location"`
	RoomId      string `json:"room_id"`
	Cover       string `json:"cover"`
	LiveUrl     string `json:"live_url"`
	Viewer      int    `json:"viewer"`
	GameType    int    `json:"game_type"`
	RoomName    string `json:"room_name"`
	Statue      int    `json:"statue"`
	FlvUrl      string `json:"flv_url"`
}

type OutUserInfo2 struct {
	Uid         int    `json:"uid"`
	Image       string `json:"image"`
	Signature   string `json:"signature"`
	Sex         int    `json:"sex"`
	Userlevel   int    `json:"user_level"`
	AnchorLevel int    `json:"anchor_level"`
	NickName    string `json:"nick_name"`
	Location    string `json:"location"`
	RoomId      string `json:"room_id"`
	Cover       string `json:"cover"`
	LiveUrl     string `json:"live"`
	Viewer      int    `json:"viewer"`
	GameType    int    `json:"game_type"`
	RoomName    string `json:"room_name"`
	Statue      int    `json:"statue"`
	Score       int64  `json:"score"`
	FlvUrl      string `json:"flv_url"`
}

func (self *Focus) Update() bool {
	_, err := orm.Id(self.Id).Update(self)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return false
	}
	return true
}
func (self *Focus) UpdateFront(filed string) bool {
	_, err := orm.Id(self.Id).MustCols(filed).Update(self)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return false
	}
	return true
}

func FocusOther(selfid, otherid int) int {
	var user1 int
	var user2 int
	one_focus := 0
	two_focus := 0

	if selfid == otherid {
		return common.ERR_FOCUS_TO_SELF
	}
	if selfid > otherid {
		user1 = selfid
		user2 = otherid
		one_focus = 1
	} else {
		user1 = otherid
		user2 = selfid
		two_focus = 1
	}
	focus := &Focus{}
	has, err := orm.Where("user1=? and user2 =?", user1, user2).Get(focus)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}

	if one_focus == common.FOCUS_STATUE_MAIN {
		if focus.OneFocus == 1 {
			return common.ERR_FOCUS_ALREADY
		}

		focus.OneFocus = one_focus
		focus.FocusTime1 = time.Now()
	} else if two_focus == common.FOCUS_STATUE_MAIN {
		if focus.TwoFocus == 1 {
			return common.ERR_FOCUS_ALREADY
		}

		focus.TwoFocus = two_focus
		focus.FocusTime2 = time.Now()
	}
	if has {
		focus.Update()
	} else {
		focus.User1 = user1
		focus.User2 = user2
		_, err := orm.InsertOne(focus)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
	}
	return common.ERR_SUCCESS
}

func CancleFocusByOid(selfid, otherid int) int {
	var user1 int
	var user2 int
	one_focus := 1
	two_focus := 1

	if selfid > otherid {
		user1 = selfid
		user2 = otherid
		one_focus = 0
	} else {
		user1 = otherid
		user2 = selfid
		two_focus = 0
	}

	focus := &Focus{}
	has, err := orm.Where("user1=? and user2 =?", user1, user2).Get(focus)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}

	if has {
		if one_focus == common.FOCUS_STATUE_NONE {
			focus.OneFocus = one_focus
			if focus.UpdateFront("one_focus") {
				return common.ERR_SUCCESS
			}
		} else if two_focus == common.FOCUS_STATUE_NONE {
			focus.TwoFocus = two_focus
			if focus.UpdateFront("two_focus") {
				return common.ERR_SUCCESS
			}
		}

	} else {
		return common.ERR_FOCUS_EXIST
	}
	return common.ERR_UNKNOWN
}

/*
func CancleFocus(id, uid int) int {
	focus := &Focus{}
	has, err := orm.Id(id).Get(focus)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has == false {
		return common.ERR_FOCUS_EXIST
	}
	if uid > focus.User2 {
		if focus.TwoFocus == common.FOCUS_STATUE_NONE {
			_, err = orm.Id(id).Delete(focus)
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())
				return common.ERR_UNKNOWN
			}
		} else {
			focus.OneFocus = common.FOCUS_STATUE_NONE
			focus.Update()
		}
	} else {
		if focus.OneFocus == common.FOCUS_STATUE_NONE {
			_, err = orm.Id(id).Delete(focus)
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())
				return common.ERR_UNKNOWN
			}
		} else {
			focus.TwoFocus = common.FOCUS_STATUE_NONE
			focus.Update()
		}
	}
	return common.ERR_SUCCESS
}
*/
func GetFocusCount(uid int) (focus_ int, fans_ int) {
	var focus string
	var fans string
	var err error
	sql := fmt.Sprintf("select count(*) as focus from  (select * from go_focus where user1=%d and one_focus=1 union all select * from go_focus where user2=%d and two_focus=1) as b ", uid, uid)

	//sql := fmt.Sprintf("select  (select count(*) from focus where otherid=%d) as fans,  (select count(*) from focus where uid=%d) as focus", uid, uid)
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return 0, 0
	}
	if len(rowArray) != 0 {
		focus = common.BytesToString(rowArray[0]["focus"])
		focus_, _ = strconv.Atoi(focus)
	}

	sql = fmt.Sprintf("select count(*) as fans from  (select * from go_focus where user1=%d and two_focus=1 union all select * from go_focus where user2=%d and one_focus=1) as b", uid, uid)
	rowArray, err = orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return 0, 0
	}
	if len(rowArray) != 0 {
		fans = common.BytesToString(rowArray[0]["fans"])
		fans_, _ = strconv.Atoi(fans)
	}

	return focus_, fans_
}

type UserFocusCache struct {
	Uid    int
	Statue int
	RoomId string
}

func GetFocusListWithCache(uid, index int) (out_users []OutUserInfo2, ret int) {
	//res, err := orm.Query("select owner_id as uid,room_id,statue from go_room_list where  owner_id in (  select user2  as uid from go_focus where user1=? and one_focus=1 union all select user1  as uid from go_focus where user2=? and two_focus=1) ", uid, uid)
	res, err := orm.Query("select user2  as uid from go_focus where user1=? and one_focus=1 union all select user1  as uid from go_focus where user2=? and two_focus=1", uid, uid)
	//res, err := orm.Query("SELECT id,uid FROM  (SELECT id,user2  AS uid ,focus_time1 AS focus_time FROM go_focus WHERE user1=? AND one_focus=1 UNION ALL SELECT id,user1  AS uid,focus_time2 AS focus_time FROM go_focus WHERE user2=? AND two_focus=1) a ORDER BY a.focus_time DESC", uid, uid)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return
	}

	all_user_v2 := make(map[int]*UserFocusCache, 0)
	all_user := make([]int, 0)

	for _, row := range res {
		uid := row["uid"]
		uid_ := common.BytesToInt(uid)

		all_user = append(all_user, uid_)

	}

	res, err = orm.Query("SELECT a.uid ,  b.room_id FROM (SELECT user2  AS uid FROM go_focus WHERE user1=? AND one_focus=1 UNION ALL SELECT user1  AS uid FROM go_focus WHERE user2=? AND two_focus=1) a LEFT JOIN go_room_list b ON a.uid=b.owner_id WHERE b.statue=1", uid, uid)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return
	}
	for _, row := range res {
		s := &UserFocusCache{}
		uid := row["uid"]
		uid_ := common.BytesToInt(uid)
		s.Uid = uid_
		s.Statue = common.USER_STATUE_LIVE
		rid := row["room_id"]
		s.RoomId = common.BytesToString(rid)
		all_user_v2[s.Uid] = s
	}

	users := make([]User, 0)
	if len(all_user) == 0 {
		return
	}
	err = orm.In("uid", all_user).Limit(common.FOCUS_LIST_PAGE_COUNT, index*common.FOCUS_LIST_PAGE_COUNT).Find(&users)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return
	}

	for _, v := range users {
		var out OutUserInfo2
		out.Uid = v.Uid
		out.AnchorLevel = v.AnchorLevel
		out.Cover = v.Image
		out.Image = v.Image
		out.Location = v.Location
		out.NickName = v.NickName
		out.Signature = v.Signature
		out.Score = v.Score
		out.Userlevel = v.UserLevel
		out.AnchorLevel = v.AnchorLevel

		c, ok := all_user_v2[v.Uid]
		if !ok {
			out.Statue = common.USER_STATUE_LEAVE
			out_users = append(out_users, out)
			continue
		} else {
			out.RoomId = c.RoomId
			room := &RoomList{}
			has, err := orm.Where("room_id=?", c.RoomId).Get(room)
			if err != nil {
				common.Log.Errf("db err %s", err.Error())
				return
			}

			if has {
				out.LiveUrl = room.LiveUrl
				out.GameType = room.GameType
				out.RoomName = room.RoomName
				out.Statue = room.Statue
				out.FlvUrl = room.FlvUrl
			} else {
				continue
			}

			out_users = append(out_users, out)
		}

	}
	return
}

// 从GetFocusListWithCache拷贝而来，区别是屏蔽掉GameType==3的（后期本函数可能会去掉）
func GetFocusListWithCache2(uid, index int) (out_users []OutUserInfo2, ret int) {
	//res, err := orm.Query("select id,user2  as uid from go_focus where user1=? and one_focus=1 union all select id,user1  as uid from go_focus where user2=? and two_focus=1 ", uid, uid)

	res, err := orm.Query("SELECT id,uid FROM  (SELECT id,user2  AS uid ,focus_time1 AS focus_time FROM go_focus WHERE user1=? AND one_focus=1 UNION ALL SELECT id,user1  AS uid,focus_time2 AS focus_time FROM go_focus WHERE user2=? AND two_focus=1) a ORDER BY a.focus_time DESC", uid, uid)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return
	}

	all_user := make([]int, 0)
	live_rooms := make(map[int]*CacheUser, 0)
	for _, row := range res {
		uid := row["uid"]
		uid_ := common.BytesToInt(uid)
		all_user = append(all_user, uid_)
		u, err := GetCacheUser(uid_)
		if err == redis.Nil {
			continue
		} else if err != nil {
			return
		} else {
			if u.Status == common.USER_STATUE_LIVE {
				live_rooms[uid_] = u
			}
		}
	}

	users := make([]User, 0)
	if len(all_user) == 0 {
		return
	}
	err = orm.In("uid", all_user).Limit(common.FOCUS_LIST_PAGE_COUNT, index*common.FOCUS_LIST_PAGE_COUNT).Find(&users)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return
	}

	for _, v := range users {
		var out OutUserInfo2
		out.Uid = v.Uid
		out.AnchorLevel = v.AnchorLevel
		out.Cover = v.Image
		out.Image = v.Image
		out.Location = v.Location
		out.NickName = v.NickName
		out.Signature = v.Signature
		out.Score = v.Score
		out.Userlevel = v.UserLevel
		out.AnchorLevel = v.AnchorLevel
		c, ok := live_rooms[v.Uid]
		if !ok {
			out.Statue = common.USER_STATUE_LEAVE
			out_users = append(out_users, out)
			continue
		} else {
			out.RoomId = c.RoomId
			room := &RoomList{}
			has, err := orm.Where("room_id=?", c.RoomId).Get(room)
			if err != nil {
				common.Log.Errf("db err %s", err.Error())
				return
			}
			if has {
				out.LiveUrl = room.LiveUrl
				out.GameType = room.GameType
				out.RoomName = room.RoomName
				out.Statue = room.Statue
				out.FlvUrl = room.FlvUrl
			} else {
				continue
			}

			if room.GameType != common.GAME_TYPE_GOLDEN_FLOWER {
				out_users = append(out_users, out)
			}
		}
	}
	return
}

/*
func GetFocusList(uid, index int) ([]map[string]string, int) {
	sql := fmt.Sprintf("select a.id,a.uid,b.image,b.signature,b.sex,b.user_level,b.anchor_level,b.nick_name,b.statue,b.live,b.room_id,b.location,b.cover,b.score from  (select id,user2  as uid from go_focus where user1=%d and one_focus=1 union all select id,user1  as uid from go_focus where user2=%d and two_focus=1)  a left join go_user b on a.uid=b.uid limit %d,%d", uid, uid, index*common.FOCUS_LIST_PAGE_COUNT, common.FOCUS_LIST_PAGE_COUNT)
	//sql := fmt.Sprintf("select a.id,a.uid,b.image,b.signature,b.sex,b.user_level,b.anchor_level,b.nick_name,b.statue,b.live,b.room_id,b.location,b.cover,b.score,c.game_type from  (select id,user2  as uid from go_focus where user1=%d and one_focus=1 union all select id,user1  as uid from go_focus where user2=%d and two_focus=1)  a left join go_user b on a.uid=b.uid left join go_room_list c on b.room_id=c.room_id limit %d,%d", uid, uid, index*common.FOCUS_LIST_PAGE_COUNT, common.FOCUS_LIST_PAGE_COUNT)

	rowArray, err := orm.Query(sql)
	retMap := make([]map[string]string, 0)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return retMap, 0
	}
	//path := common.Cfg.MustValue("video", "video_addr")
	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		uid := ss["uid"]
		uid_, _ := strconv.Atoi(uid)
		room_id := ss["room_id"]
		if ss["statue"] == "1" && room_id != "" {
			sess := GetUserSessByUid(uid_)
			if sess == nil {
				ss["statue"] = "0"
			} else {
				room := &RoomList{}
				has, err := orm.Where("room_id=?", room_id).Get(room)
				if err != nil {
					common.Log.Errf("db err %s", err.Error())
					return retMap, common.ERR_UNKNOWN
				}
				if has {
					ss["live"] = room.LiveUrl
					ss["game_type"] = strconv.Itoa(room.GameType)
				}
			}
		} else {
			ss["statue"] = "0"
			ss["room_id"] = ""
		}
		retMap = append(retMap, ss)
	}

	return retMap, common.ERR_SUCCESS
}

*/

//获取关注的人在直播列表
func GetLiveFocusWithCache(uid, index int) (out_users []OutUserInfo, ret int) {
	live_users := make([]int, 0)
	//map key保存临时uid，value保存临时房间id
	rooms := make(map[int]string, 0)
	res, err := orm.Query("select owner_id ,room_id from go_room_list where  statue=1 and owner_id in (  select user2  as uid from go_focus where user1=? and one_focus=1 union all select user1  as uid from go_focus where user2=? and two_focus=1) limit ?,?", uid, uid, index*common.FOCUS_LIST_PAGE_COUNT, common.FOCUS_LIST_PAGE_COUNT)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}

	for _, v := range res {
		b, ok := v["owner_id"]
		if ok {
			uid := common.BytesToInt(b)
			live_users = append(live_users, uid)

			b, ok = v["room_id"]

			if ok {
				room_id := common.BytesToString(b)
				rooms[uid] = room_id
			}
		}
	}

	if len(live_users) == 0 {
		ret = common.ERR_SUCCESS
		return
	}
	users := make([]User, 0)
	err = orm.In("uid", live_users).Limit(common.FOCUS_LIST_PAGE_COUNT, index*common.FOCUS_LIST_PAGE_COUNT).Find(&users)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}

	//out_users=make([]OutUserInfo,0)

	for _, v := range users {

		var out OutUserInfo
		out.Uid = v.Uid
		out.AnchorLevel = v.AnchorLevel
		out.Cover = v.Image
		out.Image = v.Image
		//out.Location = v.Location
		out.NickName = v.NickName
		out.Signature = v.Signature
		c, ok := rooms[v.Uid]
		if !ok {
			continue
		}

		out.RoomId = c
		room := &RoomList{}
		has, err := orm.Where("room_id=?", c).Get(room)
		if err != nil {
			common.Log.Errf("db err %s", err.Error())
			ret = common.ERR_UNKNOWN
			return
		}
		if has {
			out.LiveUrl = room.LiveUrl
			out.GameType = room.GameType
			out.RoomName = room.RoomName
			out.Statue = room.Statue
			out.FlvUrl = room.FlvUrl
			out.Location = room.Location
		} else {
			continue
		}

		chat := GetChatRoom(c)
		if chat != nil {
			out.Viewer = chat.GetCount() + chat.VRobotNumber
		} else {
			continue
		}

		out_users = append(out_users, out)
	}
	ret = common.ERR_SUCCESS
	return
}

//屏蔽掉游戏3[砸金花]的函数【后续客户版本都升级后应该会停掉】
func GetLiveFocusWithCache2(uid, index int) (out_users []OutUserInfo, ret int) {
	live_users := make([]int, 0)

	//map key保存临时uid，value保存临时房间id
	rooms := make(map[int]string, 0)
	res, err := orm.Query("select owner_id ,room_id from go_room_list where  statue=1 and owner_id in (  select user2  as uid from go_focus where user1=? and one_focus=1 union all select user1  as uid from go_focus where user2=? and two_focus=1) limit ?,?", uid, uid, index*common.FOCUS_LIST_PAGE_COUNT, common.FOCUS_LIST_PAGE_COUNT)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}

	for _, v := range res {
		b, ok := v["owner_id"]
		if ok {
			uid := common.BytesToInt(b)
			live_users = append(live_users, uid)

			b, ok = v["room_id"]

			if ok {
				room_id := common.BytesToString(b)
				rooms[uid] = room_id
			}
		}
	}

	if len(live_users) == 0 {
		ret = common.ERR_SUCCESS
		return
	}
	users := make([]User, 0)
	err = orm.In("uid", live_users).Limit(common.FOCUS_LIST_PAGE_COUNT, index*common.FOCUS_LIST_PAGE_COUNT).Find(&users)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}

	//out_users=make([]OutUserInfo,0)

	for _, v := range users {

		var out OutUserInfo
		out.Uid = v.Uid
		out.AnchorLevel = v.AnchorLevel
		out.Cover = v.Image
		out.Image = v.Image
		//out.Location = v.Location
		out.NickName = v.NickName
		out.Signature = v.Signature
		c, ok := rooms[v.Uid]
		if !ok {
			continue
		}

		out.RoomId = c
		room := &RoomList{}
		has, err := orm.Where("room_id=?", c).Get(room)
		if err != nil {
			common.Log.Errf("db err %s", err.Error())
			ret = common.ERR_UNKNOWN
			return
		}
		if has {
			out.LiveUrl = room.LiveUrl
			out.GameType = room.GameType
			out.RoomName = room.RoomName
			out.Statue = room.Statue
			out.FlvUrl = room.FlvUrl
			out.Location = room.Location
		} else {
			continue
		}

		chat := GetChatRoom(c)
		if chat != nil {
			out.Viewer = chat.GetCount() + chat.VRobotNumber
		} else {
			continue
		}

		if out.GameType != common.GAME_TYPE_GOLDEN_FLOWER {
			out_users = append(out_users, out)
		}
	}
	ret = common.ERR_SUCCESS
	return
}

/*
func GetLiveFocusList(uid, index int) ([]map[string]string, int) {
	sql := fmt.Sprintf("select a.id,a.uid,b.image,b.signature,b.sex,b.user_level,b.anchor_level,b.nick_name,b.statue,b.live,b.room_id,b.location,b.cover  from  (select id,user2  as uid from go_focus where user1=%d and one_focus=1 union all select id,user1  as uid from go_focus where user2=%d and two_focus=1)  a left join go_user b on a.uid=b.uid where statue=1 limit %d,%d", uid, uid, index*common.FOCUS_LIST_PAGE_COUNT, common.FOCUS_LIST_PAGE_COUNT)
	rowArray, err := orm.Query(sql)
	retMap := make([]map[string]string, 0)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return retMap, common.ERR_UNKNOWN
	}

	count := 0
	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		oid := ss["uid"]
		oid_, _ := strconv.Atoi(oid)
		//room_id := ss["room_id"]
		//	if ss["statue"] == "1" && room_id != "" {
		if ss["statue"] == "1" {
			sess := GetUserSessByUid(oid_)
			if sess == nil {
				ss["statue"] = "0"
				user, hasuser := GetUserByUid(oid_)
				if hasuser {
					//user.SetStatue(common.USER_STATUE_LEAVE)
					user.LeaveRoom()
					count++
				}
				break
			} else {
				room := &RoomList{}
				has, err := orm.Where("room_id=?", sess.Roomid).Get(room)
				if err != nil {
					common.Log.Errf("db err %s", err.Error())
					return retMap, common.ERR_UNKNOWN
				}
				if has {
					ss["live"] = room.LiveUrl
					ss["game_type"] = strconv.Itoa(room.GameType)
				}
			}
		} else {
			ss["statue"] = "0"
			ss["room_id"] = ""
		}
		retMap = append(retMap, ss)
	}


	if count > 0 {
		return GetLiveFocusList(uid, index)
	}
	retRoomMap := make([]map[string]string, 0)

	for _, value := range retMap {
		bUserMap := make(map[string]string)
		rid := value["room_id"]
		room, _ := GetRoomById(rid)

		bUserMap["room_id"] = rid
		bUserMap["room_name"] = room.RoomName
		bUserMap["owner_id"] = value["uid"]
		bUserMap["location"] = room.Location
		bUserMap["cover"] = room.Cover
		bUserMap["live_url"] = room.LiveUrl
		//user, has := GetUserByUidStr(bid)
		bUserMap["nick_name"] = value["nick_name"]
		bUserMap["sex"] = value["sex"]
		bUserMap["statue"] = "1"
		bUserMap["image"] = value["image"]
		bUserMap["game_type"] = value["game_type"]
		chat := GetChatRoom(rid)
		if chat != nil {
			bUserMap["viewer"] = strconv.Itoa(chat.GetCount())
		} else {
			continue
		}

		retRoomMap = append(retRoomMap, bUserMap)
	}
	return retRoomMap, common.ERR_SUCCESS
}
*/
/*
type ClientRoomInfo struct {
	Rid      string
	RoomName string
	OwnerId  int
	Location string
	Cover    string
	LiveUrl  string
	NickName string
	Sex      int
	Statue   int
	Image    string
}
*/

func GetFansList(uid, index int) ([]map[string]string, int) {
	sql := fmt.Sprintf("select a.id,a.uid,b.image,b.signature,b.sex,b.user_level,b.anchor_level,b.nick_name,a.relation from  (select id, user1 as uid,two_focus  as relation ,focus_time2 as focus_time from go_focus where user2=%d and one_focus=1 union all select id,user2  as uid,one_focus as relation,focus_time1 as focus_time from go_focus where user1=%d and two_focus=1)  a left join go_user b on a.uid=b.uid  order by a.focus_time  desc limit %d,%d", uid, uid, index*common.FOCUS_LIST_PAGE_COUNT, common.FOCUS_LIST_PAGE_COUNT)

	rowArray, err := orm.Query(sql)
	retMap := make([]map[string]string, 0)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return retMap, 0
	}

	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		retMap = append(retMap, ss)
	}

	return retMap, common.ERR_SUCCESS
}

func GetNotFocusLiveCount(uid int) int {
	sql := fmt.Sprintf("select count(*) count_num from go_room_list LEFT JOIN go_user ON go_room_list.owner_id = go_user.uid  where  statue=1 and owner_id not in ( select user2  as uid from go_focus where user1=%d and one_focus=1 union all select user1  as uid from go_focus where user2=%d and two_focus=1)", uid, uid)

	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return 0
	}

	count_num := common.BytesToString(rowArray[0]["count_num"])
	count_num_, _ := strconv.Atoi(count_num)

	return count_num_
}

func GetLiveNotFocusWithCache(uid, index int) (out_users []OutUserInfo, ret int) {
	live_users := make([]int, 0)

	//map key保存临时uid，value保存临时房间id
	rooms := make(map[int]string, 0)
	res, err := orm.Query("select owner_id ,room_id,go_user.nick_name,go_user.signature from go_room_list LEFT JOIN go_user ON go_room_list.owner_id = go_user.uid  where  statue=1 and owner_id not in ( select user2  as uid from go_focus where user1=? and one_focus=1 union all select user1  as uid from go_focus where user2=? and two_focus=1)limit ?,?", uid, uid, index*common.FOCUS_LIST_PAGE_COUNT, common.FOCUS_LIST_PAGE_COUNT)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}

	for _, v := range res {
		b, ok := v["owner_id"]
		if ok {
			uid := common.BytesToInt(b)
			live_users = append(live_users, uid)

			b, ok = v["room_id"]

			if ok {
				room_id := common.BytesToString(b)
				rooms[uid] = room_id
			}
		}
	}

	if len(live_users) == 0 {
		ret = common.ERR_SUCCESS
		return
	}
	users := make([]User, 0)
	err = orm.In("uid", live_users).Limit(common.FOCUS_LIST_PAGE_COUNT, index*common.FOCUS_LIST_PAGE_COUNT).Find(&users)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}

	//out_users=make([]OutUserInfo,0)
	user, _ := GetUserByUid(uid)
	acountType := user.AccountType
	if acountType == 1 { //测试号
		for _, v := range users {

			var out OutUserInfo
			out.Uid = v.Uid
			out.AnchorLevel = v.AnchorLevel
			out.Cover = v.Image
			out.Image = v.Image
			//out.Location = v.Location
			out.NickName = v.NickName
			out.Signature = v.Signature
			c, ok := rooms[v.Uid]
			if !ok {
				continue
			}

			out.RoomId = c
			room := &RoomList{}
			has, err := orm.Where("room_id=?", c).Get(room)
			if err != nil {
				common.Log.Errf("db err %s", err.Error())
				ret = common.ERR_UNKNOWN
				return
			}
			if has {
				out.LiveUrl = room.LiveUrl
				out.GameType = room.GameType
				out.RoomName = room.RoomName
				out.Statue = room.Statue
				out.FlvUrl = room.FlvUrl
				out.Location = room.Location
			} else {
				continue
			}

			chat := GetChatRoom(c)
			if chat != nil {
				out.Viewer = chat.GetCount() + chat.VRobotNumber
			} else {
				continue
			}

			out_users = append(out_users, out)
		}
		ret = common.ERR_SUCCESS
		return
	} else { //正常号
		for _, v := range users {
			if v.AccountType == 0 {
				var out OutUserInfo
				out.Uid = v.Uid
				out.AnchorLevel = v.AnchorLevel
				out.Cover = v.Image
				out.Image = v.Image
				//out.Location = v.Location
				out.NickName = v.NickName
				out.Signature = v.Signature
				c, ok := rooms[v.Uid]
				if !ok {
					continue
				}

				out.RoomId = c
				room := &RoomList{}
				has, err := orm.Where("room_id=?", c).Get(room)
				if err != nil {
					common.Log.Errf("db err %s", err.Error())
					ret = common.ERR_UNKNOWN
					return
				}
				if has {
					out.LiveUrl = room.LiveUrl
					out.GameType = room.GameType
					out.RoomName = room.RoomName
					out.Statue = room.Statue
					out.FlvUrl = room.FlvUrl
					out.Location = room.Location
				} else {
					continue
				}

				chat := GetChatRoom(c)
				if chat != nil {
					out.Viewer = chat.GetCount() + chat.VRobotNumber
				} else {
					continue
				}

				out_users = append(out_users, out)
			}
		}
		ret = common.ERR_SUCCESS
		return
	}
}
