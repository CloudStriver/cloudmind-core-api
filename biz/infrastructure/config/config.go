package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"os"

	"github.com/zeromicro/go-zero/core/service"

	"github.com/zeromicro/go-zero/core/conf"
)

var config *Config

type Auth struct {
	SecretKey    string
	PublicKey    string
	AccessExpire int64
}

type OauthConf struct {
	ClientId string
	Secret   string
}

type KqConfig struct {
	Brokers []string
	Topic   string
}

type Config struct {
	service.ServiceConf
	ListenOn              string
	Auth                  Auth
	GithubConf            OauthConf
	GiteeConf             OauthConf
	Redis                 redis.RedisConf
	CreateNotificationsKq KqConfig
	ReadNotificationsKq   KqConfig
	CleanNotificationKq   KqConfig
}

func NewConfig() (*Config, error) {
	c := new(Config)
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "etc/config.yaml"
	}
	err := conf.Load(path, c)
	if err != nil {
		return nil, err
	}
	err = c.SetUp()
	if err != nil {
		return nil, err
	}
	config = c
	return c, nil
}

func GetConfig() *Config {
	return config
}
