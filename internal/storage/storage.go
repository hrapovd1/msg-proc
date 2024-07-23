package storage

import (
	"context"

	"github.com/hrapovd1/msg-proc/internal/types"
)

// Option тип для модификации хранилища MemStorage
type Option func(mem *MemStorage) *MemStorage

// MemStorage тип реализации хранения в памяти
type MemStorage struct {
	buffer map[string]any
}

// NewMemStorage создает хранилище MemStorage
func NewMemStorage(opts ...Option) *MemStorage {
	buffer := make(map[string]any)
	ms := &MemStorage{
		buffer: buffer,
	}

	for _, opt := range opts {
		ms = opt(ms)
	}

	return ms
}

// Save сохраняет новое значение сообщения
func (ms *MemStorage) Save(ctx context.Context, msgID, message string) {
	select {
	case <-ctx.Done():
		return
	default:
		ms.buffer[msgID] = types.Message{
			Msg:    message,
			Status: "",
		}
	}
}

// Update устанавливает статус сообщения после обработки
func (ms *MemStorage) Update(ctx context.Context, msgID, status string) {
	select {
	case <-ctx.Done():
		return
	default:
		if _, ok := ms.buffer[msgID]; ok {
			ms.buffer[msgID] = types.Message{
				Status: status,
			}
		}
	}
}

// Close для реализации интерфейса Storager
func (ms *MemStorage) Close() error { return nil }

// Ping для реализации интерфейса Storager
func (ms *MemStorage) Ping(ctx context.Context) bool { return false }

// WitBuffer модифицирует MemStorage позволяя передать внешний буфер
// как внутреннее хранилище, используется в тестах
func WithBuffer(buffer map[string]interface{}) Option {
	return func(mem *MemStorage) *MemStorage {
		mem.buffer = buffer
		return mem
	}
}
