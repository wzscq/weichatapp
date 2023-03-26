package public

const (
	//消息类型
	MsgTypeText = "text"
	MsgTypeSubscribe = "subscribe"
	MsgTypeUnsubscribe = "unsubscribe"
)

type MessageRequest struct {
	ToUserName string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime string `xml:"CreateTime"`
	MsgType string `xml:"MsgType"`
	Content string `xml:"Content"`
	MsgId string `xml:"MsgId"`
	MsgDataId string `xml:"MsgDataId"`
	Idx string `xml:"Idx"`
}

type MessageHandler interface {
	HandleMessage(msg *MessageRequest)
}

func DealMessage(
	msg *MessageRequest,
	chatCompletionHandler ChatCompletionHandler,
	redirectClient RedirectClient){
	//调用openai进行问答
	msgHandler:=GetMessageHandler(msg,chatCompletionHandler,redirectClient)
	if(msgHandler!=nil){
		msgHandler.HandleMessage(msg)
	}
}

func GetMessageHandler(
	msg *MessageRequest,
	chatCompletionHandler ChatCompletionHandler,
	redirectClient RedirectClient)(MessageHandler){
	if msg.MsgType==MsgTypeText {
		return &TextMessageHandler{
			ChatCompletionHandler:chatCompletionHandler,
		}
	} else if msg.MsgType==MsgTypeSubscribe||
		msg.MsgType==MsgTypeUnsubscribe	 {
		return &RedirectMessageHandler{
			Client:redirectClient,
		}
	}
	return nil
}