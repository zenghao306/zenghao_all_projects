package client

import (
	//"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"golang.org/x/net/context"
	//. "hdt_app_go/gateway_app/client_msg"
	. "hdt_app_go/gateway_app/conf"
	. "hdt_app_go/gateway_app/log"
	//. "hdt_app_go/gateway_app/model"
	proto "hdt_app_go/protcol"
	//"github.com/astaxie/beego/context/param"
)

type RegisterRPCCli struct {
	conn proto.UserServceRpcService
}

func NewRegisterRPCCli(serviceName string) (c *RegisterRPCCli, err error) {

	options := registry.Addrs(Cfg.MustValue("etcd", "addr"))

	registry := etcdv3.NewRegistry(
		options,
	)

	service := micro.NewService(
		micro.Name("service.client"),
		micro.Registry(registry),
		//micro.WrapClient(wrapper),

	)
	service.Init()

	c = new(RegisterRPCCli)
	c.conn = proto.NewUserServceRpcService(serviceName, service.Client())

	return
}

//(req.Tel, req.Code)
func (s *RegisterRPCCli) QianXunSnsVerify(tel, code string) int32 {
	v := &proto.QianxunReq{}
	v.Tel = tel
	v.Code = code

	rsp, err := s.conn.QianXunSnsVerify(context.TODO(), v)
	if err != nil {
		Log.Println(err)
		return 500
	}

	return rsp.ErrCode
}

func (s *RegisterRPCCli) AddQianXunCode(param *proto.QianxunReq) int32 {
	rsp, err := s.conn.AddQianXunCode(context.TODO(), param)
	if err != nil {
		Log.Println(err)
		return 500
	}

	return rsp.ErrCode
}

func (s *RegisterRPCCli) CreateAccountByTel(tel, pwd, ip string, registerFrom int32) int32 {
	v := &proto.RegisterReq{}
	v.Tel = tel
	v.Pwd = pwd
	v.Ip = ip
	v.RegisterFrom = registerFrom

	rsp, err := s.conn.Register(context.TODO(), v)
	if err != nil {
		Log.Println(err)
		return 500
	}

	return rsp.ErrCode
}

func (s *RegisterRPCCli) SetUserToken(tel, token string) int32 {
	v := &proto.TokenReq{}
	v.Tel = tel
	v.Token = token

	rsp, err := s.conn.SetUserToken(context.TODO(), v)
	if err != nil {
		Log.Println(err)
		return 500
	}

	return rsp.ErrCode
}

func (s *RegisterRPCCli) GetUserToken(tel string) (int32, string) {
	v := &proto.TelReq{}
	v.Tel = tel

	rsp, err := s.conn.GetUserToken(context.TODO(), v)
	if err != nil {
		Log.Println(err)
		return 500, ""
	}

	return rsp.ErrCode, rsp.Token
}

func (s *RegisterRPCCli) Login(tel, pwd string) (int32, *proto.UserInfo) {
	v := &proto.LoginReq{}
	v.Tel = tel
	v.Pwd = pwd

	rsp, err := s.conn.Login(context.TODO(), v)
	if err != nil {
		Log.Println(err)
		return 500, nil
	}

	return rsp.ErrCode, rsp.UserInfo
}

func (s *RegisterRPCCli) ModifyPwdByTel(tel, pwd string) int32 {
	v := &proto.LoginReq{}
	v.Tel = tel
	v.Pwd = pwd

	rsp, err := s.conn.ModifyPwdByTel(context.TODO(), v)
	if err != nil {
		Log.Println(err)
		return 500
	}

	return rsp.ErrCode
}

func (s *RegisterRPCCli) GetUserRankingInfo(tel string) (int32, *proto.RankingInfoRes) {
	v := &proto.TelReq{}
	v.Tel = tel

	rsp, err := s.conn.GetUserRankingInfo(context.TODO(), v)
	if err != nil {
		Log.Println(err)
		return 500, nil
	}

	return rsp.ErrCode, rsp
}

func (s *RegisterRPCCli) GetUseRankingHdtDig(tel string) (int32, *proto.RankingInfoRes) {
	v := &proto.TelReq{}
	v.Tel = tel

	rsp, err := s.conn.GetUseRankingHdtDig(context.TODO(), v)
	if err != nil {
		Log.Println(err)
		return 500, nil
	}

	return rsp.ErrCode, rsp
}

func (s *RegisterRPCCli) GetMinePoolInfo(tel string) (int32, *proto.MinePoolRes) {
	v := &proto.TelReq{}
	v.Tel = tel

	rsp, err := s.conn.GetMinePoolInfo(context.TODO(), v)
	if err != nil {
		Log.Println(err)
		return 500, nil
	}

	return rsp.ErrCode, rsp
}

func (s *RegisterRPCCli) GetMinePoolTaskList(token string) (int32, *proto.MinePoolTaskListRes) {
	v := &proto.TokenReq{}
	rsp, err := s.conn.GetMinePoolTaskList(context.TODO(), v)
	if err != nil {
		Log.Println(err)
		return 500, nil
	}

	return rsp.ErrCode, rsp
}

func (s *RegisterRPCCli) AppList(index int) (int32, []*proto.AppListRes_AppNameIcon) {
	v := &proto.IndexReq{}
	v.Index = int32(index)

	rsp, err := s.conn.AppList(context.TODO(), v)
	if err != nil {
		Log.Println(err)
		return 500, nil
	}

	return rsp.ErrCode, rsp.Applist
}

func (s *RegisterRPCCli) AppDetailInfo(tel string,appId int64) (int32, *proto.AppDetailInfoRes) {
	v := &proto.AppDetailReq{}
	v.AppId = appId
	v.Tel = tel
	rsp, err := s.conn.AppDetailInfo(context.TODO(), v)
	if err != nil {
		Log.Println(err)
		return 500, nil
	}

	return rsp.ErrCode, rsp
}

//func (s *RegisterRPCCli) GetUser7DaysHdtList(param *proto.AccessCodeReq) (errcode int32, t []*proto.UidHdtListRes_HdtPerDay) {
//	v := &proto.AccessCodeReq{}
//	v.Appid = param.Appid
//	v.Uid = param.Uid
//	rsp, err := s.conn.GetUser7DaysHdtList(context.TODO(), v)
//	if err != nil {
//		Log.Println(err)
//		fmt.Println(err.Error())
//		return 500, nil
//	}
//
//	return rsp.ErrCode, rsp.Hdt
//}
//
////GetUserHdtTotal(ctx context.Context, in *AccessCodeReq, opts ...client.CallOption) (*AppUserHdtRes, error)
//func (s *RegisterRPCCli) GetUserHdtTotal(param *proto.AccessCodeReq) (errcode uint32, f float64) {
//	v := &proto.AccessCodeReq{}
//	v.Appid = param.Appid
//	v.Uid = param.Uid
//	rsp, err := s.conn.GetUserHdtTotal(context.TODO(), v)
//	if err != nil {
//		Log.Println(err)
//		fmt.Println(err.Error())
//		return 500, 0
//	}
//
//	return rsp.ErrCode, rsp.HdtNumber
//}
