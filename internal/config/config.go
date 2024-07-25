// Модуль config определяет типы и методы для формирования
// конфигурации приложения через флаги и переменные среды.
package config

import (
	"github.com/caarlos0/env/v6"
)

// environ содержит значения переменных среды
type environ struct {
	Address      string `env:"ADDRESS" envDefault:"localhost:8080"`
	DatabaseDSN  string `env:"DATABASE_DSN" envDefault:"postgres://postgres:postgres@localhost:5432/postgres"`
	KafkaBrokers string `env:"KAFKA_BROKERS" envDefault:"localhost:29092"`
	KafkaTopic   string `env:"KAFKA_TOPIC" envDefault:"messages"`
}

// Config тип итоговой конфигурации приложения
type Config struct {
	ServerAddress string `json:"address,omitempty"`
	DatabaseDSN   string `json:"database_dsn,omitempty"`
	KafkaBrokers  string `json:"kafka_brokers,omitempty"`
	KafkaTopic    string `json:"kafka_topic,omitempty"`
}

// NewAppConf генерирует рабочую конфигурацию приложения
func NewAppConf() (*Config, error) {
	var err error
	var envs environ
	// Разбираю переменные среды и проверяю значение тегов на значение по умолчанию
	if err = env.Parse(&envs); err != nil {
		return nil, err
	}
	cfg := Config{
		ServerAddress: envs.Address,
		DatabaseDSN:   envs.DatabaseDSN,
		KafkaBrokers:  envs.KafkaBrokers,
		KafkaTopic:    envs.KafkaTopic,
	}
	return &cfg, err
}
