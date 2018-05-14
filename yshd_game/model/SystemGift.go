package model

import (
	"github.com/yshd_game/common"
	"time"
)

type SystemGiftRecord struct {
	Id         int64
	UserId     int
	AdminId    int
	GiftId     int
	MoneyType  int
	Num        int
	RecordTime int64
}

func AddSysGiftRecord(user_id, admin_id, gift_id, num int, money_type int) {
	m := &SystemGiftRecord{
		UserId:     user_id,
		AdminId:    admin_id,
		GiftId:     gift_id,
		Num:        num,
		RecordTime: time.Now().Unix(),
		MoneyType:  money_type,
	}
	_, err := orm.Insert(m)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return
	}
	return
}
