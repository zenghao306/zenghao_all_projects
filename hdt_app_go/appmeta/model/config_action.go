package model

import (
	//"fmt"
	//"github.com/go-xorm/xorm"
	//"hdt_app_go/common"
	//proto "hdt_app_go/protcol"
	//. "hdt_app_go/appmeta/log"
	//"strconv"
	"sync"
)

const (
	LIMIT_ACTION_NONE       = iota
	LIMIT_ACTION_ZHUCE      = 1
	LIMIT_ACTION_DENGLU     = 2
	LIMIT_ACTION_FAYAN      = 3
	LIMIT_ACTION_PVFANGJIAN = 4
	LIMIT_ACTION_SONGLI     = 5
	LIMIT_ACTION_FENXIANG   = 6
	LIMIT_ACTION_GUANZHU    = 7
	LIMIT_ACTION_PVZHIBO    = 8
	LIMIT_ACTION_CHONGZHI   = 9
)

type ConfigAction struct {
	Id         int64
	Action     string
	ActionName string
	Power      float64
	Desc       string
	Limit      string
	LimitValue int
}

type ConfigBaseData struct {
	action_mutex  sync.Mutex
	config_action map[int64]ConfigAction

	limit_mutex  sync.Mutex
	config_limit map[string]int
}

func NewConfigBaseData() *ConfigBaseData {
	return &ConfigBaseData{
		config_action: make(map[int64]ConfigAction),
		config_limit:  make(map[string]int),
	}
}
