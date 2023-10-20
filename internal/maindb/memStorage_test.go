package maindb

import (
	"testing"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_SetMetric(t *testing.T) {
	type args struct {
		metric data.Metrics
	}
	m := data.NewMetric()
	m.MType = "counter"
	bad := data.NewMetric()
	bad.MType = "wrong type"
	store := NewMemStorage()
	tests := []struct {
		name    string
		rs      *MemStorage
		args    args
		wantErr bool
	}{
		{
			name:    "pozitive case",
			rs:      store,
			args:    args{metric: m},
			wantErr: false,
		},
		{
			name:    "negative case",
			rs:      store,
			args:    args{metric: bad},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rs.SetMetric(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.SetMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_GetMetrics(t *testing.T) {
	store := NewMemStorage()
	store.Gauges["test"] = 3.14

	tests := []struct {
		name string
		rs   *MemStorage
	}{
		{
			name: "pozitive case",
			rs:   store,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rs.GetMetrics()
			// if got := tt.rs.GetMetrics(); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("MemStorage.GetMetrics() = %v, want %v", got, tt.want)
			// }
			assert.NotNil(t, got)
		})
	}
}

func TestMemStorage_GetMetric(t *testing.T) {
	type args struct {
		met data.Metrics
	}
	good := data.NewMetric()
	good.MType = "gauge"
	good.ID = "test"
	bad := data.NewMetric()
	bad.MType = "bad type"

	store := NewMemStorage()
	store.Gauges["test"] = 3.14
	tests := []struct {
		name    string
		rs      *MemStorage
		args    args
		want    data.Metrics
		wantErr bool
	}{
		{
			name:    "pozitive case",
			rs:      store,
			args:    args{met: good},
			wantErr: false,
		},
		{
			name:    "negative case",
			rs:      store,
			args:    args{met: bad},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.rs.GetMetric(tt.args.met)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.GetMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
