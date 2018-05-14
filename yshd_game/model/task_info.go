package model

import (
	"github.com/yshd_game/common"
	"github.com/yshd_game/confdata"
	//"strconv"
	"time"
)

type TaskInfo struct {
	Id     int64
	Uid    int   `xorm:"int(11) not null UNIQUE(TASK)"` //用户ID
	TaskId int64 `xorm:"int(11) not null UNIQUE(TASK)"` //任务ID
	//Name          string
	Status     uint8 //任务状态 1未完成 2已完成 3已领取奖励
	BeginTime  int64 //接受任务时间
	FinishTime int64 //完成时间
	//TargetType    int32 //目标类型 0赢游戏 1打赏
	//TargetNum     int32 //目标数量
	TargetCurrent int32 //当前进度
}

type TaskInfoSend struct {
	Id         int64
	Uid        int   //用户ID
	TaskId     int64 //任务ID
	Status     uint8 //任务状态 1未完成 2已完成 3已领取奖励
	BeginTime  int64 //接受任务时间
	FinishTime int64 //完成时间
	//TargetNum     int32 //目标数量
	TargetCurrent int32 //当前进度
}

//接受任务
func AccecptTask(uid int, taskid int64) int {
	m := &TaskInfo{}
	//entry, ok := confdata.ConfigData.TaskById[taskid]
	//if ok {
	m.Uid = uid
	m.Status = common.TASK_STATUS_ACCECPT
	m.TaskId = taskid
	//m.TargetType = int32(entry.Target)
	m.BeginTime = time.Now().Unix()
	//m.Name = entry.Name
	//m.TargetNum = entry.TargetParam
	//}

	aff_row, err := orm.Insert(m)
	if err != nil {
		common.Log.Errf("orm is err %s", err.Error())
		return common.ERR_UNKNOWN
	}

	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	return common.ERR_SUCCESS
}

//触发任务
func TriggerTask(target confdata.TargetType, uid int, args int) {
	//m := &TaskInfo{}
	m := make(map[int]*TaskInfo)
	err := orm.Where("uid=? and status=?", uid, common.TASK_STATUS_ACCECPT).Find(m)
	if err != nil {
		common.Log.Errf("orm is err %s", err.Error())
		return
	}
	flag := false

	for _, v := range m {
		if v.Status == common.TASK_STATUS_ACCECPT {
			entry, ok := confdata.ConfigData.TaskById[v.TaskId]

			if ok {

				if entry.Target != target {
					continue
				}
				v.TargetCurrent += int32(args)
				if entry.TargetParam == v.TargetCurrent {
					v.Status = common.TASK_STATUS_POST
					aff_row, err := orm.Where("uid=? and task_id=?", uid, v.TaskId).Update(v)
					if err != nil {
						common.Log.Errf("orm is err %s", err.Error())
						return
					}

					if aff_row == 0 {
						return
					}
					flag = true
				} else if entry.TargetParam > v.TargetCurrent {

					aff_row, err := orm.Where("uid=? and task_id=?", uid, v.TaskId).Update(v)
					if err != nil {
						common.Log.Errf("orm is err %s", err.Error())
						return
					}
					if aff_row == 0 {
						return
					}
					flag = true

				}
			}
		}
	}

	if flag == true {
		sess := GetUserSessByUid(uid)
		if sess != nil {
			var rep ResponseTaskInfo
			rep.Tasks = make([]TaskInfoSend, 0)
			rep.Tasks, rep.Current, rep.Finish, rep.TaskNum, _ = ListTask(uid)
			rep.AllFinish = CheckFinish(uid)
			rep.MType = common.MESSAGE_TYPE_TASK_INFO
			SendMsgToUser(uid, rep)
		}
	}

}

//提交任务
func PostTask(uid int, taskid int64) int {
	m := &TaskInfo{}
	has, err := orm.Where("uid=? and task_id=? ", uid, taskid).Get(m)
	if err != nil {
		common.Log.Errf("orm is err %s", err.Error())
		return common.ERR_UNKNOWN
	}

	if !has {
		return common.ERR_TASK_NONE
	}
	u, ret := GetUserByUid(uid)

	if ret == common.ERR_SUCCESS {
		if m.Status == common.TASK_STATUS_POST {
			entry, ok := confdata.ConfigData.TaskById[taskid]
			if ok {
				if entry.Reward1.MoneyType != 0 {
					//u.AddMoney(entry.Reward1.MoneyType, entry.Reward1.Num)
					ret := u.AddMoney(nil, entry.Reward1.MoneyType, entry.Reward1.Num, false)
					if ret != common.ERR_SUCCESS {
						return ret
					}
				} else {
					common.Log.Errf("config err reward type is %s", entry.Reward1.MoneyType)
					return common.ERR_UNKNOWN
				}

				/*
					if entry.Reward2 != 0 {
						u.AddUserExp(int(entry.Reward2))
					}
				*/
				m.Status = common.TASK_STATUS_FINISH
				_, err = orm.Where("uid=? and task_id=?", uid, taskid).Update(m)
				if err != nil {
					common.Log.Errf("orm is err %s", err.Error())
					return common.ERR_UNKNOWN
				}
				return common.ERR_SUCCESS
			}
		} else {
			return common.ERR_TASK_STATUS
		}

	}
	return common.ERR_ACCOUNT_EXIST
}

//任务列表
func ListTask(uid int) (m []TaskInfoSend, current int, finish int, max int, ret int) {

	l := make([]TaskInfo, 0)
	err := orm.Where("uid=?", uid).Find(&l)
	if err != nil {
		common.Log.Errf("orm is err %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}

	if len(l) == 0 {
		ret = common.ERR_TASK_NONE
	} else {
		m = make([]TaskInfoSend, len(l))
		max = len(l)
		for k, v := range l {
			if v.Status == common.TASK_STATUS_POST {
				current++
			} else if v.Status == common.TASK_STATUS_FINISH {
				finish++
			}
			var s TaskInfoSend
			s.Id = v.Id
			s.Uid = v.Uid
			s.Status = v.Status
			s.TargetCurrent = v.TargetCurrent
			s.TaskId = v.TaskId
			s.BeginTime = v.BeginTime
			s.FinishTime = v.FinishTime
			m[k] = s
		}
		ret = common.ERR_SUCCESS
	}
	finish = finish + current
	return
}

func DelTask(uid int, taskid int64) int {
	if taskid == 0 {
		return common.ERR_PARAM
	}
	aff_row, err := orm.Where("uid=? and task_id=?", uid, taskid).Delete(&TaskInfo{})
	if err != nil {
		common.Log.Errf("orm is err %s", err.Error())

		return common.ERR_UNKNOWN
	}

	if aff_row == 0 {
		return common.ERR_DB_DEL
	}
	return common.ERR_SUCCESS
}

func CheckFinish(uid int) int {
	_, cur, finish, max, err := ListTask(uid)
	if err == common.ERR_SUCCESS {
		if finish == max {
			if cur == 0 {
				return 1
			}
		}
	}
	return 0
}
