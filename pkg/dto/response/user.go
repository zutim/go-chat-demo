package response

import (
	"chat/pkg/model/entity"
)

// UserAuthResponse 用户校验响应体
type UserAuthResponse struct {
	// 登录令牌
	Token string `json:"token"`
	UserId int `json:"userid"`
}

type UserChat struct {
	entity.User
	Online int `json:"online"`
	UnRead int `json:"unread"`
}

type UserReady struct {
	Users []entity.User
	ChatRecord []string
}


