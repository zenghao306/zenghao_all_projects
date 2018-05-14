package model

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/yshd_game/common"
	"strconv"
	"time"
)

type RoomList struct {
	RoomId     string    `xorm:" varchar(64) pk not null"` //房间ID
	RoomName   string    `xorm:"varchar(128) not null"`    //房间名字
	OwnerId    int       `xorm:"not null "`                //主播ID
	CreateTime time.Time //创建时间
	FinishTime time.Time //结束时间
	Location   string    `xorm:"varchar(64) not null"`  //定位
	Cover      string    `xorm:"varchar(180) not null"` //封面图片
	Statue     int       //房间状态
	LiveUrl    string    `xorm:"varchar(128) not null"` //直播流
	FlvUrl     string    `xorm:"varchar(128) not null"` //直播流
	MobileUrl  string    `xorm:"varchar(128) not null"` //移动流
	Rice       int       `xorm:"int(11)  not null"`     //收到的米粒
	Count      int       `xorm:"int(11)  not null"`     //人数
	Weight     int       `xorm:"not null default(0)"`   //排序权重
	Playback   string    `xorm:"varchar(128) not null"`
	Roomtype   int       `xorm:"not null "`         //0正常 1测试
	GameType   int       `xorm:"int(11)  not null"` //0默认 1牛牛
	Moon       int       `xorm:"int(11)  not null"`
	Device     string    `xorm:"varchar(40) "`
	UserAgent  string    `xorm:"varchar(40) "`
}

type ReRoom struct {
	llist RoomList
	url   string
}

var (
	prefix = "go_"
)

func GetRoomById(rid string) (*RoomList, bool) {
	room := &RoomList{}
	has, err := orm.Where("room_id=?", rid).Get(room)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}
	return room, has
}

func GetRoomList(index int) ([]map[string]string, int) {
	//sql := fmt.Sprintf("select a.room_id,a.room_name,a.owner_id,a.location,a.cover,a.live_url,b.nick_name,b.sex,a.statue,b.image  from room_list a  left join user b on  a.owner_id=b.uid   where a.statue=1  order by weight,create_time  desc limit %d,%d ", index*common.ROOM_LIST_PAGE_COUNT, common.ROOM_LIST_PAGE_COUNT)

	sql := fmt.Sprintf("select a.room_id,a.room_name,a.owner_id,a.location,a.cover,a.live_url,b.nick_name,b.sex,a.statue,b.image  from go_room_list a  left join go_user b on  a.owner_id=b.uid   where a.statue=1  order by weight,create_time  desc limit %d,%d ", index*common.ROOM_LIST_PAGE_COUNT, common.ROOM_LIST_PAGE_COUNT)
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil, common.ERR_UNKNOWN
	}

	retMap := make([]map[string]string, 0)
	flag := true
	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			if colName == "room_id" {
				chat := GetChatRoom(value)
				if chat != nil {
					ss["viewer"] = strconv.Itoa(chat.GetCount() + chat.GetVRobotCount())
				} else {
					flag = false
					break
				}
			}
			//value := common.BytesToString(colValue)
			ss[colName] = value
		}
		if flag {
			retMap = append(retMap, ss)
		}
		flag = true
	}

	return retMap, common.ERR_SUCCESS
}

/*
func CloseRoom(count, value int) interface{} {

	c := redigo.Get()
	defer c.Close()
	_, err := c.Do("lrem", "list.room", count, value)
	if err != nil {
		common.Log.Errf("close  room error %s", err.Error())
		return "failed"
	}
	return "ok"
}
*/
func CreateRoom(uid int, room_name, cover, location string, weight int, rtype, gameType int, device string, user_agent string) (string, int) {
	/*
		reg := regexp.MustCompile(`[0-9]+`)
		roomid := reg.FindAllString(live, -1)
		if len(roomid) == 0 {
			return "", common.ERR_UNKNOWN
		}

		rid := roomid[0]
	*/
	var rid string
	u, err := GetCacheUser(uid)
	if err == redis.Nil {
		return "", common.ERR_OPERATOR_EXPIRE
	} else if err != nil {
		return "", common.ERR_UNKNOWN
	} else {
		rid = u.PreRoomId
		u.PreRoomId = ""
		SetCacheUser(uid, u)
	}

	room_, has := GetRoomById(rid)
	if has == false {
		return "", common.ERR_PRE_LIVE_RID
	}

	if room_.Statue == common.ROOM_READY {
		return room_.RoomId, common.ERR_ROOM_READY
	}
	now_time := time.Now()

	if room_name == "" {
		room_name = fmt.Sprintf("房间%d", time.Now().Unix())
	}
	room_.RoomName = room_name
	room_.OwnerId = uid
	room_.CreateTime = now_time
	room_.Statue = common.ROOM_READY
	room_.Cover = cover
	room_.Location = location
	room_.Weight = weight
	room_.Roomtype = rtype
	room_.GameType = gameType
	room_.Device = device
	room_.UserAgent = user_agent

	aff_row, err := orm.Where("room_id=?", rid).Update(room_)
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return "", common.ERR_UNKNOWN
	}

	if aff_row == 0 {
		return room_.RoomId, common.ERR_DB_UPDATE
	}
	//sql := fmt.Sprintf("insert into go_room_list (room_name,owner_id,create_time,statue,cover,location) values ('%s',%d,'%s',%d,'%s','%s');select LAST_INSERT_ID() as id", room_name, uid, time.Now().Format("2006-01-02 15:04:05"), common.ROOM_ONLIVE, cover, location)

	return room_.RoomId, common.ERR_SUCCESS
}

/*
func CloseRoomStatus(roomid string, uid int) int {
	chat := GetChatRoom(roomid)
	if chat == nil {
		return common.ERR_NOT_CHAT_EXIST
	}
	close_room, has := GetRoomById(roomid)
	if has == true {
		if close_room.Statue == common.ROOM_FINISH {
			return common.ERR_ALREAD_FINSIH
		}
		close_room.FinishTime = time.Now()
		close_room.Statue = common.ROOM_FINISH
		close_room.Rice = chat.GetRice()
		close_room.Count = chat.GetCount()

		close_room.Moon = chat.GetMoon() + GetDump(roomid)
		_, err := orm.Where("room_id=? and owner_id=?", roomid, uid).MustCols("finish_time", "statue", "rice", "count").Update(close_room)
		if err != nil {
			common.Log.Errf("close  room error %s", err.Error())
			return common.ERR_UNKNOWN
		}
		return common.ERR_SUCCESS
	}
	return common.ERR_ROOM_EXIST
}
*/
//func CloseRoom(roomid string, uid int) int {
func CloseRoom(chat *ChatRoomInfo, uid int) int {
	//godump.Dump("CloseRoom() @@ called")
	//chat := GetChatRoom(roomid)
	if chat == nil {
		return common.ERR_NOT_CHAT_EXIST
	}
	roomid := chat.room.Rid
	if chat.Save == 1 {
		chat.Save = 0

		if ret := CheckPlayList(uid); ret == common.ERR_PLAYBACK_MAX {

			m := &PlayBack{}
			_, err := orm.Where("uid=?", uid).Asc("save_time").Limit(1, 0).Get(m)
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())

				return common.ERR_UNKNOWN
			}
			DelPlayBackList(uid, m.RoomId)
		}
		SaveM3u8File(roomid)

	}

	close_room, has := GetRoomById(roomid)
	if has == true {
		if close_room.Statue == common.ROOM_FINISH {
			return common.ERR_ALREAD_FINSIH
		}
		close_room.FinishTime = time.Now()
		close_room.Statue = common.ROOM_FINISH
		close_room.Rice = chat.GetRice()
		close_room.Count = chat.GetCount()

		close_room.Moon = chat.GetMoon() + GetDump(roomid)
		aff_row, err := orm.Where("room_id=? and owner_id=?", roomid, uid).MustCols("finish_time", "statue", "rice", "count").Update(close_room)
		if err != nil {
			common.Log.Errf("close  room error %s", err.Error())
			return common.ERR_UNKNOWN
		}

		if aff_row == 0 {
			return common.ERR_DB_UPDATE
		}
		return common.ERR_SUCCESS
	}
	return common.ERR_ROOM_EXIST
}

func GenPullAddr(uid, line int, rid string) (push string, pull string) {
	var key string
	var pullKey string
	if line == 1 {
		key = "video_addr1"
		pullKey = "rtmp_addr1"
	} else if line == 2 {
		key = "video_addr2"
		pullKey = "rtmp_addr2"
	} else {
		return "", ""
	}

	pushPath1 := common.Cfg.MustValue("video", key)
	sign := GenSercertUrl(rid)
	expire := time.Now().Unix() + int64(ExpireTime)

	pullPath := common.Cfg.MustValue("video", pullKey)

	if line == 1 {
		push = fmt.Sprintf("rtmp://%s/%s?e=%d&token=%s", pushPath1, rid, expire, sign)
		pull = fmt.Sprintf("rtmp://%s/%s", pullPath, rid)

	} else {
		push = fmt.Sprintf("rtmp://%s/%s", pushPath1, rid)
	}

	return
}

func PreCreateLiveUrl(user *User, line int) (string, int) {

	var key, pullkey, mobilekey, mobile, flvkey string
	if line == 1 {
		key = "video_addr1"
		pullkey = "rtmp_addr1"
		mobilekey = "mobile_addr1"
		flvkey = "pull_addr1"
	} else if line == 2 {
		key = "video_addr2"
		pullkey = "rtmp_addr2"
		mobilekey = "mobile_addr2"
		flvkey = "pull_addr1"
	} else {
		return "", common.ERR_PARAM
	}

	rv2 := &RoomList{}
	has, err := orm.Where("owner_id=? and statue=?", user.Uid, common.ROOM_PRE_V2).Get(rv2)
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return "", common.ERR_UNKNOWN
	}
	if has {
		path1 := common.Cfg.MustValue("video", key)
		var push string
		expire := time.Now().Unix() + int64(ExpireTime)
		sign := GenSercertUrl(rv2.RoomId)
		if line == 1 {
			push = fmt.Sprintf("rtmp://%s/%s?e=%d&token=%s", path1, rv2.RoomId, expire, sign)

		} else {
			push = fmt.Sprintf("rtmp://%s/%s", path1, rv2.RoomId)
		}

		u, err := GetCacheUser(user.Uid)
		if err == redis.Nil {
			s := &CacheUser{
				Uid:       user.Uid,
				PreRoomId: rv2.RoomId,
			}
			SetCacheUser(user.Uid, s)
		} else if err != nil {
			return push, common.ERR_UNKNOWN
		} else {
			u.PreRoomId = rv2.RoomId
			SetCacheUser(user.Uid, u)
		}
		return push, common.ERR_SUCCESS
	}

	room := &RoomList{OwnerId: user.Uid, Statue: common.ROOM_PRE_V2, RoomId: common.GetOnlyId(user.Uid)}
	_, err = orm.InsertOne(room)
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return "", common.ERR_UNKNOWN
	}

	path1 := common.Cfg.MustValue("video", key)

	sign := GenSercertUrl(room.RoomId)

	var push string
	expire := time.Now().Unix() + int64(ExpireTime)
	if line == 1 {
		push = fmt.Sprintf("rtmp://%s/%s?e=%d&token=%s", path1, room.RoomId, expire, sign)

	} else {
		push = fmt.Sprintf("rtmp://%s/%s", path1, room.RoomId)
	}

	path2 := common.Cfg.MustValue("video", pullkey)

	pull := fmt.Sprintf("rtmp://%s/%s", path2, room.RoomId)

	room.LiveUrl = pull

	path3 := common.Cfg.MustValue("video", flvkey)
	flv := fmt.Sprintf("http://%s/%s.flv", path3, room.RoomId)
	room.FlvUrl = flv
	//生成移动端播放地址
	//port := common.Cfg.MustValue("video", portkey)
	mobilepath := common.Cfg.MustValue("video", mobilekey)
	if line == 1 {
		mobile = fmt.Sprintf("http://%s/%s.m3u8", mobilepath, room.RoomId)
	} else if line == 2 {
		mobile = fmt.Sprintf("http://%s/%s/playlist.m3u8", mobilepath, room.RoomId)
	}
	room.MobileUrl = mobile

	_, err = orm.Where("room_id=?", room.RoomId).MustCols("live_url").Update(room)
	if err != nil {
		common.Log.Errf("close  room error %s", err.Error())
	}

	u, err := GetCacheUser(user.Uid)
	if err == redis.Nil {
		s := &CacheUser{
			Uid:       user.Uid,
			PreRoomId: room.RoomId,
		}
		SetCacheUser(user.Uid, s)
	} else if err != nil {
		return push, common.ERR_UNKNOWN
	} else {
		u.PreRoomId = room.RoomId
		SetCacheUser(user.Uid, u)
	}

	return push, common.ERR_SUCCESS
	/*
		sql := fmt.Sprintf("select LAST_INSERT_ID() as id")
		var roomid int
		roomid = 0
		rowArray, err := orm.Query(sql)
		if err != nil {
			common.Log.Errf("mysql error is %s", err.Error())
			return 0, common.ERR_UNKNOWN
		}
		if len(rowArray) == 1 {
			roomid = common.BytesToInt(rowArray[0]["id"])
			return roomid, common.ERR_SUCCESS
		}
		return roomid, common.ERR_UNKNOWN
	*/
}

func GetRecommandMutipleWithPlayUrl(index int) ([]map[string]string, int) {
	sql := fmt.Sprintf("SELECT  c.room_id,c.room_name,c.owner_id,c.location,c.cover,c.play_url AS live_url,d.nick_name,d.sex,d.image,c.count AS viewer FROM (SELECT * FROM (SELECT  f.play_back_recommand_rid,g.play_url ,g.weight  AS WEI  FROM go_user_extra f LEFT JOIN go_multiple_play_back g ON f.play_back_recommand_rid=g.multiple_record_id WHERE play_back_recommand_type=1 AND  play_back_recommand_rid!='' AND g.hidden=0 ) a LEFT JOIN go_multiple_room_record b ON a.play_back_recommand_rid=b.id) c LEFT JOIN go_user d ON c.owner_id=d.uid")
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil, common.ERR_UNKNOWN
	}

	retMap := make([]map[string]string, 0)
	flag := true
	for _, row := range rowArray {

		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			if colName == "owner_id" {
				sess := GetUserSessByUid(common.BytesToInt(colValue))
				if sess != nil {
					flag = false
					break
				}
			}
			ss[colName] = value
		}
		if flag {
			//godump.Dump(ss)
			retMap = append(retMap, ss)
		}
		flag = true
	}
	return retMap, common.ERR_SUCCESS
}

func GetRecommandWithPlayUrl(index int) ([]map[string]string, int) {

	//retMap, _ := GetRecommandList(index)

	sql := fmt.Sprintf("SELECT c.room_id,c.room_name,c.owner_id,c.location,c.cover,c.play_url AS live_url,d.nick_name,d.sex,d.image,c.count as viewer FROM  (SELECT * FROM  (SELECT f.play_back_recommand_rid,g.play_url ,g.weight  as WEI FROM go_user_extra f LEFT JOIN go_play_back g ON f.play_back_recommand_rid=g.room_id WHERE play_back_recommand_rid!='' and g.hidden=0 ) a LEFT JOIN go_room_list b ON a.play_back_recommand_rid=b.room_id) c LEFT JOIN go_user d ON c.owner_id =d.uid ORDER BY c.WEI DESC,d.coupons DESC ,c.create_time  DESC")
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil, common.ERR_UNKNOWN
	}
	retMap := make([]map[string]string, 0)
	flag := true
	for _, row := range rowArray {

		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			if colName == "owner_id" {
				sess := GetUserSessByUid(common.BytesToInt(colValue))
				if sess != nil {
					flag = false
					break
				}
			}
			ss[colName] = value
		}
		if flag {
			//godump.Dump(ss)
			retMap = append(retMap, ss)
		}
		flag = true
	}

	return retMap, common.ERR_SUCCESS
}

func GetRecommandListWithGameType(test int, game_type ...int) ([]map[string]string, int) {
	var param string
	for k, v := range game_type {
		if k >= 1 {
			param = param + " or "
			param = fmt.Sprintf("%s game_type=%d", param, v)
		} else {
			param = fmt.Sprintf("game_type=%d", v)
		}
	}

	retMap := make([]map[string]string, 0)
	rowArray, err := orm.Query("select  a.room_id,a.room_name,a.owner_id,a.location,a.cover,a.live_url ,a.flv_url ,b.nick_name,b.sex,a.statue,b.image,a.roomtype,a.game_type from go_room_list a  left join go_user b on  a.owner_id=b.uid   where a.statue=1 and roomtype=? and ?  order by a.weight desc,b.coupons desc ,create_time  desc", test, param)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil, common.ERR_UNKNOWN
	}
	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			if colName == "room_id" {
				value := common.BytesToString(colValue)
				chat := GetChatRoom(value)
				if chat != nil {
					ss["viewer"] = strconv.Itoa(chat.GetCount() + chat.GetVRobotCount())
				} else {
					common.Log.Errf("room status is err please check rid= %s", value)
					break
				}
				ss[colName] = value
			}
		}
		retMap = append(retMap, ss)
	}
	return retMap, common.ERR_SUCCESS
}

func GetRecommandListV3(rids_uid []string) ([]map[string]string, []map[string]string, int) {
	a, b, c, k, l := GetRecommandList(0)
	ret := make([]map[string]string, 0)
	ret2 := make([]map[string]string, 0)
	t := make(map[string]bool, 0)

	t2 := make(map[string]bool, 0)
	for _, v := range rids_uid {
		s, ok := k[v]
		if ok {
			ret = append(ret, s)
			t[v] = true

		} else {
			//godump.Dump("not find ")
		}
	}

	for _, n := range a {
		s, ok := n["owner_id"]
		if ok {
			_, ok2 := t[s]
			if ok2 == false {
				ret = append(ret, n)
			}
		}

	}

	for _, v := range rids_uid {
		s, ok := l[v]
		if ok {
			ret2 = append(ret2, s)
			t2[v] = true

		} else {
			//godump.Dump("not find 2")
		}
	}

	for _, n := range b {
		s, ok := n["owner_id"]
		if ok {
			_, ok2 := t2[s]
			if ok2 == false {
				ret2 = append(ret2, n)
			}
		}
	}
	return ret, ret2, c
}

func GetRecommandListV4(rids_uid []string) ([]map[string]string, []map[string]string, int) {

	a, b, c, k, l := GetRecommandList2(0)
	ret := make([]map[string]string, 0)
	ret2 := make([]map[string]string, 0)
	t := make(map[string]bool, 0)

	t2 := make(map[string]bool, 0)
	for _, v := range rids_uid {
		s, ok := k[v]
		if ok {
			ret = append(ret, s)
			t[v] = true

		} else {
			//godump.Dump("not find ")
		}
	}

	for _, n := range a {
		s, ok := n["owner_id"]
		if ok {
			_, ok2 := t[s]
			if ok2 == false {
				ret = append(ret, n)
			}
		}

	}

	for _, v := range rids_uid {
		s, ok := l[v]
		if ok {
			ret2 = append(ret2, s)
			t2[v] = true

		} else {
			//godump.Dump("not find 2")
		}
	}

	for _, n := range b {
		s, ok := n["owner_id"]
		if ok {
			_, ok2 := t2[s]
			if ok2 == false {
				ret2 = append(ret2, n)
			}
		}
	}
	return ret, ret2, c
}

func GetRecommandList(index int) ([]map[string]string, []map[string]string, int, map[string]map[string]string, map[string]map[string]string) {
	//sql := fmt.Sprintf("select a.room_id,a.room_name,a.owner_id,a.location,a.cover,a.live_url,b.nick_name,b.sex,a.statue,b.image from room_list a  left join user b on  a.owner_id=b.uid   where a.statue=1  order by a.weight desc,b.coupons desc ,create_time  desc limit %d,%d ", index*common.ROOM_LIST_PAGE_COUNT, common.ROOM_LIST_PAGE_COUNT)
	sql := fmt.Sprintf("select a.room_id,a.room_name,a.owner_id,a.location,a.cover,a.live_url ,a.flv_url ,b.nick_name,b.sex,a.statue,b.image,a.roomtype,a.game_type  from go_room_list a  left join go_user b on  a.owner_id=b.uid   where a.statue=1   order by a.weight desc,b.coupons desc ,create_time  desc")
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil, nil, common.ERR_UNKNOWN, nil, nil
	}

	retMap := make([]map[string]string, 0)

	retTestMap := make([]map[string]string, 0)

	keyMap := make(map[string]map[string]string, 0)
	keyMapTest := make(map[string]map[string]string, 0)

	flag2 := true
	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			if colName == "room_id" {
				chat := GetChatRoom(value)
				if chat != nil {
					ss["viewer"] = strconv.Itoa(chat.GetCount() + chat.GetVRobotCount())

				} else {
					common.Log.Errf("room status is err please check rid= %s", value)
					flag2 = false
					break
				}
				ss[colName] = value
			} else if colName == "roomtype" {
				if value == "1" {
					flag2 = false
				}
			} else {
				ss[colName] = value
			}

		}

		if !flag2 {
			retTestMap = append(retTestMap, ss)
			keyMapTest[ss["owner_id"]] = ss

		} else {
			retMap = append(retMap, ss)
			retTestMap = append(retTestMap, ss)

			keyMap[ss["owner_id"]] = ss
			keyMapTest[ss["owner_id"]] = ss
		}

		flag2 = true
	}
	return retMap, retTestMap, common.ERR_SUCCESS, keyMap, keyMapTest
}

//本函数跟GetRecommandList的区别是去掉了game_type为3的直播间
func GetRecommandList2(index int) ([]map[string]string, []map[string]string, int, map[string]map[string]string, map[string]map[string]string) {
	//sql := fmt.Sprintf("select a.room_id,a.room_name,a.owner_id,a.location,a.cover,a.live_url,b.nick_name,b.sex,a.statue,b.image from room_list a  left join user b on  a.owner_id=b.uid   where a.statue=1  order by a.weight desc,b.coupons desc ,create_time  desc limit %d,%d ", index*common.ROOM_LIST_PAGE_COUNT, common.ROOM_LIST_PAGE_COUNT)
	sql := fmt.Sprintf("select a.room_id,a.room_name,a.owner_id,a.location,a.cover,a.live_url ,a.flv_url ,b.nick_name,b.sex,a.statue,b.image,a.roomtype,a.game_type  from go_room_list a  left join go_user b on  a.owner_id=b.uid   where a.statue=1  order by a.weight desc,b.coupons desc ,create_time  desc")
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil, nil, common.ERR_UNKNOWN, nil, nil
	}
	retMap := make([]map[string]string, 0)

	retTestMap := make([]map[string]string, 0)

	keyMap := make(map[string]map[string]string, 0)
	keyMapTest := make(map[string]map[string]string, 0)
	flag2 := true
	for _, row := range rowArray {
		gameType3 := false
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			if colName == "room_id" {
				chat := GetChatRoom(value)
				if chat != nil {
					ss["viewer"] = strconv.Itoa(chat.GetCount() + chat.GetVRobotCount())

				} else {
					common.Log.Errf("room status is err please check rid= %s", value)
					break
				}
				ss[colName] = value
			} else if colName == "roomtype" {
				if value == "1" {
					flag2 = false
				}
			} else if colName == "game_type" {
				if value == "3" {
					gameType3 = true
				} else {
					gameType3 = false
				}
				ss[colName] = value
			} else {
				ss[colName] = value
			}

		}

		if !gameType3 {
			if !flag2 {
				retTestMap = append(retTestMap, ss)
				keyMapTest[ss["owner_id"]] = ss
			} else {
				retMap = append(retMap, ss)
				retTestMap = append(retTestMap, ss)

				keyMap[ss["owner_id"]] = ss
				keyMapTest[ss["owner_id"]] = ss
			}
		}

		flag2 = true
	}
	return retMap, retTestMap, common.ERR_SUCCESS, keyMap, keyMapTest
}

func GetRoomRealUserCountList() ([]map[string]string, int, int) {
	total := 0
	sql := fmt.Sprintf("select a.room_id,a.room_name,a.owner_id,a.location,a.cover,a.live_url ,a.flv_url ,b.nick_name,b.sex,a.statue,b.image,a.roomtype,a.game_type  from go_room_list a  left join go_user b on  a.owner_id=b.uid   where a.statue=1  order by a.weight desc,b.coupons desc ,create_time  desc")
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil, total, common.ERR_UNKNOWN
	}

	retTestMap := make([]map[string]string, 0)

	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			if colName == "room_id" {
				chat := GetChatRoom(value)
				if chat != nil {
					num := chat.GetChatRealUserCount()
					ss["real_viewer"] = strconv.Itoa(num)
					total += num
				} else {
					common.Log.Errf("room status is err please check rid= %s", value)
					break
				}
				ss[colName] = value
			} else if colName == "room_name" || colName == "owner_id" || colName == "nick_name" {
				ss[colName] = value
			}
		}
		retTestMap = append(retTestMap, ss)

	}

	return retTestMap, total, common.ERR_SUCCESS
}

// func GetMultipleRoomList(index int) ([]map[string]string, int) {
// 	sql := fmt.Sprintf("select a.room_id,a.room_name,a.owner_id,a.location,a.cover,a.live_url,a.statue,b.nick_name,b.sex,b.image from multiple_room_list a  left join user b on  a.owner_id=b.uid   where a.statue=1 and a.owner_id!=0 order by a.weight desc,b.coupons desc ,create_time  desc ")
// 	rowArray, err := orm.Query(sql)
// 	if err != nil {
// 		common.Log.Errf("db err %s", err.Error())
// 		return nil, common.ERR_UNKNOWN
// 	}

// 	retMap := make([]map[string]string, 0)
// 	flag := true
// 	for _, row := range rowArray {
// 		ss := make(map[string]string)
// 		for colName, colValue := range row {
// 			value := common.BytesToString(colValue)
// 			if colName == "room_id" {
// 				chat := GetChatRoom(value)
// 				if chat != nil {
// 					ss["viewer"] = strconv.Itoa(chat.GetCount())
// 				} else {
// 					flag = false
// 					if !common.RefreshRecommndSwitch {
// 						break
// 					}
// 					_, err = orm.Exec("update multiple_room_list set status=? where room_id=?", common.MULTIPLE_ROOM_LOCK, value)
// 					if err != nil {
// 						common.Log.Errf("db err %s", err.Error())
// 						return nil, common.ERR_UNKNOWN
// 					}
// 					break
// 				}
// 			}
// 			ss[colName] = value
// 		}
// 		if flag {
// 			retMap = append(retMap, ss)
// 		}
// 		flag = true
// 	}

// 	return retMap, common.ERR_SUCCESS
// }
func GetMultipleRoomList(index int) ([]map[string]string, []map[string]string, int) {
	sql := fmt.Sprintf("select a.room_id,a.room_name,a.owner_id,a.location,a.cover,a.live_url,a.statue,b.nick_name,b.sex,b.image from go_multiple_room_list a  left join go_user b on  a.owner_id=b.uid   where a.statue=1 and a.owner_id!=0 order by a.weight desc,b.coupons desc ,create_time  desc")
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil, nil, common.ERR_UNKNOWN
	}

	retMap := make([]map[string]string, 0)
	retTestMap := make([]map[string]string, 0)
	flag := true
	flag2 := true
	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			if colName == "room_id" {
				chat := GetChatRoom(value)
				if chat != nil {
					ss["viewer"] = strconv.Itoa(chat.GetCount() + chat.GetVRobotCount())
				} else {
					flag = false
					if !common.RefreshRecommndSwitch {
						break
					}
					_, err = orm.Exec("update go_multiple_room_list set status=? where room_id=?", common.MULTIPLE_ROOM_LOCK, value)
					if err != nil {
						common.Log.Errf("db err %s", err.Error())
						return nil, nil, common.ERR_UNKNOWN
					}
					break
				}
			}
			ss[colName] = value

			if colName == "owner_id" {
				uid_, _ := strconv.Atoi(colName)
				user, _ := GetUserByUid(uid_)
				if user.AccountType > 0 {
					flag2 = true
				} else {
					flag2 = false
				}
			}
		}

		if flag {
			if flag2 {
				retMap = append(retMap, ss)
				retTestMap = append(retTestMap, ss)
			} else {
				retTestMap = append(retTestMap, ss)
			}
		}
		flag = true
		flag = true
	}

	return retMap, retTestMap, common.ERR_SUCCESS
}

type MonitorData struct {
	Uid     string `json:"userId"`
	LiveUrl string `json:"playUrl"`
}

func GetMonitorBaseRoom() []MonitorData {
	//sql := fmt.Sprintf("select    select a.live_url,b.nick_name,b.sex,a.statue,b.image from room_list a  left join user b on  a.owner_id=b.uid   where a.statue=1  order by a.weight desc,b.coupons desc ,create_time  desc ")

	r := make([]RoomList, 0)
	err := orm.Where("statue=?", common.USER_STATUE_LIVE).Find(&r)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil
	}
	monitor := make([]MonitorData, 0)
	for _, v := range r {
		sess := GetUserSessByUid(v.OwnerId)
		if sess == nil {
			continue
		}
		d := MonitorData{}
		d.Uid = strconv.Itoa(v.OwnerId)
		d.LiveUrl = v.LiveUrl
		monitor = append(monitor, d)
	}
	return monitor
}

func RefreshRoomStatus() int {
	rlist := make([]RoomList, 0)
	err := orm.Where("statue=?", common.USER_STATUE_LIVE).Find(&rlist)

	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN
	}
	for _, v := range rlist {
		chat := GetChatRoom(v.RoomId)
		if chat == nil {
			v.Statue = common.ROOM_MODIFY
			aff_row, err := orm.Where("room_id=?", v.RoomId).Update(&v)
			if err != nil {
				common.Log.Errf("db err %s", err.Error())
				return common.ERR_UNKNOWN
			}

			if aff_row == 0 {
				return common.ERR_DB_UPDATE
			}
		}
	}
	return common.ERR_SUCCESS
}

func CheckRoomStatus(uid int) int {
	rlist := make([]RoomList, 0)
	err := orm.Where("statue=? and owner_id=?", common.USER_STATUE_LIVE, uid).Find(&rlist)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN
	}
	flag := false
	for _, v := range rlist {
		chat := GetChatRoom(v.RoomId)
		if chat != nil {
			common.Log.Debugf("pre live close room uid=%d,rid=%s", v.OwnerId, v.RoomId)
			user, _ := GetUserByUid(v.OwnerId)

			v.Statue = common.ROOM_MODIFY
			v.FinishTime = time.Now()
			_, err := orm.Where("room_id=?", v.RoomId).Update(&v)
			if err != nil {
				common.Log.Errf("db err %s", err.Error())
				return common.ERR_UNKNOWN
			}

			DelUserSession(user.Uid)
			common.Log.Infof("record close room status uid=%d,rid=%s", user.Uid, v.RoomId)
			CloseChat(user.Uid, v.RoomId)

			if sess := GetUserSessByUid(user.Uid); sess != nil {
				//sess.CloseSesion()
				sess.Sess.Close()
				return common.ERR_USER_SESS
			}
		} else {
			v.Statue = common.ROOM_MODIFY
			v.FinishTime = time.Now()
			aff_row, err := orm.Where("room_id=?", v.RoomId).Update(&v)
			if err != nil {
				common.Log.Errf("db err %s", err.Error())
				return common.ERR_UNKNOWN
			}

			if aff_row == 0 {
				return common.ERR_DB_UPDATE
			}
		}
		flag = true
	}

	if flag {
		return common.ERR_USER_SESS
	}
	return common.ERR_SUCCESS
}

func CheckLeader(uid int) int {

	res, err := orm.Query("select * from php_user_admin_link WHERE user_id=?", uid)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if len(res) == 0 {
		return common.ERR_SUCCESS
	}
	return common.ERR_PRE_LIVE_ADMIN

}

func InitConsistData() int {
	//format_time := time.Now().Local().Format("2006-01-02 15:04:05")
	/*
		t := time.Now()
		cur := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
		cur_tm := cur.AddDate(0, 0, -7)

		arr:=make([]RoomList,0)
		//err:=orm.Where("statue=? or statue=? and create_time>? ",common.ROOM_ONLIVE,common.ROOM_RESTART,cur_tm ).Find(&arr)
		err:=orm.Where("statue=? or statue=? and create_time>? ",common.ROOM_ONLIVE,common.ROOM_RESTART,cur_tm ).Find(&arr)
		if err!=nil {
			godump.Dump(err)
			return common.ERR_UNKNOWN
		}

		for _,v:=range arr  {
			DelAudienceKey(v.RoomId)
		}
	*/
	res, err := orm.Exec("update go_room_list set statue=? ,finish_time=now() where statue=?", common.ROOM_RESTART, common.ROOM_ONLIVE)

	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN
	}
	aff_row, err := res.RowsAffected()
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_UPDATE
	}

	return common.ERR_SUCCESS
	//return AddFinishTime()
}

func AddFinishTime() int {
	t := time.Now()
	cur_tm := t.AddDate(0, 0, -1)
	s := make([]RoomList, 0)
	err := orm.Where("create_time<? and statue=? and finish_time='0001-01-01 00:00:00'", cur_tm, common.ROOM_RESTART).Find(&s)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN
	}

	for _, v := range s {
		timeNow := v.CreateTime.Format("2006-01-02 15:04:05")
		limit_time := v.CreateTime.Add(3600 * 6 * time.Second)
		u2 := limit_time.Format("2006-01-02 15:04:05")
		res, err := orm.Query("select * from go_room_list where owner_id=? and create_time>? and create_time<=? and statue=? order by create_time limit 1", v.OwnerId, timeNow, u2, common.ROOM_FINISH)
		if err != nil {
			common.Log.Errf("db err %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if len(res) == 1 {
			b, _ := res[0]["create_time"]
			ctime := common.BytesToString(b)
			t, _ = time.Parse("2006-01-02 15:04:05", ctime)
			_, err := orm.Exec("update go_room_list set finish_time=? where room_id=?", t, v.RoomId)
			if err != nil {
				common.Log.Errf("db err %s", err.Error())
				return common.ERR_UNKNOWN
			}

		} else {
			_, err := orm.Exec("update go_room_list set finish_time=? where room_id=?", u2, v.RoomId)
			if err != nil {
				common.Log.Errf("db err %s", err.Error())
				return common.ERR_UNKNOWN
			}
		}
	}
	return common.ERR_SUCCESS
}

func CloseAllGameRoom() {
	sql := fmt.Sprintf("select room_id  from go_room_list  where game_type != 0 AND statue=1 ")
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return
	}

	for _, row := range rowArray {
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			if colName == "room_id" {
				var data ResponseSys
				data.MType = common.MESSAGE_TYPE_ADMIN_CLOSE
				data.Notice = "游戏直播间已经被管理员关闭如有问题联系管理员"
				SendMsgToRoom(value, data)
				//room := GetChatRoom(value)
				//if room != nil {
				//	CloseRoom(room, room.GetChatInfo().Uid)
				//}
				fmt.Print("\n CloseAllGameRoom() 1100行 room_id=%s", value)
				DirectCloseRoom(1, value)
			}
		}
	}
}
