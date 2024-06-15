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
