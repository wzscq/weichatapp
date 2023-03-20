package public

import (
	openai "github.com/sashabaranov/go-openai"
	"weichatapp/common"
	"log"
	"net/http"
	"io/ioutil"
	"bytes"
	"encoding/json"
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

func ChatCompletion(question string,openaiproxyConf common.OpenaiproxyConf) string {
	reqBody:=&OpenAIRequestBody{
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question,
			},
		},
		MaxTokens:openaiproxyConf.MaxTokens,
	}

	//发送http请求
	postJson,_:=json.Marshal(reqBody)
	postBody:=bytes.NewBuffer(postJson)
	log.Println("http.Post ",openaiproxyConf.Url,string(postJson))
	resp,err:=http.Post(openaiproxyConf.Url,"application/json",postBody)

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

	return chatRsp.Result.Choices[0].Message.Content
}