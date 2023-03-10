package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"openai-api-go/config"
	"openai-api-go/msg"
	"openai-api-go/service"
)

func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hi, AI.")
}

func ChatTalk(c *gin.Context) {
	var req msg.ChatReqMsg
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if req.ID == "" || req.Ts <= 0 || req.Sign == "" || req.SenderId == "" || req.ChatId == "" || req.Content == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "params error"})
		return
	}
	if !config.Cfg.IsSceneAllow(req.Scene) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "scene error"})
		return
	}
	if config.Cfg.ApiSignEnable {
		key, ok := config.Cfg.GetClientKey(req.ID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "request id error"})
			return
		}
		if !req.SignVerified(key) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "sign error"})
			return
		}
	}
	answer, err := service.Talk(req.SenderId, req.ChatId, req.Content, req.Scene)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := msg.ChatResMsg{
		Code:    200,
		Msg:     "success",
		Content: answer,
	}
	c.JSON(http.StatusOK, res)
}
