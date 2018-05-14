package common

/*
var my_session *sessions.session

func SetSession() *sessions.session {
	my_session = new(sessions.session)
	return my_session
}

func InitSession() {
	sessions.session.Set("roomAllocID", 1)
}

func GetRoomAllocId() int {
	id := sessions.session.Get("roomAllocID")
	id++
	return id
}

func SetRoomAllocID(id int) {
	sessions.session.Set("roomAllocID", id)
}
*/

func GetRoomAllocId(roomid int) int {
	return roomid + 1
}
