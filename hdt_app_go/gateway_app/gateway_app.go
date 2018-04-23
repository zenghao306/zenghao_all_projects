package main

import (
	"hdt_app_go/gateway_app/conf"
	"hdt_app_go/gateway_app/http"
	. "hdt_app_go/gateway_app/log"
	"hdt_app_go/gateway_app/rpc"
)

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func init() {

}

func main() {
	flag.Parse()

	conf.SetConfig()

	LogPath := conf.Cfg.MustValue("", "log_path")
	conf.LocalHost = conf.Cfg.MustValue("", "localhost")
	InitLog(LogPath)

	serviceName := conf.Cfg.MustValue("micro", "service_name")
	rpc.InitRpc(serviceName)

	http.NewHttp()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	sig := <-c
	Log.Infof("server is finishd sig is %v", sig)
}
