/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     rose-mqtt
 * @Date:        2024-06-14 16:24
 * @Description:
 */

package rmqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

func (c *MQTTClient) sub(topic string, qos QosType, callback mqtt.MessageHandler) error {
	token := c.client.Subscribe(topic, byte(qos), callback)

	success := false
	if c.waitTimeout > 0 {
		success = token.WaitTimeout(c.waitTimeout)
	} else {
		success = token.Wait()
	}
	if !success {
		return fmt.Errorf("wait timeout")
	}

	if err := token.Error(); err != nil {
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

func (c *MQTTClient) Subscribe(topic string, qos QosType, callback mqtt.MessageHandler) error {
	//c.mu.Lock()
	//c.subHandlers[SubCallbackKey(topic, qos)] = callback
	//c.Topics = append(c.Topics, topic)
	//c.mu.Unlock()
	return c.sub(topic, qos, callback)
}

type Consumer struct {
	Topic    string
	QosType  QosType
	Callback mqtt.MessageHandler
}

func (c *MQTTClient) subConsumer(consumer *Consumer) {
	err := c.Subscribe(consumer.Topic, consumer.QosType, consumer.Callback)
	if err != nil {
		log.Printf("[Error] subscribe topic [%v]", consumer.Topic)
	}
	log.Printf("[Info] subscribe topic [%v]", consumer.Topic)
}

// RegisterConsumers 批量注册消费者
func (c *MQTTClient) RegisterConsumers(consumers []*Consumer) {
	for _, consumer := range consumers {
		go c.subConsumer(consumer)
	}
}

func (c *MQTTClient) RegisterConsumer(consumer *Consumer) {
	go c.subConsumer(consumer)
}
