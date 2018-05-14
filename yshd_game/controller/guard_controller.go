package controller

import (
	//"github.com/liudng/godump"
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
	"time"
)

type OpenGuard struct {
	Uid   int    `form:"uid" binding:"required"`
	Token string `form:"token" binding:"required"`
	Oid   int    `form:"oid" binding:"required"`
	Index int    `form:"index"`
}

//curl -d 'uid=8&oid=7&token=1' 'http://shangtv.cn:3003/guard/open'  curl -d 'uid=8&oid=7&token=1' 'http://192.168.1.12:3003/guard/open'
func OpenGuardController(r render.Render, d OpenGuard) {
	ret_value := make(map[string]interface{})

	ret_value[ServerTag] = model.OpenGuardAnchor(d.Uid, d.Oid)
	guard := model.GetGuard(d.Uid, d.Oid)
	if guard != nil {
		nowtime := time.Now().Unix()
		if guard.FinishTime > nowtime {
			ret_value["left"] = guard.FinishTime - nowtime
		} else {
			ret_value["left"] = 0
		}
	}
	u, eret := model.GetUserByUid(d.Uid)
	if eret == common.ERR_SUCCESS {
		ret_value["diamond"] = u.Diamond
	}
	r.JSON(http.StatusOK, ret_value)
}

//curl  'http://t1.shangtv.cn:3003/guard/anchor_info?uid=100008&oid=100008&token=1&index=0'
func ListAnchorGuardController(r render.Render, d OpenGuard) {
	ret_value := make(map[string]interface{})

	ret_value["guard"], ret_value[ServerTag] = model.ListAnchorGuard(d.Oid, d.Index)
	nowtime := time.Now().Unix()
	guard := model.GetGuard(d.Uid, d.Oid)
	if guard != nil {
		if guard.FinishTime > nowtime {
			ret_value["left"] = guard.FinishTime - nowtime
		} else {
			ret_value["left"] = 0
		}

	} else {
		ret_value["left"] = 0
	}
	ret_value["count"] = model.GetGuardCount(d.Oid)

	r.JSON(http.StatusOK, ret_value)
}

//curl  'http://t1.shangtv.cn:3003/guard/self_info?uid=3&oid=2&token=1&index=1'
func ListSelfGuardController(r render.Render, d OpenGuard) {
	ret_value := make(map[string]interface{})

	ret_value["guard"], ret_value[ServerTag] = model.ListSelfGuard(d.Oid, d.Uid, d.Index)
	nowtime := time.Now().Unix()
	guard := model.GetGuard(d.Uid, d.Oid)
	if guard != nil {
		if guard.FinishTime > nowtime {
			ret_value["left"] = guard.FinishTime - nowtime
		} else {
			ret_value["left"] = 0
		}

	} else {
		ret_value["left"] = 0
	}
	ret_value["count"] = model.GetGuardCount(d.Oid)

	allnum, has := model.GetConfigGuardPrice(1)
	if has == false {
		allnum = 0
	}
	ret_value["money"] = allnum
	ret_value["day"] = common.GUARD_KEEP_DAY
	r.JSON(http.StatusOK, ret_value)
}
