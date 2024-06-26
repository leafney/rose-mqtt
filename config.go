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
	clientId      string        // 客户端 Id
	cleanSession  bool          // 断线后是否清理会话
	username      string        // 用户名
	password      string        // 用户密码
	keepAlive     time.Duration // 保活时间
	waitTimeout   time.Duration // 等待时间
	reconnectType ReConnType    // 自动重连选项

	//willEnabled     bool          // 遗嘱消息
	willTopic         string
	willPayload       []byte
	willQos           QosLevel
	willPayloadIsByte bool
	willRetained      bool

	// tls
	tlsVerify         bool   // 安全性验证
	tlsCaCert         string // CA证书
	tlsCaCertFile     string
	tlsClientCert     string
	tlsClientCertFile string
	tlsClientKey      string
	tlsClientKeyFile  string

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

func WithCleanSession(clean bool) Option {
	return func(c *Config) {
		c.cleanSession = clean
	}
}

func WithWaitTimeout(t time.Duration) Option {
	return func(c *Config) {
		c.waitTimeout = t
	}
}

func WithWaitTimeoutSec(sec int64) Option {
	return func(c *Config) {
		c.waitTimeout = time.Duration(sec) * time.Second
	}
}

func NewConfig(options ...Option) *Config {
	// 默认值
	cfg := &Config{
		cleanSession:  true,
		tlsVerify:     false, // default skipVerify
		reconnectType: ReConnTypeDefault,
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

func (c *Config) SetCleanSession(clean bool) *Config {
	c.cleanSession = clean
	return c
}

func (c *Config) SetWaitTimeout(t time.Duration) *Config {
	c.waitTimeout = t
	return c
}

func (c *Config) SetReconnectType(t ReConnType) *Config {
	c.reconnectType = t
	return c
}

//func (c *Config) SetWillEnabled(enabled bool) *Config {
//	c.willEnabled = enabled
//	return c
//}

// SetWill 设置遗嘱消息
func (c *Config) SetWill(topic, payload string, qos QosLevel, retained bool) *Config {
	c.willTopic = topic
	c.willPayload = []byte(payload)
	c.willPayloadIsByte = false
	c.willQos = qos
	c.willRetained = retained
	return c
}

// SetWillByte 设置遗嘱消息
func (c *Config) SetWillByte(topic string, payload []byte, qos QosLevel, retained bool) *Config {
	c.willTopic = topic
	c.willPayload = payload
	c.willPayloadIsByte = true
	c.willQos = qos
	c.willRetained = retained
	return c
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

func (c *Config) SetTlsCaCertFile(pem string) *Config {
	c.tlsCaCertFile = pem
	return c
}

func (c *Config) SetTlsClientCertFile(cert string, key string) *Config {
	c.tlsClientCertFile = cert
	c.tlsClientKeyFile = key
	return c
}

func (c *Config) SetTlsConfig(verify bool) *Config {
	c.tlsVerify = verify

	// TODO 其他 tls 配置
	//tls.NoClientCert
	//tls.RequestClientCert
	//tls.RequireAnyClientCert
	//tls.RequireAndVerifyClientCert

	return c

}
