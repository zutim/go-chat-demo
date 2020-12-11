package handler

import (
	"chat/http/response"
	"chat/pkg/dto/request"
	"chat/pkg/service"
	"github.com/gin-gonic/gin"
)


func UserHandler(ctx *gin.Context) {
	var req request.AuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.WrapContext(ctx).Error(1001, "参数错误:" + err.Error())
		return
	}

	// 调用service的Auth方法，获得结果
	res, err := service.User().Auth(req)

	// 有错就抛panic
	//egu.SecurePanic(err)
	if err != nil {
		response.WrapContext(ctx).Error(1001, err.Error())
	}

	response.WrapContext(ctx).Success(res)
	// 自定义的response模块，更多使用方法见[1.1.3.响应]
	//response.WrapContext(ctx).Success(nil)
}