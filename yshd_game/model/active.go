package model

import "github.com/yshd_game/common"

type ChargeActive struct {
	ItemId     string `xorm:" varchar(64) pk not null"` //充值道具ID
	MoneyType  int    `xorm:"not null "`                //赠送金钱类型
	ExtraNum   int64  `xorm:"not null "`                //赠送金钱数量
	Status     int    `xorm: int(11) `                  //开关0关 1开
	BeginTime  int64  //开始时间
	FinishTime int64  //结束时间
}

func GetChargeActive(item_id string) *ChargeActive {
	if item_id == "" {
		return nil
	}
	m := &ChargeActive{}
	_, err := orm.Where("item_id=?", item_id).Get(m)
	if err != nil {
		common.Log.Err(err.Error())
		return nil
	}
	return m
}
