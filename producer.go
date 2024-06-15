/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     rose-mqtt
 * @Date:        2024-06-14 16:24
 * @Description:
 */

package rmqtt

import "log"

func (c *MQTTClient) Publish(topic string, qos QosType, payload interface{}) error {
	token := c.client.Publish(topic, byte(qos), false, payload)
	if token.Wait() && token.Error() != nil {
		err := token.Error()
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
