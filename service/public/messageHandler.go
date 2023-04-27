package public

import (
	"log"
)

const (
	//消息类型
	MsgTypeText = "text"
	MsgTypeEvent = "event"
	MsgTypeNews = "news"
)

const (
	EventTypeSubscribe = "subscribe"
	EventTypeUnsubscribe = "unsubscribe"
)

type MessageRequest struct {
	ToUserName string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime string `xml:"CreateTime"`
	MsgType string `xml:"MsgType"`
	Event string `xml:"Event"`
	EventKey string `xml:"EventKey"`
	Content string `xml:"Content"`
	MsgId string `xml:"MsgId"`
	MsgDataId string `xml:"MsgDataId"`
	Idx string `xml:"Idx"`
	CustomerInfo *UserInfoResponse `json:"CustomerInfo"`
}

type MessageHandler interface {
	HandleMessage(msg *MessageRequest)
}

func DealMessage(
	msg *MessageRequest,
	chatCompletionHandler ChatCompletionHandler,
	redirectClient RedirectClient,
	customerRepo CustomerRepo){
	//调用openai进行问答
	msgHandler:=GetMessageHandler(msg,chatCompletionHandler,redirectClient,customerRepo)
	if(msgHandler!=nil){
		msgHandler.HandleMessage(msg)
	}
}

func GetMessageHandler(
	msg *MessageRequest,
	chatCompletionHandler ChatCompletionHandler,
	redirectClient RedirectClient,
	customerRepo CustomerRepo)(MessageHandler){
	if msg.MsgType==MsgTypeText {
		return &TextMessageHandler{
			ChatCompletionHandler:chatCompletionHandler,
		}
	} else if msg.MsgType==MsgTypeEvent {
		//if msg.Event==EventTypeSubscribe||msg.Event==EventTypeUnsubscribe	 {
			return &RedirectMessageHandler{
			Client:redirectClient,
			CustomerRepo:customerRepo,
			}
		//}
	}
		
	log.Printf("not supported message type: %s",msg.MsgType)

	return nil
}