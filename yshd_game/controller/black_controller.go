package controller

import (
	"github.com/martini-contrib/render"
	"github.com/yshd_game/model"
	"net/http"
)

type AddBlackReq struct {
	Uid     int    `form:"uid" binding:"required"`
	Token   string `form:"token" binding:"required"`
	BlackId int    `form:"blackid" binding:"required"`
}

type DelBlackReq struct {
	Uid     int    `form:"uid" binding:"required"`
	Token   string `form:"token" binding:"required"`
	BlackId int    `form:"blackid" binding:"required"`
}

type BlackListReq struct {
	Uid   int    `form:"uid" binding:"required"`
	Token string `form:"token" binding:"required"`
	Index int    `form:"index" `
}

//curl -d 'token=15492786411&blackid=113&uid=328' 'http://192.168.1.12:3000/black/add_black'
func AddBlackFunc(req *http.Request, r render.Render, d AddBlackReq) {
	ret_value := make(map[string]interface{})
	/*
		uid := req.FormValue("uid")
		uid_, _ := strconv.Atoi(uid)

		blackid := req.FormValue("blackid")
		blackid_, _ := strconv.Atoi(blackid)
	*/

	ret_value[ServerTag] = model.AddBlack(d.Uid, d.BlackId)

	//ret_value[ServerTag] = model.AddBlack(uid_, blackid_)
	r.JSON(http.StatusOK, ret_value)
}

func DelBlackFunc(req *http.Request, r render.Render, d DelBlackReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = model.DelBlack(d.Uid, d.BlackId)
	r.JSON(http.StatusOK, ret_value)
}

//curl 'http://192.168.1.12:3000/black/black_list?token=d9770f85bf111493593c359dc92d2c97&index=0&uid=5'
func BlackListFunc(req *http.Request, r render.Render, d BlackListReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag], ret_value["black"] = model.BlackList(d.Uid, d.Index)
	r.JSON(http.StatusOK, ret_value)
}
