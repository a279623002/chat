package util

import (
	"chat/models"
)

type MsgLogin struct {
	State bool   `json:"state"`
	Msg   string `json:"msg"`
	User  models.User
}

type Msg struct {
	State bool   `json:"state"`
	Msg   string `json:"msg"`
}

type MsgHead struct {
	State bool   `json:"state"`
	Msg   string `json:"msg"`
	Err   error
}
