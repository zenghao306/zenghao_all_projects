package model

import (
	//	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	//"net/http"
	"time"
	//"encoding/json"
)

func MonitorClose(uid int, close_type int) int {
	sess := GetUserSessByUid(uid)
	if sess == nil {
		ForbidUser(uid, close_type)
		return common.ERR_NOT_CHAT_EXIST
	}
	chat := GetChatRoom(sess.Roomid)
	if chat == nil {
		return common.ERR_NOT_CHAT_EXIST
	}

	if chat.room.Uid != uid {
		return common.ERR_CLOSE_ROOM
	}
	res := &ResponseSys{MType: common.MESSAGE_TYPE_ADMIN, Notice: "主播账号违规现在已经被关闭如有问题联系管理员"}

	//sess.Sess.SendMsg(res)

	//time.Sleep(5)
	err := sess.Sess.CloseWithMsgAndJson(res)
	if err != nil {
		common.Log.Errf("err is %s", err.Error())
	}

	return ForbidUser(uid, close_type)

}

func ChartRoomClose(uid int) int {

	sess := GetUserSessByUid(uid)
	if sess == nil {
		return common.ERR_NOT_CHAT_EXIST
	}
	chat := GetChatRoom(sess.Roomid)
	if chat == nil {
		return common.ERR_NOT_CHAT_EXIST
	}

	if chat.room.Uid != uid {
		return common.ERR_CLOSE_ROOM
	}
	res := &ResponseSys{MType: common.MESSAGE_TYPE_ADMIN, Notice: "房间已被管理员关闭"}

	sess.Sess.CloseWithMsgAndJson(res)
	//SendMsgToUserWithClose(chat.room.Uid, res)

	//uid := chat.room.Uid
	//sess := GetUserSessByUid(uid)
	//	sess.CloseDirect()
	return common.ERR_SUCCESS

}

func ForbidUser(uid int, close_type int) int {
	//user, _ := GetUserByUid(uid)
	var keep_time int64
	switch close_type {
	case 1:
		keep_time = common.FORBID_ACCOUNT_KEEP_TIME
	case 2:
		keep_time = 24 * 3600
	case 3:
		keep_time = 1
	case 4:
		keep_time = 123600
	default:
		keep_time = common.FORBID_ACCOUNT_KEEP_TIME
	}
	a := time.Now().Unix() + keep_time
	//godump.Dump(a)
	res, err := orm.Exec("update go_user set forbid=1 ,forbid_time=? where uid=?", a, uid)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	aff_row, err := res.RowsAffected()
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	return common.ERR_SUCCESS
}

func ForbidUserPower(uid int, forbid int, forbid_time int) int {
	if forbid < common.FORBID_POWERS_NONE && forbid > common.FORBID_POWER_MAX {
		return common.ERR_PARAM
	}

	if forbid == common.FORBID_POWERS_FOREVER {

		res, err := orm.Exec("update go_user set forbid_powers=? where uid=?", common.FORBID_POWERS_FOREVER, uid)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}

		aff_row, err := res.RowsAffected()
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff_row == 0 {
			return common.ERR_DB_UPDATE
		}
	} else if forbid == common.FORBID_POWERS_TIME {

		add_time := time.Now().Unix() + int64(forbid_time)
		res, err := orm.Exec("update go_user set forbid_powers=? and forbid_powers_time=? where uid=?", common.FORBID_POWERS_TIME, add_time, uid)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}

		aff_row, err := res.RowsAffected()
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff_row == 0 {
			return common.ERR_DB_UPDATE
		}
	} else if forbid == common.FORBID_POWERS_NONE {

		res, err := orm.Exec("update go_user set forbid_powers=? and forbid_powers_time=? where uid=?", common.FORBID_POWERS_NONE, 0, uid)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}

		aff_row, err := res.RowsAffected()
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff_row == 0 {
			return common.ERR_DB_UPDATE
		}
		return common.ERR_SUCCESS
	} else {
		return common.ERR_PARAM
	}

	sess := GetUserSessByUid(uid)

	if sess != nil {
		//chat := GetChatRoom(sess.Roomid)

		//if chat != nil {
		res := &ResponseSys{MType: common.MESSAGE_TYPE_ADMIN_FORBID, Notice: "你的账号违规现在已经被关闭如有问题联系管理员"}
		//SendMsgToUserWithClose(uid, res)
		sess.Sess.CloseWithMsgAndJson(res)
		//}
	}
	return common.ERR_SUCCESS
}
