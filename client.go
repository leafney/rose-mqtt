/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     rose-mqtt
 * @Date:        2024-06-14 16:23
 * @Description:
 */

package rmqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"sync"
	"time"
)

type MQTTClient struct {
	debug          bool
	waitTimeout    time.Duration
	client         mqtt.Client
	Ops            *mqtt.ClientOptions
	mu             *sync.Mutex
	tls            tlsConfig
	allTopics      []string                       // topic 集合
	subHandlers    map[string]mqtt.MessageHandler // 单独订阅时 key:topic#qos value:handler
	subMutHandlers map[string]mqtt.MessageHandler // 多个 topic 同时订阅时 key: json([{topic1:qos1},{topic2:qos2}]) value: handle
}

type tlsConfig struct {
	enabled    bool
	ca         string
	verify     bool
	clientKey  string
	clientCert string
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
		subHandlers: map[string]mqtt.MessageHandler{},
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
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURI)

	if !IsStrEmpty(cfg.username) {
		opts.SetUsername(cfg.username)
	}
	if !IsStrEmpty(cfg.password) {
		opts.SetPassword(cfg.password)
	}
	if !IsStrEmpty(cfg.clientId) {
		opts.SetClientID(cfg.clientId)
	}
	if cfg.keepAlive > 0 {
		opts.SetKeepAlive(cfg.keepAlive)
	}

	opts.SetCleanSession(cfg.cleanSession)

	// 自动重连配置
	switch cfg.reconnectType {
	case ReConnTypeAutomatic:
		//	自动重连
		opts.SetAutoReconnect(true)
		opts.SetConnectRetry(true)
		opts.SetConnectRetryInterval(5 * time.Second)
		opts.SetMaxReconnectInterval(60 * time.Second)
	case ReConnTypeManual:
		//	手动重连
		opts.SetAutoReconnect(false)
		opts.SetConnectionLostHandler(ReconnectManualHandler)
	default:
		//	默认
	}

	if cfg.defaultHandler != nil {
		opts.SetDefaultPublishHandler(cfg.defaultHandler) // 当接收数据没有匹配的处理函数时触发
	}

	if cfg.connHandler != nil {
		opts.SetOnConnectHandler(cfg.connHandler) // 连接回调
	}
	if cfg.connLostHandler != nil {
		opts.SetConnectionLostHandler(cfg.connLostHandler) // 连接意外中断回调
	}

	// 遗嘱消息
	if !IsStrEmpty(opts.WillTopic) || len(opts.WillPayload) > 0 {
		if cfg.willPayloadIsByte {
			opts.SetBinaryWill(cfg.willTopic, cfg.willPayload, byte(cfg.willQos), cfg.willRetained)
		} else {
			opts.SetWill(cfg.willTopic, string(cfg.willPayload), byte(cfg.willQos), cfg.willRetained)
		}
	}

	// tls
	tlsCfg := tlsConfig{
		enabled:    false,
		ca:         cfg.tlsCaCertFile,
		verify:     cfg.tlsVerify,
		clientCert: cfg.tlsClientCertFile,
		clientKey:  cfg.tlsClientKeyFile,
	}
	if !IsStrEmpty(cfg.tlsCaCertFile) {
		tlsCfg.enabled = true
	}

	return &MQTTClient{
		debug:          cfg.debug,
		waitTimeout:    cfg.waitTimeout,
		Ops:            opts,
		tls:            tlsCfg,
		mu:             &sync.Mutex{},
		subHandlers:    map[string]mqtt.MessageHandler{},
		subMutHandlers: map[string]mqtt.MessageHandler{},
	}
}

func (c *MQTTClient) Connect() (err error) {
	if c.Ops.OnConnect == nil {
		c.Ops.OnConnect = c.DefaultOnConnect
	}

	// tls
	if c.tls.enabled {
		tCfg := &tls.Config{
			RootCAs:            nil,
			ClientAuth:         tls.NoClientCert,
			InsecureSkipVerify: !c.tls.verify, // 注意这里的值
			Certificates:       nil,
		}

		certPool := x509.NewCertPool()
		pemCerts, err := os.ReadFile(c.tls.ca)
		if err != nil {
			if c.debug {
				log.Printf("[Error] tls config load CaCert file error [%v]", err)
			}
			return err
		}
		certPool.AppendCertsFromPEM(pemCerts)
		tCfg.RootCAs = certPool

		if !IsStrEmpty(c.tls.clientCert) && !IsStrEmpty(c.tls.clientKey) {
			cert, err := tls.LoadX509KeyPair(c.tls.clientCert, c.tls.clientKey)
			if err != nil {
				if c.debug {
					log.Printf("[Error] tls config load ClientCert and ClientKey error [%v]", err)
				}
				return err
			}
			tCfg.Certificates = []tls.Certificate{cert}
		}

		c.Ops.SetTLSConfig(tCfg)
	}

	c.client = mqtt.NewClient(c.Ops)

	// 连接，默认
	//if token := c.client.Connect(); token.Wait() && token.Error() != nil {
	//	err = token.Error()
	//	return err
	//}

	// 连接，自定义超时时间
	token := c.client.Connect()
	waitRes := false
	if c.waitTimeout > 0 {
		waitRes = token.WaitTimeout(c.waitTimeout)
	} else {
		waitRes = token.Wait()
	}
	if !waitRes {
		return fmt.Errorf("wait timeout for %v", c.waitTimeout)
	}
	if err := token.Error(); err != nil {
		if c.debug {
			log.Printf("[Error] ** Connect ** ClientID [%v] error [%v]", c.Ops.ClientID, err)
		}
		return err
	}
	if c.debug {
		log.Printf("[Info] ** Connect ** ClientID [%v] success", c.Ops.ClientID)
	}

	// 连接，自动重试
	//err = AutoRetry(func() error {
	//	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
	//		err = token.Error()
	//		return err
	//	}
	//	return nil
	//}, 3, 5*time.Second)

	//// 连接，自动重试 优化
	//err = AutoRetry(c.tryConnect, 3, 5*time.Second)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (c *MQTTClient) tryConnect() error {
	token := c.client.Connect()
	waitRes := false
	if c.waitTimeout > 0 {
		waitRes = token.WaitTimeout(c.waitTimeout)
	} else {
		waitRes = token.Wait()
	}
	if !waitRes {
		return fmt.Errorf("wait timeout for %v", c.waitTimeout)
	}

	if err := token.Error(); err != nil {
		if c.debug {
			log.Printf("[Error] tryConnect error [%v]", err)
		}
		return err
	}
	return nil
}

func (c *MQTTClient) Close() {
	if len(c.allTopics) > 0 {
		c.client.Unsubscribe(c.allTopics...)
		c.allTopics = nil
	}
	c.subHandlers = nil
	c.client.Disconnect(1000)
}
