package service

import (
	"chat/pkg/model/data"
	"github.com/zutim/ego/app"
)

type wsService struct {
}

// User 用户服务
func Ws() *wsService {
	return &wsService{}
}

func (ws *wsService)ParseToken(token string)(int,error) {
	var claim data.UserClaims
	err := app.Jwt().ParseTokenWithClaims(token,&claim)
	if err != nil {
		return 0,err
	}
	return claim.User.Id,nil
}