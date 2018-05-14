package controller

import (
	//"github.com/liudng/godump"
	"github.com/martini-contrib/render"
	"github.com/yshd_game/model"
	"net/http"
	"strconv"
)

type GoodsShowReq struct {
	Uid   int    `form:"uid" binding:"required"`
	Index int    `form:"index" `
	Token string `form:"token" binding:"required"`
}

type UserSaleGoodsListReq struct {
	Uid      int    `form:"uid" binding:"required"`
	AnchorID int    `form:"anchor_id" binding:"required"`
	Index    int    `form:"index" `
	Token    string `form:"token" binding:"required"`
}

type SelectShowClickReq struct {
	Uid   int    `form:"uid" binding:"required"`
	Token string `form:"token" binding:"required"`
	Id    int    `form:"id" `
}

//curl http://localhost:3000/goods/list?uid=1233&token=856efa5902fbfcb8277c0c16201c0dd2&index=0
func UserSaleGoodsListFunc(req *http.Request, r render.Render, d UserSaleGoodsListReq) {
	ret_value := make(map[string]interface{})

	uid := req.FormValue("anchor_id")
	uid_, _ := strconv.Atoi(uid)

	index := req.FormValue("index")
	index_, _ := strconv.Atoi(index)

	ret_value[ServerTag], ret_value["number"] = model.GetUserSaleGoodsNumbers(uid_)
	ret_value[ServerTag], ret_value["list"] = model.UserSaleGoodsList(uid_, index_)
	r.JSON(http.StatusOK, ret_value)
}

func ShowListFunc(req *http.Request, r render.Render, d GoodsShowReq) {
	ret_value := make(map[string]interface{})

	index := req.FormValue("index")
	index_, _ := strconv.Atoi(index)

	ret_value[ServerTag], ret_value["list"] = model.GetShowList(index_)
	r.JSON(http.StatusOK, ret_value)
}

func SelectShowClick(req *http.Request, r render.Render, d SelectShowClickReq) {
	ret_value := make(map[string]interface{})

	selectShowId := req.FormValue("id")
	selectShowId_, _ := strconv.Atoi(selectShowId)

	ret_value[ServerTag] = model.SelectShowClick(selectShowId_)
	r.JSON(http.StatusOK, ret_value)
}
