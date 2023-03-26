package openaiproxy

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"weichatapp/common"
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

	result:=opc.ChatCompletionHandler.ProxyChatCompletion(req.Sessionid,req.Messages)
	//处理消息
	log.Println("ChatCompletion")
	//返回结果
	rsp:=common.CreateResponse(nil,result)
	c.IndentedJSON(http.StatusOK, rsp)
}


func (opc *OpenaiproxyController) Bind(router *gin.Engine) {
	router.POST("/openaiproxy/ChatCompletion", opc.ChatCompletion)
}
