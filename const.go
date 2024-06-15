/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     rose-mqtt
 * @Date:        2024-06-14 19:32
 * @Description:
 */

package rmqtt

type QosType byte

const (
	Qos0 QosType = 0
	Qos1 QosType = 1
	Qos2 QosType = 2
)
