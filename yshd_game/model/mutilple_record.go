package model

import (
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	//"github.com/yshd_game/timer"
	//"strconv"
	//"fmt"
	//"sync"
	"time"
)

type MultipleRoomRecord struct {
	Id         int64
	RoomId     string    `xorm:"varchar(128)   not null"` //房间ID
	RoomName   string    `xorm:"varchar(255)"  `          //房间名字
	OwnerId    int       `xorm:"not null "`               //主播ID
	CreateTime time.Time //创建时间
	FinishTime time.Time //结束时间
	Location   string    `xorm:"varchar(128)"  ` //定位
	Cover      string    `xorm:"varchar(255)"  ` //封面图片
	LiveUrl    string    `xorm:"varchar(255)"  ` //直播流
	MobileUrl  string    `xorm:"varchar(255)"  ` //移动流
	Rice       int       //收到的米粒
	Count      int       //人数
	Weight     int       `xorm:"not null default(0)"` //排序权重
	LockTime   int64
	Playback   string
}

func GetMutipleRoomRecordByID(id int64) (*MultipleRoomRecord, bool) {
	m := &MultipleRoomRecord{}
	ok, err := orm.Id(id).Get(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return nil, false
	}
	return m, ok
}

func AddMutipleRecord(rid, rname, location, cover, liveurl, mobileurl string, uid int) int64 {
	m := &MultipleRoomRecord{
		RoomId:     rid,
		RoomName:   rname,
		OwnerId:    uid,
		Location:   location,
		Cover:      cover,
		LiveUrl:    liveurl,
		MobileUrl:  mobileurl,
		CreateTime: time.Now(),
	}
	_, err := orm.Insert(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return 0
	}
	return m.Id
}

func CloseMutipleRecord(id int64, rice, count int) int {
	m := &MultipleRoomRecord{}
	_, err := orm.Id(id).Get(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	m.Rice = rice
	m.Count = count
	m.FinishTime = time.Now()
	_, err = orm.Id(id).Update(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS
}
