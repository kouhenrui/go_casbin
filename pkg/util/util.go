package util

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/google/uuid"
)

// generateTraceID 生成追踪ID
func GenerateTraceID() string {
	return "trace-" + time.Now().Format("20060102150405") + "-" + RandomInt(6)
}

// randomInt 生成随机整数
func RandomInt(length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 退化为时间种子
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		} else {
			b[i] = charset[num.Int64()]
		}
	}
	return string(b)
}

// randomString 生成随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// randomUUID 生成随机UUID
func RandomUUID() string {
	return uuid.New().String()
}
