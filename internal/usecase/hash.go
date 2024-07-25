package usecase

import (
	"time"
)

const timestampFormat = "20060102T150405.000"

// GenMsgID генерирует hash сообщения как его ID
func GenMsgID(msg string) string {
	return time.Now().Format(timestampFormat)
}
