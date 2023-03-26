package public

import (
	"encoding/json"
)

type RedirectClient interface {
	SendMessage(sessionid,msg string)
}

type RedirectMessageHandler struct {
	Client RedirectClient
}

func (h *RedirectMessageHandler)HandleMessage(msg *MessageRequest){
	//将msg参数转换为json
	msgJson,_:=json.Marshal(msg)
	h.Client.SendMessage(msg.FromUserName,string(msgJson))
}