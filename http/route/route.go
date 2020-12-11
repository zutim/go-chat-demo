package route

import (
	"chat/http/handler"
	wshandler "chat/http/handler/wsmanager"
	"chat/http/response"
	"chat/pkg/app"
	"chat/pkg/model/data"
	"github.com/gin-gonic/gin"
	"github.com/zutim/ego/http/middleware"
)
func LoadRoute(router *gin.Engine) {
	//导入静态地址
	router.LoadHTMLGlob("./html/view/*")

	//router.LoadHTMLFiles("./html/view/index.html", "./html/view/login.html")
	//加载静态资源，例如网页的css、js
	router.Static("html/static", "./html/static")

	router.GET("/user/login", func(ctx *gin.Context) {
		ctx.HTML(200,"login.html",nil)
	})

	router.GET("/index", func(ctx *gin.Context) {
		ctx.HTML(200,"index.html",nil)
	})

	router.GET("/index/chat", func(ctx *gin.Context) {
		ctx.HTML(200,"chat.html",nil)
	})

	router.POST("/user/auth", handler.UserHandler)

	// 定义需要token校验的路由
	index := router.Group("/ws").Use(middleware.JWT(&data.UserClaims{}))
	{
		// 获取用户信息
		index.GET("index", func(ctx *gin.Context) {
			response.WrapContext(ctx).Success("ok")
		})
	}


	// websocket
	ws :=app.Websocket()

	router.GET("/ws", wshandler.WebsocketHandler)

	go ws.Start()

}
