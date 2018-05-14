package controller

import (
	//"github.com/liudng/godump"
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
)

type CancleFocusReq struct {
	Uid   int    `form:"uid" `
	Token string `form:"token" binding:"required"`
	Oid   int    `form:"oid" binding:"required"`
}

type LiveFocusReq struct {
	Uid        int    `form:"uid" `
	Token      string `form:"token" binding:"required"`
	Index      int    `form:"index" `
	AppVersion string `form:"app_version"`
	Os         int    `form:"os"`
}

type FocusOtherReq struct {
	Uid   int    `form:"uid" `
	Token string `form:"token" binding:"required"`
	Oid   int    `form:"oid" binding:"required"`
}

type FansListReq struct {
	Uid   int    `form:"uid" `
	Token string `form:"token" binding:"required"`
	Index int    `form:"index" `
}

type FocusInfoReq struct {
	Uid   int    `form:"uid" `
	Token string `form:"token" binding:"required"`
	Oid   int    `form:"oid" binding:"required"`
}

type FocusLiveListReq struct {
	Uid        int    `form:"uid" `
	Token      string `form:"token" binding:"required"`
	Index      int    `form:"index" `
	AppVersion string `form:"app_version"`
	Os         int    `form:"os"`
}

//curl -d 'token=15492786411&oid=5' 'http://192.168.1.12:3000/focus/focus_other'
func FocusOther(req *http.Request, r render.Render, d FocusOtherReq) {
	ret_value := make(map[string]interface{})
	//token := req.FormValue("token")
	//user, _ := model.GetUserByToken(token)

	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	//oid := req.FormValue("oid")
	//oid_, _ := strconv.Atoi(oid)
	_, ret := model.GetUserByUid(d.Oid)
	if ret == common.ERR_SUCCESS {
		ret_value[ServerTag] = user.FocusOtherPublic(d.Oid)
	} else {
		ret_value[ServerTag] = common.ERR_ACCOUNT_EXIST
	}

	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'token=477387272886999058c9acf3ddc29f19&oid=3' 'http://shangtv.cn:3003/focus/cancel_focus'
func CancleFocus(req *http.Request, r render.Render, d CancleFocusReq) {
	ret_value := make(map[string]interface{})
	//user, _ := model.GetUserByToken(d.Token)

	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}
	ret_value[ServerTag] = user.CancleFocusPublic(d.Oid)
	r.JSON(http.StatusOK, ret_value)
}

//curl 'http://192.168.1.12:3003/focus/focus_list??uid=100008&token=6a41050adeb870bb58fab4cc8fa06eff&index=0&app_version=1.1'
//curl 'http://t1.shangtv.cn:3003/focus/focus_list?uid=100008&token=6a41050adeb870bb58fab4cc8fa06eff&index=0&app_version=1.1'
func GetFocusList(req *http.Request, r render.Render, d LiveFocusReq) {
	ret_value := make(map[string]interface{})
	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}
	//user, _ := model.GetUserByToken(d.Token)

	v := make([]model.OutUserInfo2, 0)

	if d.AppVersion != "" {
		v, ret_value[ServerTag] = model.GetFocusListWithCache(user.Uid, d.Index)
		if len(v) == 0 {
			ret_value["focus"] = make([]model.OutUserInfo2, 0)
		} else {
			ret_value["focus"] = v
		}

		//user.GetFocusList(d.Index)
		ret_value["count"], _ = model.GetFocusCount(user.Uid)
		r.JSON(http.StatusOK, ret_value)
	} else {
		v, ret_value[ServerTag] = model.GetFocusListWithCache2(user.Uid, d.Index)
		if len(v) == 0 {
			ret_value["focus"] = make([]model.OutUserInfo2, 0)
		} else {
			ret_value["focus"] = v
		}

		//user.GetFocusList(d.Index)
		ret_value["count"], _ = model.GetFocusCount(user.Uid)
		r.JSON(http.StatusOK, ret_value)
	}
}

//curl 'shangtv.cn:3003/focus/focus_live_list?uid=442404&token=6ad67d19a6c96ed1dd08ae68f7763fca&index=0&app_version=1.1'
func GeFocustLiveList(req *http.Request, r render.Render, d FocusLiveListReq) {
	ret_value := make(map[string]interface{})
	//user, _ := model.GetUserByToken(d.Token)

	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	//v, ret := model.GetLiveFocusWithCache(user.Uid, d.Index)
	if d.AppVersion != "" {
		v := make([]model.OutUserInfo, 0)
		v, ret := model.GetLiveFocusWithCache(user.Uid, d.Index)

		if len(v) == 0 {
			ret_value["focus"] = make([]model.OutUserInfo, 0)
		} else {
			ret_value["focus"] = v
		}

		ret_value[ServerTag] = ret
		ret_value["count"], _ = model.GetFocusCount(user.Uid)
		r.JSON(http.StatusOK, ret_value)
	} else {
		v := make([]model.OutUserInfo, 0)
		v, ret := model.GetLiveFocusWithCache2(user.Uid, d.Index)
		if len(v) == 0 {
			ret_value["focus"] = make([]model.OutUserInfo, 0)
		} else {
			ret_value["focus"] = v
		}

		ret_value[ServerTag] = ret
		ret_value["count"], _ = model.GetFocusCount(user.Uid)
		r.JSON(http.StatusOK, ret_value)
	}

}

func RecommendLiveList(req *http.Request, r render.Render, d FocusLiveListReq) {
	ret_value := make(map[string]interface{})

	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	v := make([]model.OutUserInfo, 0)
	v, ret := model.GetLiveNotFocusWithCache(user.Uid, d.Index)

	if len(v) == 0 {
		ret_value["recommend"] = make([]model.OutUserInfo, 0)
	} else {
		ret_value["recommend"] = v
	}

	ret_value[ServerTag] = ret
	ret_value["count"] = model.GetNotFocusLiveCount(user.Uid)
	r.JSON(http.StatusOK, ret_value)
}

//curl 'http://192.168.1.12:3000/focus/fans_list?token=15492786411'
func GetFansList(req *http.Request, r render.Render, d FansListReq) {
	ret_value := make(map[string]interface{})

	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	//token := req.FormValue("token")
	//user, _ := model.GetUserByToken(token)
	//index := req.FormValue("index")
	//index_, _ := strconv.Atoi(index)
	ret_value["fans"], ret_value[ServerTag] = user.GetFansList(d.Index)
	_, ret_value["count"] = model.GetFocusCount(user.Uid)
	r.JSON(http.StatusOK, ret_value)
}

//curl 'http://192.168.1.12:3000/focus/search_focus?token=15492786411&oid=5'
func GetFocusInfo(req *http.Request, r render.Render, d FocusInfoReq) {
	ret_value := make(map[string]interface{})
	//token := req.FormValue("token")
	//user, _ := model.GetUserByToken(token)

	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	//oid := req.FormValue("oid")

	//oid_, _ := strconv.Atoi(oid)
	ret_value[ServerTag] = user.GetFocusInfo(d.Oid)
	r.JSON(http.StatusOK, ret_value)
}
