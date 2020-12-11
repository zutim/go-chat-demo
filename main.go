package main

import (
	"chat/http/route"
	"github.com/zutim/ego"
	"github.com/zutim/ego/app"
	"github.com/zutim/egu"
)

func init() {
	// 加载配置
	egu.SecurePanic(app.Config().LoadFile("app.yaml"))

	// 初始化数据库
	egu.SecurePanic(app.InitDB())


	//egu.SecurePanic(app.DB().Use(mysql.Resolver().DB))

	//egu.SecurePanic(app.Redis().Connect())
	egu.SecurePanic(app.Redis().ConnectCluster())

}
func main() {
	s := ego.HttpServer()

	// 加载路由
	route.LoadRoute(s.Router)

	// 启动
	egu.SecurePanic(s.Start())
}
