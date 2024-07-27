package metric

import (
	"context"
	"sync"
	"testing"

	"github.com/hrapovd1/msg-proc/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMetrics(t *testing.T) {
	mStorager := storgr{int64(1), int64(2)}
	type args struct {
		db types.Storager
	}
	tests := []struct {
		name string
		args args
		want *Metrics
	}{
		{
			name: "simple check",
			args: args{db: mStorager},
			want: &Metrics{
				Mem: map[string]int64{},
				Ds:  mStorager,
				Mu:  &sync.Mutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMetrics(mStorager)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMetrics_IncrementInput(t *testing.T) {
	mStorager := storgr{}
	mtrcs := NewMetrics(mStorager)
	tests := []struct {
		name string
		want int64
	}{
		{
			name: "simple test",
			want: int64(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := int64(0); i < tt.want; i++ {
				mtrcs.IncrementInput()
			}
			assert.Equal(t, tt.want, mtrcs.Mem[types.InputMetric])
		})
	}
}

func TestMetrics_IncrementProcessed(t *testing.T) {
	mStorager := storgr{}
	mtrcs := NewMetrics(mStorager)
	tests := []struct {
		name string
		want int64
	}{
		{
			name: "simple test",
			want: int64(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := int64(0); i < tt.want; i++ {
				mtrcs.IncrementProcessed()
			}
			assert.Equal(t, tt.want, mtrcs.Mem[types.ProcessMetric])
		})
	}
}

func TestMetrics_sync(t *testing.T) {
	mStorager := storgr{int64(2), int64(5)}
	mtrcs := NewMetrics(mStorager)
	tests := []struct {
		name          string
		wantInput     int64
		wantProcessed int64
	}{
		{
			name:          "simple check",
			wantInput:     int64(2),
			wantProcessed: int64(5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mtrcs.sync(context.Background())
			assert.Equal(t, tt.wantInput, mtrcs.Mem[types.InputMetric])
			assert.Equal(t, tt.wantProcessed, mtrcs.Mem[types.ProcessMetric])
		})
	}
}

// Mock тип интерфейса Storager
type storgr []int64

func (s storgr) Close() error {
	return nil
}

func (s storgr) Ping(ctx context.Context) bool {
	return true
}

func (s storgr) GetCount(ctx context.Context, processed bool) int64 {
	if processed {
		return s[1]
	}
	return s[0]
}
