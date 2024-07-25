// Модуль usecase содержит общие для проекта методы.
package usecase

import (
	"context"

	"github.com/hrapovd1/msg-proc/internal/types"
)

// AddMessageID добавляет внутренний ID к полученному сообщению
func AddMessageID(data *types.Message) {
	messageID := GenMsgID(data.Msg)
	data.ID = messageID
}

// WriteJSONMessage сохраняет сообщение в Repository полученное в
// JSON формате POST запроса.
func WriteJSONMessage(ctx context.Context, data types.Message, repo types.Repository) {
	repo.Save(ctx, data.ID, data.Msg)
}

// SendJSONMessage отправляет сообщение в шину сообщений в JSON формате из POST запроса
func SendJSONMessage(ctx context.Context, data types.Message, bus types.BusMessenger) {
	bus.Write(ctx, data)
}
