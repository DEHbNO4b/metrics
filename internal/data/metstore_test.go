package data

import (
	"reflect"
	"testing"
)

func TestMetStore_GetMetric(t *testing.T) {
	store := NewMetStore(StoreConfig{})
	type args struct {
		met Metrics
	}
	tests := []struct {
		name    string
		ms      *MetStore
		args    args
		want    Metrics
		wantErr bool
	}{
		{
			name:    "empty id",
			ms:      store,
			args:    args{met: Metrics{}},
			want:    Metrics{},
			wantErr: true,
		},
		{
			name:    "wrong type",
			ms:      store,
			args:    args{met: Metrics{MType: "wrongType"}},
			want:    Metrics{},
			wantErr: true,
		},
		{
			name:    "wrong type",
			ms:      store,
			args:    args{met: Metrics{MType: "gauge", ID: "someName"}},
			want:    Metrics{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ms.GetMetric(tt.args.met)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetStore.GetMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetStore.GetMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
