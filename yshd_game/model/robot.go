package model

import (
	"container/list"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
)

var robot_manager map[int]*User

var robot_arr []User

func InitRobot() map[int]*User {
	robot_manager = make(map[int]*User)

	err := orm.Where("robot=?", true).Find(&robot_manager)
	if err != nil {
		common.Log.Panicf("orm err is %s", err.Error())
	}

	err = orm.Where("robot=?", true).Find(&robot_arr)
	if err != nil {
		common.Log.Panicf("orm err is %s", err.Error())
	}
	return robot_manager
}

func AddRobot(num int, l *list.List) {
	rlen := len(robot_manager)
	if l.Len() > 1 {
		return
	}
	if rlen <= num {
		return
	}
	s := make([]int, 0)
	flag := false

	for count := 0; count < 3*num; count++ {
		r := common.RadnomRange(0, rlen)

		for i := 0; i < len(s); i++ {
			if s[i] == r {
				flag = true
				break
			}
		}
		if flag {
			flag = false
			continue
		}
		s = append(s, r)
		if len(s) == num {
			break
		}
	}

	for i := 0; i < len(s); i++ {
		v := robot_manager[s[i]]

		if v == nil {
			common.Log.Errf("robot err index is %d,%d", i, s[i])
			continue
		}
		u := &UserInfo{}
		v.GetChatUser(u)
		u.IsRobit = true

		l.PushBack(u)
	}
}

/*
//机器人人数至少100以上
func AddRobotFirst(l *list.List) {
	rlen := len(robot_manager)
	if l.Len() > 1 {
		return
	}
	if rlen >= 50 {
		return
	}

	flag := false
	num := common.RandnomRange64(1, 18)
	s := common.RandomRangeArr(1, 50, num)
	for i := 0; i < len(s); i++ {
		v := robot_manager[s[i]]
		if v == nil {
			common.Log.Errf("robot err index is %d,%d", i, s[i])
			continue
		}
		u := &UserInfo{}
		v.GetChatUser(u)
		l.PushBack(u)
	}
}
*/
/*
func AddRobotSecond(l *list.List) {

	if l.Len() >= 100 {
		return
	}
	num := common.RandnomRange64(1, 18)
	s := common.RandomRangeArr(50, 100, num)

	rlen := len(robot_manager)
	if l.Len() > 1 {
		return
	}
	if rlen >= 100 {
		return
	}

	flag := false

	for i := 0; i < len(s); i++ {
		v := robot_manager[s[i]]
		if v == nil {
			common.Log.Errf("robot err index is %d,%d", i, s[i])
			continue
		}
		u := &UserInfo{}
		v.GetChatUser(u)
		l.PushBack(u)
	}
}
*/

/*
func AddRobotFirst(l *list.List, min, max int) {
	//godump.Dump(len(robot_arr))
	if len(robot_arr) < max {
		return
	}
	num := common.RadnomRange(min, max)
	s := common.RandomRangeArr(1, len(robot_arr), num)
	for i := 0; i < len(s); i++ {
		v := robot_arr[i]
		u := &UserInfo{}
		v.GetChatUser(u)
		flag := false

		for e := l.Front(); e != nil; e = e.Next() {
			u := e.Value.(*UserInfo)
			if u.Chat.Uid == v.Uid {
				flag = true
				break
			}
		}
		if !flag {
			u.IsRobit = true

			l.PushBack(u)
		}

	}
}
*/
/*
func AddRobotOfNumber(l *list.List, num int) int {
	numberFlag := 0                                  //用来记录实际增加的机器人
	randNum := common.RadnomRange(1, len(robot_arr)) //取一个随机数

	for i := 0; i < len(robot_arr) && numberFlag < num; i++ {
		temp := (i + randNum) % len(robot_arr)
		v := robot_arr[temp]
		u := &UserInfo{}
		v.GetChatUser(u) //取机器人有用的信息
		flag := false

		for e := l.Front(); e != nil; e = e.Next() {
			u := e.Value.(*UserInfo)
			if u.Chat.Uid == v.Uid {
				flag = true
				break
			}
		}
		if !flag {
			u.IsRobit = true

			// //根据配置文件随机更改机器人的昵称
			// if u.Chat.Sex == 0 {
			// 	//如果是女的
			// 	index := common.RandnomRange64(0, int64(len(RobotWomanNickName)-1))
			// 	u.Chat.NickName = RobotWomanNickName[index]
			// } else {
			// 	//如果是男的
			// 	index := common.RandnomRange64(0, int64(len(RobotManNickName)-1))
			// 	u.Chat.NickName = RobotManNickName[index]
			// }

			l.PushBack(u)
			numberFlag++
		}
	}

	if numberFlag < num { //返回虚拟机器人数量
		return num - numberFlag
	} else {
		return 0
	}
}

func RemoveRobotByNumber(l *list.List, num int) {
	if num < 1 {
		return
	}

	var next *list.Element
	for e := l.Front(); e != nil && num > 0; {
		u := e.Value.(*UserInfo)

		if u.IsRobit {
			next = e.Next()
			l.Remove(e)
			e = next
			num--
			if num < 1 {
				break
			}
		} else {
			e = e.Next()
		}
	}

}
*/

//使用前外部不能加锁
func AddRobotOfNumber(roomid string, num int) int {
	numberFlag := 0                                  //用来记录实际增加的机器人
	randNum := common.RadnomRange(1, len(robot_arr)) //取一个随机数
	chat := GetChatRoom(roomid)
	if chat == nil {
		return 0
	}
	mutex_chat_guardv2.Lock()
	defer mutex_chat_guardv2.Unlock()
	for i := 0; i < len(robot_arr) && numberFlag < num; i++ {
		temp := (i + randNum) % len(robot_arr)
		v := robot_arr[temp]
		u := &UserInfo{}
		v.GetChatUser(u) //取机器人有用的信息

		u.IsRobit = true

		if chat.CheckAudienceByUid(u.Chat.Uid) == false {
			chat.AddAudience(u, -1)
			numberFlag++
		}

	}

	if numberFlag < num { //返回虚拟机器人数量
		return num - numberFlag
	} else {
		return 0
	}
}

//调用前必须在外部保证加锁
func RemoveRobotByNumber(chat *ChatRoomInfo, num int) {
	if num < 1 {
		return
	}

	uids := make([]int, 0)
	//chat := GetChatRoom(roomid)

	for _, u := range chat.UserInfoMap { //遍历房间里的UserInfoMap看User是否是机器人
		if u.IsRobit { //如果是机器人，记录到uids里面
			//chat.DelAudience(u.Chat.Uid)
			uids = append(uids, u.Chat.Uid)
			num--
		}

		if num < 1 { //如果小于1，没有机器人可以删了
			break
		}
	}

	for _, uid := range uids {
		chat.DelAudience(uid) //根据uids保存的uid，删机器人
	}
}
