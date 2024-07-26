package usecase

import (
	"context"
	"sync"
	"testing"

	"github.com/hrapovd1/msg-proc/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestAddMessageID(t *testing.T) {
	tests := []struct {
		name string
		data *types.Message
	}{
		{"Not nil", new(types.Message)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddMessageID(tt.data)
			assert.NotNil(t, tt.data.ID)
			assert.IsType(t, "", tt.data.ID)
		})
	}
}

func TestWriteJSONMessage(t *testing.T) {
	mockRepo := NewStorBus()
	type args struct {
		ctx  context.Context
		data types.Message
		repo types.Repository
	}
	tests := []struct {
		name string
		args args
		want types.Message
	}{
		{
			"simple message",
			args{context.Background(), types.Message{"simple", "message"}, mockRepo},
			types.Message{"simple", "message"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteJSONMessage(tt.args.ctx, tt.args.data, tt.args.repo)
			result := types.Message{mockRepo.mm[tt.args.data.ID].message, tt.args.data.ID}
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestSendJSONMessage(t *testing.T) {
	mockBus := NewStorBus()
	type args struct {
		ctx  context.Context
		data types.Message
		bus  types.BusMessenger
	}
	tests := []struct {
		name string
		args args
		want types.Message
	}{
		{
			"simple message",
			args{context.Background(), types.Message{"simple", "message"}, mockBus},
			types.Message{"simple", "message"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendJSONMessage(tt.args.ctx, tt.args.data, tt.args.bus)
			assert.Equal(t, tt.want, mockBus.mb[0])
		})
	}
}

// Тип реализующий интерфейсы Repository, BusMessenger
type storBus struct {
	mu sync.Mutex
	mm map[string]struct {
		message string
		status  string
	}
	mb []types.Message
}

func NewStorBus() *storBus {
	return &storBus{
		mm: make(map[string]struct {
			message string
			status  string
		}),
		mb: make([]types.Message, 0),
	}
}

func (sb *storBus) Save(c context.Context, msgID string, msg string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.mm[msgID] = struct {
		message string
		status  string
	}{msg, ""}
}

func (sb *storBus) Update(c context.Context, msgID string, status string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	msg := sb.mm[msgID]
	msg.status = status
	sb.mm[msgID] = msg
}

func (sb *storBus) Write(c context.Context, msg types.Message) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.mb = append(sb.mb, msg)
}

func (sb *storBus) Read(c context.Context) types.Message {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	if len(sb.mb) > 0 {
		return sb.mb[len(sb.mb)-1]
	}
	return types.Message{}
}
