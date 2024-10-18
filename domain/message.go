package domain

import "time"

type MessageBody struct {
	Message string    `json:"message"`
	ID      string    `json:"message_id"`
	Date    time.Time `json:"date"`
	UserID  string    `json:"user_id"`
}

const UserIDAttributeName string = "UserID"
