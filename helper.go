/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     rose-mqtt
 * @Date:        2024-06-15 21:02
 * @Description:
 */

package rmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"slices"
	"strings"
	"time"
)

// ReconnectManualHandler 手动实现自动重连机制
func ReconnectManualHandler(client mqtt.Client, err error) {
	log.Printf("MQTT connection lost: [%v]", err)

	maxReties := 10
	minReties := 3
	minDelay := 5
	maxDelay := 5

	for i := 0; i < maxReties; i++ {
		log.Printf("Attempt [%v] to reconnect...", i+1)

		//	每次重连尝试 3 次，间隔时间递增
		for j := 0; j < minReties; j++ {
			interval := time.Duration(minDelay+j*minDelay) * time.Second
			log.Printf("Reconnecting in %v ...", interval)
			time.Sleep(interval)

			token := client.Connect()
			token.Wait()
			if err := token.Error(); err == nil {
				log.Println("Reconnected to MQTT broker success")
				return
			}
		}

		//	如 3 次重连失败，等待 10 分钟后再次尝试
		log.Printf("Reconnection failed, waiting %v minutes before next attempt...", maxDelay)
		time.Sleep(time.Duration(maxDelay) * time.Minute)
	}

	log.Println("Maximum reconnection attempts reached, giving up.")
}

// DefaultOnConnect 重连后自动注册订阅
func (c *MQTTClient) DefaultOnConnect(cli mqtt.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//
	for key, handler := range c.subHandlers {
		go func(topic string, cb mqtt.MessageHandler) {
			split := strings.Split(topic, "#")

			if len(split) == 2 {
				var qos QosLevel
				switch split[1] {
				case "0":
					qos = Qos0
				case "1":
					qos = Qos1
				case "2":
					qos = Qos2
				default:
					qos = Qos0
				}

				// 判断当前 topic 是否已被取消订阅
				if !slices.Contains(c.allTopics, split[0]) {
					return
				}

				if err := c.sub(split[0], qos, cb); err != nil {
					log.Printf("[Error] topic [%v] reconnect register error [%v]", split[0], err)
					return
				}
				if c.debug {
					log.Printf("[Info] topic [%v] reconnect register success", split[0])
				}
			}
		}(key, handler)
	}

	//
	for key, handler := range c.subMutHandlers {
		go func(topic string, cb mqtt.MessageHandler) {
			//
			var result []TopicQosPair
			if err := JsonUnMarshal(topic, &result); err != nil {
				if c.debug {
					log.Printf("[Error] topic [%v] unmarshal error [%v]", topic, err)
				}
				return
			}

			// 判断当前的 topic 是否已被取消订阅
			newResult := make([]TopicQosPair, 0)
			for _, tq := range result {
				if slices.Contains(c.allTopics, tq.Topic) {
					newResult = append(newResult, tq)
				}
			}

			resMap := ConvOrderedArrayToMap(newResult)

			if err := c.subMultiple(resMap, cb); err != nil {
				log.Printf("[Error] topic [%v] reconnect register error [%v]", resMap, err)
			}
			if c.debug {
				log.Printf("[Info] topic [%v] reconnect register success", resMap)
			}
		}(key, handler)
	}
}
