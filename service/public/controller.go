package public

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type SignatureRequest struct {
	Signature string `json:"signature"`
	Timestamp string `json:"timestamp"`
	Nonce string `json:"nonce"`
	Echostr string `json:"echostr"`
}

type PublicController struct {
	Token string
}

func (pc *PublicController)checkSignature(c *gin.Context){
	//获取请求体中携带的消息
	var req SignatureRequest
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusOK, nil)
		return 
  }
	
	//校验签名
	if CheckSignature(req.Signature,req.Timestamp,req.Nonce,pc.Token) {
		log.Println("checkSignature success")
		c.IndentedJSON(http.StatusOK, req.Echostr)
		return
	}
	
	//返回结果
	log.Println("checkSignature failed")
	c.IndentedJSON(http.StatusOK, nil)
}

func (opc *PublicController) Bind(router *gin.Engine) {
	router.POST("/public/", opc.checkSignature)
}
