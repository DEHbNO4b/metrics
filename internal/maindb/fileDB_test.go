package maindb

import (
	"testing"

	"github.com/DEHbNO4b/metrics/internal/data"
)

func TestFileDB_WriteMetrics(t *testing.T) {
	f := NewFileDB("../../test.json")
	type args struct {
		data []data.Metrics
	}
	metricsData := make([]data.Metrics, 0, 3)
	metr := data.NewMetric()
	metr.ID = "some_new_id"
	metr.MType = "gauge"
	*metr.Value = 3.14
	metricsData = append(metricsData, metr)

	tests := []struct {
		name    string
		f       *FileDB
		args    args
		wantErr bool
	}{
		{
			name: "pozitive case",
			f:    f,
			args: args{
				data: metricsData,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.WriteMetrics(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("FileDB.WriteMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileDB_ReadMetrics(t *testing.T) {
	f := NewFileDB("../../test.json")
	tests := []struct {
		name    string
		f       *FileDB
		want    []data.Metrics
		wantErr bool
	}{
		{
			name:    "pozitive case",
			f:       f,
			wantErr: false,
		},
		{
			name:    "negative case",
			f:       NewFileDB("wrong file path"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.f.ReadMetrics()

			if (err != nil) != tt.wantErr {
				t.Errorf("FileDB.ReadMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func BenchmarkFileWriteMetrics(b *testing.B) {
	f := NewFileDB("../../test.json")
	// metricsData := make([]data.Metrics, 0, 3)
	// metr := data.NewMetric()
	// metr.ID = "test"
	// metr.MType = "gauge"
	// *metr.Value = 3.14
	// metricsData = append(metricsData, metr)
	metricsData := getRundomMetrics(27)

	for i := 0; i < b.N; i++ {
		f.WriteMetrics(metricsData)
	}
}
func BenchmarkFileReadMetrics(b *testing.B) {

	f := NewFileDB("../../test.json")
	// for _, size := range []int{1, 10, 100, 1000, 10000} {
	metricsData := getRundomMetrics(20)
	f.WriteMetrics(metricsData)

	b.ResetTimer()
	// name := fmt.Sprintf("Contains-%d", size)
	// b.Run(name, func(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f.ReadMetrics()
	}
	// })
	// }
}
