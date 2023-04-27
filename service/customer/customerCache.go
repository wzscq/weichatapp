package customer

import (
	"github.com/go-redis/redis/v8"
	"log"
	"encoding/json"
)

type CustomerCache struct {
	Server string
	Password string 
	DB int 
	client *redis.Client
}

func (this *CustomerCache)GetRedisClient()(*redis.Client){
	//if this.client==nil{
		this.client=redis.NewClient(&redis.Options{
			Addr:     this.Server,
			Password: this.Password, 
			DB:       this.DB,  
		})
	//}
	//初始化redis
	return this.client
}

func (this *CustomerCache)GetCustomerInfo(openid string)(*CustomerInfo){
	client:=this.GetRedisClient()
	defer client.Close()

	//获取当前已使用的token
	usedInfo, err:=client.Get(client.Context(), openid).Result()
	if err != nil {
		log.Println(err)
		return nil
	}

	//usedInfo转换为CustomerInfo
	customerInfo:=&CustomerInfo{}
	err=json.Unmarshal([]byte(usedInfo),customerInfo)
	if err != nil {
		log.Println(err)
		return nil
	}

	return customerInfo
}

func (this *CustomerCache)SaveCustomerInfo(customerInfo *CustomerInfo){
	client:=this.GetRedisClient()
	defer client.Close()

	//customerInfo转换为JSON字符串
	customerInfoStr,err:=json.Marshal(customerInfo)
	if err != nil {
		log.Println(err)
		return
	}

	//将字符串保存到redis
	err=client.Set(client.Context(), customerInfo.OpenID, string(customerInfoStr), 0).Err()
}