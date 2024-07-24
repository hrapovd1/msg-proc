// Модуль msgbus реализует типы и методы для работы с
// шиной сообщений - Kafka
package msgbus

import (
	"context"
	"log"
	"strings"

	"github.com/hrapovd1/msg-proc/internal/config"
	"github.com/hrapovd1/msg-proc/internal/types"
	"github.com/segmentio/kafka-go"
)

const (
	WriterAsync           = true //Не блокирующая запись
	WriterAutoTopicCreate = true //Авто создание топика
)

// KafkaBus реализует шину сообщений на Kafka
type KafkaBus struct {
	Reader *kafka.Reader
	Writer *kafka.Writer
}

// NewKfkBus правильно инициализирует kafka reader/writer
func NewKfkBus(appConf config.Config, logger *log.Logger) *KafkaBus {
	kfkBus := KafkaBus{
		Reader: newKafkaReader(appConf, logger),
		Writer: newKafkaWriter(appConf, logger),
	}
	return &kfkBus
}

// newKafkaReader правильно инициализирует kafka Reader
func newKafkaReader(appConf config.Config, logger *log.Logger) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     strings.Split(appConf.KafkaBrokers, ","),
		Topic:       appConf.KafkaTopic,
		Logger:      logger,
		ErrorLogger: nil, // will use logger
	})
}

// newKafkaWriter правильно инициализирует kafka Writer
func newKafkaWriter(appConf config.Config, logger *log.Logger) *kafka.Writer {
	brokers := strings.Split(appConf.KafkaBrokers, ",")
	return &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  appConf.KafkaTopic,
		Async:                  WriterAsync,
		AllowAutoTopicCreation: WriterAutoTopicCreate,
		Logger:                 logger,
		ErrorLogger:            nil, // will use logger
	}
}

// Write реализует запись сообщения в kafka topic
func (kb *KafkaBus) Write(ctx context.Context, msg types.Message) {
	if err := kb.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.ID),
		Value: []byte(msg.Msg),
	}); err != nil {
		kb.Writer.Logger.Printf("Error when write in kafka: %v\n", err)
	}
}

// Read реализует чтение сообщения из топика по текущему offset
func (kb *KafkaBus) Read(ctx context.Context) types.Message {
	select {
	case <-ctx.Done():
		return types.Message{}
	default:
		msg, err := kb.Reader.ReadMessage(ctx)
		if err != nil {
			kb.Reader.Config().Logger.Printf("Error when read from kafka: %v\n", err)
			return types.Message{}
		}
		return types.Message{
			Msg: string(msg.Key),
			ID:  string(msg.Value),
		}
	}
}

// Consume реализует интерфейс MsgConsumer для непрерывного чтения
// сообщений из шины
func (kb *KafkaBus) Consume(ctx context.Context, stor types.Repository) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			msg := kb.Read(ctx)
			stor.Update(ctx, msg.ID, "processed")
		}
	}
}
