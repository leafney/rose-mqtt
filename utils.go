/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     rose-mqtt
 * @Date:        2024-06-14 19:51
 * @Description:
 */

package rmqtt

import (
	"fmt"
	"strings"
	"time"
)

// AutoRetry 失败后重试指定次数
func AutoRetry(callback func() error, maxRetries int, interval time.Duration) (err error) {
	for i := 0; i < maxRetries; i++ {
		if err = callback(); err != nil {
			time.Sleep(interval)
			continue
		}
		return
	}
	return
}

func SubCallbackKey(topic string, qos QosType) string {
	return fmt.Sprintf("%s#%v", topic, qos)
}

func IsStrEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
