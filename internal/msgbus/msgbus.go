// Модуль msgbus реализует типы и методы для работы с
// шиной сообщений - Kafka
package msgbus

import (
	"log"

	"github.com/hrapovd1/msg-proc/internal/config"
	"github.com/segmentio/kafka-go"
)

const (
	WriterAsync           = true
	WriterAutoTopicCreate = true
)

type KafkaBus struct {
	Reader *kafka.Reader
	Writer *kafka.Writer
}

func NewKfkBus(appConf *config.Config, logger *log.Logger) *KafkaBus {
	kfkBus := KafkaBus{
		Reader: NewKafkaReader(appConf, logger),
		Writer: NewKafkaWriter(appConf, logger),
	}
	return &kfkBus
}

func NewKafkaReader(appConf *config.Config, logger *log.Logger) *kafka.Reader {}

func NewKafkaWriter(appConf *config.Config, logger *log.Logger) *kafka.Writer {}
