/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     rose-mqtt
 * @Date:        2024-06-15 11:33
 * @Description:
 */

package rmqtt

import "time"

type Config struct {
	debug        bool
	clientId     string // 客户端 Id
	cleanSession bool   // 断线后是否清理会话
	username     string // 用户名
	password     string // 用户密码
	keepAlive    time.Duration
	waitTime     time.Duration // 等待时间
}

type Option func(c *Config)

func WithDebug(debug bool) Option {
	return func(c *Config) {
		c.debug = debug
	}
}

func WithUserAndPwd(userName, passWord string) Option {
	return func(c *Config) {
		c.username = userName
		c.password = passWord
	}
}

func NewConfig(options ...Option) *Config {
	cfg := &Config{}
	for _, opt := range options {
		opt(cfg)
	}
	return cfg
}

func (c *Config) SetClientID(id string) *Config {
	c.clientId = id
	return c
}

func (c *Config) SetKeepAlive(t time.Duration) *Config {
	c.keepAlive = t
	return c
}

func (c *Config) SetKeepAliveSec(sec int64) *Config {
	c.keepAlive = time.Duration(sec) * time.Second
	return c
}
