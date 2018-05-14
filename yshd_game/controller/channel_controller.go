package controller

import (
	"github.com/martini-contrib/render"
	"github.com/yshd_game/model"
	"net/http"
)

type ChannelVersionConfigReq struct {
	Uid       int    `form:"uid" binding:"required"`
	Token     string `form:"token" binding:"required"`
	ChannelId string `form:"channel_id"`
	Version   string `form:"version"`
}

func ChannelVersionConfig(req *http.Request, r render.Render, d ChannelVersionConfigReq) {
	ret_value := make(map[string]interface{})

	ret_value["ErrCode"], ret_value["channel_version_info"] = model.ChannelVersion(d.ChannelId, d.Version)
	r.JSON(http.StatusOK, ret_value)
}
