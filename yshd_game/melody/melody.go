package melody

import (
	"github.com/gorilla/websocket"
	//"github.com/liudng/godump"
	"net/http"
	//"time"
	"encoding/json"
	"errors"
	"github.com/yshd_game/common"
	"sync"
)

const (
	CloseNormalClosure           = 1000
	CloseGoingAway               = 1001
	CloseProtocolError           = 1002
	CloseUnsupportedData         = 1003
	CloseNoStatusReceived        = 1005
	CloseAbnormalClosure         = 1006
	CloseInvalidFramePayloadData = 1007
	ClosePolicyViolation         = 1008
	CloseMessageTooBig           = 1009
	CloseMandatoryExtension      = 1010
	CloseInternalServerErr       = 1011
	CloseServiceRestart          = 1012
	CloseTryAgainLater           = 1013
	CloseTLSHandshake            = 1015
)

// Duplicate of codes from gorilla/websocket for convenience.
var validReceivedCloseCodes = map[int]bool{
	// see http://www.iana.org/assignments/websocket/websocket.xhtml#close-code-number

	CloseNormalClosure:           true,
	CloseGoingAway:               true,
	CloseProtocolError:           true,
	CloseUnsupportedData:         true,
	CloseNoStatusReceived:        false,
	CloseAbnormalClosure:         false,
	CloseInvalidFramePayloadData: true,
	ClosePolicyViolation:         true,
	CloseMessageTooBig:           true,
	CloseMandatoryExtension:      true,
	CloseInternalServerErr:       true,
	CloseServiceRestart:          true,
	CloseTryAgainLater:           true,
	CloseTLSHandshake:            false,
}

type handleMessageFunc func(*Session, []byte)
type handleErrorFunc func(*Session, error)
type handleCloseFunc func(*Session, int, string) error
type handleSessionFunc func(*Session)
type filterFunc func(*Session) bool

type AfterRunFunc func(*Session)

type ExpectionFunc func(*Session)

type Melody struct {
	Config                   *Config
	Upgrader                 *websocket.Upgrader
	messageHandler           handleMessageFunc
	messageHandlerBinary     handleMessageFunc
	messageSentHandler       handleMessageFunc
	messageSentHandlerBinary handleMessageFunc
	errorHandler             handleErrorFunc
	closeHandler             handleCloseFunc
	connectHandler           handleSessionFunc
	disconnectHandler        handleSessionFunc
	pongHandler              handleSessionFunc
	hub                      *hub
	ExpectionHandler         handleSessionFunc
	PongChan                 chan string
}

// Returns a new melody instance with default Upgrader and Config.
func New() *Melody {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	hub := newHub()

	go hub.run()

	return &Melody{
		Config:                   newConfig(),
		Upgrader:                 upgrader,
		messageHandler:           func(*Session, []byte) {},
		messageHandlerBinary:     func(*Session, []byte) {},
		messageSentHandler:       func(*Session, []byte) {},
		messageSentHandlerBinary: func(*Session, []byte) {},
		errorHandler:             func(*Session, error) {},
		closeHandler:             nil,
		connectHandler:           func(*Session) {},
		disconnectHandler:        func(*Session) {},

		ExpectionHandler: func(*Session) {},
		hub:              hub,
		PongChan:         make(chan string, 6000),
	}
}

func (m *Melody) HandleExpection(fn func(*Session)) {
	m.ExpectionHandler = fn
}

// Fires fn when a session connects.
func (m *Melody) HandleConnect(fn func(*Session)) {
	m.connectHandler = fn
}

// Fires fn when a session disconnects.
func (m *Melody) HandleDisconnect(fn func(*Session)) {
	m.disconnectHandler = fn
}

func (m *Melody) HandlePong(fn func(*Session)) {
	m.pongHandler = fn
}

// Callback when a text message comes in.
func (m *Melody) HandleMessage(fn func(*Session, []byte)) {
	m.messageHandler = fn
}

// Callback when a binary message comes in.
func (m *Melody) HandleMessageBinary(fn func(*Session, []byte)) {
	m.messageHandlerBinary = fn
}

// HandleSentMessage fires fn when a text message is successfully sent.
func (m *Melody) HandleSentMessage(fn func(*Session, []byte)) {
	m.messageSentHandler = fn
}

// HandleSentMessageBinary fires fn when a binary message is successfully sent.
func (m *Melody) HandleSentMessageBinary(fn func(*Session, []byte)) {
	m.messageSentHandlerBinary = fn
}

// Fires when a session has an error.
func (m *Melody) HandleError(fn func(*Session, error)) {
	m.errorHandler = fn
}

// Handles http requests and upgrades them to websocket connections.
func (m *Melody) HandleClose(fn func(*Session, int, string) error) {
	if fn != nil {
		m.closeHandler = fn
	}
}
func (m *Melody) HandleRequest(w http.ResponseWriter, r *http.Request) error {
	defer common.PrintPanicStack()
	if m.hub.closed() {
		return errors.New("melody instance is closed")
	}

	conn, err := m.Upgrader.Upgrade(w, r, nil)

	if err != nil {
		return err
	}

	session := &Session{
		Request: r,
		conn:    conn,
		output:  make(chan *envelope, m.Config.MessageBufferSize),
		melody:  m,
		timeout: 0,
		times:   0,
		open:    true,
		rwmutex: &sync.RWMutex{},
	}

	m.hub.register <- session

	m.hub.wgConns.Add(1)

	go m.connectHandler(session)

	go session.writePump()

	session.readPump()

	if m.hub.open {
		m.hub.unregister <- session
	}
	session.close()
	m.disconnectHandler(session)

	m.hub.wgConns.Done()
	return nil
}

// Broadcasts a text message to all sessions.
func (m *Melody) Broadcast(msg []byte) {
	message := &envelope{t: websocket.TextMessage, msg: msg}
	m.hub.broadcast <- message
}

//

func (m *Melody) BroadcastFilterByJson(data interface{}, fn func(*Session) bool) {
	if b, err := json.Marshal(data); err == nil {
		msg := b
		m.BroadcastFilter(msg, fn)
	}
}

// Broadcasts a text message to all sessions that fn returns true for.

func (m *Melody) BroadcastFilter(msg []byte, fn func(*Session) bool) {
	message := &envelope{t: websocket.TextMessage, msg: msg, filter: fn}
	if m.hub.open {
		m.hub.broadcast <- message
	}
}

// Broadcasts a text message to all sessions except session s.
func (m *Melody) BroadcastOthers(msg []byte, s *Session) {
	m.BroadcastFilter(msg, func(q *Session) bool {
		return s != q
	})
}

func (m *Melody) SendToSelf(msg []byte, s *Session) {
	message := &envelope_to_self{t: websocket.TextMessage, msg: msg, sess: s}
	if m.hub.open {
		m.hub.selfmsg <- message
	}
}

func (m *Melody) SendToSelfWithFunc(msg []byte, s *Session, a AfterRunFunc) {
	message := &envelope_to_self{t: websocket.TextMessage, msg: msg, sess: s, afunc: a}
	if m.hub.open {
		m.hub.selfmsg <- message
	}
}

// Broadcasts a binary message to all sessions.
func (m *Melody) BroadcastBinary(msg []byte) {
	message := &envelope{t: websocket.BinaryMessage, msg: msg}
	m.hub.broadcast <- message
}

// Broadcasts a binary message to all sessions that fn returns true for.
func (m *Melody) BroadcastBinaryFilter(msg []byte, fn func(*Session) bool) {
	message := &envelope{t: websocket.BinaryMessage, msg: msg, filter: fn}
	m.hub.broadcast <- message
}

// Broadcasts a binary message to all sessions except session s.
func (m *Melody) BroadcastBinaryOthers(msg []byte, s *Session) {
	m.BroadcastBinaryFilter(msg, func(q *Session) bool {
		return s != q
	})
}

// Closes the melody instance and all connected sessions.
func (m *Melody) Close() error {
	if m.hub.closed() {
		return errors.New("melody instance is already closed")
	}
	m.hub.exit <- &envelope{t: websocket.CloseMessage, msg: []byte{}}
	m.hub.wgLn.Wait()

	m.hub.wgConns.Wait()
	return nil
}
