package model

import (
	"openai-api-go/config"
	"time"
)

type Model struct {
	ID        int       `json:"id" form:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Message struct {
	Model
	SessionId        string `json:"sessionId"` // 会话ID，生成策略 SessionId = Scene - SenderId - ChatId
	Scene            string `json:"scene"`     // 使用场景
	SenderId         string `json:"senderId"`  // 发送者ID
	ChatId           string `json:"chatId"`    // 聊天ID：用于区分频道、话题等
	Question         string `json:"question"`
	Answer           string `json:"answer"`
	PromptTokens     int    `json:"promptTokens"`
	CompletionTokens int    `json:"completionTokens"`
	TotalTokens      int    `json:"totalTokens"`
}

type HistoryMessage struct {
	Message
}

type LarkEvent struct {
	Model
	AppId    string `json:"appId" gorm:"index"`
	EventId  string `json:"eventId" gorm:"index"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func InitDB() (err error) {
	err = config.DB.AutoMigrate(&Message{}, &HistoryMessage{}, &LarkEvent{})
	return
}
