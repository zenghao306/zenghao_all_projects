package model

import (
	"fmt"
	//"github.com/liudng/godump"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
)

type PlayBack struct {
	Id       int64
	Uid      int    `xorm:"int(11) not null "`
	PlayUrl  string `xorm:"varchar(128)`
	SaveTime int64
	RoomId   string `xorm:"varchar(48) UNIQUE(PLAY_ROOM_ID)`
	Hidden   int    `xorm:"int(11) not null "`
	Weight   int    `xorm:"not null default(0)"`
}

//新增回播记录
func AddNewPlayBack(uid int, url, rid string, rtime int64) int {

	if ret := CheckPlayList(uid); ret == common.ERR_PLAYBACK_MAX {

		m := &PlayBack{}
		has, err := orm.Where("uid=?", uid).Asc("save_time").Limit(1, 0).Get(m)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if has {
			DelPlayBackList(uid, m.RoomId)
		}

	}

	allurl := fmt.Sprintf("http://%s/%s", DomainVod, url)
	res, err := orm.Exec("insert into go_play_back (`uid`,`play_url`,`save_time`,`room_id`) values  (?,?,?,?)", uid, allurl, rtime, rid)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return common.ERR_UNKNOWN
	}
	aff_row, err := res.RowsAffected()
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	return common.ERR_SUCCESS
}

//获取回播列表
func GetPlayBackList(uid int) []map[string]string {
	rowArray, err := orm.Query("select a.uid,a.play_url,a.save_time ,b.room_name,b.count,b.room_id  from go_play_back a left join go_room_list b on a.room_id=b.room_id where a.uid=? order by a.weight", uid)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return nil
	}

	retMap := make([]map[string]string, 0)
	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		retMap = append(retMap, ss)
	}
	return retMap
}

//删除回播记录同时删除文件
func DelPlayBackList(uid int, roomid string) int {
	m := &PlayBack{}
	_, err := orm.Where("uid=? and room_id=?", uid, roomid).Get(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	will_del_file := m.PlayUrl

	aff_row, err := orm.Where("uid=? and room_id=?", uid, roomid).Delete(&PlayBack{})
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	if aff_row == 0 {
		return common.ERR_DB_DEL
	}
	//DelQiNiuFile3(DomainVod,will_del_file)

	extra, has := GetUserExtraByUid(uid)
	if has {
		if extra.PlayBackRecommandRid == roomid {
			CancelRecommandFlag(uid)
		}
	}

	DelQiNiuFile3(bucket_vod, DomainVod, will_del_file)
	return common.ERR_SUCCESS
}

//通过ID查找回播记录
func GetPlayBack(rid string) (*PlayBack, bool) {
	v := &PlayBack{}
	has, err := orm.Where("room_id=?", rid).Get(v)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return nil, false
	}
	return v, has
}

//检查回播记录限制
func CheckPlayList(uid int) int {
	rowArray, err := orm.Query("select count(*) as count from go_play_back where uid= ?", uid)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	count := common.BytesToInt(rowArray[0]["count"])
	if count >= 5 {
		return common.ERR_PLAYBACK_MAX
	}

	rowArray, err = orm.Query("select count(*) as count from go_multiple_play_back where uid=?", uid)

	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	count2 := common.BytesToInt(rowArray[0]["count"])

	if count+count2 >= 5 {
		return common.ERR_PLAYBACK_MAX
	}
	return common.ERR_SUCCESS

	/*
		user, ok := GetUserExtraByUid(uid)
		if ok == false {
			return common.ERR_ACCOUNT_EXIST
		}
		if user.PlayBackCount >= 5 {
			return common.ERR_PLAYBACK_MAX
		}
		return common.ERR_SUCCESS
	*/
}

//设置回播推荐
func UpdateRecommandFlag(uid int, rid string) int {
	extra, has := GetUserExtraByUid(uid)
	if has {
		has2, err := orm.Where("uid=? and room_id=?", uid, rid).Get(&PlayBack{})
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if has2 {
			willclose := extra.PlayBackRecommandRid
			if extra.PlayBackRecommandRid == rid {
				return common.ERR_SUCCESS
			}
			extra.PlayBackRecommandRid = rid
			_, err := extra.UpdateByColS("play_back_recommand_rid")
			if err != nil {
				return common.ERR_UNKNOWN
			}

			if willclose != "" {
				common.Log.Debugf("delete chat room when update flag uid=%d,rid=%s", uid, rid)
				DelChatRoom(willclose)
			}
			c := NewChatRoomInfo()

			o := c.GetChatInfo()
			o.Rid = rid
			o.Uid = uid

			room, has3 := GetRoomById(rid)
			if has3 == true {
				o.Image = room.Cover
			}
			c.Statue = common.ROOM_PLAYBACK
			AddChatRoom(c)
			return common.ERR_SUCCESS
		}
		return common.ERR_ROOM_EXIST
	}
	return common.ERR_UNKNOWN
}

func CancelRecommandFlag(uid int) int {
	extra, has := GetUserExtraByUid(uid)
	if has {
		delRid := extra.PlayBackRecommandRid
		extra.PlayBackRecommandRid = ""
		_, err := extra.UpdateByColS("play_back_recommand_rid")
		if err != nil {
			return common.ERR_UNKNOWN
		}
		if extra.PlayBackRecommandType == common.SIGNEL_ROOM {
			common.Log.Debugf("delete chat room when cancle uid=%d", uid)
			DelChatRoom(delRid)
		}
		return common.ERR_SUCCESS
	}
	return common.ERR_UNKNOWN
}

func HiddenPlayBack(rid string) int {
	res, err := orm.Exec("update go_play_back set hidden=1 where room_id=?", rid)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	aff_row, err := res.RowsAffected()
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_UPDATE
	}
	return common.ERR_SUCCESS
}
