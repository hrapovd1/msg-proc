// Модуль config определяет типы и методы для формирования
// конфигурации приложения через флаги и переменные среды.
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAppConf(t *testing.T) {
	envVars := map[string]string{
		"ADDRESS":       "test:8888",
		"DATABASE_DSN":  "postgres:5555",
		"KAFKA_BROKERS": "kafka1,kafka2",
		"KAFKA_TOPIC":   "topic1",
	}
	for k, v := range envVars {
		os.Setenv(k, v)
	}
	t.Run("check config", func(t *testing.T) {
		got, err := NewAppConf()
		require.NoError(t, err)
		assert.Equal(t, envVars["ADDRESS"], got.ServerAddress)
		assert.Equal(t, envVars["DATABASE_DSN"], got.DatabaseDSN)
		assert.Equal(t, envVars["KAFKA_BROKERS"], got.KafkaBrokers)
		assert.Equal(t, envVars["KAFKA_TOPIC"], got.KafkaTopic)
	})
}
