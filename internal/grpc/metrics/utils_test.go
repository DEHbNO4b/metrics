package metrics

import (
	"testing"

	pb "github.com/DEHbNO4b/metrics/proto"
)

func Test_validate(t *testing.T) {
	type args struct {
		m *pb.Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
