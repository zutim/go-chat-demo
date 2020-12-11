package chat

import (
	"encoding/json"
	"fmt"
	"github.com/ebar-go/ego/app"
	"github.com/ebar-go/egu"
	"strconv"
	"time"
)
//
//type ChatServer interface {
//	setChatRecord(chat *Chat)
//	getChatRecord(chat *Chat)
//	getUnreadMsgCount()
//}

type Data struct {
	From int
	To int
	Message string
	Sent time.Time
	recd string
	EnTo string //接收方加密字段
}
var data *Data

func init()  {
	data = newData()
}

func newData()  *Data{
	return &Data{
		From:0,
		To:0,
		Message:"",
		Sent:time.Now(),
		recd:"0_/",
		EnTo:"",
	}
}
func DataServe()  *Data{
	return data
}

/**
发送消息时保存聊天记录
 */
func (d *Data)SetChatRecord() (int64) {
	keyName := fmt.Sprintf("rec:%s",getRecKeyName(d.From,d.To))

	msg,_ := json.Marshal(d)

	res,err := app.Redis().LPush(keyName,string(msg)).Result()
	if err== nil {
		egu.SecurePanic(err)
	}

	if b:=chechUserIsOnLine(d.To);!b{
		fmt.Println("接收方不在线")
		d.cacheUnreadMsg()
	}

	return res
}

/**
获取聊天记录
 */
func (d *Data)GetChatRecord(num int64) ([]string) {
	keyName := fmt.Sprintf("rec:%s",getRecKeyName(d.From,d.To))
	recList,err := app.Redis().LRange(keyName,0,num).Result()
	if err !=nil{
		egu.SecurePanic(err)
	}
	return recList
}

/**
获取未读消息的内容
 */
func (d *Data)GetUnreadMsg() ([]string){
	count := getUnreadMsgCount(d.To)
	keyName := fmt.Sprintf("rec:%s",getRecKeyName(d.From,d.To))
	s := strconv.Itoa(count)
	int64obj,_:= strconv.ParseInt(s,10,64)
	fmt.Println(getRecKeyName(d.From,d.To))
	msg, err := app.Redis().LRange(keyName,0,int64obj).Result()
	if err != nil {
		egu.SecurePanic(err)
	}
	return msg
}

/**
* 将消息设为已读
 */
func (d *Data)SetUnreadToRead() (int)  {
	to := strconv.Itoa(d.To)
	_,err :=app.Redis().HDel(fmt.Sprintf("unread_%d",d.From),to).Result()
	if err != nil {
		egu.SecurePanic(err)
	}
	return 1
}

/**
当用户不在线时，或者当前没有立刻接收消息时，缓存未读消息,将未读消息的数目和发送者信息存到一个与接受者关联的hash数据中
 */
func (d *Data)cacheUnreadMsg(){
	from2 := strconv.Itoa(d.From)
	app.Redis().HIncrBy(fmt.Sprintf("unread_%d",d.To),from2,1)
}

func (d *Data)GetUnreadMsgs()  (map[string]string){
	data,err := app.Redis().HGetAll(fmt.Sprintf("unread_%d",d.To)).Result()
	if err != nil {
		egu.SecurePanic(err)
	}
	return data
}

/**
生成聊天记录的键名，即按大小规则将两个数字排序
 */
func getRecKeyName(from int,to int) (string)  {
	res :=""
	if from>to {
		res = fmt.Sprintf("%d_%d",to,from)
	}else {
		res = fmt.Sprintf("%d_%d",from,to)
	}
	return res
}


/**
当用户上线时，或点开聊天框时，获取未读消息的数目
*/
func getUnreadMsgCount(userid int) (int){
	unRead,err:=app.Redis().HGetAll(fmt.Sprintf("unread_%d",userid)).Result()
	if err != nil {
		egu.SecurePanic(err)
	}
	return len(unRead)
}

/**
验证用户是否是在线状态
 */
func chechUserIsOnLine(to int) (bool){
	online := false
	key := RoomServe().User[to]
	if  key!=""{
		online = true
	}
	return online
}


