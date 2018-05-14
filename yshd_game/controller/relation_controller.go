package controller

import (
	///"github.com/liudng/godump"
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
	"strconv"
)

//curl  'http://shangtv.cn:3003/look_info?uid=56'
func LookInfo(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	uid := req.FormValue("uid")
	uid_, _ := strconv.Atoi(uid)

	info := &model.UserIndexInfo{}
	user, ret := model.GetUserByUid(uid_)
	if ret != common.ERR_SUCCESS {
		ret_value["ErrCode"] = ret
		r.JSON(http.StatusOK, ret_value)
		return
	}

	//godump.Dump(user)
	user.GetUserIndex(info)
	focus, fans := model.GetFocusCount(uid_)
	info.Focus = focus
	info.Fans = fans

	ret_value["ErrCode"] = common.ERR_SUCCESS
	ret_value["user"] = info
	r.JSON(http.StatusOK, ret_value)
}

//查看个人全部信息
//curl  'http://192.168.1.12:3003/look_info_all?oid=11&token=e5dabb1159058d63780b96ea116d1835&uid=2'
func LookInfoAll(req *http.Request, r render.Render, d CommonWithOid2Req) {
	ret_value := make(map[string]interface{})

	info := &model.UserIndexInfo{}
	other, ret := model.GetUserByUid(d.Oid)
	if ret != common.ERR_SUCCESS {
		ret_value["ErrCode"] = ret
		r.JSON(http.StatusOK, ret_value)
		return
	}
	other.GetUserIndex(info)
	focus, fans := model.GetFocusCount(d.Oid)
	info.Focus = focus
	info.Fans = fans

	sess := model.GetUserSessByUid(d.Uid)
	if sess != nil {
		gag := model.GetGagByUidAndRoomID(d.Oid, sess.Roomid)
		if gag != nil {
			info.Gag = 1
		} else {
			info.Gag = 0
		}
	} else {
		info.Gag = 0
	}

	/*
		uid := req.FormValue("uid")

		if uid == "" {
			ret_value["ErrCode"] = common.ERR_SUCCESS
			info.IsFocus = false
			ret_value["user"] = info
			r.JSON(http.StatusOK, ret_value)
			return
		}
		user, _ := model.GetUserByUidStr(uid)
	*/
	user, _ := model.GetUserByUid(d.Uid)
	ret_value["ErrCode"] = common.ERR_SUCCESS
	if user.GetFocusInfo(d.Oid) == common.ERR_SUCCESS {
		info.IsFocus = true
	} else {
		info.IsFocus = false
	}
	ret_value["user"] = info

	r.JSON(http.StatusOK, ret_value)
}
