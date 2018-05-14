package controller

import (
	//"github.com/liudng/godump"
	//"github.com/olahol/melody"
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
	"regexp"
	"strconv"
)

type ChatInfoReq struct {
	Rid   string `form:"rid" binding:"required"`
	Index int    `form:"index" `
}

type AudienceInfoReq struct {
	Rid   string `form:"rid" binding:"required"`
	Index int    `form:"index" `
}

func CheckChannelController(w http.ResponseWriter, r *http.Request, q render.Render) {
	ret_value := make(map[string]interface{})
	reg := regexp.MustCompile(`[0-9]+`)
	roomid := reg.FindAllString(r.URL.Path, -1)
	if len(roomid) == 0 {
		common.Log.Errf("get room panic")
		ret_value[ServerTag] = common.ERR_PARAM
		q.JSON(http.StatusOK, ret_value)
		return
	}
	roomid_ := roomid[0]
	c_room := model.GetChatRoom(roomid_)
	if c_room == nil {
		ret_value[ServerTag] = common.ERR_ROOM_EXIST
		q.JSON(http.StatusOK, ret_value)
		return
	}

	if c_room.Statue == common.ROOM_PRE_V2 || c_room.Statue == common.ROOM_ONLIVE || c_room.Statue == common.ROOM_PLAYBACK || c_room.Statue == common.ROOM_READY {
		token := r.FormValue("token")
		if token != "" {
			user, has := model.GetUserByToken(token)
			if !has {
				ret_value[ServerTag] = common.ERR_TOEKN_EXPIRE
				q.JSON(http.StatusOK, ret_value)
				return
			}

			if forbid := user.CheckPowerAccount(); forbid == true {
				ret_value[ServerTag] = common.ERR_FORBID
				q.JSON(http.StatusOK, ret_value)
				return
			}

			if sess := model.GetUserSessByUid(user.Uid); sess != nil {
				ret_value[ServerTag] = common.ERR_USER_SESS
				q.JSON(http.StatusOK, ret_value)
				return
			}
			if c_room.GetChatInfo().Uid == user.Uid {
				if c_room.Statue == common.ROOM_PLAYBACK {
					ret_value[ServerTag] = common.ERR_PLAYERBACK_OWNER
					q.JSON(http.StatusOK, ret_value)
					return
				}
			}
			ret_value[ServerTag] = common.ERR_SUCCESS
			q.JSON(http.StatusOK, ret_value)
			return
		}
	}
	ret_value[ServerTag] = common.ERR_CHAT_FINISH
	q.JSON(http.StatusOK, ret_value)

	return
}

//curl -d '' 'ws://192.168.1.12:3000/chat/join/11/ws'
func Join(w http.ResponseWriter, r *http.Request, q render.Render) {
	isWebChannel := false
	//ret_value := make(map[string]interface{})
	reg := regexp.MustCompile(`[0-9]+`)
	roomid := reg.FindAllString(r.URL.Path, -1)

	token := r.FormValue("token")
	//"channel" "web"
	uid := r.FormValue("uid")
	uid_, _ := strconv.Atoi(uid)
	channel := r.FormValue("channel")

	if channel == "web" {
		isWebChannel = true
	}

	if len(roomid) == 0 {
		common.Log.Errf("get room panic")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	roomid_ := roomid[0]
	c_room := model.GetChatRoom(roomid_)
	if c_room == nil {
		common.Log.Errf("room null panic")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if token != "" && !isWebChannel {
		user, ret := model.GetUserByUidStr(uid)
		if ret != common.ERR_SUCCESS {
			return
		} else {
			if token != user.Token {
				return
			}

			if forbid := user.CheckPowerAccount(); forbid == true {
				return
			}

			if sess := model.GetUserSessByUid(user.Uid); sess != nil {
				common.Log.Errf("sess is not nill 已存在@,uid= %d", user.Uid)
				ChecOtherkLoginIn(user.Uid)
				return
			}

			if c_room.GetChatInfo().Uid == user.Uid {

				if c_room.Statue == common.ROOM_PLAYBACK {
					return
				}
				if c_room.Statue == common.ROOM_PRE_V2 || c_room.Statue == common.ROOM_READY {
					a := model.GetChat()

					model.ReportActionDate(c_room.GetChatInfo().Uid, "view", uid_)
					a.HandleRequest(w, r)
					return
				}
			}
		}
	} else if isWebChannel {
		a := model.GetChat()
		model.ReportActionDate(c_room.GetChatInfo().Uid, "view", uid_)
		a.HandleRequest(w, r)
		return
	} else {
		return
	}

	if c_room.Statue == common.ROOM_ONLIVE || c_room.Statue == common.ROOM_PLAYBACK || c_room.Statue == common.ROOM_READY {
		//controller.Join(w, r)
		a := model.GetChat()
		model.ReportActionDate(c_room.GetChatInfo().Uid, "view", uid_)
		a.HandleRequest(w, r)
	} else {
		//common.Log.Errf("room is close")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

}

//
//curl 'http://t1.shangtv.cn:3003/chat_info?rid=100008150356878311979'
////curl 'http://shangtv.cn:3000/chat_info?rid=1222148421175495534'
func GetChatInfo(w http.ResponseWriter, req *http.Request, r render.Render, d ChatInfoReq) {
	ret_value := make(map[string]interface{})

	ret_value["info"] = model.GetChatBaseInfo(d.Rid, d.Index)
	chat := model.GetChatRoom(d.Rid)
	if chat == nil {
		ret_value["ErrCode"] = common.ERR_ROOM_EXIST
		r.JSON(http.StatusOK, ret_value)
		return
	}
	ret_value[ServerTag] = common.ERR_SUCCESS
	room := chat.GetChatInfo()
	ret_value["room"] = room

	owner, _ := model.GetUserByUid(room.Uid)
	if owner == nil {
		ret_value["ErrCode"] = common.ERR_ACCOUNT_EXIST
		r.JSON(http.StatusOK, ret_value)
		return
	}
	info := &model.UserInfo{}
	owner.GetChatUser(info)
	ret_value["owner"] = info.Chat
	ret_value["count"] = chat.GetCount()
	ret_value["vr_count"] = chat.GetVRobotCount()

	sum, _ := model.GetSendMoneyNum(room.Uid, common.MONEY_TYPE_DIAMOND)
	//sum_guard, _ := model.GetOpenGuardMoneyNum(room.Uid)
	ret_value["rice"] = sum

	sum2, ret := model.GetSendMoneyNum(room.Uid, common.MONEY_TYPE_SCORE)
	if ret == common.ERR_SUCCESS {
		ret_value["moon"] = sum2
	}

	r.JSON(http.StatusOK, ret_value)
}

//curl 'http://192.168.1.12:3000/audience_info?rid=18'
//curl 'http://shangtv.cn/audience_info?rid=18'
//curl 'http://120.76.156.177:3003/chat_info?rid=10005614992321889804'
//curl 'http://120.76.156.177:3003/audience_info?rid=7150097520846666'
//curl 'http://192.168.1.11:3003/audience_info?rid=7150097520846666'
func GetAudienceInfo(w http.ResponseWriter, req *http.Request, r render.Render, d AudienceInfoReq) {
	ret_value := make(map[string]interface{})

	ret_value["info"] = model.GetChatBaseInfo(d.Rid, d.Index)

	chat := model.GetChatRoom(d.Rid)
	if chat != nil {
		user, ok := model.GetUserByUid(chat.GetChatInfo().Uid)
		if ok != common.ERR_SUCCESS {
			ret_value[ServerTag] = ok
			return
		}
		sum, ret := model.GetSendMoneyNum(user.Uid, common.MONEY_TYPE_DIAMOND)
		if ret == common.ERR_SUCCESS {
			ret_value["rice"] = sum
		}
		sum2, ret := model.GetSendMoneyNum(user.Uid, common.MONEY_TYPE_SCORE)
		if ret == common.ERR_SUCCESS {
			ret_value["moon"] = sum2
		}
		//ret_value["rice"] = user.Coupons
		//ret_value["moon"] = user.Moon
		ret_value["count"] = chat.GetCount()
		ret_value["vr_count"] = chat.GetVRobotCount()
		ret_value[ServerTag] = common.ERR_SUCCESS
	} else {
		ret_value[ServerTag] = common.ERR_NOT_CHAT_EXIST
		ret_value["rice"] = 0
		ret_value["count"] = 0
		ret_value["moon"] = 0
	}
	//ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

type GagReq struct {
	Uid   int    `form:"uid"`   //禁言者（比如主播或超级用户）UID
	Oid   int    `form:"oid"`   //被禁言者ID
	Token string `form:"token"` //禁言者（比如主播或超级用户）TOKEN
	Type  int    `form:"type"`  //禁言类型
}

//curl -d 'token=1f59dd2bce2c73d041556ae8f85f9341&oid=5' 'http://192.168.1.12:3000/chat/gag'
//func GagUser(w http.ResponseWriter, req *http.Request, r render.Render) {
//	common.Log.Info("GagUser() called@@@@@@")
//
//	ret_value := make(map[string]interface{})
//	token := req.FormValue("token")
//	user, _ := model.GetUserByToken(token)
//	oid := req.FormValue("oid")
//	oid_, _ := strconv.Atoi(oid)
//	ret_value[ServerTag] = user.GagUser(oid_)
//	r.JSON(http.StatusOK, ret_value)
//}
func GagUser(req *http.Request, r render.Render, d GagReq) {
	ret_value := make(map[string]interface{})
	user, _ := model.GetUserByToken(d.Token)
	oid := req.FormValue("oid")
	oid_, _ := strconv.Atoi(oid)
	ret_value[ServerTag] = user.GagUser(oid_)
	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'token=1f59dd2bce2c73d041556ae8f85f9341&oid=5' 'http://192.168.1.12:3000/chat/cancel_gag'
func CancelGagUser(w http.ResponseWriter, req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	token := req.FormValue("token")
	user, _ := model.GetUserByToken(token)
	oid := req.FormValue("oid")
	oid_, _ := strconv.Atoi(oid)
	ret_value[ServerTag] = user.CancelGagUser(oid_)
	r.JSON(http.StatusOK, ret_value)
}

//  curl  'http://120.76.156.177:3003/chat/gag_status?token=a6cce661a6cafeed28caad7e6cc98a93&uid=17?rid=7150106086699649'
func GagStatusUser(r render.Render, d CommonReqWithRid) {
	ret_value := make(map[string]interface{})
	user, ret := model.GetUserByUid(d.Uid)
	if ret == common.ERR_SUCCESS {
		ret_value[ServerTag] = user.GagStatusUser(d.Rid)
	} else {
		ret_value[ServerTag] = ret
	}
	r.JSON(http.StatusOK, ret_value)
}
