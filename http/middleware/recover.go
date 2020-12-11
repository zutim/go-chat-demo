package middleware

import (
	"chat/http/response"
	"github.com/gin-gonic/gin"
)

// Recover
func Recover(ctx *gin.Context) {
	defer func() {
		// TODO 定制化error
		if r := recover(); r != nil {
			response.WrapContext(ctx).Error(500, "System Error")

			//app.Logger().Error("system_error", map[string]interface{}{
			//	"error": r,
			//})
		}
	}()
	ctx.Next()
}
