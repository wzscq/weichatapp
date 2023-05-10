package openaiproxy

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"weichatapp/common"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAIProxyRequestBody struct {
	Messages []openai.ChatCompletionMessage `json:"messages"`
	Sessionid string `json:"sessionid"`
}

type OpenaiproxyController struct {
	ChatCompletionHandler *GPTChatCompletionHandler
}

func (opc *OpenaiproxyController)ChatCompletion(c *gin.Context){
	//获取请求体中携带的消息
	var req OpenAIProxyRequestBody
	if err := c.ShouldBind(&req); err != nil {
		log.Println(err)
		c.String(http.StatusOK, "不知道你在说什么")
		return 
	}
	//处理消息
	log.Println("ChatCompletion Sessionid:",req.Sessionid)
	result:=opc.ChatCompletionHandler.ProxyChatCompletion(req.Sessionid,req.Messages)
	
	//返回结果
	rsp:=common.CreateResponse(nil,result)
	c.IndentedJSON(http.StatusOK, rsp)
}

func (opc *OpenaiproxyController)StreamCompletion(c *gin.Context){
	log.Println("openAIChatStreamGPT4 start...")
	w := c.Writer
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

  flusher, ok := w.(http.Flusher)
	if !ok {
		log.Println("server not support") //浏览器不兼容
		return
	}

	//获取请求体中携带的消息
	var req OpenAIProxyRequestBody
	if err := c.ShouldBind(&req); err != nil {
		log.Println(err)
		fmt.Fprintf(w, "data: 不知道你在说什么\n\n")
		fmt.Fprintf(w, "event: close")
		return 
	}

	//处理消息
	log.Println("ChatCompletion Sessionid:",req.Sessionid)
	opc.ChatCompletionHandler.ProxyStreamCompletion(req.Sessionid,req.Messages,flusher)
}

func (opc *OpenaiproxyController) Bind(router *gin.Engine) {
	router.POST("/openaiproxy/ChatCompletion", opc.ChatCompletion)
	router.POST("/openaiproxy/StreamCompletion", opc.StreamCompletion)
}
