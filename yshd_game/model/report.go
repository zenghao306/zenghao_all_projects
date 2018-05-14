package model

import (
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	//"strconv"
	"time"
)

type Report struct {
	Id         int64
	Uid        int    `xorm:"int(11) not null"`
	Oid        int    `xorm:"int(11) not null"`
	Desc       string `xorm:"varchar(255) not null"`
	CreateTime time.Time
}

func AddReport(uid, oid int, desc string) int {
	aff_row, err := orm.Insert(&Report{Uid: uid, Oid: oid, Desc: desc, CreateTime: time.Now()})
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	return common.ERR_SUCCESS
}
