package models

import "time"

type EventMessage struct {
	UserID    string    `json:"userid"`
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"timestamp"`
}
