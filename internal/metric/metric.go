package metric

import (
	"context"
	"time"

	"github.com/hrapovd1/msg-proc/internal/types"
)

// Metrics тип для хранения метрик в памяти
type Metrics struct {
	Mem map[string]int64
	Ds  types.Storager
}

// NewMetrics инициализирует хранилище метрик
func NewMetrics(db types.Storager) *Metrics {
	return &Metrics{
		Mem: make(map[string]int64, 2),
		Ds:  db,
	}
}

// IncrementInput увеличивает счетчик входящих сообщений
func (m *Metrics) IncrementInput() {
	m.Mem[types.InputMetric] = m.Mem[types.InputMetric] + 1
}

// IncrementProcessed увеличивает счетчик обработанных сообщений
func (m *Metrics) IncrementProcessed() {
	m.Mem[types.ProcessMetric] = m.Mem[types.ProcessMetric] + 1
}

// SyncWithDB синхронизирует счетчики в памяти с данными из базы
func (m *Metrics) SyncWithDB(ctx context.Context) {
	ticker := time.NewTicker(types.SyncPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			m.Mem[types.InputMetric] = m.Ds.GetCount(ctx, false)
			m.Mem[types.ProcessMetric] = m.Ds.GetCount(ctx, true)
		}
	}
}
