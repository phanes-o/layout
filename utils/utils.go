package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"math"
)

func Throw(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	saltSize = 16
)

var SnowGen = &snowflake.Node{}

func ToJsonString(data interface{}) string {
	return string(toJson(data))
}

func ToJsonBytes(data interface{}) []byte {
	return toJson(data)
}

func toJson(data interface{}) []byte {
	if data == nil {
		return nil
	}
	marshal, _ := json.Marshal(data)
	return marshal
}

// PageUtil 分页工具
func PageUtil(size, index int) (int, int) {
	if index == 0 {
		index = 1
	}
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}
	return size, index
}

func PageCount(total, pageSize int64) int64 {
	return int64(math.Ceil(float64(total) / float64(pageSize)))
}

// GenerateSalt 生成密码盐值
func GenerateSalt() (string, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}
