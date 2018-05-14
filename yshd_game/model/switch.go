package model

import (
	"github.com/yshd_game/common"
)

func GetChannelSwitch(channel_id string) (status int, ret int) {
	res, err := orm.Query("select status from php_cash_moon_switch where channel_id=?", channel_id)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}

	if len(res) == 0 {
		ret = common.ERR_CONFIG_SWITCH
		return
	}

	b := res[0]["status"]
	status = common.BytesToInt(b)
	ret = common.ERR_SUCCESS
	return
}
