/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     rose-mqtt
 * @Date:        2024-06-14 19:51
 * @Description:
 */

package rmqtt

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// generateRandomClientID 生成随机 ClientID
func generateRandomClientID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := r.Intn(clientIDMaxLength-clientIDMinLength+1) + clientIDMinLength

	clientID := make([]byte, length)
	for i := range clientID {
		clientID[i] = clientIDCharset[r.Intn(len(clientIDCharset))]
	}

	return clientIDPrefix + string(clientID)
}

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

// ConvMapToOrderedArray 将 map 转换为有序集合
func ConvMapToOrderedArray(m map[string]QosType) ([]string, []TopicQosPair) {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}

	// 对键排序
	sort.Strings(keys)

	orderedArray := make([]TopicQosPair, 0)
	for _, k := range keys {
		orderedArray = append(orderedArray, TopicQosPair{
			Topic: k,
			Qos:   m[k],
		})
	}

	return keys, orderedArray
}

// ConvOrderedArrayToMap 将有序集合转换为 map
func ConvOrderedArrayToMap(array []TopicQosPair) map[string]QosType {
	m := make(map[string]QosType, len(array))

	for _, kv := range array {
		if IsStrEmpty(kv.Topic) {
			continue
		}
		m[kv.Topic] = kv.Qos
	}
	return m
}

func JsonMarshal(v interface{}) string {
	bt, _ := json.Marshal(v)
	return string(bt)
}

func JsonUnMarshal(s string, v interface{}) {
	json.Unmarshal([]byte(s), v)
}
