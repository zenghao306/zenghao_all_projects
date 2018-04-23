package http

import (
	"fmt"
	"github.com/kataras/iris"
	"hdt_app_go/gateway_app/conf"
	"hdt_app_go/gateway_app/http/controller"
	"hdt_app_go/gateway_app/http/session"
	//"hdt_app_go/gateway_app/rpc"
)

func NewHttp() {

	app := iris.New()
	session.InitWebScoket()

	/**
	 * @api {post} /action/upload Post Action information
	 * @apiName Action
	 * @apiGroup None
	 *
	 * @apiParam {Appid} 应用ID.
	 *
	 * @apiSuccess {String} firstname Firstname of the User.
	 * @apiSuccess {String} lastname  Lastname of the User.
	 */

	app.Get("/hour/hdt/list", controller.HourHdtList) //过去一小时HDT排行榜

	app.Post("/register", controller.Register)
	app.Post("/login", controller.Login)
	app.Post("/modify/pwd", controller.FindPwdByTel)
	app.Post("/ranking/info", controller.GetUserRankingInfo)
	app.Post("/app/list", controller.AppList)

	//短信验证
	app.Post("/send/sns", controller.QianXunSnsController)

	http_port := conf.Cfg.MustValue("", "http_port")
	addr := fmt.Sprint(":", http_port)
	fmt.Printf("\naddr:%s\n", addr)
	if addr == "" {
		addr = ":3000"
	}

	go app.Run(iris.Addr(addr))

}
