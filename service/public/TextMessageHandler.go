package public

import (
	"weichatapp/common"
)

const (
	MSG_GETID="我的ID"
)

type TextMessageHandler struct {
	OpenaiproxyConf common.OpenaiproxyConf
	ChatCompletionHandler ChatCompletionHandler
}

type ChatCompletionHandler interface {
	ChatCompletion(sessionid,question string) string
}

func (h *TextMessageHandler)HandleMessage(msg *MessageRequest){
	var answer string
	if msg.Content==MSG_GETID {
		answer=msg.FromUserName
	} else {
		answer=h.ChatCompletionHandler.ChatCompletion(msg.FromUserName,msg.Content)
	}
	//发送微信客户消息
	PostCSMessage(msg.FromUserName,answer)
}