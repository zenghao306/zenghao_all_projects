package model

import (
	"github.com/yshd_game/common"
)

type GroupWeight struct {
	GroupId int `xorm:"pk not null"` //组别ID
	Weight  int `xorm:"not null "`   //权重
}

var weight_mgr_map map[int]GroupWeight

func InitWeightMgr() {
	weight_mgr_map = make(map[int]GroupWeight)
	err := orm.Find(&weight_mgr_map)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return
	}
}

func GetWeightByGroupId(id int) int {
	t, ok := weight_mgr_map[id]
	if ok {
		return t.Weight
	}
	return 0
}

func GetWeightByGroupId2(id int) int {
	s := &GroupWeight{}
	has, err := orm.Where("group_id=?", id).Get(s)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return 0
	}
	if has {
		return s.Weight
	}
	return 0
}
