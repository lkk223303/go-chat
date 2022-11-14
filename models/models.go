package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"time"
)

var (
	MsgCache = "line_income_msg_cache"
)

type EventMessage struct {
	UserID    string    `json:"userid"`
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"timestamp"`
}

func (e EventMessage) MarshalBinary() ([]byte, error) {

	b := bytes.NewBuffer(nil)
	b.WriteString(e.UserID)
	b.WriteByte('\t')
	b.WriteString(e.Message)
	b.WriteByte('\t')
	b.WriteString(e.TimeStamp.String())
	return b.Bytes(), nil
}

func (e EventMessage) UnmarshalBinary(data []byte) (EventMessage, error) {
	// 先按照 json 解析
	err := json.Unmarshal(data, &e)
	if err == nil {
		// 解析成功直接返回
		log.Println("json parsed")
		return e, nil
	}

	// 前面失败，按照 \t 分割解析
	results := bytes.Split(data, []byte{'\t'})
	if len(results) != 2 {
		return e, errors.New("not valid format")
	}
	e.UserID = string(results[0])
	e.Message = string(results[1])
	e.TimeStamp, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(results[2]))
	if err != nil {
		return e, err
	}
	return e, nil
}
