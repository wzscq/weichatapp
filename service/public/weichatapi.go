package public

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"crypto/sha1"
	"fmt"
)

type SignatureRequest struct {
	Signature string `form:"signature"`
	Timestamp string `form:"timestamp"`
	Nonce string `form:"nonce"`
	Echostr string `form:"echostr"`
}

type CSMessageText struct {
	Content string `json:"content"`
}

type CSMessageRequest struct {
	ToUser string `json:"touser"`
	MsgType string `json:"msgtype"`
	Text *CSMessageText `json:"text"`
}

type CSMessageResponse struct {
	ErrCode int `json:"errcode"`
	ErrMsg string `json:"errmsg"`
}

type Scene struct {
	SceneID int64 `json:"scene_id"`
}

type ActionInfo struct {
	Scene  Scene  `json:"scene"`
}

type GetTicketRequest struct {
	ExpireSeconds int `json:"expire_seconds"`
	ActionName string `json:"action_name"`
	ActionInfo  ActionInfo `json:"action_info"`
}

type GetTicketResponse struct {
	Ticket string `json:"ticket"`
	ExpireSeconds int `json:"expire_seconds"`
	URL string `json:"url"`
	SceneID int64 `json:"sceneID"`
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

func GetTicket(reqBody *GetTicketRequest)(*GetTicketResponse,error){
	reqBody.ActionInfo.Scene.SceneID=GetSceneID()
	//发送http请求
	ticketUrl:="https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="+GetAccessToken()
	postJson,_:=json.Marshal(reqBody)
	postBody:=bytes.NewBuffer(postJson)
	log.Println("http.Post ",ticketUrl,string(postJson))
	resp,err:=http.Post(ticketUrl,"application/json",postBody)

	if err != nil || resp==nil || resp.StatusCode != 200 { 
		log.Println("http.Post error",err)
		return nil,err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error",err)
		return nil,err
	}

	var ticketResponse GetTicketResponse
	err=json.Unmarshal(body,&ticketResponse)
	if err != nil {
		log.Println("json.Unmarshal error",err)
		return nil,err
	}

	ticketResponse.SceneID=reqBody.ActionInfo.Scene.SceneID

	return &ticketResponse,nil
}

func PostCSMessage(openid,message string){
	reqBody:=&CSMessageRequest{
		ToUser:openid,
		MsgType:MsgTypeText,
		Text:&CSMessageText{
			Content:message,
		},
	}
	//发送http请求
	csUrl:="https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="+GetAccessToken()
	postJson,_:=json.Marshal(reqBody)
	postBody:=bytes.NewBuffer(postJson)
	log.Println("http.Post ",csUrl,string(postJson))
	resp,err:=http.Post(csUrl,"application/json",postBody)

	if err != nil || resp==nil || resp.StatusCode != 200 { 
		log.Println("http.Post error",err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error",err)
	}
	log.Println("http.Post response",string(body))
}

type UserAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID string `json:"openid"`
	Scope string `json:"scope"`
	Errcode int `json:"errcode"`
	Errmsg string `json:"errmsg"`
}

func RefreshAccessToken(refreshAccessToken string)(*UserAccessTokenResponse){
	//发送http请求
	accessTokenUrl:="https://api.weixin.qq.com/sns/oauth2/refresh_token?appid="+GetAppID()+"&grant_type=refresh_token&refresh_token="+refreshAccessToken
	log.Println("http.Get ",accessTokenUrl)
	resp,err:=http.Get(accessTokenUrl)

	if err != nil || resp==nil || resp.StatusCode != 200 { 
		log.Println("http.Get error",err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error",err)
		return nil
	}

	var accessTokenResponse UserAccessTokenResponse
	err=json.Unmarshal(body,&accessTokenResponse)
	if err != nil {
		log.Println("json.Unmarshal error",err)
		return nil
	}

	return &accessTokenResponse
}

type UserInfoResponse struct {
	OpenID string `json:"openid"`
	NickName string `json:"nickname"`
	Sex int `json:"sex"`
	Province string `json:"province"`
	City string `json:"city"`
	Country string `json:"country"`
	HeadImgURL string `json:"headimgurl"`
	Privilege []string `json:"privilege"`
	UnionID string `json:"unionid"`
}

func RefreshCustomerInfo(openid,accessToken string)(*UserInfoResponse){
	//发送http请求
	customerInfoUrl:="https://api.weixin.qq.com/sns/userinfo?access_token="+accessToken+"&openid="+openid+"&lang=zh_CN"
	log.Println("http.Get ",customerInfoUrl)
	resp,err:=http.Get(customerInfoUrl)

	if err != nil || resp==nil || resp.StatusCode != 200 { 
		log.Println("http.Get error",err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error",err)
		return nil
	}

	var customerInfo UserInfoResponse
	err=json.Unmarshal(body,&customerInfo)
	if err != nil {
		log.Println("json.Unmarshal error",err)
		return nil
	}

	return &customerInfo
}

type Article struct {
	Title string `json:"title"`
	Description string `json:"description"`
	URL string `json:"url"`
	PicURL string `json:"picurl"`
}

type NewsMessage struct {
	Articles []Article `json:"articles"`
}

type CSNewsMessageRequest struct {
	ToUser string `json:"touser"`
	MsgType string `json:"msgtype"`
  News *NewsMessage `json:"news"`
}

func PostAuthMessage(openid string){
	redirectUrl:="http://vipgpt.top/weichatauth"
	reqBody:=&CSNewsMessageRequest{
		ToUser:openid,
		MsgType:MsgTypeNews,
		News:&NewsMessage{
			Articles:[]Article{
				Article{
					Title:"微信授权",
					Description:"点击进入微信授权，授权后即可开始使用",
					URL:"https://open.weixin.qq.com/connect/oauth2/authorize?appid="+GetAppID()+"&redirect_uri="+redirectUrl+"&response_type=code&scope=snsapi_userinfo&state=123#wechat_redirect",
				},
			},
		},
	}
	//发送http请求
	csUrl:="https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="+GetAccessToken()
	postJson,_:=json.Marshal(reqBody)
	postBody:=bytes.NewBuffer(postJson)
	log.Println("http.Post ",csUrl,string(postJson))
	resp,err:=http.Post(csUrl,"application/json",postBody)

	if err != nil || resp==nil || resp.StatusCode != 200 { 
		log.Println("http.Post error",err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error",err)
	}
	log.Println("http.Post response",string(body))
}

func GetUserAccessTokenByCode(code string)(*UserAccessTokenResponse){
	//发送http请求
	accessTokenUrl:="https://api.weixin.qq.com/sns/oauth2/access_token?appid="+GetAppID()+"&secret="+GetAppSecret()+"&code="+code+"&grant_type=authorization_code"
	log.Println("http.Get ",accessTokenUrl)
	resp,err:=http.Get(accessTokenUrl)

	if err != nil || resp==nil || resp.StatusCode != 200 { 
		log.Println("http.Get error",err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error",err)
		return nil
	}

	var accessTokenResponse UserAccessTokenResponse
	err=json.Unmarshal(body,&accessTokenResponse)
	if err != nil {
		log.Println("json.Unmarshal error",err)
		return nil
	}

	if accessTokenResponse.Errcode != 0 {
		log.Println("Errcode:",accessTokenResponse.Errcode)
		log.Println("Errmsg:",accessTokenResponse.Errmsg)
		return nil
	}

	return &accessTokenResponse
}

