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

const (
	clientIDCharset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	clientIDMinLength = 8        // 随机字符串最小长度
	clientIDMaxLength = 14       // 随机字符串最大长度
	clientIDPrefix    = "rmqtt-" // 前缀
)

// 重试方式
type ReConnType string

const (
	ReConnTypeDefault   ReConnType = ""       // 默认
	ReConnTypeAutomatic ReConnType = "auto"   // 自动
	ReConnTypeManual    ReConnType = "manual" // 手动
)

type TopicQosPair struct {
	Topic string
	Qos   QosType
}
