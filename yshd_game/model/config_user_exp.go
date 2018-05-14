package model

import (
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
)

var user_exp_config_map map[int]ConfigUserExp

type ConfigUserExp struct {
	Level int `xorm:"int(11) pk not null "`
	Exp   int `xorm:"int(11) not null"`
}

func LoadConfigUserExp() map[int]ConfigUserExp {
	user_exp_config_map = make(map[int]ConfigUserExp)
	err := orm.Find(&user_exp_config_map)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return nil
	}
	return user_exp_config_map
}

func GetUserExpByLevel(level int) (ConfigUserExp, bool) {
	comsmer, exist := user_exp_config_map[level]
	return comsmer, exist
}
