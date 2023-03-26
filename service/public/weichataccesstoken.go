package public

import (
	"net/http"
	"time"
	"log"
	"encoding/json"
	"io/ioutil"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
}

var gAccessToken AccessTokenResponse

func UpdateTokenRoutine(appID,appSecret string){
	for{
		UpdateAccessToken(appID,appSecret)
		if(gAccessToken.ExpiresIn==0){
			return
		}
		duration := time.Duration(gAccessToken.ExpiresIn-20) * time.Second
		time.Sleep(duration)
	}
}

func UpdateAccessToken(appID,appSecret string){
	//GET 
	GetAccessTokenUrl := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid="+appID+"&secret="+appSecret
	resp,err:=http.Get(GetAccessTokenUrl)
	if err!=nil{
		log.Println(err)
		return
	}

	//解析结果
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println("resp",string(body))
	//识别返回的类型，如果是二进制流，则直接转发到前端

	at:=AccessTokenResponse{}
	if err := json.Unmarshal(body, &at); err != nil {
			log.Println(err)
			return
	}

	gAccessToken=at
}

func GetAccessToken() string{
	return gAccessToken.AccessToken
}