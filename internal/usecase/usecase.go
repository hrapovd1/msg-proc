// Модуль usecase содержит общие для проекта методы.
package usecase

import (
	"context"

	"github.com/hrapovd1/msg-proc/internal/types"
)

// WriteJSONMessage сохраняет сообщение в Repository полученное в
// JSON формате POST запроса.
func WriteJSONMessage(ctx context.Context, data types.Message, repo types.Repository) error {
	messageID, err := GenMsgID(data.Msg)
	if err != nil {
		return err
	}
	repo.Save(ctx, messageID, data.Msg)
	return nil
}
