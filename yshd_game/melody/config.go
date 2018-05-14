package melody

import "time"

// Melody configuration.
type Config struct {
	WriteWait         time.Duration // Milliseconds until write times out.
	PongWait          time.Duration // Timeout for waiting on pong.
	PingPeriod        time.Duration // Milliseconds between pings.
	MaxMessageSize    int64         // Maximum size in bytes of a message.
	MessageBufferSize int           // The max amount of messages that can be in a sessions buffer before it starts dropping them.
	SelfPongPeriod    time.Duration
}

func newConfig() *Config {

	return &Config{
		WriteWait:         10 * time.Second,
		PongWait:          20 * time.Second,
		PingPeriod:        (20 * time.Second * 9) / 10,
		MaxMessageSize:    512,
		MessageBufferSize: 256,
		SelfPongPeriod:    20 * time.Second,
	}

	/*
		return &Config{
			WriteWait:         10 * time.Second,
			PongWait:          50 * time.Second,
			PingPeriod:        4 * time.Second,
			MaxMessageSize:    512,
			MessageBufferSize: 256,
		}
	*/
}
