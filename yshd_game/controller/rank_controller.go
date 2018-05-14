package controller

import (
	// "github.com/liudng/godump"
	//"github.com/olahol/melody"
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
	"strconv"
)

type CommonWithOidReq struct {
	Uid   int    `form:"uid" `
	Token string `form:"token" binding:"required"`
	Index int    `form:"index"`
	Oid   int    `form:"oid" binding:"required" `
}

type CouponsRankReq struct {
	Uid    int    `form:"uid" `
	Token  string `form:"token" binding:"required"`
	Index  int    `form:"index"`
	RevUid int    `form:"rev_uid" binding:"required"`
}

type GameRankReq struct {
	Uid   int    `form:"uid" `
	Token string `form:"token" binding:"required"`
	Index int    `form:"index"`
}

//送礼排行榜
//curl 'http://shangtv.cn:3003/gift/send_rank?token=e72b101275a1b15a7ae54eeee4c3a1b3&index=0&oid=2354'
func SendGiftRank(req *http.Request, r render.Render, d CommonWithOidReq) {
	ret_value := make(map[string]interface{})

	ret_value["rank"], ret_value[ServerTag] = model.GetSendDiamonGiftRank(d.Oid, d.Index)

	r.JSON(http.StatusOK, ret_value)
}

//curl 'http://192.168.1.12:3003/gift/gain_rank?token=81863f3a4cea09644bb90a68b64bc7b5&index=0&oid=100056'
//获得礼物排行
func GainGiftRank(req *http.Request, r render.Render, d CommonWithOidReq) {
	ret_value := make(map[string]interface{})

	sum, ret := model.GetSendMoneyNum(d.Oid, common.MONEY_TYPE_DIAMOND)
	if ret == common.ERR_SUCCESS || ret == common.ERR_DB_FIND {
		ret_value["num"] = sum
	}
	ret_value["rank"], ret_value[ServerTag] = model.GetGainDiamonGiftRank(d.Oid, d.Index)
	r.JSON(http.StatusOK, ret_value)
}

//curl 'http://192.168.1.12:3000/all_rank?&index=0'
func AllGiftRank(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	index := req.FormValue("index")
	index_, _ := strconv.Atoi(index)
	ret_value["rank"], ret_value[ServerTag] = model.AllGiftRank(index_)
	r.JSON(http.StatusOK, ret_value)
}

//curl 'http://120.76.156.177:3003/gift/coupons/rank/list?token=fc36c9994a6d63519aac1caf9df8bf41&index=0&rev_uid=52'
func CouponsRankList(req *http.Request, r render.Render, d CouponsRankReq) {
	ret_value := make(map[string]interface{})
	sum, ret := model.GetSendMoneyNum(d.RevUid, common.MONEY_TYPE_DIAMOND)
	if ret == common.ERR_SUCCESS || ret == common.ERR_DB_FIND {
		ret_value["num"] = sum
	}
	ret_value["list"], ret_value[ServerTag] = model.GetCouponsRankList(d.RevUid, d.Index)
	r.JSON(http.StatusOK, ret_value)
}

func CouponsRankWeekList(r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})

	ret_value["list"] = model.GetWeekRevRank()
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func CouponsRankMonthList(r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})

	ret_value["list"] = model.GetMonthRevRank()
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func CouponsRankAllList(r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})

	ret_value["list"] = model.GetAllRevRank()
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func SendGameGiftRank(req *http.Request, r render.Render, d CommonWithOidReq) {
	ret_value := make(map[string]interface{})
	ret_value["rank"], ret_value[ServerTag] = model.GetSendGameGiftRank(d.Oid, d.Index)
	r.JSON(http.StatusOK, ret_value)
}

func GainMoonGiftRank(req *http.Request, r render.Render, d CommonWithOidReq) {
	ret_value := make(map[string]interface{})

	sum, ret := model.GetSendMoneyNum(d.Oid, common.MONEY_TYPE_SCORE)
	if ret == common.ERR_SUCCESS || ret == common.ERR_DB_FIND {
		ret_value["num"] = sum
	}
	ret_value["rank"], ret_value[ServerTag] = model.GetGainGameGiftRank(d.Oid, d.Index)
	r.JSON(http.StatusOK, ret_value)
}

func MoonRankList(req *http.Request, r render.Render, d CouponsRankReq) {
	ret_value := make(map[string]interface{})
	sum, ret := model.GetSendMoneyNum(d.RevUid, common.MONEY_TYPE_SCORE)
	if ret == common.ERR_SUCCESS || ret == common.ERR_DB_FIND {
		ret_value["num"] = sum
	}

	ret_value["list"], ret_value[ServerTag] = model.GetMoonRankList(d.RevUid, d.Index)
	r.JSON(http.StatusOK, ret_value)
}

//

func GetSendGiftRankAll(req *http.Request, r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})

	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["rank"] = model.GetAllSendRank()

	r.JSON(http.StatusOK, ret_value)

}

// t1.shangtv.cn:3003/gift/send_rank_week?uid=3&token=894235f80c32c1323e57fcd550345db2
func GetSendGiftRankWeek(req *http.Request, r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})

	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["rank"] = model.GetWeekSendRank()

	r.JSON(http.StatusOK, ret_value)

}

// 192.168.1.12:3003/gift/send_rank_month?uid=3&token=894235f80c32c1323e57fcd550345db2
func GetSendGiftRankMonth(req *http.Request, r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["rank"] = model.GetMonthSendRank()

	r.JSON(http.StatusOK, ret_value)

}

func GetGameWinScoreRankWeek(req *http.Request, r render.Render, d GameRankReq) {
	ret_value := make(map[string]interface{})

	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["rank"] = model.GetGameWeekRank()

	r.JSON(http.StatusOK, ret_value)

}

func GetGameWinScoreRankMonth(req *http.Request, r render.Render, d GameRankReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["rank"] = model.GetMonthRankOfGame()

	r.JSON(http.StatusOK, ret_value)

}

func GetGameWinScoreRankAll(req *http.Request, r render.Render, d GameRankReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["rank"] = model.GetAllRankOfGame()

	r.JSON(http.StatusOK, ret_value)

}
