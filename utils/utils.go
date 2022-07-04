package utils

import "encoding/json"

func Throw(err error) {
	if err != nil {
		panic(err)
	}
}

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
