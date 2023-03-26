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
}

func (this *MessageCache)GetRedisClient()(*redis.Client){
	//初始化redis
	return redis.NewClient(&redis.Options{
		Addr:     this.Server,
		Password: this.Password, 
		DB:       this.DB,  
	})
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