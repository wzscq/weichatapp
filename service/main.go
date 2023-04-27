package main

import (
	"log"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"weichatapp/common"
	"weichatapp/mqtt"
	"weichatapp/bill"
	"weichatapp/public"
	"weichatapp/openaiproxy"
	"weichatapp/customer"
	"time"
)

func main() {
	//设置log打印文件名和行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//初始化时区
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone

	confFile:="conf/conf.json"
	if len(os.Args)>1 {
			confFile=os.Args[1]
			log.Println(confFile)
	}

	//初始化配置
	conf:=common.InitConfig(confFile)

	//初始化消息重定向处理器
	redirectClient:=&mqtt.MQTTRedirectClient{
		Topic:conf.MQTT.RedirectTopic,
		Broker:conf.MQTT.Broker,
		User:conf.MQTT.User,
		Password:conf.MQTT.Password,
		ClientID:conf.MQTT.ClientID,
	}

	expired,_:=time.ParseDuration(conf.MessageCache.Expire)

	accountCache:=&openaiproxy.AccountCache{
		Server:conf.AccountCache.Server,
		Password:conf.AccountCache.Password, 
		DB:conf.AccountCache.DB,
	}

	customerCache:=&customer.CustomerCache{
		Server:conf.CustomerCache.Server,
		Password:conf.CustomerCache.Password, 
		DB:conf.CustomerCache.DB,
	}

	customerRepo:=&customer.CustomerRepo{
		CustomerCache:customerCache,
	}

	customerController:=&customer.CustomerController{
		CustomerCache:customerCache,
	}

	messageCache:=&openaiproxy.MessageCache{
		Server:conf.MessageCache.Server,
		Password:conf.MessageCache.Password, 
		DB:conf.MessageCache.DB,
		Expire:expired,
		Count:conf.MessageCache.Count,
	}

	billSender:=&bill.BillSender{
		Topic:conf.MQTT.BillTopic,
		Broker:conf.MQTT.Broker,
		User:conf.MQTT.User,
		Password:conf.MQTT.Password,
		ClientID:conf.MQTT.ClientID,
	}

	chatCompletionHandler:=&openaiproxy.GPTChatCompletionHandler{
		OpenaiproxyConf:conf.Openaiproxy,
		MessageCache:messageCache,
		BillingHandler:billSender,
		AccountCache:accountCache,
	}

	proxyController:=&openaiproxy.OpenaiproxyController{
		ChatCompletionHandler:chatCompletionHandler,
	}

	//初始化openai代理控制器
	publicController:=public.PublicController{
		Token:conf.Public.Token,
		ChatCompletionHandler:chatCompletionHandler,
		RedirectClient:redirectClient,
		CustomerRepo:customerRepo,
	}

	router := gin.Default()
	//允许跨域访问
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:true,
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	//绑定路由
	publicController.Bind(router)
	proxyController.Bind(router)
	customerController.Bind(router)
	//启动token刷新任务
	//go public.UpdateTokenRoutine(conf.Public.AppID,conf.Public.Secret)
	//启动服务
	router.Run(conf.Service.Port)
}