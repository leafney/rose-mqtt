/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     rose-mqtt
 * @Date:        2024-06-14 16:24
 * @Description:
 */

package rmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"testing"
)

const (
	brokerURL = "tcp://broker.emqx.io:1883"
)

func TestNewMClient(t *testing.T) {

	// 第一版
	//c := NewMQTTClient("tcp://broker.emqx.io:1883", WithDebug(true), WithClientID("93949"))
	//.SetClientID("009988")

	//c = c.SetClientID("0088885")
	//c.SetClientID("123456")
	//.
	//	SetKeepAliveSec(60)

	// 第二版
	cfg := NewConfig(
		WithDebug(true),
		//WithUserAndPwd("", ""),
	)
	//cfg.SetClientID("3499587")

	c := NewMQTTClient(brokerURL, cfg)

	if err := c.Connect(); err != nil {
		t.Fatal(err)
	}

	//cs := []*Consumer{
	//	{
	//		Topic:   "yili/1",
	//		QosType: Qos0,
	//		Callback: func(c mqtt.Client, msg mqtt.Message) {
	//			t.Logf("topic [%v] msg [%v]", msg.Topic(), string(msg.Payload()))
	//		},
	//	},
	//	{
	//		Topic:   "yili/2",
	//		QosType: Qos0,
	//		Callback: func(c mqtt.Client, msg mqtt.Message) {
	//			t.Logf("topic [%v] msg [%v]", msg.Topic(), string(msg.Payload()))
	//		},
	//	},
	//}

	//c.RegisterConsumers(cs)

	//go func() {
	//	if err := c.Sub("topic/hello", Qos1, func(_ mqtt.Client, msg mqtt.Message) {
	//		log.Printf("topic [%v] msg [%v]", msg.Topic(), string(msg.Payload()))
	//	}); err != nil {
	//		t.Error(err)
	//	}
	//	//log.Println("等待 subscribe")
	//}()

	//for i := 0; i < 10; i++ {
	//	topic := fmt.Sprintf("yili/%v", i%2+1)
	//	c.Publish(topic, Qos1, fmt.Sprintf("hello-%v", i))
	//	time.Sleep(1 * time.Second)
	//}

	select {}
	//	c.Close()
}

func TestSecond(t *testing.T) {
	cfg := NewConfig(
		WithDebug(true),
		//WithUserAndPwd("", ""),

	)
	//cfg.SetClientID("3499587")

	c := NewMQTTClient(brokerURL, cfg)

	if err := c.Connect(); err != nil {
		t.Fatal(err)
	}

	// 接收
	c.RegisterConsumer(&Consumer{
		Topic:   "topic/hello",
		QosType: 0,
		Callback: func(_ mqtt.Client, m mqtt.Message) {
			log.Printf("Subscribe333 callback topic [%v] msg [%v]", m.Topic(), string(m.Payload()))
		},
	})

	select {}
}
