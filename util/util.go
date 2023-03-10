package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

func Md5Hex(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

func GetSignature(params map[string]string, key string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var values []string
	for _, k := range keys {
		v := params[k]
		if v == "" || v == "0" {
			continue
		}
		values = append(values, v)
	}
	values = append(values, key)
	text := strings.Join(values, "#")
	return Md5Hex(text)
}

// MakeSessionId 生成Chat会话 SessionId = Scene - SenderId - ChatId
func MakeSessionId(scene, senderId, chatId string) string {
	return fmt.Sprintf("%s-%s-%s", scene, senderId, chatId)
}
