package domain

import "time"


type MessageBody struct {
    Message string `json:"message"`
    Date    time.Time `json:"date"`
}


const UserIDAttributeName string = "UserID"