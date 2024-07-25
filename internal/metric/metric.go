package metric

import "github.com/hrapovd1/msg-proc/internal/types"

// Metrics тип для хранения метрик в памяти
type Metrics struct {
	Mem map[string]int64
	Ds  types.Repository
}

// NewMetrics инициализирует хранилище метрик
func NewMetrics(db *types.Repository) *Metrics {
	return &Metrics{
		Mem: make(map[string]int64, 2),
		Ds:  *db,
	}
}

// IncrementInput увеличивает счетчик входящих сообщений
func (m *Metrics) IncrementInput() {}

// IncrementProcessed увеличивает счетчик обработанных сообщений
func (m *Metrics) IncrementProcessed() {}

// SyncWithDB синхронизирует счетчики в памяти с данными из базы
func (m *Metrics) SyncWithDB() {}
