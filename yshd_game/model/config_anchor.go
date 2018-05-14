package model

import (
	"github.com/yshd_game/common"
)

var anchor_config_map map[int]ConfigAnchorExp

type ConfigAnchorExp struct {
	Level int `xorm:"int(11) pk not null"`
	Exp   int `xorm:"int(11) not null "`
}

func LoadAnchorExp() map[int]ConfigAnchorExp {
	anchor_config_map = make(map[int]ConfigAnchorExp)
	err := orm.Find(&anchor_config_map)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return nil
	}
	return anchor_config_map
}

func GetAnchorById(cid int) (ConfigAnchorExp, bool) {
	anchor, exist := anchor_config_map[cid]
	return anchor, exist
}
