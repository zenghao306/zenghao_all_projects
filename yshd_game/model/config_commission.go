package model

import (
	"github.com/yshd_game/common"
)

type ConfigCommossion struct {
	//CommossionType int `xorm:"int(11) not null pk autoincr"` //分配类型

	/*
		CommonDiamonGift      uint32 //普通礼物分配比例
		CommonDiamonValue uint32
		ExtravagantDiamonGift uint32 //奢华礼物分配比例
		CommonGameGift        uint32 //普通游戏礼物分配比例
		ExtravagantGameGift   uint32 //奢华游戏礼物分配比例
		GameCommssion         uint32 //系统抽成
	*/
	GiftCommossion   int     `xorm:"int(11) not null"` //关联礼物分配类型
	UserCommossion   int     `xorm:"int(11) not null"` //关联用户分配类型
	OwnerProportion  float32 `xorm:" not null"`        //自己分成
	AdminProportion  float32 `xorm:"int(11) not null"` //管理者分成
	SystemProportion float32 `xorm:"int(11) not null"` //系统分成

	Pump float32 `xorm:"int(11) not null"` //抽水
}

var CommmossionMgr map[int]ConfigCommossion

func LoadCommossion() {
	CommmossionMgr = make(map[int]ConfigCommossion)
	err := orm.Find(&CommmossionMgr)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return
	}
}

/*
func GetCommossion(gift_alloc, user_alloc int) *ConfigCommossion {
	var c ConfigCommossion
	has, err := orm.Where("gift_commossion=? and user_commossion=?", gift_alloc, user_alloc).Get(&c)
	if err != nil {
		common.Log.Errf("db error:", err.Error())

		return nil
	}
	if has {
		return &c
	}
	return nil
}
*/

//这个是web定义表结构

type UserGiftPercent struct {
	//UserId int     `xorm:"int(11) not null"`
	//GiftType int     `xorm:"int(11) not null"`
	OwnerPercent  float32 `xorm:" not null"`
	SystemPercent float32 `xorm:" not null"`
	LeaderPercent float32 `xorm:" not null"`
}

type UserGuardPercent struct {
	OwnerPercent  float32 `xorm:" not null"`
	SystemPercent float32 `xorm:" not null"`
	LeaderPercent float32 `xorm:" not null"`
}

func GetGamePercent(uid int) (r float32, has bool) {
	has = true
	res, err := orm.Query("select * from php_user_game_percent where user_id=? ", uid)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		has = false
		return
	}

	if len(res) == 0 {
		has = false
		return
	}
	c, ok := res[0]["percent"]
	if ok {
		r = common.BytesToFloat32(c)
	} else {
		has = false
	}
	return
}

func GetFamilyPercent(uid int, category int) (r UserGiftPercent, has bool) {
	has = true
	res, err := orm.Query("select * from php_user_gift_percent where user_id=? and gift_type=?", uid, category)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		has = false
		return
	}

	if len(res) == 0 {
		has = false
		return
	}

	c, ok := res[0]["owner_percent"]
	if ok {
		r.OwnerPercent = common.BytesToFloat32(c)
	} else {

		has = false
	}
	c, ok = res[0]["system_percent"]
	if ok {

		r.SystemPercent = common.BytesToFloat32(c)
	} else {

		has = false
	}
	c, ok = res[0]["leader_percent"]
	if ok {

		r.LeaderPercent = common.BytesToFloat32(c)
	} else {

		has = false
	}
	return
}

func GetGroupId(admin_id int) (ret int, bind_user_id int, bind_group_id int) {

	res2, err := orm.Query("select group_id from php_auth_group_access where uid=?", admin_id)
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}

	if len(res2) == 0 {
		ret = common.ERR_CONFIG_GROUP
		return
	}
	b, ok := res2[0]["group_id"]
	if !ok {
		ret = common.ERR_UNKNOWN
		return
	}
	bind_group_id = common.BytesToInt(b)

	if bind_group_id == 11 {
		ret = common.ERR_SUCCESS
		return
	}

	res, err := orm.Query("SELECT user_id  FROM php_user_admin_link WHERE admin_id=?", admin_id)
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}
	if len(res) == 0 {
		ret = common.ERR_CONFIG_ADMIN
		return
	}
	b, ok = res[0]["user_id"]
	if !ok {
		ret = common.ERR_UNKNOWN
		return
	}

	bind_user_id = common.BytesToInt(b)
	ret = common.ERR_SUCCESS
	return
}

func GetGuardPercent(uid int, guard_type int) (r UserGuardPercent, has bool) {
	has = true
	res, err := orm.Query("select * from php_user_guard_percent where user_id=? and guard_type=?", uid, guard_type)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		has = false
		return
	}

	if len(res) == 0 {
		has = false
		return
	}
	c, ok := res[0]["owner_percent"]
	if ok {
		r.OwnerPercent = common.BytesToFloat32(c)
	} else {

		has = false
	}
	c, ok = res[0]["system_percent"]
	if ok {
		r.SystemPercent = common.BytesToFloat32(c)
	} else {

		has = false
	}
	c, ok = res[0]["leader_percent"]
	if ok {

		r.LeaderPercent = common.BytesToFloat32(c)
	} else {

		has = false
	}
	return
}

func GetConfigGuardPrice(guard int) (price int, has bool) {
	res, err := orm.Query("select * from php_guard where id=?", guard)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		has = false
		return
	}
	c, ok := res[0]["price"]
	if ok {
		has = true
		price = common.BytesToInt(c)
	} else {
		has = false
	}
	return
}

func GetBindUser(admin_id int) (ret int, bind_user_id int) {

	res, err := orm.Query("SELECT user_id  FROM php_user_admin_link WHERE admin_id=?", admin_id)
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}
	if len(res) == 0 {
		ret = common.ERR_CONFIG_ADMIN
		return
	}
	b, ok := res[0]["user_id"]
	if !ok {
		ret = common.ERR_UNKNOWN
		return
	}
	bind_user_id = common.BytesToInt(b)
	ret = common.ERR_SUCCESS
	return
}
