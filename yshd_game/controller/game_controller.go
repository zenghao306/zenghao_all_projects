package controller

import (
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
)

type UserGameRaiseReq struct {
	CheckUid  int `form:"check_id" binding:"required"`
	TimeStart int `form:"time_start" binding:"required"`
}

func UserGameRaise(req *http.Request, r render.Render, d UserGameRaiseReq) {
	ret_value := make(map[string]interface{})

	ret_value["userGain"] = common.ERR_SUCCESS
	ret_value["userGain"], ret_value["userNotGain"], ret_value["totalraise"] = model.CheckUserBetReward(d.CheckUid, d.TimeStart)

	r.JSON(http.StatusOK, ret_value)
}
