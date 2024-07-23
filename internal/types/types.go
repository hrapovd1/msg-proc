// Модуль types содержит общие для проетка типы и интерфейсы.
package types

import (
	"context"
)

// Префикс в названиях таблиц базы
const DBtableName = "app_messages"

// Metric тип JSON формата метрики
type Message struct {
	Msg    string `json:"msg"` // сообщение
	Status string
}

// Repository основной интерфейс хранилища сообщений
type Repository interface {
	Save(ctx context.Context, msgID string, message string)
	Update(ctx context.Context, msgID string, status string)
}

// Storager вспомогательный интерфейс хранилища сообщений
type Storager interface {
	Close() error
	Ping(ctx context.Context) bool
}

// MessageModel модель таблицы для хранения сообщения в базе
type MessageModel struct {
	Timestamp int64 `gorm:"primaryKey;autoCreateTime"`
	ID        string
	Message   string
	Status    string
}
