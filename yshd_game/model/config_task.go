package model

import (
	"github.com/yshd_game/common"
	"github.com/yshd_game/confdata"
	"strconv"
	//"sort"
	"strings"
)

type ConfigTask struct {
	TaskId int64 `xorm:" pk not null unique"`

	//任务名字
	Name string `xorm:"varchar(128) not null"`

	//任务类型
	Type int // 0日常 1限时任务

	//任务目标
	Target int // 0胜利场

	TargetParam int32

	//奖励1
	Reward1 string
	//奖励2
	Reward2 int64
}

var task_config_cache []*confdata.TaskDefine

func LoadConfigTask() {
	m := make(map[int64]ConfigTask)

	var err error
	err = orm.Find(m)
	if err != nil {
		common.Log.Errf("orm is err %s", err.Error())
		return
	}

	for _, v := range m {
		var task confdata.TaskDefine
		//task=new(confdata.TaskDefine)
		task.Reward1 = new(confdata.RewardProp)
		task.Id = v.TaskId

		task.Name = v.Name
		task.Target = confdata.TargetType(v.Target)
		task.TargetParam = v.TargetParam
		task.Type = confdata.TaskType(v.Type)
		task.Reward2 = v.Reward2
		task.Reward1.MoneyType = 1
		//godump.Dump(v.Reward1)
		s := strings.Split(v.Reward1, ";")
		if len(s) == 2 {
			//godump.Dump(s)

			mtype, err := strconv.ParseInt(s[0], 10, 32)
			//strconv.p
			if err != nil {
				common.Log.Panicf("config task err is %s ", err.Error())
			}
			/*
				mtype, err := strconv.Atoi(s[0])

				if err != nil {
					common.Log.Errf("config task err is %s ",err.Error())
					return
				}
			*/

			//godump.Dump(mtype)

			task.Reward1.MoneyType = int32(mtype)

			num, err := strconv.Atoi(s[1])

			if err != nil {
				common.Log.Panicf("config task err is %s ", err.Error())
				return
			}

			task.Reward1.Num = int64(num)

		} else if len(s) == 0 {

		} else {
			common.Log.Panic("config  task err ")
		}

		confdata.ConfigData.TaskById[task.Id] = &task
	}

	task_config_cache = make([]*confdata.TaskDefine, 0)
	for _, v := range confdata.ConfigData.TaskById {
		task_config_cache = append(task_config_cache, v)
	}
}

func GetTaskConfig() []*confdata.TaskDefine {
	return task_config_cache
}
