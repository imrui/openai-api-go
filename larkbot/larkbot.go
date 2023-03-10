package larkbot

import (
	"context"
	"encoding/json"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"openai-api-go/config"
	"openai-api-go/service"
	"time"
)

type LarkBot struct {
	Cli             *lark.Client
	editInterval    time.Duration
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
			Cli:             cli,
			EventDispatcher: eventDispatcher,
		}
		eventDispatcher.OnP2MessageReceiveV1(bot.handleOnP2MessageReceiveV1)
		Bots[c.UrlPath] = bot
	}
}

func (b *LarkBot) handleOnP2MessageReceiveV1(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	openId := *event.Event.Sender.SenderId.OpenId
	content := *event.Event.Message.Content
	chatId := *event.Event.Message.ChatId
	var text ContentText
	err := json.Unmarshal([]byte(content), &text)
	if err != nil {
		_, _ = b.Send(openId, "消息识别异常，快去请瑞神！")
		return err
	}
	answer, err := service.Talk(openId, chatId, content, "lark")
	if err != nil {
		_, _ = b.Send(openId, "我emo了，快去请瑞神！")
		return err
	}
	_, err = b.Send(openId, answer)
	return err
}

func (b *LarkBot) Send(openId string, text string) (*larkim.CreateMessageResp, error) {
	content := &ContentText{
		Text: text,
	}
	data, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	return b.Cli.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(openId).
			Content(string(data)).
			Build()).
		Build())
}
