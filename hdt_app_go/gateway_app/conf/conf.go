package conf

import (
	cfg "github.com/Unknwon/goconfig"
)

var (
	Cfg *cfg.ConfigFile
)
var LocalHost = "http://192.168.1.12:3000/"

func SetConfig() {
	c, err := cfg.LoadConfigFile("./gateway_app_config.ini")
	if err != nil {
		//Cfg, err = cfg.LoadConfigFile("../config.ini")
		//Log.Panic("load ini config")
	}
	Cfg = c

}
