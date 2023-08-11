package logger

import (
	"io"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

func fileOutputWriter(filePath string, fileName string, maxSize, maxAge int) io.Writer {
	var name string
	if !exist(filePath) {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return nil
		}
	}
	name = filePath + "/" + fileName
	if !strings.HasSuffix(fileName, ".log") {
		name = filePath + "/" + fileName + ".log"
	}
	writer := &lumberjack.Logger{
		Filename:   name,
		MaxSize:    maxSize, // megabytes default 50
		MaxBackups: 3,
		MaxAge:     maxAge, // days default 3
		Compress:   true,
	}
	return writer
}

func exist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			return true
		}
	}
	return false
}
