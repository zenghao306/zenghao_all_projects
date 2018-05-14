package model

import (
	"github.com/yshd_game/common"
	"time"
)

type PvInfo struct {
	Id         int64
	Uid        int
	Rid        string
	RecordTime int64
	Ip         string
}

func StatisticsRoom(rid string, uid int, ip string) int {
	m := &PvInfo{
		Uid:        uid,
		Rid:        rid,
		RecordTime: time.Now().Unix(),
		Ip:         ip,
	}
	aff_row, err := orm.Insert(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	if aff_row == 0 {
		return common.ERR_DB_UPDATE
	}
	return common.ERR_SUCCESS
}
