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

func (c *MQTTClient) sub(topic string, qos QosLevel, callback mqtt.MessageHandler) error {
	token := c.client.Subscribe(topic, byte(qos), callback)

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
			log.Printf("[Error] ** Subscribe ** ClientID [%v] topic [%v] qos [%v] error [%v]", c.Ops.ClientID, topic, qos, err)
		}
		return err
	}

	if c.debug {
		log.Printf("[Info] ** Subscribe ** ClientID [%v] topic [%v] qos [%v] success", c.Ops.ClientID, topic, qos)
	}
	return nil
}

func (c *MQTTClient) Subscribe(topic string, qos QosLevel, callback mqtt.MessageHandler) error {
	c.mu.Lock()
	c.subHandlers[SubCallbackKey(topic, qos)] = callback
	c.allTopics = append(c.allTopics, topic)
	c.mu.Unlock()

	return c.sub(topic, qos, callback)
}

func (c *MQTTClient) subMultiple(filters map[string]QosLevel, callback mqtt.MessageHandler) error {
	topics := map[string]byte{}
	for t, q := range filters {
		topics[t] = byte(q)
	}

	token := c.client.SubscribeMultiple(topics, callback)

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
			log.Printf("[Error] ** SubscribeMultiple ** ClientID [%v] topics [%v] error [%v]", c.Ops.ClientID, topics, err)
		}
		return err
	}

	if c.debug {
		log.Printf("[Info] ** SubscribeMultiple ** ClientID [%v] topics [%v] success", c.Ops.ClientID, topics)
	}
	return nil
}

func (c *MQTTClient) SubscribeMultiple(topics map[string]QosLevel, callback mqtt.MessageHandler) error {
	c.mu.Lock()
	keys, arr := ConvMapToOrderedArray(topics)
	arrStr := JsonMarshal(arr)
	c.subMutHandlers[arrStr] = callback
	c.allTopics = append(c.allTopics, keys...)
	c.mu.Unlock()

	return c.subMultiple(topics, callback)
}

type (
	Consumer struct {
		Topic    string
		QosType  QosLevel
		CallBack mqtt.MessageHandler
	}

	MultipleConsumer struct {
		Topics   map[string]QosLevel
		CallBack mqtt.MessageHandler
	}
)

func (c *MQTTClient) subConsumer(consumer *Consumer) {
	err := c.Subscribe(consumer.Topic, consumer.QosType, consumer.CallBack)
	if err != nil {
		if c.debug {
			log.Printf("[Error] ** subConsumer ** topic [%v] error [%v]", consumer.Topic, err)
		}
		return
	}
	if c.debug {
		log.Printf("[Info] ** subConsumer ** topic [%v] success", consumer.Topic)
	}
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

func (c *MQTTClient) subMultipleConsumer(mConsumer *MultipleConsumer) {
	err := c.SubscribeMultiple(mConsumer.Topics, mConsumer.CallBack)
	if err != nil {
		if c.debug {
			log.Printf("[Error] ** subMultipleConsumer ** topics [%v] error [%v]", mConsumer.Topics, err)
		}
		return
	}
	if c.debug {
		log.Printf("[Info] ** subMultipleConsumer ** topics [%v] success", mConsumer.Topics)
	}
}

func (c *MQTTClient) RegisterMultipleConsumer(mConsumer *MultipleConsumer) {
	go c.subMultipleConsumer(mConsumer)
}

func (c *MQTTClient) RegisterMultipleConsumers(mConsumers []*MultipleConsumer) {
	for _, mConsumer := range mConsumers {
		go c.subMultipleConsumer(mConsumer)
	}
}

func (c *MQTTClient) RegisterOnlyTopic(topic string, qos QosLevel) {
	c.RegisterConsumer(&Consumer{
		Topic:   topic,
		QosType: qos,
	})
}

func (c *MQTTClient) UnSubscribe(topics ...string) error {
	// 取消订阅
	token := c.client.Unsubscribe(topics...)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	//	移除 topics
	c.allTopics = SliceRmvSubSlice(c.allTopics, topics)

	return nil
}
