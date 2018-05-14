package model

import (
	"fmt"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"strconv"
)

type MultiplePlayBack struct {
	//Id              int64
	MultipleRecordId int64  `xorm:"varchar(48) pk `
	Uid              int    `xorm:"int(11) not null "`
	PlayUrl          string `xorm:"varchar(128)`
	SaveTime         int64
	Hidden           int `xorm:"int(11) not null "`
	Weight           int `xorm:"not null default(0)"`
}

func AddMutiplePlayBack(uid int, rid, url string, rtime, recordId int64) int {
	user, ok := GetUserExtraByUid(uid)

	if ok == false {
		return common.ERR_ACCOUNT_EXIST
	}

	if user.PlayBackCount >= 5 {
		return common.ERR_PLAYBACK_MAX
	}

	allurl := fmt.Sprintf("http://%s/%s", DomainVod, url)
	res, err := orm.Exec("insert into go_multiple_play_back (`uid`,`play_url`,`save_time`,`multiple_record_id`) values  (?,?,?,?)", uid, allurl, rtime, recordId)
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

func DelMutiplePlayBack(uid int, recordId string) int {
	_, err := orm.Where("uid=? and multiple_record_id=?", uid, recordId).Delete(&MultiplePlayBack{})
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	extra, has := GetUserExtraByUid(uid)
	if has {
		if extra.PlayBackRecommandType == common.MUTIPLE_ROOM {
			if extra.PlayBackRecommandRid == recordId {
				CancelRecommandFlag(uid)
			}
		}
	}
	return common.ERR_SUCCESS
}

func UpdateMutipleRecommandFlag(uid int, id string) int {
	extra, has := GetUserExtraByUid(uid)
	if has {
		has2, err := orm.Where("multiple_record_id=?", id).Get(&MultiplePlayBack{})
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if !has2 {
			return common.ERR_PARAM
		}

		extra.PlayBackRecommandType = common.MUTIPLE_ROOM
		extra.PlayBackRecommandRid = id
		_, err = extra.UpdateByColS("play_back_recommand_rid, play_back_recommand_type")
		if err != nil {
			return common.ERR_UNKNOWN
		}

		c := NewChatRoomInfo()

		o := c.GetChatInfo()
		o.Rid = id
		o.Uid = uid

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return common.ERR_UNKNOWN
		}
		room, has3 := GetMutipleRoomRecordByID(intId)
		if has3 == true {
			o.Image = room.Cover
		}
		c.Statue = common.ROOM_PLAYBACK
		AddChatRoom(c)
		return common.ERR_SUCCESS
	}
	return common.ERR_PARAM
}

func GetMultiplePlayBackList(uid int) []map[string]string {
	rowArray, err := orm.Query("SELECT a.uid,a.play_url,a.save_time ,b.room_name,b.count,b.id as room_id FROM go_multiple_play_back a LEFT JOIN go_multiple_room_record b ON a.multiple_record_id=b.id WHERE a.uid=? ORDER BY a.weight", uid)
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
