package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

const hashKey = "fjsdfll87sf.sdfwsLJDL:FJ"

// GenMsgID генерирует hash сообщения как его ID
func GenMsgID(msg string) (string, error) {
	h := hmac.New(sha256.New, []byte(hashKey))
	_, err := h.Write([]byte(msg))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
