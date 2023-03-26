package openaiproxy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"weichatapp/common"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAIRequestBody struct {
	Messages []openai.ChatCompletionMessage `json:"messages"`
	MaxTokens int `json:"maxTokens"`
}

type OpenAIProxyResponse struct {
	Result openai.ChatCompletionResponse `json:"result"`
	ErrorCode int `json:"errorCode"`
	Message string `json:"message"`
	Error bool `json:"error"`
	Params map[string]interface{} `json:"params"`
}

type BillingHandler interface {
	BillingRecord(sessionid string,amount int)
}

type GPTChatCompletionHandler struct {
	OpenaiproxyConf common.OpenaiproxyConf
	MessageCache *MessageCache
	BillingHandler BillingHandler
	AccountCache *AccountCache
}

func (h *GPTChatCompletionHandler)GTPChatCompletion(maxTokens int, sessionid string,messages []openai.ChatCompletionMessage) string {
	reqBody:=&OpenAIRequestBody{
		Messages:messages,
		MaxTokens:maxTokens,
	}

	//发送http请求
	postJson,_:=json.Marshal(reqBody)
	postBody:=bytes.NewBuffer(postJson)
	log.Println("http.Post ",h.OpenaiproxyConf.Url,string(postJson))
	resp,err:=http.Post(h.OpenaiproxyConf.Url,"application/json",postBody)

	if err != nil || resp==nil || resp.StatusCode != 200 { 
		log.Println("http.Post error",err)
		return "服务器太忙了，请稍后再试"
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println("resp",string(body))
	//识别返回的类型，如果是二进制流，则直接转发到前端

	chatRsp:=&OpenAIProxyResponse{}
	if err := json.Unmarshal(body, chatRsp); err != nil {
			log.Println(err)
			return "服务器太忙了，请稍后再试"
	}

	
	if h.BillingHandler!=nil {
		h.BillingHandler.BillingRecord(sessionid,chatRsp.Result.Usage.TotalTokens)
	}
	choice:=chatRsp.Result.Choices[0]
	return choice.Message.Content
}

func (h *GPTChatCompletionHandler)ProxyChatCompletion(
	sessionid string,messages []openai.ChatCompletionMessage)(string){
	//检查账户余额
	maxTokens:=h.AccountCache.GetToken(sessionid)
	if maxTokens<=0 {
		return "您的账户余额不足，请充值后再试"
	}

	if maxTokens>h.OpenaiproxyConf.MaxTokens {
		maxTokens=h.OpenaiproxyConf.MaxTokens
	}

	answer:=h.GTPChatCompletion(maxTokens,sessionid,messages)
	return answer
}

func (h *GPTChatCompletionHandler)ChatCompletion(sessionid,question string)(string){
	//检查账户余额
	maxTokens:=h.AccountCache.GetToken(sessionid)
	if maxTokens<=0 {
		return "您的账户余额不足，请充值后再试"
	}

	if maxTokens>h.OpenaiproxyConf.MaxTokens {
		maxTokens=h.OpenaiproxyConf.MaxTokens
	}

	//用openid查找缓存中的历史消息记录
	messages:=h.MessageCache.GetMessages(sessionid)
	if messages==nil {
		messages=[]openai.ChatCompletionMessage{}
	}
	//将新消息加入到历史消息记录中
	messages=append(messages,openai.ChatCompletionMessage{
		Content:question,
		Role:openai.ChatMessageRoleUser,
	})
	answer:=h.GTPChatCompletion(maxTokens,sessionid,messages)
	//将新消息加入到缓存中
	messages=append(messages,openai.ChatCompletionMessage{
		Content:answer,
		Role:openai.ChatMessageRoleAssistant,
	})
	h.MessageCache.SetMessages(sessionid,messages)
	return answer
}