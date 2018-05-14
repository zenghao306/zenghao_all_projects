package controller

import (
	"github.com/martini-contrib/render"
	"github.com/yshd_game/model"
	"net/http"
)

type AddRealNameReq struct {
	Uid      int    `form:"uid" binding:"required"`
	Token    string `form:"token" binding:"required"`
	RealName string `form:"name" binding:"required"`
	RealSex  int    `form:"sex" `
	//BirthDay       string `form:"birthday"  binding:"required"`
	//IdentifyType   int    `form:"identify_type" `
	Identification string `form:"identification" binding:"required"`
	//School         string `form:"school"  `
	Tel     string `form:"tel"  binding:"required" `
	AdminId int    `form:"admin_id" `
}

func AuthRealName(req *http.Request, r render.Render, d AddRealNameReq) {
	ret_value := make(map[string]interface{})
	user, _ := model.GetUserByUid(d.Uid)
	ret_value["ErrCode"] = model.AddRealNameInfo(user, d.RealName, d.RealSex, d.Identification, d.Tel, d.AdminId)
	r.JSON(http.StatusOK, ret_value)
}
