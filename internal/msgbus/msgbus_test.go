// Модуль msgbus реализует типы и методы для работы с
// шиной сообщений - Kafka
package msgbus

import (
	"log"
	"testing"
	"unsafe"

	"github.com/hrapovd1/msg-proc/internal/config"
	"github.com/hrapovd1/msg-proc/internal/metric"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestNewKfkBus(t *testing.T) {
	type args struct {
		appConf config.Config
		logger  *log.Logger
		metrics *metric.Metrics
	}
	tests := []struct {
		name string
		args args
		want *KafkaBus
	}{
		{
			"size check",
			args{
				config.Config{KafkaBrokers: "", KafkaTopic: "topic"},
				log.Default(),
				&metric.Metrics{},
			},
			&KafkaBus{
				&kafka.Reader{},
				&kafka.Writer{},
				metric.Metrics{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewKfkBus(tt.args.appConf, tt.args.logger, tt.args.metrics)
			_ = got
			assert.Equal(t, unsafe.Sizeof(tt.want), unsafe.Sizeof(got))
		})
	}
}

func Test_newKafkaReader(t *testing.T) {
	type args struct {
		appConf config.Config
		logger  *log.Logger
	}
	tests := []struct {
		name string
		args args
		want *kafka.Reader
	}{
		{
			"size check",
			args{
				config.Config{KafkaBrokers: "", KafkaTopic: "topic"},
				log.Default(),
			},
			&kafka.Reader{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newKafkaReader(tt.args.appConf, tt.args.logger)
			_ = got
			assert.Equal(t, unsafe.Sizeof(tt.want), unsafe.Sizeof(got))
		})
	}
}

func Test_newKafkaWriter(t *testing.T) {
	type args struct {
		appConf config.Config
		logger  *log.Logger
	}
	tests := []struct {
		name string
		args args
		want *kafka.Writer
	}{
		{
			"size check",
			args{
				config.Config{KafkaBrokers: "", KafkaTopic: "topic"},
				log.Default(),
			},
			&kafka.Writer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newKafkaWriter(tt.args.appConf, tt.args.logger)
			_ = got
			assert.Equal(t, unsafe.Sizeof(tt.want), unsafe.Sizeof(got))
		})
	}
}
