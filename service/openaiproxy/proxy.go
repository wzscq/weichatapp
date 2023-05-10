package openaiproxy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"weichatapp/common"
	"fmt"
	"io"
	"context"
	"errors"
	"strings"
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

func (h *GPTChatCompletionHandler)ProxyStreamCompletion(sessionid string,messages []openai.ChatCompletionMessage,flusher http.Flusher){
	w:=flusher.(io.Writer)
	//检查账户余额
	maxTokens:=h.AccountCache.GetToken(sessionid)
	if maxTokens<=0 {
		fmt.Fprintf(w, "data: 您的账户余额不足，请充值后再试\n\n")
		fmt.Fprintf(w, "event: close")
		return
	}
	
	if maxTokens>h.OpenaiproxyConf.MaxTokens {
		maxTokens=h.OpenaiproxyConf.MaxTokens
	}

	//调用openai接口
	client := openai.NewClient(h.OpenaiproxyConf.Key)
	ctx := context.Background()
	gptReq := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		Messages: messages,
		MaxTokens: maxTokens,
		Stream: true,
	}
	log.Println("ProxyStreamCompletion CreateChatCompletionStream")
	stream, err := client.CreateChatCompletionStream(ctx, gptReq)
	if err != nil {
		fmt.Printf("ProxyStreamCompletion error: %v\n", err)
		fmt.Fprintf(w, "data: 访问GPT接口出错，请稍后重试，或联系管理员处理\n\n")
		fmt.Fprintf(w, "event: close")
		return
	}
	defer stream.Close()

	totalTokens:=0
	log.Println("ProxyStreamCompletion receive response")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			break;
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			break;
		}

		totalTokens=totalTokens+1
		contentStr:=response.Choices[0].Delta.Content
		//将内容中的\n替换为<br/>
		contentStr=strings.Replace(contentStr,"\n","<br/>",-1)
		log.Println(contentStr)
		fmt.Fprintf(w, "data: "+contentStr+"\n\n")
		flusher.Flush()
	}

	fmt.Fprintf(w, "event: close")

	//
	if h.BillingHandler!=nil {
		h.BillingHandler.BillingRecord(sessionid,totalTokens)
	}
}