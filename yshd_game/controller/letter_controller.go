package controller

import (
	//"github.com/olahol/melody"
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
	"strconv"
)

//写私信
//curl -d 'token=502672622e6fb0bd159bcbd5bb947aa9&oid=9&msg=uuuuu&uid=8' 'http://shangtv.cn:3003/letter/write_letter'
//curl -d 'token=c7432a1a50dc22256f029e784738800e&oid=9&msg=uuuuu&uid=8' 'http://192.168.1.12:3003/letter/write_letter'
//curl -d 'token=c7432a1a50dc22256f029e784738800e&oid=9&msg=uuuuu&uid=8' 'http://t1.shangtv.cn:3003/letter/write_letter'
func WriteLetter(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	uid := req.FormValue("uid")
	user, _ := model.GetUserByUidStr(uid)
	oid := req.FormValue("oid")
	oid_, _ := strconv.Atoi(oid)

	if oid_ == 0 {
		ret_value[ServerTag] = common.ERR_PARAM
		r.JSON(http.StatusOK, ret_value)
		return
	}
	msg := req.FormValue("msg")
	if user.AuthRealInfo == true || user.UserLevel >= 5 {
		ret := model.SendLetter(user.Uid, oid_, msg)
		if ret == common.ERR_SUCCESS {
			sess := model.GetUserSessByUid(oid_)
			if sess != nil {
				m := &model.ResponseLetterUnread{}
				m.MType = common.MESSAGE_TYPE_LETTER_UNREAD
				ret, m.Num = model.GetLetterUnreadNum(oid_)
				if ret == common.ERR_SUCCESS {
					sess.Sess.SendMsg(m)
				}
			}
			ret_value[ServerTag] = ret
		}
	} else {
		ret_value[ServerTag] = common.ERR_LETTER_OPT
	}

	r.JSON(http.StatusOK, ret_value)
}

//私信列表
//curl 'http://t1.shangtv.cn:3003/letter/letter_list?token=88353ed07092b83b9c3489d37fe84d52&uid=100&index=0'
//curl 'http://192.168.1.12:3003/letter/letter_list?token=e58d767414cc646e85681e567759d775&uid=2&index=0'
func LetterList(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	uid := req.FormValue("uid")
	user, _ := model.GetUserByUidStr(uid)
	index := req.FormValue("index")
	index_, _ := strconv.Atoi(index)

	ret_value["letter"], ret_value[ServerTag] = model.ShowLetterList(user.Uid, index_)

	r.JSON(http.StatusOK, ret_value)
}

//私信详情
//curl 'http://192.168.1.12:3000/letter/letter_open?token=15411111701&uid=122&index=0&session_id=1'
func OpenLetter(req *http.Request, r render.Render) {
	/*
		ret_value := make(map[string]interface{})

		session_id := req.FormValue("session_id")
		session_id_, _ := strconv.Atoi(session_id)

		index := req.FormValue("index")
		index_, _ := strconv.Atoi(index)
		ret_value["letter"], ret_value[ServerTag] = model.ShowLetterDeatil(session_id_, index_, user.Uid, 1)

		r.JSON(http.StatusOK, ret_value)
	*/
}

//curl 'http://120.76.156.177:3003/letter/letter_show?token=528945af676a234cb4769d785dd4d1ed&uid=4&index=0&oid=3'
//curl 'http://192.168.1.12:3003/letter/letter_show?token=528945af676a234cb4769d785dd4d1ed&uid=4&index=0&oid=3'

func ShowLetter(req *http.Request, r render.Render) {

	ret_value := make(map[string]interface{})
	oid := req.FormValue("oid")
	uid := req.FormValue("uid")
	index := req.FormValue("index")

	index_, _ := strconv.Atoi(index)
	oid_, _ := strconv.Atoi(oid)
	user, _ := model.GetUserByUidStr(uid)

	if oid_ == 0 {
		ret_value[ServerTag] = common.ERR_PARAM
		r.JSON(http.StatusOK, ret_value)
		return
	}
	session_id, has := model.GetLetterSessionId(user.Uid, oid_)
	if has {
		letter := make([]model.LetterMsg, 0)
		letter, ret_value[ServerTag] = model.ShowLetterDeatil(session_id, index_, user.Uid, oid_)
		letter_resq := make([]model.LetterMsgResq, 0)

		for _, v := range letter {
			var resqtemp model.LetterMsgResq
			resqtemp.SetByOut(&v)
			letter_resq = append(letter_resq, resqtemp)
		}
		ret_value["letter"] = letter_resq
	} else {
		ret_value[ServerTag] = common.ERR_SUCCESS
	}
	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'token=c0805fb1372e31d644a0eeed2d862f92&uid=6&oid=366'  'http://192.168.1.12:3000/letter/del_session'
func DelSession(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	oid := req.FormValue("oid")
	uid := req.FormValue("uid")
	oid_, _ := strconv.Atoi(oid)
	user, _ := model.GetUserByUidStr(uid)
	ret_value[ServerTag] = user.DelSessionById(oid_)
	r.JSON(http.StatusOK, ret_value)
}

//curl  'http://192.168.1.12:3003/letter/letter_unread?token=6c0c1960bc9919af9b13917e393fd569&uid=10&oid=9'
func UnreadLetter(req *http.Request, r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag], ret_value["unread"] = model.GetLetterUnreadNum(d.Uid)
	r.JSON(http.StatusOK, ret_value)
}
