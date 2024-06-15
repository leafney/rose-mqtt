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
	"log"
)

func (c *MQTTClient) Publish(topic string, qos QosType, retained bool, payload interface{}) error {
	token := c.client.Publish(topic, byte(qos), retained, payload)

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
			log.Printf("[Error] ** Publish ** ClientID [%v] topic [%v] qos [%v] error [%v]", c.Ops.ClientID, topic, qos, err)
		}
		return err
	}
	if c.debug {
		log.Printf("[Info] ** Publish ** ClientID [%v] topic [%v] qos [%v] success", c.Ops.ClientID, topic, qos)
	}
	return nil
}
