package svr

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"golang.org/x/net/context"
	. "hdt_app_go/appmeta/conf"
	"hdt_app_go/appmeta/db"
	"hdt_app_go/common"
	proto "hdt_app_go/protcol"
	//	. "hdt_app_go/appmeta/log"
	. "hdt_app_go/appmeta/model"
	"strconv"
	"time"
)

type UserServceRpc struct {
	base *ConfigBaseData
	dao  *db.Dao
	//dao  *xorm.Engine
	//redis *redis.Client
}

func (g *UserServceRpc) LoadConfig() {

}

func (g *UserServceRpc) Login(ctx context.Context, req *proto.LoginReq, rsp *proto.LoginRes) (err error) {
	errCode, userInfo := g.dao.GetUserAccountInfo(req.Tel, req.Pwd)
	rsp.ErrCode = errCode
	rsp.UserInfo = userInfo

	if errCode == proto.ERR_OK { //更新下登陆时间
		g.dao.UpdateUserLoginTime(req.Tel)
	}
	return
}

func (g *UserServceRpc) ModifyPwdByTel(ctx context.Context, req *proto.LoginReq, rsp *proto.ReturnRes) (err error) {
	errCode := g.dao.ModifyPwdByTel(req.Tel, req.Pwd)
	rsp.ErrCode = errCode

	return
}

func (g *UserServceRpc) Register(ctx context.Context, req *proto.RegisterReq, rsp *proto.ReturnRes) (err error) {
	errCode := g.dao.CreateAccountByTel(req.Tel, req.Pwd, req.Ip, int(req.RegisterFrom))
	rsp.ErrCode = errCode
	return
}

func (g *UserServceRpc) AddQianXunCode(ctx context.Context, req *proto.QianxunReq, rsp *proto.ReturnRes) (err error) {
	errCode := g.dao.AddQianXunCode(req.Tel, req.Code)
	rsp.ErrCode = errCode
	return
}

func (g *UserServceRpc) QianXunSnsVerify(ctx context.Context, req *proto.QianxunReq, rsp *proto.ReturnRes) (err error) {
	errCode := g.dao.QianXunSnsVerify(req.Tel, req.Code)
	rsp.ErrCode = errCode
	return
}

func (g *UserServceRpc) SetUserToken(ctx context.Context, req *proto.TokenReq, rsp *proto.ReturnRes) (err error) {
	errCode := g.dao.SetUserToken(req.Tel, req.Token)
	rsp.ErrCode = errCode
	return
}

func (g *UserServceRpc) GetUserToken(ctx context.Context, req *proto.TelReq, rsp *proto.TokenRes) (err error) {
	errCode, token := g.dao.GetUserToken(req.Tel)
	if errCode == proto.ERR_OK {
		rsp.ErrCode = errCode
		rsp.Token = token
	} else { //如果上次执行错误了再执行一次，发生错误的几率接近于零
		errCode, token = g.dao.GetUserToken(req.Tel)
		rsp.ErrCode = errCode
		rsp.Token = token
	}
	return
}

func (g *UserServceRpc) GetUserRankingInfo(ctx context.Context, req *proto.TelReq, rsp *proto.RankingInfoRes) (err error) {
	res := g.dao.GetHourRankingOfHdtDig()
	i := 0
	rsp.MiningIndex = 0
	rsp.HdtMiningLast = 0
	for k, v := range res {
		i++
		fValue, _ := strconv.ParseFloat(v, 64)
		if k == req.Tel {
			rsp.MiningIndex = int32(i) //记录这家伙的排名
			rsp.HdtMiningLast = fValue //上次挖矿获取的互动币数量
			if i > 10 {
				break
			}
		}

		//var r proto.RankingInfoRes_HdtDigInfo
		//r.Tel = k
		//r.Hdt = fValue
		//
		//a := &proto.RankingInfoRes_HdtDigInfo{
		//	Tel: k,
		//	Hdt: fValue,
		//}
		//if i <= 10 { //只统计10个
		//	rsp.RankingOfHdtDig = append(rsp.RankingOfHdtDig, a) //挖矿排名，进行记录
		//}
	}

	_, rsp.HdtMiningTotal = g.dao.GetUserHdtMiningTotalByTel(req.Tel) //已挖到的总的HDT数量

	rsp.DegreeOfDifficulty = g.dao.GetHdtDegreeOfDifficulty()

	rsp.ErrCode = proto.ERR_OK
	return
}

func (g *UserServceRpc) GetUseRankingHdtDig(ctx context.Context, req *proto.TelReq, rsp *proto.RankingInfoRes) (err error) {
	res := g.dao.GetHourRankingOfHdtDig()
	i := 0
	for k, v := range res {
		i++
		fValue, _ := strconv.ParseFloat(v, 64)

		f := fmt.Sprintf("%.5f", fValue)
		fValue2 := common.ParseFloat(f)

		if k == req.Tel {
			rsp.MiningIndex = int32(i) //记录这家伙的排名
			rsp.HdtMiningLast = fValue //上次挖矿获取的互动币数量
			if i > 10 {
				break
			}
		}

		var r proto.RankingInfoRes_HdtDigInfo
		r.Tel = k
		r.Hdt = fValue2

		a := &proto.RankingInfoRes_HdtDigInfo{
			Tel: k,
			Hdt: fValue2,
		}
		if i <= 10 { //只统计10个
			rsp.RankingOfHdtDig = append(rsp.RankingOfHdtDig, a) //挖矿排名，进行记录
		}else{
			break
		}
	}

	rsp.ErrCode = proto.ERR_OK
	return
}

func NewRpcSvr() {
	options := registry.Addrs(Cfg.MustValue("etcd", "addr"))

	registry := etcdv3.NewRegistry(
		options,
	)

	serviceName := Cfg.MustValue("kite", "service_name")
	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version("1.0.0"),
		micro.Registry(registry),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
	)

	service.Init()
	s := &UserServceRpc{
		dao:  db.NewDao(),
		base: NewConfigBaseData(),
	}
	s.LoadConfig()
	proto.RegisterUserServceRpcHandler(service.Server(), s)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
