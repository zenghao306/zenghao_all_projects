package model

import (
	"github.com/yshd_game/common"
)

type ConfigScoreExchange struct {
	Diamond int   `xorm:"int(11) default(0) pk"`
	Score   int64 `xorm:" default(0)"`
}

var score_map map[int]ConfigScoreExchange

var socre_cahce []ConfigScoreExchange

func LoadScoreExchange() map[int]ConfigScoreExchange {
	score_map = make(map[int]ConfigScoreExchange)
	socre_cahce = make([]ConfigScoreExchange, 0)

	err := orm.Find(&score_map)
	if err != nil {
		common.Log.Panic("orm err is %s", err.Error())
	}

	err = orm.Find(&socre_cahce)
	if err != nil {
		common.Log.Panic("orm err is %s", err.Error())
	}
	/*
		for _, v := range score_map {
			socre_cahce = append(socre_cahce, v)
		}
	*/
	return score_map
}

func GetScoreById(id int) *ConfigScoreExchange {
	n, ok := score_map[id]
	if ok {
		return &n
	}
	return nil
}

func GetScoreConfig() []ConfigScoreExchange {
	return socre_cahce
}
