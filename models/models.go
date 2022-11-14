package models

import (
	"time"

	"github.com/vmihailenco/msgpack"
)

type EventMessage struct {
	UserID    string    `json:"userid"`
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"timestamp"`
}

var (
	MsgCache = "line_income_msg_cache"
)

func (s EventMessage) MarshalBinary() ([]byte, error) {
	return msgpack.Marshal(s)
}

func (s EventMessage) UnmarshalBinary(data []byte) error {
	return msgpack.Unmarshal(data, s)
}
