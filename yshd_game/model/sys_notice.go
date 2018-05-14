package model

import (
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	//"strconv"
)

type ConfigSysNotice struct {
	NoticeId int    `xorm:"int(11) pk not null autoincr"`
	Content  string `xorm:"varchar(256) not null"`
	Status   bool
}

var sys_notice_map map[int]ConfigSysNotice

var sys_notice_ad string

func InitSysNotice() {
	sys_notice_map = make(map[int]ConfigSysNotice)
	LoadSysNotice()

	for _, v := range sys_notice_map {
		if v.Status == true {
			sys_notice_ad = v.Content
			return
		}
	}
	/*
		if sys_notice_ad_, ok := sys_notice_map[1]; ok {
			sys_notice_ad = sys_notice_ad_.Content
			return
		}
	*/
	common.Log.Info("ad err is nil")
}

func LoadSysNotice() map[int]ConfigSysNotice {
	err := orm.Find(&sys_notice_map)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return nil
	}
	return sys_notice_map
}

func GetSysNotice() string {
	return sys_notice_ad
}

func ResetNotice() {
	sys_notice_map_tmp := make(map[int]ConfigSysNotice)
	err := orm.Find(&sys_notice_map_tmp)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return
	}
	for _, v := range sys_notice_map_tmp {
		if v.Status {
			sys_notice_ad = v.Content
			return
		}
	}
	sys_notice_ad = ""
}
