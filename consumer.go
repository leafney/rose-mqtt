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
)

type Consumer struct {
	Topic    string
	QosType  QosType
	Callback mqtt.MessageHandler
}

func (c *MQTTClient) Sub(topic string, qos QosType, callback mqtt.MessageHandler) error {
	token := c.client.Subscribe(topic, byte(qos), callback)
	if token.Wait() && token.Error() != nil {
		err := token.Error()
		if c.debug {
			log.Printf("[Error] ** Subscribe ** ClientID [%v] topic [%v] qos [%v] error [%v]", c.Ops.ClientID, topic, qos, err)
		}
		return err

	}
	if c.debug {
		log.Printf("[Info] ** Subscribe ** ClientID [%v] topic [%v] qos [%v] success", c.Ops.ClientID, topic, qos)
	}
	return nil
}

//
//func (c *MQTTClient) Subscribe(topic string, qos QosType, callback mqtt.MessageHandler) error {
//	//c.mu.Lock()
//	//c.msgHandlers[SubCallbackKey(topic, qos)] = callback
//	//c.Topics = append(c.Topics, topic)
//	//c.mu.Unlock()
//	return c.sub(topic, qos, callback)
//}
//
//func (c *MQTTClient) goSubConsumer(consumer *Consumer) {
//
//	err := c.Subscribe(consumer.Topic, consumer.QosType, consumer.Callback)
//	if err != nil {
//		log.Printf("[Error] ** goSubConsumer ** ClientID [%v] topic [%v] qos [%v] error [%v]", c.Ops.ClientID, consumer.Topic, consumer.QosType, err)
//		return
//	}
//	if c.debug {
//		log.Printf("[Info] ** goSubConsumer ** ClientID [%v] topic [%v] qos [%v] success", c.Ops.ClientID, consumer.Topic, consumer.QosType)
//	}
//}
//
//// RegisterConsumers 批量注册消费者
//func (c *MQTTClient) RegisterConsumers(consumers []*Consumer) {
//	for _, consumer := range consumers {
//		go c.goSubConsumer(consumer)
//	}
//}
