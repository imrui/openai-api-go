package msg

import (
	"openai-api-go/util"
	"strconv"
)

type ChatResMsg struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Content string `json:"content"`
}

type ChatReqMsg struct {
	SenderId string `json:"senderId"` // 发送者ID
	ChatId   string `json:"chatId"`   // 聊天ID：用于区分频道、话题等
	MsgId    string `json:"msgId"`    // 消息ID
	Content  string `json:"content"`  // 消息内容
	ID       string `json:"id"`       // 请求者ID
	Scene    string `json:"scene"`    // 使用场景：wx/lark/qq/web
	Ts       int64  `json:"ts"`       // 时间戳 秒
	Sign     string `json:"sign"`     // 签名
}

func (r *ChatReqMsg) SignVerified(key string) bool {
	params := map[string]string{
		"senderId": r.SenderId,
		"chatId":   r.ChatId,
		"msgId":    r.MsgId,
		"content":  r.Content,
		"id":       r.ID,
		"scene":    r.Scene,
		"ts":       strconv.FormatInt(r.Ts, 10),
	}
	signature := util.GetSignature(params, key)
	return r.Sign == signature
}

func (r *ChatReqMsg) MakeSessionId() string {
	return util.MakeSessionId(r.Scene, r.SenderId, r.ChatId)
}
