package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"openai-api-go/config"
	"openai-api-go/handler"
	"openai-api-go/model"
)

func main() {
	err := model.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.GET("/", handler.Index)
	r.GET("/ai/chat/", handler.Index)
	r.POST("/ai/chat/api/talk", handler.ChatTalk)
	log.Fatal(r.Run(config.Cfg.ServerAddr))
}
