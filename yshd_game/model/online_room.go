package model

import (
//"fmt"
//"sync"
)

type RoomInfo struct {
	Rid       string
	Name      string
	Begintime int
	Uid       int
	Image     string
	Mtype     int
}

/*
type OwnerInfo struct {
	Uid int
}

type onlineRoom struct {
	//room_list map[uint32]RoomInfo
	mutex sync.RWMutex
}

var _instance *onlineRoom

func InstanceOnlineRoom() *onlineRoom {
	if _instance == nil {
		_instance = new(onlineRoom)
	}
	return _instance
}

func (this *onlineRoom) AddNewRoom(room RoomInfo) {
	//this.room_list
	fmt.Println("hello")
}

func (this *onlineRoom) Hello() {
	fmt.Println("hello")
}
*/
/*
package singleton

import (
	"fmt"
	"sync"
)

type RoomInfo struct {
	room_id        uint32
	room_begintime uint32
	room_uid       uint32
	romm_type      uint32
}

var (
	mutex     sync.RWMutex
	room_data = make(map[uint32]interface{})
)

func Set(key uint32, val interface{}) {
	mutex.Lock()
	room_data[key] = val
	mutex.Unlock()
}

func Hello() {
	fmt.Println("hello")
}
*/
