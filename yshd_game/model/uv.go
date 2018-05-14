package model

import (
	"github.com/yshd_game/common"
	"time"
)

type UvInfo struct {
	Imei       string `xorm:"varchar(128)  pk"`
	CreateTime int64
	Device     string `xorm:"varchar(128) "`
	ChannelId  string `xorm:"varchar(128) "`
}

func AddImei(imei string, device string, channel string) int {

	has, err := orm.Where("imei=?", imei).Get(&UvInfo{})
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has {
		common.Log.Err(imei)
		return common.ERR_DUPLICATE_UV
	}
	var info UvInfo
	info.Imei = imei
	info.Device = device
	info.CreateTime = time.Now().Unix()
	info.ChannelId = channel
	_, err = orm.InsertOne(&info)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS
}
