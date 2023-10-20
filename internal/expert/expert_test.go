package expert

import (
	"testing"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/internal/maindb"
)

func TestExpert_SetMetric(t *testing.T) {
	store := maindb.NewMemStorage()
	expert := NewExpert(WithRAM(store))
	good := data.NewMetric()
	good.MType = "gauge"
	good.ID = "some_id"
	bad := data.NewMetric()
	bad.MType = "bad type"

	type args struct {
		m data.Metrics
	}
	tests := []struct {
		name    string
		e       *Expert
		args    args
		wantErr bool
	}{
		{
			name:    "positive case",
			e:       expert,
			args:    args{m: good},
			wantErr: false,
		},
		{
			name:    "negative case",
			e:       expert,
			args:    args{m: bad},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.SetMetric(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("Expert.SetMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExpert_setMetricToRAM(t *testing.T) {

	store := maindb.NewMemStorage()
	store.Counters["test"] = 10
	expert := NewExpert(WithRAM(store))
	good := data.NewMetric()
	good.MType = "counter"
	good.ID = "test"
	*good.Delta = 10

	bad := data.NewMetric()
	bad.MType = "bad type"

	type args struct {
		m data.Metrics
	}
	tests := []struct {
		name    string
		e       *Expert
		args    args
		wantErr bool
	}{
		{
			name:    "positive case",
			e:       expert,
			args:    args{m: good},
			wantErr: false,
		},
		{
			name:    "negative case",
			e:       expert,
			args:    args{m: bad},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.setMetricToRAM(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("Expert.setMetricToRAM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
