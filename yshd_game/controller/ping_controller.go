package controller

import (
	"github.com/martini-contrib/render"
	"github.com/yshd_game/model"
	"net/http"
)

func PingReconnect(req *http.Request, r render.Render, d CommonReqOnlyUid) {
	ret_value := make(map[string]interface{})
	//ret_value["ErrCode"] = model.AnchorMgr.Ping(d.Uid)
	model.GetChat().Close()
	model.AnchorMgr.Close()
	r.JSON(http.StatusOK, ret_value)
}
