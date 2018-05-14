package model

import (
	//"fmt"
	//"github.com/go-xorm/xorm"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	//"math"
	//"strconv"
	"time"
)

type WatchRecord struct {
	Id            int64
	Uid           int    `xorm:"int(11) not null "` //用户ID
	JoinChatTime  int64  `xorm:" default(0)"`       //进入房间时间
	LevelChatTime int64  `xorm:" default(0)"`       //离开房间时间
	Rid           string `xorm:"vchar(255) "`       //房间ID
	RoomType      int    `xorm:" default(0)"`
	Ip            string `xorm:"vchar(128) "`
}

func JoinChat(uid int, rid string, room_type int, ip string) (bool, int64) {
	if uid == 0 || rid == "" {
		common.Log.Errf("join chat uid=%d rid=%s,room_type", uid, rid, room_type)
		return false, 0
	}
	watch := &WatchRecord{Uid: uid, JoinChatTime: time.Now().Unix(), Rid: rid, RoomType: room_type, Ip: ip}
	aff_row, err := orm.Insert(watch)

	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return false, 0
	}
	if aff_row == 0 {
		return false, watch.Id
	}
	return true, watch.Id
}

//离开时候进行观看时间统计和经验的处理
func LeaveChat(id int64) bool {
	watch := &WatchRecord{}
	has, err := orm.Id(id).Get(watch)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return false
	}
	nowtime := time.Now()
	if has {
		watch.LevelChatTime = nowtime.Unix()
		_, err = orm.Id(id).Update(watch)
		if err != nil {
			common.Log.Errf("db err %s", err.Error())
			return false
		}
		extre, has := GetUserExtraByUid(watch.Uid)
		if !has {
			return false
		}
		nowtime_ := nowtime.Format("20060102")
		dff := watch.LevelChatTime - watch.JoinChatTime

		var owner int
		if watch.RoomType == 0 {
			room, has := GetRoomById(watch.Rid)
			if has {
				owner = room.OwnerId
			} else {
				common.Log.Errf("watch_id is null wid=%s  time=%d ", id, time.Now().Unix())
				return false
			}
		} else if watch.RoomType == 1 {
			has, room := GetMultipleRoomByRid(watch.Rid)
			if has {
				owner = room.OwnerId
			} else {
				common.Log.Errf("watch_id is null wid=%s now=%d ", id, time.Now().Unix())
				return false
			}
		}

		if owner == watch.Uid {
			add := extre.DayAnchorTime / 3600
			if nowtime_ == extre.LastWatchTime {
				extre.DayAnchorTime += int(dff)
			} else {
				extre.DayAnchorTime = int(dff)
				extre.LastWatchTime = nowtime_
				extre.AddAnchorTimes = 0
			}

			if add > 5 {
				add = 5
			}
			diff := add - extre.AddAnchorTimes
			if diff > 0 {
				anchor_user, _ := GetUserByUid(watch.Uid)
				anchor_user.AddUserExp(nil, 2*diff, false)
				extre.AddAnchorTimes = add
			}
		} else {
			if nowtime_ == extre.LastWatchTime {
				extre.DayWatchTime += int(dff)
			} else {
				extre.DayWatchTime = int(dff)
				extre.AddWacthTimes = 0
				extre.LastWatchTime = nowtime_
			}

			add := extre.DayWatchTime / 3600
			if add >= 4 {
				add = 4
			}
			diff := add - extre.AddWacthTimes
			if diff > 0 {
				main_user, _ := GetUserByUid(watch.Uid)
				main_user.AddUserExp(nil, 5*diff, false)
				extre.AddWacthTimes = add
			}
		}
		extre.UpdateByColS("day_anchor_time", "day_watch_time", "add_wacth_times", "last_watch_time", "add_anchor_times")
		return true

	}
	return false
}
