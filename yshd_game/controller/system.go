package controller

import (
	"fmt"
	"github.com/martini-contrib/render"
	"github.com/yshd_game/model"
	"net/http"
)

func BannerList(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value["ErrCode"], ret_value["list"] = model.BannerList()
	r.JSON(http.StatusOK, ret_value)
}

//重新加载系统变量值【从数据库里读取哦】
func SystemVariableInitOrReset(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value["ErrCode"] = model.SystemVariableInitOrReset()
	fmt.Printf("\n NiuNiuPem= %f", model.NiuNiuPem)
	r.JSON(http.StatusOK, ret_value)
}

func GetToyConfig(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})

	ret_value[ServerTag], ret_value["item"] = model.GetToyInfoConfig()

	r.JSON(http.StatusOK, ret_value)
}
