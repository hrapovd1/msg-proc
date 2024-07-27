// Модуль config определяет типы и методы для формирования
// конфигурации приложения через флаги и переменные среды.
package config

import (
	"reflect"
	"testing"
)

func TestNewAppConf(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAppConf()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAppConf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAppConf() = %v, want %v", got, tt.want)
			}
		})
	}
}
