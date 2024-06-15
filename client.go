/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     rose-mqtt
 * @Date:        2024-06-14 16:23
 * @Description:
 */

package rmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sync"
	"time"
)

type MQTTClient struct {
	debug       bool
	ClientId    string
	client      mqtt.Client
	Ops         *mqtt.ClientOptions
	mu          *sync.Mutex
	Topics      []string                       // topic 集合
	msgHandlers map[string]mqtt.MessageHandler // key:topic#qos value:handler
}

/*
type Option func(c *MQTTClient)

func WithUserPwd(userName, passWord string) Option {
	return func(c *MQTTClient) {
		c.Ops.SetUsername(userName)
		c.Ops.SetPassword(passWord)
	}
}

func WithDebug(debug bool) Option {
	return func(c *MQTTClient) {
		c.debug = debug
	}
}

func WithClientID(id string) Option {
	return func(c *MQTTClient) {
		c.Ops.SetClientID(id)
	}
}

// NewMQTTClient brokerURI: tcp://foobar.com:1883
func NewMQTTClient(brokerURI string, options ...Option) *MQTTClient {

	ops := mqtt.NewClientOptions()
	ops.AddBroker(brokerURI)

	// 1. 测试直接在 NewClient 的时候设置ClientID，经测试后正常
	ops.SetClientID("9983300")
	ops.SetAutoReconnect(true)
	//ops.SetCleanSession(true)

	mc := &MQTTClient{
		Ops:         ops,
		mu:          &sync.Mutex{},
		msgHandlers: map[string]mqtt.MessageHandler{},
	}

	for _, opt := range options {
		opt(mc)
	}

	return mc
}

func (c *MQTTClient) SetClientID(clientId string) *MQTTClient {
	// 2. 测试通过这种方式来设置 ClientID ,经测试后，无法正常接收数据
	// 3. 但是可以通过 WithClientID 的方式来设置
	c.Ops.SetClientID(clientId)
	return c
}

*/

func NewMQTTClient(brokerURI string, cfg *Config) *MQTTClient {
	ops := mqtt.NewClientOptions()
	ops.AddBroker(brokerURI)

	//ops.SetUsername(cfg.username)
	//ops.SetPassword(cfg.password)
	ops.SetClientID(cfg.clientId)

	if cfg.keepAlive > 0 {
		ops.SetKeepAlive(cfg.keepAlive)
	}

	return &MQTTClient{
		debug:       cfg.debug,
		Ops:         ops,
		mu:          &sync.Mutex{},
		msgHandlers: map[string]mqtt.MessageHandler{},
	}
}

func (c *MQTTClient) Connect() (err error) {
	//if c.Ops.OnConnect == nil {
	//	//c.Ops.OnConnect=c.
	//}

	//fmt.Printf("ops [%v]", rose.JsonMarshalStr(&c.Ops))

	//
	c.client = mqtt.NewClient(c.Ops)
	// 连接，自动重试
	err = AutoRetry(func() error {
		if token := c.client.Connect(); token.Wait() && token.Error() != nil {
			err = token.Error()
			return err
		}
		return nil
	}, 3, 5*time.Second)
	if err != nil {
		return err
	}

	//if token := c.client.Connect(); token.Wait() && token.Error() != nil {
	//	err = token.Error()
	//	return err
	//}

	return nil
}

func (c *MQTTClient) Close() {
	if len(c.Topics) > 0 {
		c.client.Unsubscribe(c.Topics...)
		c.Topics = nil
	}
	c.msgHandlers = nil
	c.client.Disconnect(1000)
}
