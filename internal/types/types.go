// Модуль types содержит общие для проетка типы и интерфейсы.
package types

import (
	"context"
	"time"
)

const (
	DBtableName   = "app_messages" // Имя таблицы в базе
	InputMetric   = "totalInput"
	ProcessMetric = "processed"
	SyncPeriod    = time.Second * 30
)

// Message тип JSON формата сообщения
type Message struct {
	Msg string `json:"msg"` // сообщение
	ID  string `json:"-"`
}

// Metrics тип JSON формата для handler
type Metrics struct {
	Total     int64 `json:"total"`
	Processed int64 `json:"processed"`
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
	GetCount(ctx context.Context, processed bool) int64
}

// BusMessenger интерфейс шины сообщений
type BusMessenger interface {
	Write(ctx context.Context, msg Message)
	Read(ctx context.Context) Message
}

// MsgConsumer интерфейс потребителя сообщений из шины
type MsgConsumer interface {
	Consume(ctx context.Context, stor Repository)
}

// MessageModel модель таблицы для хранения сообщения в базе
type MessageModel struct {
	Timestamp int64 `gorm:"primaryKey;autoCreateTime"`
	ID        string
	Message   string
	Status    string
}
