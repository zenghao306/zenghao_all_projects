package melody

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"

	"github.com/yshd_game/common"
	"net/http"
	"sync"
	"time"
)

// A melody session.
type Session struct {
	Request *http.Request
	Keys    map[string]interface{}
	conn    *websocket.Conn
	output  chan *envelope
	melody  *Melody
	timeout int64
	times   int
	open    bool
	rwmutex *sync.RWMutex
}

func (s *Session) writeMessage(message *envelope) {
	if s.closed() {
		s.melody.errorHandler(s, errors.New("tried to write to closed a session"))
		return
	}
	select {
	case s.output <- message:
	default:
		s.melody.errorHandler(s, errors.New("Message buffer full"))
	}
}

func (s *Session) writeRaw(message *envelope) error {
	if s.closed() {
		return errors.New("tried to write to a closed session")
	}
	s.conn.SetWriteDeadline(time.Now().Add(s.melody.Config.WriteWait))
	err := s.conn.WriteMessage(message.t, message.msg)

	if err != nil {
		return err
	}
	/*
		if message.t == websocket.CloseMessage {
			err := s.conn.Close()

			if err != nil {
				return err
			}
		}
	*/
	return nil
}

func (s *Session) closed() bool {

	s.rwmutex.RLock()
	defer s.rwmutex.RUnlock()

	return !s.open
}

func (s *Session) close() {
	if !s.closed() {
		s.rwmutex.Lock()
		s.open = false
		s.conn.Close()
		close(s.output)
		s.rwmutex.Unlock()
	}
}

func (s *Session) ping() {
	s.writeRaw(&envelope{t: websocket.PingMessage, msg: []byte{}})
}

func (s *Session) writePump() {
	defer common.PrintPanicStack()
	//defer s.conn.Close()

	ticker := time.NewTicker(s.melody.Config.PingPeriod)
	defer ticker.Stop()

	keep_ticker := time.NewTicker(s.melody.Config.SelfPongPeriod)
	defer keep_ticker.Stop()
loop:
	for {
		select {
		case msg, ok := <-s.output:
			if !ok {
				//s.close()
				break loop
			}
			if err := s.writeRaw(msg); err != nil {
				s.melody.errorHandler(s, err)
				break loop
			}

			if msg.t == websocket.CloseMessage {
				//time.Sleep(1*time.Second)
				break loop
			}

			if msg.t == websocket.TextMessage {
				s.melody.messageSentHandler(s, msg.msg)
			}

			if msg.t == websocket.BinaryMessage {
				s.melody.messageSentHandlerBinary(s, msg.msg)
			}
		case <-ticker.C:
			s.ping()

		case <-keep_ticker.C:
			{
				if time.Now().Unix() > s.timeout {
					if s.times >= 3 {
						s.close()
						break loop
					} else {
						s.times++
					}
				}

			}
		}
	}
}

func (s *Session) readPump() {
	//	defer s.conn.Close()

	s.conn.SetReadLimit(s.melody.Config.MaxMessageSize)
	s.conn.SetReadDeadline(time.Now().Add(s.melody.Config.PongWait))

	s.conn.SetPongHandler(func(string) error {
		s.conn.SetReadDeadline(time.Now().Add(s.melody.Config.PongWait))
		s.timeout = time.Now().Unix() + int64(s.melody.Config.PongWait/time.Second)

		s.times = 0
		uid := s.Request.FormValue("uid")
		s.melody.PongChan <- uid
		return nil
	})

	for {
		t, message, err := s.conn.ReadMessage()
		if err != nil {
			s.melody.errorHandler(s, err)
			break
		}
		if t == websocket.TextMessage {
			s.melody.messageHandler(s, message)
		}

		if t == websocket.BinaryMessage {
			s.melody.messageHandlerBinary(s, message)
		}
	}
}

// Write message to session.
func (s *Session) Write(msg []byte) error {
	if s.closed() {
		return errors.New("session is closed")
	}
	s.writeMessage(&envelope{t: websocket.TextMessage, msg: msg})
	return nil
}

// Write binary message to session.
func (s *Session) WriteBinary(msg []byte) error {
	if s.closed() {
		return errors.New("session is closed")
	}
	s.writeMessage(&envelope{t: websocket.BinaryMessage, msg: msg})
	return nil
}

// Close session.
func (s *Session) Close() error {
	if s.closed() {
		return errors.New("session is already closed")
	}
	s.writeMessage(&envelope{t: websocket.CloseMessage, msg: []byte{}})
	return nil
}

/*
func (s *Session) WriteAndClose(msg []byte) {
	/*
		_, ok := <-s.output
		if !ok {
			s.close()
			common.Log.Debug("channel is err")
			return
		}

	s.writeRaw(&envelope{t: websocket.TextMessage, msg: msg})
	s.close()
}


*/
func (s *Session) SendMsg(msg interface{}) (err error) {
	if b, err := json.Marshal(msg); err == nil {
		return s.Write(b)
	}
	return err
}

func (s *Session) CloseWithMsgAndJson(msg interface{}) (err error) {
	s.SendMsg(msg)
	time.Sleep(2 * time.Second)
	if b, err := json.Marshal(msg); err == nil {
		return s.CloseWithMsg(b)
	}
	return err
}

func (s *Session) CloseWithMsg(msg []byte) error {
	if s.closed() {
		return errors.New("session is already closed")
	}

	s.writeMessage(&envelope{t: websocket.CloseMessage, msg: msg})

	return nil
}

func (s *Session) Set(key string, value interface{}) {
	if s.Keys == nil {
		s.Keys = make(map[string]interface{})
	}

	s.Keys[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (s *Session) Get(key string) (value interface{}, exists bool) {
	if s.Keys != nil {
		value, exists = s.Keys[key]
	}

	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (s *Session) MustGet(key string) interface{} {
	if value, exists := s.Get(key); exists {
		return value
	}

	panic("Key \"" + key + "\" does not exist")
}

// IsClosed returns the status of the connection.
func (s *Session) IsClosed() bool {
	return s.closed()
}

/*
func (s *Session) SendMsgToSelf(res interface{}) bool {
	if b, err := json.Marshal(res); err == nil {
		s.Write(b)
		return true
	}
	return false
}


func (s *Session) SendMsgToSelf(res interface{}) bool {
	s.BroadcastFilter(msg, func(q *melody.Session) bool {
		return s == q
	})
}
*/
