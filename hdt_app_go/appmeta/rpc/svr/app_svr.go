package svr

import (
	"golang.org/x/net/context"
	proto "hdt_app_go/protcol"
)

func (g *UserServceRpc) AppList(ctx context.Context, req *proto.IndexReq, rsp *proto.AppListRes) (err error) {
	rsp.ErrCode, rsp.Applist = g.dao.GetAPPIconNameList(req.Index)
	return
}

func (g *UserServceRpc) AppDetailInfo(ctx context.Context, req *proto.AppDetailReq, rsp *proto.AppDetailInfoRes) (err error) {
	var errCode2, errCode3, errCode4 int32
	rsp.ErrCode, rsp.UserAppHdt = g.dao.GetUserAppHdt(req.AppId, req.Tel) //获取用户在该平台挖到的HDT数量

	errCode2, rsp.AppHdtTotal = g.dao.GetAppHdtTotal(req.AppId) //获取APP开发者投放的HDT
	if errCode2 != proto.ERR_OK {
		rsp.ErrCode = errCode2
		return
	}

	errCode3, rsp.AppContent, rsp.IosAddress, rsp.AndroidAddress = g.dao.GetAppContent(req.AppId) //获取app内容、IOS下载地址、Android下载地址
	if errCode3 != proto.ERR_OK {
		rsp.ErrCode = errCode3
		return
	}

	errCode4, rsp.AppImg = g.dao.GeAPPImageList(req.AppId) //获取APP图片
	if errCode3 != proto.ERR_OK {
		rsp.ErrCode = errCode4
	}

	return
}

func (g *UserServceRpc) GetMinePoolInfo(ctx context.Context, req *proto.TelReq, rsp *proto.MinePoolRes) (err error) {
	rsp.DegreeOfDifficulty = g.dao.GetHdtDegreeOfDifficulty()

	rsp.HdtSupplyLimit, rsp.HdtTotalSupply, _ = g.dao.GetMinedInfo()

	rsp.ErrCode, rsp.AppHdtBalanceTotal = g.dao.GetAppHdtBalanceTotal()

	return
}

//rpc MinePoolTaskList(TokenReq) returns (MinePoolTaskListRes) {}
func (g *UserServceRpc) GetMinePoolTaskList(ctx context.Context, req *proto.TokenReq, rsp *proto.MinePoolTaskListRes) (err error) {
	//lists := make([]*proto.MinePoolTaskListRes_MinePoolTask, 0)
	rsp.ErrCode, rsp.MinePoolTasklist = g.dao.GetMinePoolTaskReleaseList()

	list2 := g.dao.GetMinePoolTaskDigInfo()
	for _, v := range list2 {
		temp := &proto.MinePoolTaskListRes_MinePoolTask{}
		temp.HdtTaskBalance = v.HdtTaskBalance
		temp.Hdt = v.Hdt
		temp.Time = v.Time
		temp.AppId = v.AppId
		temp.AppIcoPath = v.AppIcoPath
		temp.Style = v.Style
		temp.AppName = v.AppName
		rsp.MinePoolTasklist = append(rsp.MinePoolTasklist, temp)
	}

	return
}
