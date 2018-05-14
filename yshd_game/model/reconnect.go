package model

import (
	"container/list"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"sync"
	"time"
)

type AnchorInfo struct {
	Uid     int
	Rid     string
	Timeout int64
}

type AnchorRoom struct {
	reconnects   map[int]*AnchorInfo
	closeconnect chan *AnchorInfo
	mutex_anchor sync.RWMutex
	exit         chan bool
	//wgClose      sync.WaitGroup
	wgReconnect sync.WaitGroup
}

var AnchorMgr *AnchorRoom

func InitAchorRoom() {
	AnchorMgr = newAnchorRoom()
}

func newAnchorRoom() *AnchorRoom {
	return &AnchorRoom{
		reconnects:   make(map[int]*AnchorInfo, 200),
		closeconnect: make(chan *AnchorInfo, 200),
		exit:         make(chan bool, 200),
	}
}

func (r *AnchorRoom) ReconnectRoom() {
	defer common.PrintPanicStack()
	keep_ticker := time.NewTicker(1 * time.Second)
	defer keep_ticker.Stop()
	var tl list.List
	tl.Init()
	go func() {
		for m := range r.closeconnect {
			common.Log.Debugf("reconnect close room uid=%d,rid=%s", m.Uid, m.Rid)
			CloseChat(m.Uid, m.Rid)
			//r.wgClose.Done()
		}
	}()

loop:
	for {
		select {
		case <-keep_ticker.C:
			ntime := time.Now().Unix()

			/*
				for _, s := range r.reconnects {
					if s.Timeout < ntime {
						//godump.Dump("time out")
						r.AddCloseChannel(s)
						tl.PushBack(s.Uid)
					}
				}

				for e := tl.Front(); e != nil; e = e.Next() {
					uid := e.Value.(int)
					r.DelReconnectMap(uid)
					tl.Remove(e)
				}
			*/
			for _, s := range r.reconnects {
				if s.Timeout < ntime {
					r.SetCloseRoom(s.Uid)
				}
			}
		case <-r.exit:
			for _, s := range r.reconnects {
				r.AddCloseChannel(s)
				r.DelReconnectMap(s.Uid)
			}
			break loop
		}
	}
}

func (r *AnchorRoom) Ping(uid int) int {
	v, ok := r.reconnects[uid]
	if ok {
		ntime := time.Now().Unix()
		v.Timeout = ntime + common.RECONNECT_TIMEOUT
		return common.ERR_SUCCESS
	}
	return common.ERR_RECONNECT_ROOM
}

func (r *AnchorRoom) AddAnchor(uid int, rid string) bool {
	anchor := &AnchorInfo{
		Uid:     uid,
		Rid:     rid,
		Timeout: time.Now().Unix() + common.RECONNECT_TIMEOUT,
	}
	_, ok := r.reconnects[uid]
	if ok {
		common.Log.Errf("wait reconnect uid=%d,rid=%s", uid, rid)
		return false
	}
	r.AddReconnectMap(anchor)
	return true
}

func (r *AnchorRoom) SetCloseRoom(uid int) bool {
	s, ok := r.reconnects[uid]
	if ok {
		r.AddCloseChannel(s)
		r.DelReconnectMap(uid)
	}

	return ok
}

/*
func (r *AnchorRoom) CheckRes(uid int) bool {
	_, ok := r.reconnects[uid]
	if ok {
		//r.SetCloseRoom(uid)
		r.AddCloseChannel(s)
	}
	return ok
}


*/

func (r *AnchorRoom) Close() {
	for k, v := range chat_room_manager {
		DirectCloseRoom(v.room.Uid, k)
	}

	/*
		r.exit <- true
		r.wgReconnect.Wait()
	*/
	//r.wgClose.Wait()
}

func (r *AnchorRoom) AddCloseChannel(s *AnchorInfo) {
	//r.wgClose.Add(1)
	r.closeconnect <- s
}

func (r *AnchorRoom) AddReconnectMap(s *AnchorInfo) {
	common.Log.Debugf("add anchor into reconnect uid=%d,rid=%s", s.Uid, s.Rid)
	r.mutex_anchor.Lock()
	defer r.mutex_anchor.Unlock()
	r.wgReconnect.Add(1)
	r.reconnects[s.Uid] = s
}

func (r *AnchorRoom) DelReconnectMap(uid int) bool {
	r.mutex_anchor.Lock()
	defer r.mutex_anchor.Unlock()
	_, ok := r.reconnects[uid]
	if ok {
		//r.DelReconnectMap(uid)
		defer r.wgReconnect.Done()
		delete(r.reconnects, uid)
		return true
	}
	return false
}

/*
func (r *AnchorRoom) DelReconnectMap(uid int) {
	r.mutex_anchor.Lock()
	defer r.mutex_anchor.Unlock()
	defer r.wgReconnect.Done()
	delete(r.reconnects, uid)
}
*/
