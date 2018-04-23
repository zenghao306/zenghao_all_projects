package conf

import (
	cfg "github.com/Unknwon/goconfig"
)

var (
	Cfg *cfg.ConfigFile
)
var LocalHost = "http://192.168.40.1:3003/"

func SetConfig() {
	c, err := cfg.LoadConfigFile("./app_config.ini")
	if err != nil {
		//Cfg, err = cfg.LoadConfigFile("../config.ini")
		//Log.Panic("load ini config")
	}
	Cfg = c

}
