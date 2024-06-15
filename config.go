/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     rose-mqtt
 * @Date:        2024-06-15 11:33
 * @Description:
 */

package rmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type Config struct {
	debug         bool
	clientId      string // 客户端 Id
	cleanSession  bool   // 断线后是否清理会话
	username      string // 用户名
	password      string // 用户密码
	keepAlive     time.Duration
	waitTimeout   time.Duration // 等待时间
	autoReconnect bool          // 是否自动重连

	defaultHandler  mqtt.MessageHandler        // 发布消息回调
	connHandler     mqtt.OnConnectHandler      // 连接回调
	connLostHandler mqtt.ConnectionLostHandler // 连接意外中断回调

}

type Option func(c *Config)

/* WithXXX 适用于特定的简单参数 */

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

func WithClientID(id string) Option {
	return func(c *Config) {
		c.clientId = id
	}
}

func NewConfig(options ...Option) *Config {
	// 默认值
	cfg := &Config{
		cleanSession:  true,
		autoReconnect: true,
		clientId:      generateRandomClientID(),
	}
	// 初始配置
	for _, opt := range options {
		opt(cfg)
	}
	return cfg
}

/* SetXXX 适用于参数或方法 */

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

func (c *Config) SetUserAndPwd(userName, passWord string) *Config {
	c.username = userName
	c.password = passWord
	return c
}

func (c *Config) SetTLSConfig() {

}

func (c *Config) SetDefaultPublishHandler(handler mqtt.MessageHandler) *Config {
	c.defaultHandler = handler
	return c
}

func (c *Config) SetConnHandler(handler mqtt.OnConnectHandler) *Config {
	c.connHandler = handler
	return c
}

func (c *Config) SetConnLostHandler(handler mqtt.ConnectionLostHandler) *Config {
	c.connLostHandler = handler
	return c
}
