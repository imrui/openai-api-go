package config

import (
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var (
	Cfg *Config
	DB  *gorm.DB
)

type LarkConfig struct {
	UrlPath           string `mapstructure:"url-path"`
	AppId             string `mapstructure:"app-id"`
	AppName           string `mapstructure:"app-name"`
	AppSecret         string `mapstructure:"app-secret"`
	VerificationToken string `mapstructure:"verification-token"`
	EncryptKey        string `mapstructure:"encrypt-key"`
	LarkHost          string `mapstructure:"lark-host"`
}

type Config struct {
	ServerAddr        string            `mapstructure:"server-addr"`
	OpenAiApiKey      string            `mapstructure:"openai-api-key"`
	OpenAiModel       string            `mapstructure:"openai-model"`
	OpenAiMaxToken    int               `mapstructure:"openai-max-token"`
	OpenAiTemperature float32           `mapstructure:"openai-temperature"`
	ClientIdKeys      map[string]string `mapstructure:"client-id-keys"`
	ApiSignEnable     bool              `mapstructure:"api-sign-enable"`
	Scenes            []string          `mapstructure:"scenes"`
	SceneDeleteTexts  map[string]string `mapstructure:"scene-delete-texts"`
	LarkConfigs       []LarkConfig      `mapstructure:"lark-bots"`
}

func (c *Config) GetClientKey(id string) (key string, ok bool) {
	key, ok = c.ClientIdKeys[id]
	return
}

func (c *Config) IsSceneAllow(scene string) bool {
	for _, val := range c.Scenes {
		if val == scene {
			return true
		}
	}
	return false
}

func (c *Config) GetSceneDeleteText(scene string) (text string, ok bool) {
	text, ok = c.SceneDeleteTexts[scene]
	return
}

func init() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("[err] load config err:", err)
		return
	}
	c := &Config{}
	err = viper.Unmarshal(c)
	if err != nil {
		log.Fatal("[err] parse config err:", err)
		return
	}
	Cfg = c
	log.Println("[info] load config:", Cfg)

	db, err := gorm.Open(sqlite.Open("app.db.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatal("[err] db err:", err)
	}
	DB = db
}
