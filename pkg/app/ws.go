package app

import "github.com/zutim/ego"

var ws ego.WsServer

func init()  {
	ws= ego.WebsocketServer()
}

func Websocket() ego.WsServer {
	return ws
}

