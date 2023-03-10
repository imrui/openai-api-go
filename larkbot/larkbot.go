package larkbot

import (
	"context"
	"encoding/json"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"log"
	"openai-api-go/config"
	"openai-api-go/model"
	"openai-api-go/service"
)

type LarkBot struct {
	AppId           string
	AppName         string
	Cli             *lark.Client
	EventDispatcher *dispatcher.EventDispatcher
}

type ContentText struct {
	Text string `json:"text"`
}

var Bots map[string]*LarkBot

func init() {
	Bots = make(map[string]*LarkBot)
	for _, c := range config.Cfg.LarkConfigs {
		eventDispatcher := dispatcher.NewEventDispatcher(c.VerificationToken, c.EncryptKey)
		cli := lark.NewClient(c.AppId, c.AppSecret, lark.WithLogReqAtDebug(true), lark.WithLogLevel(larkcore.LogLevelDebug))
		bot := &LarkBot{
			AppId:           c.AppId,
			AppName:         c.AppName,
			Cli:             cli,
			EventDispatcher: eventDispatcher,
		}
		eventDispatcher.OnP2MessageReceiveV1(bot.handleOnP2MessageReceiveV1)
		Bots[c.UrlPath] = bot
	}
}

func (b *LarkBot) handleOnP2MessageReceiveV1(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	eventId := event.EventV2Base.Header.EventID
	if eventId != "" {
		larkEvent, err := service.DbFindLarkEvent(b.AppId, eventId)
		if larkEvent != nil && larkEvent.ID > 0 {
			log.Println("[lark] repeat LarkEvent", err)
			return nil
		}
	}
	return b.doP2MessageReceiveV1(event)
}

func (b *LarkBot) doP2MessageReceiveV1(event *larkim.P2MessageReceiveV1) (err error) {
	openId := *event.Event.Sender.SenderId.OpenId
	messageType := *event.Event.Message.MessageType
	if messageType != larkim.MsgTypeText {
		_ = b.Send(openId, "我还不会其他类型的提问，快去找瑞神给我升级！")
		return
	}
	content := *event.Event.Message.Content
	var text ContentText
	err = json.Unmarshal([]byte(content), &text)
	if err != nil {
		_ = b.Send(openId, "消息识别异常，快去请瑞神！")
		return
	}
	eventId := event.EventV2Base.Header.EventID
	chatId := *event.Event.Message.ChatId
	chatType := *event.Event.Message.ChatType
	// 私聊消息，直接回复
	if chatType == "p2p" {
		answer, err1 := b.talkOpenAI(eventId, openId, chatId, text.Text)
		if err1 != nil {
			answer = "我emo了，快去请瑞神！"
		}
		_ = b.Send(openId, answer)
		return
	}
	// 群聊消息，需要@机器人
	if chatType == "group" {
		// 被提及用户的信息，为空表示日常消息，则忽略
		mentions := event.Event.Message.Mentions
		if mentions == nil || len(mentions) == 0 {
			log.Println("[lark] chat_type group: mentions empty.")
			return
		}
		// 未提及机器人，则忽略
		if *mentions[0].Name != b.AppName {
			log.Println("[lark] chat_type group: mention[0] not bot.")
			return
		}
		answer, err2 := b.talkOpenAI(eventId, openId, chatId, text.Text)
		if err2 != nil {
			answer = "我好像故障了，你们继续聊，我先休息一会！"
		}
		_ = b.Send(openId, answer)
	}
	return
}

func (b *LarkBot) talkOpenAI(eventId, openId, chatId, question string) (answer string, err error) {
	answer, err = service.Talk(openId, chatId, question, "lark")
	if eventId == "" {
		return
	}
	// 异步保存事件信息
	go func() {
		larkEvent := &model.LarkEvent{
			AppId:    b.AppId,
			EventId:  eventId,
			Question: question,
			Answer:   answer,
		}
		_ = service.DbAddLarkEvent(larkEvent)
	}()
	return
}

func (b *LarkBot) Send(openId string, text string) (err error) {
	content := &ContentText{
		Text: text,
	}
	data, err := json.Marshal(content)
	if err != nil {
		return
	}
	_, err = b.Cli.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(openId).
			Content(string(data)).
			Build()).
		Build())
	return
}
