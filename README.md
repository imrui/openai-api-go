# openai-api-go

基于 GPT-3.5-Turbo 模型的聊天服务，使用 OpenAI 的 AIP keys。

提供通用HTTP服务，便于集成到其他系统之中，如：Telegram、Lark/飞书、企业微信等各种机器人场景。

## HTTP API

### POST `/ai/chat/api/talk` 聊天对话

**Headers:**

Content-Type: application/json

**Body:**

| 名称       | 类型     | 备注                     |
|----------|--------|------------------------|
| senderId | string | 发送者ID                  |
| chatId   | string | 聊天ID：用于区分频道、话题等        |
| msgId    | string | 消息ID                   |
| content  | string | 消息内容                   |
| id       | string | 请求者ID                  |
| scene    | string | 使用场景：lark/tg/wx/qq/web |
| ts       | int64  | 时间戳 秒                  |
| sign     | string | 签名                     |

**签名算法：**

1. 所有不为空、不为0的参数都需要加入签名，参数值必须为`UrlEncode`之前的原始数值。如参数`中文`，作为参数传输时编码为`%e4%b8%ad%e6%96%87`，签名计算时则要用其原始值中文(注意字符集必须是`UTF-8`)
2. 对所有待签名参数，按照参数名字母升序排列。
3. 使用符号#拼装排序后的参数值，最后用#连接KEY，得到签名文本。
4. 对签名文本进行MD5计算，就是sign参数的值。MD5值32位小写。
