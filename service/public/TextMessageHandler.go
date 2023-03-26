package public

import (
	"weichatapp/common"
)

type TextMessageHandler struct {
	OpenaiproxyConf common.OpenaiproxyConf
	ChatCompletionHandler ChatCompletionHandler
}

type ChatCompletionHandler interface {
	ChatCompletion(sessionid,question string) string
}

func (h *TextMessageHandler)HandleMessage(msg *MessageRequest){
	answer:=h.ChatCompletionHandler.ChatCompletion(msg.FromUserName,msg.Content)
	//发送微信客户消息
	PostCSMessage(msg.FromUserName,answer)
}