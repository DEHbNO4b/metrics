package handlers

import (
	"bytes"
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
	//create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	//create mock object
	m := mocks.NewMockMetricsStorage(ctrl)
	// гарантируем, что заглушка при вызове с аргументом data.Metrics{} вернёт nil
	m.EXPECT().SetMetric(gomock.Any()).Return(nil)
	metrics := NewMetrics(m)

	type want struct {
		code int
	}
	type args struct {
		body io.Reader
	}

	tests := []struct {
		name string
		want want
		args args
	}{
		{
			name: "negative case empty req body",
			want: want{
				code: 400,
			},
			args: args{body: bytes.NewReader([]byte(""))},
		},
		{
			name: "negative case wrong data",
			want: want{
				code: 400,
			},
			args: args{body: bytes.NewReader([]byte("wrong data"))},
		},
		{
			name: "positive case #1",
			want: want{
				code: 200,
			},
			args: args{body: bytes.NewReader([]byte(`{"ID":"some_id"}`))},
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/", tt.args.body)
			w := httptest.NewRecorder()
			metrics.SetMetricJSON(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func TestMetrics_SetMetricsJSON(t *testing.T) {
	//create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	//create mock object
	m := mocks.NewMockMetricsStorage(ctrl)
	// гарантируем, что заглушка при вызове с аргументом data.Metrics{} вернёт nil
	m.EXPECT().SetMetric(gomock.Any()).Return(nil)
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
		want want
		args args
	}{
		{
			name: "negative case: empty body",
			ms:   &metrics,
			want: want{
				code: 400,
			},
			args: args{body: bytes.NewReader([]byte(""))},
		},
		{
			name: "negative case: wrong data",
			ms:   &metrics,
			want: want{
				code: 400,
			},
			args: args{body: bytes.NewReader([]byte("wrong data"))},
		},
		{
			name: "positive case #1",
			want: want{
				code: 200,
			},
			args: args{body: bytes.NewReader([]byte(`[{"ID":"some_id"}]`))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/updates/", tt.args.body)
			w := httptest.NewRecorder()
			metrics.SetMetricsJSON(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func TestMetrics_GetMetricJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockMetricsStorage(ctrl)
	value := data.NewMetric()
	m.EXPECT().GetMetric(gomock.Any()).Return(value, nil)

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
		want want
		args args
	}{
		{
			name: "negative case: empty body",
			ms:   &metrics,
			want: want{
				code: 400,
			},
			args: args{body: bytes.NewReader([]byte(""))},
		},
		{
			name: "negative case: wrong data",
			ms:   &metrics,
			want: want{
				code: 400,
			},
			args: args{body: bytes.NewReader([]byte("wrong data"))},
		},
		{
			name: "positive case #1",
			want: want{
				code: 200,
			},
			args: args{body: bytes.NewReader([]byte(`{"ID":"some_id"}`))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/value/", tt.args.body)
			w := httptest.NewRecorder()
			metrics.GetMetricJSON(w, r)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}
