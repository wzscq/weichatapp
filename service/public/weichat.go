package public

import (
	"crypto/sha1"
	"sort"
	"strings"
	"fmt"
	"log"
	"encoding/xml"
)

const (
	//消息类型
	MsgTypeText = "text"
)

type SignatureRequest struct {
	Signature string `form:"signature"`
	Timestamp string `form:"timestamp"`
	Nonce string `form:"nonce"`
	Echostr string `form:"echostr"`
}

type NomalMessageRequest struct {
	ToUserName string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime string `xml:"CreateTime"`
	MsgType string `xml:"MsgType"`
	Content string `xml:"Content"`
	MsgId string `xml:"MsgId"`
	MsgDataId string `xml:"MsgDataId"`
	Idx string `xml:"Idx"`
}

type TextMessageResponse struct {
	XMLName xml.Name `xml:"xml"`
	ToUserName string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime string `xml:"CreateTime"`
	MsgType string `xml:"MsgType"`
	Content string `xml:"Content"`
}

func CheckSignature(signature string, timestamp string, nonce string,token string)(bool) {
	//将token、timestamp、nonce三个参数放到一个数组中
	strs:=[]string{token,timestamp,nonce}
	//将数组中的元素按照字典序排序
	sort.Strings(strs)
	//将三个参数字符串拼接成一个字符串进行sha1加密
	joinStr:=strings.Join(strs,"")
	log.Println(joinStr)
	//生成签名
	sha1:=sha1.Sum([]byte(joinStr))
	//将签名转为字符串
	sha1Str:=fmt.Sprintf("%x",sha1)
	//将sha1加密后的字符串与signature进行对比
	if sha1Str==signature {
		return true
	}
	log.Printf("checkSignature failed,token:%s timestamp:%s nonce:%s \n",token,timestamp,nonce)
	log.Printf("checkSignature failed,sha1 %s,signature %s \n",sha1Str,signature)
	return false
}

func CreateTextResponse(req *NomalMessageRequest,answer string)(*TextMessageResponse){
	response:=&TextMessageResponse{
		ToUserName:req.FromUserName,
		FromUserName:req.ToUserName,
		CreateTime:req.CreateTime,
		MsgType:MsgTypeText,
		Content:answer,
	}
	return response
}

func TestSignature(){
	signature:="749e13ef2af4b9a97b6355b3a5f03d7b33d13f21"
	token:="2gb6g99dq52gb6g99dq5"
	timestamp:="1679280322"
	nonce:="1106204565"

	CheckSignature(signature,timestamp,nonce,token)
}