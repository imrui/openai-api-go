package service

import (
	"context"
	"encoding/json"
	"github.com/sashabaranov/go-openai"
	"log"
	"openai-api-go/config"
	"openai-api-go/model"
	"openai-api-go/util"
	"strings"
)

const (
	CommandReset = "/reset"
	CommandClear = "/clear"
)

var client *openai.Client

func init() {
	client = openai.NewClient(config.Cfg.OpenAiApiKey)
}

// CommandProcess 处理指令
func CommandProcess(question string, sessionId string) (bool, string) {
	if !strings.HasPrefix(question, CommandReset) && !strings.HasPrefix(question, CommandClear) {
		return false, ""
	}
	go func() {
		_ = DbDeleteMessageBySessionId(sessionId)
	}()
	return true, "清除上下文成功，你可以继续提问。"
}

func Talk(senderId, chatId, question, scene string) (answer string, err error) {
	sessionId := util.MakeSessionId(scene, senderId, chatId)
	question = strings.TrimSpace(question)
	if stop, stopMsg := CommandProcess(question, sessionId); stop {
		answer = stopMsg
		return
	}
	if text, ok := config.Cfg.GetSceneDeleteText(scene); ok {
		question = strings.ReplaceAll(question, text, "")
	}
	log.Println("[chat] question:", question)
	messages, err := BuildTalkMessages(sessionId, question)
	if err != nil {
		log.Println("[chat] build talk messages err:", err)
		return
	}
	body, err := json.Marshal(messages)
	log.Println("talk message: ", string(body))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    config.Cfg.OpenAiModel,
			Messages: messages,
		},
	)
	if err != nil {
		log.Println("[chat] req openai err:", err)
		return
	}
	answer = strings.TrimSpace(resp.Choices[0].Message.Content)
	log.Println("[chat] answer:", answer)
	m := &model.Message{
		SessionId:        sessionId,
		Scene:            scene,
		SenderId:         senderId,
		ChatId:           chatId,
		Question:         question,
		Answer:           answer,
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
		TotalTokens:      resp.Usage.TotalTokens,
	}
	err = DbAddMessage(m)
	return
}

func BuildTalkMessages(sessionId, question string) (items []openai.ChatCompletionMessage, err error) {
	messages, err := DbFindAllMessages(sessionId)
	if err != nil {
		return
	}
	tokens := 0
	dbMaxId := 0
	i := len(messages) - 1
	for ; i >= 0; i-- {
		m := messages[i]
		if m == nil {
			continue
		}
		tokens += m.TotalTokens
		if tokens >= config.Cfg.OpenAiMaxToken {
			dbMaxId = m.ID
			break
		}
	}
	for _, m := range messages[i+1:] {
		items = append(items, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: m.Question,
		})
		items = append(items, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: m.Answer,
		})
	}
	items = append(items, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: question,
	})
	if dbMaxId > 0 {
		go func() {
			_ = DbDeleteMessageLeID(sessionId, dbMaxId)
		}()
	}
	return
}

func DbFindMessage(id int) (message *model.Message, err error) {
	err = config.DB.First(&message, id).Error
	return
}

func DbFindAllMessages(sessionId string) (messages []*model.Message, err error) {
	err = config.DB.Where(&model.Message{SessionId: sessionId}).Find(&messages).Error
	return
}

func DbAddMessage(message *model.Message) (err error) {
	err = config.DB.Create(&message).Error
	return
}

func DbSaveMessage(message *model.Message) (err error) {
	err = config.DB.Save(&message).Error
	return
}

func DbDeleteMessageBySessionId(sessionId string) (err error) {
	err = config.DB.Delete(&model.Message{SessionId: sessionId}).Error
	return
}

func DbDeleteMessageLeID(sessionId string, maxId int) (err error) {
	err = config.DB.Delete(&model.Message{SessionId: sessionId}, "id <= ?", maxId).Error
	return
}
