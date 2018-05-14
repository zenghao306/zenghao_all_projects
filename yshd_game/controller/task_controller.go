package controller

import (
	"github.com/martini-contrib/render"
	//"github.com/martini-contrib/sessions"
	//"github.com/yshd/common"
	"github.com/yshd_game/model"
	"net/http"
	//"strconv"
	//"fmt"
	//"github.com/liudng/godump"
	//	"github.com/apiguy/go-hmacauth"
	//"github.com/martini-contrib/csrf"
	//"fmt"
	"github.com/yshd_game/common"
	"time"
	//	"go/test"
	//"github.com/liudng/godump"
)

type PostTaskListReq struct {
	Uid    int    `form:"uid"`
	Token  string `form:"token" binding:"required"`
	TaskId int64  `form:"task_id"`
}

type AcceptTaskListReq struct {
	Uid    int    `form:"uid"`
	Token  string `form:"token" binding:"required"`
	TaskId int64  `form:"task_id"`
}

//shangtv.cn:3003/task/list_task?uid=100166&token=ffcd4af88157924cc28f96e4a8985712
//192.168.1.12:3003/task/list_task?uid=23&token=52a4d0239bcfb1ea1ac3e2e198fe4f5b
func DailyTaskListController(req *http.Request, r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})

	extra, ok := model.GetUserExtraByUid(d.Uid)
	if ok {
		extra.Refresh()
	}
	ret_value["task"], ret_value["current"], ret_value["finish"], ret_value["task_num"], ret_value[ServerTag] = model.ListTask(d.Uid)
	ret_value["all_finish"] = model.CheckFinish(d.Uid)

	r.JSON(http.StatusOK, ret_value)
}

// curl -d  'uid=14&token=f6a0ef36f33ac9027fb2422b5210991a&task_id=2'  shangtv.cn:3003/task/post_task
func PostTaskController(req *http.Request, r render.Render, d PostTaskListReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = model.PostTask(d.Uid, d.TaskId)
	ret_value["all_finish"] = model.CheckFinish(d.Uid)
	_, ret_value["current"], _, _, _ = model.ListTask(d.Uid)

	y, _ := model.GetUserByUid(d.Uid)
	ret_value["score"] = y.Score
	r.JSON(http.StatusOK, ret_value)
}

func AcceptedTaskController(req *http.Request, r render.Render, d AcceptTaskListReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = model.AccecptTask(d.Uid, d.TaskId)
	r.JSON(http.StatusOK, ret_value)
}

//192.168.1.12:3003/task/refresh_task?uid=1&token=a38b8f6ae3675b6e09af53070fa779ca
//curl -d 'uid=1&token=60a65e2fcf3ca743d87d913c2c4d351f' 192.168.1.12:3003/task/refresh_task
func RefreshTaskController(req *http.Request, r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})
	u, _ := model.GetUserExtraByUid(d.Uid)
	u.Refresh()
	ret_value[ServerTag] = common.ERR_SUCCESS
	//ret_value[ServerTag] = model.AccecptTask(d.Uid,d.TaskId)
	r.JSON(http.StatusOK, ret_value)
}

//192.168.1.12:3003/daily/online_time?uid=13&token=6a53d58ce311f744a375531c51456afe
func OnlineTimeController(req *http.Request, r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})
	u, ok := model.GetUserExtraByUid(d.Uid)
	if ok {
		if u.OnlineRewardTimes >= 6 {
			ret_value[ServerTag] = common.ERR_DAILY_ONLINE_TIMES
			r.JSON(http.StatusOK, ret_value)
			return
		}
		ret_value["reward_time"] = u.OnlineRewardFinishTime
		if time.Now().Unix() < u.OnlineRewardFinishTime {
			ret_value["left"] = u.OnlineRewardFinishTime - time.Now().Unix()
		} else {
			ret_value["left"] = 0
		}
		ret_value["times"] = u.OnlineRewardTimes
		ret_value["reward_num"] = common.ONLINE_TIME_SOCRE
		ret_value["reward_time_long"] = common.ONlINE_TIME_REWARD
		//ret_value["reward_times"]=u.OnlineRewardTimes
		ret_value[ServerTag] = common.ERR_SUCCESS
	} else {
		ret_value[ServerTag] = common.ERR_UNKNOWN
	}
	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'uid=8&token=c314c270e5b886342d296ecbeab4d011' 192.168.1.12:3003/daily/online_reward
func PostOnlineRewardController(req *http.Request, r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})
	u, ok := model.GetUserExtraByUid(d.Uid)
	if ok {
		ret_value[ServerTag] = u.PostOnlineReward()
		//godump.Dump(u.OnlineRewardFinishTime)
		if u.OnlineRewardFinishTime == 0 {
			ret_value["left"] = 0
		} else if time.Now().Unix() < u.OnlineRewardFinishTime {
			ret_value["left"] = u.OnlineRewardFinishTime - time.Now().Unix()
		}

		m, _ := model.GetUserByUid(d.Uid)
		ret_value["score"] = m.Score
	} else {
		ret_value[ServerTag] = common.ERR_UNKNOWN
	}
	r.JSON(http.StatusOK, ret_value)
}
