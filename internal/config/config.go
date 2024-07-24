// Модуль config определяет типы и методы для формирования
// конфигурации приложения через флаги и переменные среды.
package config

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"reflect"

	"github.com/caarlos0/env/v6"
)

// environ содержит значения переменных среды
type environ struct {
	Address     string `env:"ADDRESS" envDefault:"localhost:8080"`
	DatabaseDSN string `env:"DATABASE_DSN" envDefault:""`
	ConfigFile  string `env:"CONFIG" envDefault:""`
}

// Config тип итоговой конфигурации приложения
type Config struct {
	ServerAddress string          `json:"address,omitempty"`
	DatabaseDSN   string          `json:"database_dsn,omitempty"`
	KafkaBrokers  string          `json:"kafka_brokers,omitempty"`
	KafkaTopic    string          `json:"kafka_topic,omitempty"`
	tagsDefault   map[string]bool `json:"-"`
}

// NewAppConf генерирует рабочую конфигурацию приложения
func NewAppConf(flags Flags) (*Config, error) {
	var err error
	var cfg Config
	var fileCfg Config
	cfg.tagsDefault = make(map[string]bool)
	var envs environ
	// Разбираю переменные среды и проверяю значение тегов на значение по умолчанию
	if err = env.Parse(&envs, env.Options{OnSet: cfg.getTags}); err != nil {
		return nil, err
	}
	// Определяю файл конфигурации и использую как
	// 3й источник конфигурации
	if !cfg.tagsDefault["CONFIG"] {
		if err := fileCfg.setConfigFromFile(envs.ConfigFile); err != nil {
			return nil, err
		}
	}
	if flags.configFile != "" {
		if err := fileCfg.setConfigFromFile(flags.configFile); err != nil {
			return nil, err
		}
	}
	// Определяю адрес сервера
	if flags.address != "" && cfg.tagsDefault["ADDRESS"] {
		cfg.ServerAddress = flags.address
	} else {
		cfg.ServerAddress = envs.Address
	}
	if flags.address == "" && cfg.tagsDefault["ADDRESS"] && fileCfg.valueExists("ServerAddress") {
		cfg.ServerAddress = fileCfg.ServerAddress
	}

	return &cfg, err
}

// getTags проверка и отметка значений переменных среды что они по умолчанию или нет
func (cfg *Config) getTags(tag string, value interface{}, isDefault bool) {
	cfg.tagsDefault[tag] = isDefault
}

// setConfigFromFile устанавливает параметры конфигурации из файла в формате JSON
func (cfg *Config) setConfigFromFile(cFile string) error {
	rawJSON, err := getRawJSONConfig(cFile)
	if err != nil {
		return err
	}
	if !json.Valid(rawJSON) {
		return errors.New("JSON from " + cFile + " NOT valid.")
	}
	conf := Config{}
	if err := json.Unmarshal(rawJSON, &conf); err != nil {
		return err
	}
	*cfg = conf
	return nil
}

func (cfg *Config) UnmarshalJSON(data []byte) error {
	type ConfigAlias Config

	aliasValue := &struct {
		*ConfigAlias
	}{
		ConfigAlias: (*ConfigAlias)(cfg),
	}
	if err := json.Unmarshal(data, aliasValue); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) valueExists(val string) bool {
	rCfg := reflect.ValueOf(*cfg)
	return !rCfg.FieldByName(val).IsZero()

}

func getRawJSONConfig(fName string) ([]byte, error) {
	fileStat, err := os.Stat(fName)
	if err != nil {
		return nil, err
	}
	if fileStat.Size() > 2000 {
		return nil, errors.New(fName + " too big.")
	}
	rawJSON := make([]byte, 2000)
	cf, err := os.Open(fName)
	if err != nil {
		return nil, err
	}
	n, err := cf.Read(rawJSON)
	if err != nil {
		return nil, err
	}
	return rawJSON[:n], nil
}

// Flags содержит значения флагов переданные при запуске
type Flags struct {
	address    string
	dbDSN      string
	configFile string
}

// GetServerFlags - считывае флаги сервера
func GetAppFlags() Flags {
	flags := Flags{}
	flag.StringVar(&flags.address, "a", "", "Address of server, for example: 0.0.0.0:8000")
	flag.StringVar(&flags.dbDSN, "d", "", "Database connect source, for example: postgres://username:password@localhost:5432/database_name")
	flag.StringVar(&flags.configFile, "c", "", "(or -config) Path to config file in JSON format")
	flag.StringVar(&flags.configFile, "config", "", "(or -c) Path to config file in JSON format")
	flag.Parse()
	return flags
}
