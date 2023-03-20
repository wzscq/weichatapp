package main

import (
	"log"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"weichatapp/common"
	"weichatapp/public"
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

	router := gin.Default()
	//允许跨域访问
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:true,
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	//初始化openai代理控制器
	publicController:=public.PublicController{
		Token:conf.Public.Token,
		OpenaiproxyConf:conf.Openaiproxy,
	}

	//绑定路由
	publicController.Bind(router)

	//启动服务
	router.Run(conf.Service.Port)
}