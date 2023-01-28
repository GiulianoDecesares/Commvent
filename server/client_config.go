package server

import "time"

type ClientConfig struct {
	BufferSize     int   `yaml:"bufferSize"`
	MaxMessageSize int64 `yaml:"maxMessageSize"`

	WriteWait  time.Duration `yaml:"writeWait"`
	PongWait   time.Duration `yaml:"pongWait"`
	PingPeriod time.Duration `yaml:"pingPeriod"`
}
