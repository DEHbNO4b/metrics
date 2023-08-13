package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetrics_SetMetricJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockMetricsStorage(ctrl)
	// var val float64 = 100
	m.EXPECT().SetMetric(gomock.Any()).Return(errors.New("some error")).MinTimes(0)

	metrics := NewMetrics(m)
	b := []byte("")
	type want struct {
		code int
	}
	type args struct {
		body io.Reader
	}
	tests := []struct {
		name string
		ms   *Metrics
		args args
		want want
	}{
		{
			name: "empty request body",
			ms:   &metrics,
			args: args{
				body: bytes.NewReader(b),
			},
			want: want{
				code: 400,
			},
		},
		// {
		// 	name: "positiv test",
		// 	ms:   &metrics,
		// 	args: args{
		// 		body: bytes.NewReader([]byte(`{"id":"some","type":"gauge","value":100}`)),
		// 	},
		// 	want: want{
		// 		code: 200,
		// 	},
		// },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/", test.args.body)
			w := httptest.NewRecorder()
			test.ms.SetMetricJSON(w, request)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			res.Body.Close()
		})
	}
}

func TestMetrics_GetMetricJSON(t *testing.T) {
	// store := maindb.NewRAMStore(maindb.StoreConfig{}, maindb.NewFileDB(""))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockMetricsStorage(ctrl)
	m.EXPECT().
		GetMetric(gomock.Any()).Return(data.Metrics{}, errors.New("some error")).MinTimes(0)

	metrics := NewMetrics(m)
	type want struct {
		code int
	}

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name string
		ms   *Metrics
		args args
		want want
	}{
		{
			name: "emty body",
			ms:   &metrics,
			args: args{body: bytes.NewReader([]byte(""))},
			want: want{code: 400},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/", tt.args.body)
			w := httptest.NewRecorder()
			tt.ms.GetMetricJSON(w, request)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			res.Body.Close()
		})
	}
}
