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

	app.Get("/test", controller.Test)
	app.Get("/hour/hdt/list", controller.HourHdtList) //过去一小时HDT排行榜

	app.Post("/register", controller.Register)               //注册接口
	app.Post("/login", controller.Login)                     //登陆
	app.Post("/modify/pwd", controller.FindPwdByTel)         //修改密码
	app.Post("/ranking/info", controller.GetUserRankingInfo) //排行榜
	app.Post("/ranking/hdt/dig", controller.GetUseRankingHdtDig)

	app.Post("/app/list", controller.AppList)              //APP挖矿排名
	app.Post("/app/detail/info", controller.AppDetailInfo) //APP详情

	app.Post("/mine/pool/info", controller.MinePoolInfo)
	app.Post("/mine/pool/tast/list", controller.GetMinePoolTaskList)

	//短信验证
	app.Post("/send/sns", controller.QianXunSnsController) //短信验证

	http_port := conf.Cfg.MustValue("", "http_port")
	addr := fmt.Sprint(":", http_port)
	fmt.Printf("\n addr is:%s\n", addr)
	if addr == "" { //如果未能从文件里读取到数据则默认的端口为3000
		addr = ":3000"
	}

	go app.Run(iris.Addr(addr)) //iris开始监听

}
