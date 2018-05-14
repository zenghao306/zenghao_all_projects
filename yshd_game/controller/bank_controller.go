package controller

import (
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
)

//获取银行列表
func GetPresentBankList(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value["list"], ret_value["ErrCode"] = model.GetPresentBankList()

	r.JSON(http.StatusOK, ret_value)
}

func GetAdminList(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value["ErrCode"] = common.ERR_SUCCESS
	ret_value["admin"] = model.GetAdminList()

	r.JSON(http.StatusOK, ret_value)
}
