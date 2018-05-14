package model

import (
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	//"strconv"
)

type ConfigAdMgr struct {
	AdId  int    `xorm:"int(11) pk not null unique"`
	Image string `xorm:"varchar(128) not null"`
}

var ad_mgr_map map[int]ConfigAdMgr

var new_user_ad ConfigAdMgr

func InitAdMgr() {
	ad_mgr_map = make(map[int]ConfigAdMgr)
	LoadAdMgr()
	var ok bool
	if new_user_ad, ok = ad_mgr_map[1]; !ok {
		common.Log.Err("ad err is nil")
	}
}

func LoadAdMgr() map[int]ConfigAdMgr {
	err := orm.Find(&ad_mgr_map)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return nil
	}
	return ad_mgr_map
}

func GetAd() *ConfigAdMgr {
	return &new_user_ad
}
