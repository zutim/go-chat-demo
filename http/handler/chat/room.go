package chat

import (
	"github.com/zutim/ego"
	"sync"
)

type room struct {
	User map[int]string
	Mu      sync.Mutex
}

var ro *room
//
func newRoom()  *room{
	return &room{
		User: 	 make(map[int]string),
		Mu:      sync.Mutex{},
	}
}
//
func RoomServe() (r *room)  {
	return ro
}
//
func init() {
	ro = newRoom()
}



func (r *room)Bind(userId int,UUID string){
	r.Mu.Lock()
	r.User[userId] = UUID
	r.Mu.Unlock()
}

func (r *room)GetBind(id int,ID string) (int)  {
	res :=0
	bindUUID := r.User[id]
	if bindUUID==ID { //绑定的用户id就是当前的id
		res = 1
	}

	return res
}

func (r *room)DelBind(UUID string) {
	ego.WebsocketServer().Unregister(UUID)
	for v := range r.User {
		if r.User[v] == UUID {
			delete(r.User,v)
		}
	}
}

//func (r *room)GetUsers() (map[string]*ego.WebsocketConn){
//	conns := make(map[string]*ego.WebsocketConn)
//	//conns = ego.WebsocketServer().GetOnlines()
//	ego.WebsocketServer().GetOnlines()
//	//fmt.Println("hao:",conns)
//	return conns
//}

//type Room struct {
//	//Clients map[string]*ego.WebsocketConn
//	User map[string]int
//	Mu      sync.Mutex
//}
//


