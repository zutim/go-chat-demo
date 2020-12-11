// service 服务模块
package service

import (
	"chat/http/handler/chat"
	"chat/http/handler/msg"
	chatApp "chat/pkg/app"
	"chat/pkg/dto/request"
	"chat/pkg/dto/response"
	"chat/pkg/enum/statusCode"
	"chat/pkg/model/dao"
	"chat/pkg/model/data"
	"chat/pkg/model/entity"
	"chat/pkg/utils"
	"fmt"
	"github.com/ebar-go/ego"
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/ego/errors"
	"github.com/ebar-go/egu"
	"strconv"
	"time"
)

type userService struct {
}

// User 用户服务
func User() *userService {
	return &userService{}
}



// Auth 校验用户登录
func (service *userService) Auth(req request.AuthRequest) (*response.UserAuthResponse, error) {
	userDao := dao.User(app.DB())
	user,err := userDao.Auth(req)
	if err!=nil{
		return nil,errors.New(statusCode.DatabaseQueryFailed,"账号错误")
	}

	if user.Password != egu.Md5(req.Password) {
		return nil,errors.New(statusCode.PasswordWrong,"密码错误")
	}

	//str,err := app.Redis().Get(fmt.Sprintf("user:tel:%s:userid",req.Tel)).Result()
	//
	//if err != nil {
	//	return nil, errors.New(statusCode.DataNotFound, fmt.Sprintf("获取用户信息失败:%v", err))
	//}
	//
	//id, err := strconv.Atoi(str)
	//if err != nil {
	//	return nil, errors.New(statusCode.PasswordWrong, "输入信息错误")
	//}
	//
	//if res := chat.RoomServe().GetBind(id,"0"); res==1{
	//	return nil,errors.New(statusCode.DataNotFound, fmt.Sprintf("该用户已经登录，无法绑定"))
	//}
	//
	// 组装结构体

	res := new(response.UserAuthResponse)
	// 生成token
	userClaims := new(data.UserClaims)
	userClaims.ExpiresAt = egu.GetTimeStamp() + 3600*24*10
	userClaims.User.Id = user.Id
	userClaims.User.Tel = user.Tel
	token, err := app.Jwt().GenerateToken(userClaims)
	if err != nil {
		return nil, errors.New(statusCode.TokenGenerateFailed, err.Error())
	}
	res.Token = token
	res.UserId = user.Id
	//
	//keyStr := app.Config().GetString("server.desKey")
	//key := []byte(keyStr)
	//userid,err := utils.Encrypt(str,key)
	//if err != nil {
	//	fmt.Println(err)
	//	egu.SecurePanic(err)
	//}
	//res.UserId = userid
	return res, nil
}

func (service *userService) WsBindUserid(id int,ID string) (error)  {
	//id :=getDeUserId(myUserID)
	//id,_ := chat.RoomServe().User[]
	res := chat.RoomServe().GetBind(id,ID)
	if  res==1{
		chat.RoomServe().DelBind(ID)
	}

	chat.RoomServe().Bind(id,ID)
	return nil
}

func (service *userService) GetUsers(conn *ego.WebsocketConn,id int)  ([]*response.UserChat,error){
	//获取的数据库所有用户
	//users,_ := app.Redis().HGetAll("user").Result()
	userDao := dao.User(app.DB())
	users,err := userDao.GetUsers()
	if err != nil {
		egu.SecurePanic(err)
	}
	len := len(users)

	//返回的数据
	resUser := make([]*response.UserChat,len,len)

	//获取该用户未读的所有信息
	touserid := getUserIdByUUID(conn.ID)
	chat.DataServe().To= touserid

	unreadNums := chat.DataServe().GetUnreadMsgs()
	fmt.Println("未读信息",unreadNums)

	for i := range users {
		if users[i].Id == id {
			continue
		}
		var user response.UserChat
		user.Id = users[i].Id
		user.Username = users[i].Username
		user.Headimg = users[i].Headimg

		userid_str := strconv.Itoa(user.Id)
		for key := range unreadNums {
			if key==userid_str {
				num ,_:= strconv.Atoi(unreadNums[key])
				user.UnRead = num
			}
		}

		//判断是否在线
		if _,ok := chat.RoomServe().User[user.Id]; ok{
			uuid := chat.RoomServe().User[user.Id] //获取uuid

			if isRe :=chatApp.Websocket().IsRegist(uuid);isRe==1 {
				user.Online =1
			}
		}

		resUser[i] = &user
	}
	return resUser,nil
}

func (service *userService) CloseConnection(conn *ego.WebsocketConn)  {
	chat.RoomServe().DelBind(conn.ID)
}

func (service *userService) Chat(conn *ego.WebsocketConn,msgs *msg.Common)  {
	from := msgs.Args
	to := msgs.MsgType
	fromUserId,_ := strconv.Atoi(from)
	chat.DataServe().From = fromUserId
	toUserId,_ := strconv.Atoi(to)
	chat.DataServe().To = toUserId
	ToUUID := chat.RoomServe().User[toUserId]
	//获取to的连接信息
	ToConn := chatApp.Websocket().GetConnection(ToUUID)

	chat.DataServe().Message = msgs.Msg//string(record)
	chat.DataServe().Sent = time.Now()
	chat.DataServe().SetChatRecord()

	if res := chatApp.Websocket().IsRegist(ToUUID);res==1{
		resMsg := msg.ResComm("chat","",msgs.Msg,"",1)
		chatApp.Websocket().Send([]byte(resMsg),ToConn)
		return
	}
	fmt.Println("离线了")
}


func (service *userService) ChatReady(conn *ego.WebsocketConn,toid string) (response.UserReady) {

	//touserid := getDeUserId(toid)
	touserid,_ := strconv.Atoi(toid)

	chat.DataServe().To = touserid

	userid := getUserIdByUUID(conn.ID)
	chat.DataServe().From = userid
	chat.DataServe().SetUnreadToRead()
	res := chat.DataServe().GetChatRecord(10)
	resUsers := GetChatUser(userid,touserid)
	fmt.Println("touserid:",touserid,"userid:",userid)
	var userReady response.UserReady
	userReady.ChatRecord = res
	userReady.Users = resUsers
	return userReady
}

func (service *userService) ChatClose(conn *ego.WebsocketConn)  {

}

func (service *userService) GetEnUserId(id int) (string){
	idstr := strconv.Itoa(id)
	return getEnUserId(idstr)
}

func getUserIdByUUID(UUID string) (int) {
	touserid := 0
	for v:= range chat.RoomServe().User {
		if chat.RoomServe().User[v] == UUID {
			touserid = v
		}
	}
	return touserid
}

func GetChatUser(userid,touserid int)([]entity.User){
	userDao := dao.User(app.DB())
	users,_ := userDao.GetUserByID(userid,touserid)
	return users
}

func getDeUserId(userid string) (int) {
	keyStr := app.Config().GetString("server.desKey")
	key := []byte(keyStr)
	idStr,err := utils.Decrypt(userid,key)
	if err != nil {
		fmt.Println("加密的：",idStr)
		//egu.SecurePanic(err)
	}
	id,_ := strconv.Atoi(idStr)
	return id
}

func getEnUserId(id string) (string) {
	keyStr := app.Config().GetString("server.desKey")
	key := []byte(keyStr)
	userid,err := utils.Encrypt(id,key)
	if err != nil {
		egu.SecurePanic(err)
	}
	return userid
}
// Register 注册
//func (service *userService) Register(req request.UserRegisterRequest) error {
//	db := app.DB()
//	userDao := dao.User(db)
//	// 根据邮箱获取用户信息
//	user, err := userDao.GetByEmail(req.Email)
//	if err != nil && err != gorm.ErrRecordNotFound {
//		return errors.New(statusCode.DatabaseQueryFailed, fmt.Sprintf("获取用户信息失败:%v", err))
//	}
//
//	// 用户已存在
//	if user != nil {
//		return errors.New(statusCode.EmailRegistered, "该邮箱已被注册")
//	}
//
//	now := int(egu.GetTimeStamp())
//
//	user = new(entity.UserEntity)
//	user.Email = req.Email
//	user.Password = egu.Md5(req.Pass)
//	user.CreatedAt = now
//	user.UpdatedAt = now
//
//	if err := userDao.Create(user); err != nil {
//		return errors.New(statusCode.DatabaseSaveFailed, fmt.Sprintf("保存数据失败:%v", err))
//	}
//
//	return nil
//}
