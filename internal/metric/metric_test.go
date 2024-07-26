package metric

import (
	"reflect"
	"testing"

	"github.com/hrapovd1/msg-proc/internal/types"
)

func TestNewMetrics(t *testing.T) {
	type args struct {
		db types.Storager
	}
	tests := []struct {
		name string
		args args
		want *Metrics
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMetrics(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
