package public

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type PublicController struct {
	Token string
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
	var req NomalMessageRequest
	if err := c.ShouldBindXML(&req); err != nil {
		log.Println(err)
		c.String(http.StatusOK, "success")
		return 
	}
	
	//调用openai进行问答
	answer:="目前仅支持文本消息"
	if (req.MsgType==MsgTypeText){
		answer=ChatCompletion(req.Content)
	}

	//返回结果
	resp:=CreateTextResponse(&req,answer)
	
	//处理消息
	log.Println("normalMessage")
	c.XML(http.StatusOK, resp)
}

func (opc *PublicController) Bind(router *gin.Engine) {
	router.GET("/public/", opc.checkSignature)
	router.POST("/public/",opc.normalMessage)
}
