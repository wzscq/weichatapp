package openaiproxy

import (
	"github.com/go-redis/redis/v8"
	openai "github.com/sashabaranov/go-openai"
	"time"
	"log"
	"encoding/json"
)

type MessageCache struct {
	Server string
	Password string 
	DB int 
	Expire time.Duration
	Count int
	client *redis.Client
}

func (this *MessageCache)GetRedisClient()(*redis.Client){
	if this.client==nil{
		this.client=redis.NewClient(&redis.Options{
			Addr:     this.Server,
			Password: this.Password, 
			DB:       this.DB,  
		})
	}
	//初始化redis
	return this.client
}

func (this *MessageCache)SetMessages(key string,messages []openai.ChatCompletionMessage){
	bytes, err := json.Marshal(messages)
	if err!=nil {
		log.Println(err)
		return
	}
  // Convert bytes to string.
  jsonStr := string(bytes)

	client:=this.GetRedisClient()
	defer client.Close()
	err= client.Set(client.Context(), key,jsonStr, this.Expire).Err()
  if err!=nil {
		log.Println(err)
	}
}

func (this *MessageCache)GetMessages(key string)([]openai.ChatCompletionMessage){
	client:=this.GetRedisClient()
	defer client.Close()

	val, err := client.Get(client.Context(), key).Result()
	if err != nil {
		log.Println(err)
		return nil
	}

	var messages []openai.ChatCompletionMessage
	err = json.Unmarshal([]byte(val), &messages)
	if err != nil {
		log.Println(err)
	}
	return messages
}