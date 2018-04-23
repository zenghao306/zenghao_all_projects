package main

import (
	"flag"
	"hdt_app_go/appmeta/conf"
	. "hdt_app_go/appmeta/log"
	"hdt_app_go/appmeta/rpc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()

	conf.SetConfig()

	LogPath := conf.Cfg.MustValue("", "log_path")
	conf.LocalHost = conf.Cfg.MustValue("", "localhost")
	InitLog(LogPath)
	//db.Orm=db.SetEngine()
	rpc.InitRpc()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	sig := <-c
	Log.Infof("server is finishd sig is %v", sig)
}
