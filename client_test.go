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

func TestNewMClient(t *testing.T) {
	URL := "tcp://broker.emqx.io:1883"
	//c := NewMQTTClient("tcp://broker.emqx.io:1883", WithDebug(true), WithClientID("93949"))
	//.SetClientID("009988")

	//c = c.SetClientID("0088885")
	//c.SetClientID("123456")
	//.
	//	SetKeepAliveSec(60)

	cfg := NewConfig(
		WithDebug(true),
		//WithUserAndPwd("", ""),
	)
	//cfg.SetClientID("3499587")

	c := NewMQTTClient(URL, cfg)
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

	go func() {
		if err := c.Sub("topic/hello", Qos1, func(_ mqtt.Client, msg mqtt.Message) {
			log.Printf("topic [%v] msg [%v]", msg.Topic(), string(msg.Payload()))
		}); err != nil {
			t.Error(err)
		}
		//log.Println("等待 subscribe")
	}()

	//for i := 0; i < 10; i++ {
	//	topic := fmt.Sprintf("yili/%v", i%2+1)
	//	c.Publish(topic, Qos1, fmt.Sprintf("hello-%v", i))
	//	time.Sleep(1 * time.Second)
	//}

	select {}
	//	c.Close()
}
