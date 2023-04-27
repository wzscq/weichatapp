package public

import (
	"encoding/json"
	"log"
)

type RedirectClient interface {
	SendMessage(sessionid,msg string)
}

type CustomerRepo interface {
	GetCustomerInfo(openid string) *UserInfoResponse 
}

type RedirectMessageHandler struct {
	Client RedirectClient
	CustomerRepo CustomerRepo
}

func (h *RedirectMessageHandler)HandleMessage(msg *MessageRequest){
	//获取客户信息
	customerInfo:=h.CustomerRepo.GetCustomerInfo(msg.FromUserName)
	//如果客户信息为空，则发送授权连接地址给客户
	if customerInfo==nil {
		//发送授权连接地址给客户
		log.Println("发送授权连接地址给客户")
		PostAuthMessage(msg.FromUserName)
		return
	}
	//将msg参数转换为json
	msgJson,_:=json.Marshal(msg)
	h.Client.SendMessage(msg.FromUserName,string(msgJson))
}