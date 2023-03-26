package openaiproxy

import (
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

type AccountCache struct {
	Server string
	Password string 
	DB int 
}

func (this *AccountCache)GetRedisClient()(*redis.Client){
	//初始化redis
	return redis.NewClient(&redis.Options{
		Addr:     this.Server,
		Password: this.Password, 
		DB:       this.DB,  
	})
}

func (this *AccountCache)GetToken(key string)(int){
	client:=this.GetRedisClient()
	defer client.Close()

	totalTokenStr, err:=client.Get(client.Context(), key+":totalToken").Result()
	if err != nil {
		log.Println(err)
		return 0
	}

	//将totalToken转换为int类型
	totalToken,err:= strconv.Atoi(totalTokenStr)
	if err != nil {
		log.Println(err)
		return 0
	}

	usedTokenStr, err:=client.Get(client.Context(), key+":usedToken").Result()
	if err != nil {
		log.Println(err)
		return 0
	}

	//将totalToken转换为int类型
	usedToken,err:= strconv.Atoi(usedTokenStr)
	if err != nil {
		log.Println(err)
		return 0
	}

	return totalToken - usedToken
}