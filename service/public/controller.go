package public

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"weichatapp/common"
)

type PublicController struct {
	Token string
	ChatCompletionHandler ChatCompletionHandler
	RedirectClient RedirectClient
}

func (pc *PublicController)checkSignature(c *gin.Context){
	//获取请求体中携带的消息
	var req SignatureRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Println(err)
		c.String(http.StatusOK, "")
		return 
  }
	
	//校验签名
	if CheckSignature(req.Signature,req.Timestamp,req.Nonce,pc.Token) {
		log.Println("checkSignature success")
		c.String(http.StatusOK, req.Echostr)
		return
	}
	
	//返回结果
	log.Println("checkSignature failed")
	c.String(http.StatusOK, "")
}

func (pc *PublicController)normalMessage(c *gin.Context){
	//获取请求体中携带的消息
	var req MessageRequest
	if err := c.ShouldBindXML(&req); err != nil {
		log.Println(err)
		c.String(http.StatusOK, "success")
		return 
	}

	log.Printf("from:%s to: %s type:%s content:%s",req.FromUserName,req.ToUserName,req.MsgType,req.Content)
	
	go DealMessage(
		&req,
		pc.ChatCompletionHandler,
		pc.RedirectClient)
	
	//处理消息
	log.Println("normalMessage")
	c.String(http.StatusOK, "success")
}

func (pc *PublicController)getTicket(c *gin.Context){
	log.Println("getTicket start")
	//获取请求体中携带的消息
	var req GetTicketRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Println(err)
		c.String(http.StatusOK, "")
		return 
	}

	resp,err:=GetTicket(&req)
	if err!=nil {
		rsp:=common.CreateResponse(common.CreateError(common.ResultWeiChatAPIError,nil),resp)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	rsp:=common.CreateResponse(nil,resp)
	c.IndentedJSON(http.StatusOK, rsp)
	//处理消息
	log.Println("getTicket end")
}

func (opc *PublicController) Bind(router *gin.Engine) {
	router.GET("/public/", opc.checkSignature)
	router.POST("/public/",opc.normalMessage)
	router.POST("/public/getTicket",opc.getTicket)
}
