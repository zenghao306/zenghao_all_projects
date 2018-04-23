package client

import (
	"github.com/micro/go-micro"
	. "hdt_app_go/appmeta/conf"
	proto "hdt_app_go/protcol"
)

type RegisterRPCCli struct {
	conns map[string]*proto.GreeterClient
}

func NewRegisterRPCCli() (c *RegisterRPCCli, err error) {
	name := Cfg.MustValue("kite", "service_name")
	version := Cfg.MustValue("kite", "version")

	service := micro.NewService(micro.Name("greeter.client"))
	greeter := proto.NewGreeterClient("greeter", service.Client())

}

/*
func (s *RegisterRPCCli) BannerList() (str []map[string]string, err error) {
	response, err := s.getRandCon().Tell("bannerlist")
	if err != nil {
		Log.Err(err.Error())
		return
	}


	str = make([]map[string]string, 0)


	retSlice, err := response.Slice()
	if err != nil {
		return str, err
	}

	for _, v := range retSlice {

		a := v.MustMap()
		h := make(map[string]string)
		for k2, v2 := range a {
			b := v2.MustString()
			h[k2] = b
		}

		str = append(str, h)
	}

	return str, err
}

*/
