package model

import (
	"bufio"
	"github.com/yshd_game/common"
	"io"
	"os"
	"strings"
)

var SayRobotArr []string

var lenghtSayRobot int

var RobotManNickName []string   //男性机器人昵称
var RobotWomanNickName []string //女性机器人昵称

func ReadLine(filePth string, key_word_set map[string]bool) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer f.Close()
	bfRd := bufio.NewReader(f)
	for {
		line_, err := bfRd.ReadBytes('\n')
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				return nil
			}
			return err
		}
		line := string(line_)
		hookfn(key_word_set, line) //放在错误处理前面，即使发生错误，也会处理已经读取到的数据。
	}
	return nil
}
func hookfn(key_word_set map[string]bool, line string) {
	line = strings.Replace(line, "\n", "", -1)
	line = strings.Replace(line, "\r", "", -1)
	key_word_set[line] = true
}

func InitSayWord() {
	if len(SayRobotArr) > 0 {
		for k, _ := range SayRobotArr {
			SayRobotArr[k] = ""
		}
	}

	key_word_set := make(map[string]bool, 0)
	ReadLine("./config/robot.txt", key_word_set)
	SayRobotArr = make([]string, len(key_word_set))

	//SayRobotArr = new([]string)
	LoadRobotSay(key_word_set)
	lenghtSayRobot = len(key_word_set)

}

func LoadRobotSay(key_word_set map[string]bool) {
	i := 0
	for k, _ := range key_word_set {
		SayRobotArr[i] = k
		i++
	}
}

func GetRandomSay() string {
	if len(SayRobotArr) == 0 {
		return ""
	}
	i := common.RadnomRange(0, len(SayRobotArr)-1)
	return SayRobotArr[i]
}

func GetRobotSayByID(id int) (string, bool) {
	if id > len(SayRobotArr) {
		return "", false
	}
	return SayRobotArr[id], true
}

type RobotSayInfo struct {
	say_record map[int]bool
	count      int
}

func (self *RobotSayInfo) Init(count int) {

	self.say_record = make(map[int]bool)
	for i := 0; i < count; i++ {
		self.say_record[i] = true
	}
	self.count = count
}

func (self *RobotSayInfo) Reset() {
	for k, _ := range self.say_record {
		delete(self.say_record, k)
	}

	for i := 0; i < self.count; i++ {
		self.say_record[i] = true
	}

}
func (self *RobotSayInfo) GetRandomSayInChat() (int, bool) {
	if self.count <= 1 {
		return 0, false
	}
	i := common.RadnomRange(0, self.count-1)
	count := 0
	for k, _ := range self.say_record {
		if count == i {
			delete(self.say_record, k)
			self.count--
			return k, true
		}
		count++
	}
	return 0, false
}
func LoadAllRobotNickName() {
	if len(RobotManNickName) > 0 {
		for k, _ := range RobotManNickName {
			RobotManNickName[k] = ""
		}
	}

	key_word_set := make(map[string]bool, 0)
	ReadLine("./config/robot_man_nickname.txt", key_word_set)
	RobotManNickName = make([]string, len(key_word_set))

	i := 0
	for key_word, _ := range key_word_set {
		RobotManNickName[i] = key_word
		i++
	}

	if len(RobotWomanNickName) > 0 {
		for k, _ := range RobotWomanNickName {
			RobotWomanNickName[k] = ""
		}
	}

	key_word_set2 := make(map[string]bool, 0)
	ReadLine("./config/robot_woman_nickname.txt", key_word_set2)
	RobotWomanNickName = make([]string, len(key_word_set2))

	i = 0
	for key_word, _ := range key_word_set2 {
		RobotWomanNickName[i] = key_word
		i++
	}

}
