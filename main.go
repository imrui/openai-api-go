package main

import (
	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-gin"
	"log"
	"openai-api-go/config"
	"openai-api-go/handler"
	"openai-api-go/larkbot"
	"openai-api-go/model"
)

func main() {
	err := model.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.GET("/", handler.Index)
	r.POST("/ai/chat/api/talk", handler.ChatTalk)
	for path, bot := range larkbot.Bots {
		r.POST("/bot/lark/webhook/event/"+path, sdkginext.NewEventHandlerFunc(bot.EventDispatcher))
	}
	log.Fatal(r.Run(config.Cfg.ServerAddr))
}
