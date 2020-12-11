package wsmanager

import (
	"chat/http/handler/msg"
	"chat/pkg/app"
	"chat/pkg/service"
	"chat/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zutim/ego"
	"github.com/zutim/egu"
	"strings"
)


func WebsocketHandler(ctx *gin.Context)  {
	conn, err := app.Websocket().UpgradeConn(ctx.Writer, ctx.Request)
	if err != nil {
		egu.SecurePanic(err)
	}

	app.Websocket().Register(conn, func(message []byte){
		var msgs *msg.Common
		if err := utils.JsonDecode(message,&msgs); err != nil {
			fmt.Println(err)
		}

		option := msgs.Op

		switch option {
			case "connection":
				Connection(msgs,conn)
				break
			case "ready":
				Connection(msgs,conn)
				userReady := service.User().ChatReady(conn,msgs.Args)
				record,_ :=json.Marshal(userReady)
				resMsg := msg.ResComm("record","",string(record),"",1)
				app.Websocket().Send([]byte(resMsg),conn)
				break
			case "chat":
				service.User().Chat(conn,msgs)
			    break
			case "close":
				service.User().CloseConnection(conn)
				Broad()
				break
			case "error":
				service.User().CloseConnection(conn)
				Broad()
				break
			case "pageClose":
				service.User().ChatClose(conn)
				break

		}
	})
}

func Broad()  {
	online := app.Websocket().GetOnline()
	res := msg.ResComm("online","",fmt.Sprintf("连接成功，当前在线人数：%v", online),"string",0)
	app.Websocket().Broadcast([]byte(res),nil)
}

func ConvertStringsToBytes(stringContent []string) ([]byte){
	byteContent := "\x00"+ strings.Join(stringContent, "\x02\x00")  // x20 = space and x00 = null
	return []byte(byteContent)
}

func Connection(msgs *msg.Common,conn *ego.WebsocketConn){
	//连接成功，写入User会话
	token := msgs.Msg
	id,err := service.Ws().ParseToken(token)
	if err != nil {
		resMsg := msg.ResComm("connErr","","连接失败，token不正确，重新登录","",1)
		app.Websocket().Send([]byte(resMsg),conn)
	}

	err = service.User().WsBindUserid(id,conn.ID) //绑定
	if err != nil{
		resMsg := msg.ResComm("connErr","","连接失败，已经登录过","",1)
		app.Websocket().Send([]byte(resMsg),conn)
		service.User().CloseConnection(conn)
		return
	}
	//获取所有列表和未读信息
	users,err := service.User().GetUsers(conn,id)
	if err != nil{
		fmt.Println(err)
	}
	user,err :=json.Marshal(users)
	resMsg := msg.ResComm("connSuccess","",string(user),"",1)
	app.Websocket().Send([]byte(resMsg),conn)
	//在线广播
	Broad()

}





