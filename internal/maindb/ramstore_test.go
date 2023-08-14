package maindb

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/DEHbNO4b/metrics/internal/data"
// )

// func TestRamStore_GetMetric(t *testing.T) {
// 	store := NewRAMStore(StoreConfig{})
// 	type args struct {
// 		met data.Metrics
// 	}
// 	tests := []struct {
// 		name    string
// 		ms      *RAMStore
// 		args    args
// 		want    data.Metrics
// 		wantErr bool
// 	}{
// 		{
// 			name:    "empty id",
// 			ms:      store,
// 			args:    args{met: data.Metrics{}},
// 			want:    data.Metrics{},
// 			wantErr: true,
// 		},
// 		{
// 			name:    "wrong type",
// 			ms:      store,
// 			args:    args{met: data.Metrics{MType: "wrongType"}},
// 			want:    data.Metrics{},
// 			wantErr: true,
// 		},
// 		{
// 			name:    "positive test",
// 			ms:      store,
// 			args:    args{met: data.Metrics{MType: "gauge", ID: "someName"}},
// 			want:    data.Metrics{},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := tt.ms.GetMetric(tt.args.met)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("RamStore.GetMetric() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("RamStore.GetMetric() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
