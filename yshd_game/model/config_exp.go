package model

import (
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
)

var comsumer_config_map map[int]ConfigConsumer

type ConfigConsumer struct {
	Cid        int `xorm:"int(11) pk not null"`
	ConsumeNum int `xorm:"int(11) not null "`
	Exp        int `xorm:"int(11) not null"`
}

func LoadComsumer() map[int]ConfigConsumer {
	comsumer_config_map = make(map[int]ConfigConsumer)
	err := orm.Find(&comsumer_config_map)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return nil
	}
	return comsumer_config_map
}

func GetComsumerById(cid int) (ConfigConsumer, bool) {
	comsmer, exist := comsumer_config_map[cid]
	return comsmer, exist
}
