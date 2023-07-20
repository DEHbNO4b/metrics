package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestMetrics_SetMetrics(t *testing.T) {
	store := data.MetStore{}
	memSt := NewMetrics(&store)

	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		ms      *Metrics
		request string
		want    want
	}{
		{
			ms:      &memSt,
			name:    "positive test #1",
			request: "/update/something/somemetric/300",
			want: want{
				statusCode: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			tt.ms.SetGauge(w, req)
			result := w.Result()
			result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestMetrics_SetGauge(t *testing.T) {
	store := data.NewMetStore()
	memSt := NewMetrics(store)

	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		ms      *Metrics
		request string
		want    want
	}{
		{
			ms:      &memSt,
			name:    "positive test ",
			request: "/update/gauge/somemetric/300",
			want: want{
				statusCode: 200,
			},
		},
		{
			ms:      &memSt,
			name:    "negative test ",
			request: "/update/gauge/somemetric/k",
			want: want{
				statusCode: 400,
			},
		},
		{
			ms:      &memSt,
			name:    "zero test",
			request: "/update/gauge/somemetric/0",
			want: want{
				statusCode: 200,
			},
		},
		{
			ms:      &memSt,
			name:    "big number ",
			request: "/update/gauge/somemetric/9845649.8816513",
			want: want{
				statusCode: 200,
			},
		},
		{
			ms:      &memSt,
			name:    "big negative number ",
			request: "/update/gauge/somemetric/-9845649.8816513",
			want: want{
				statusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			tt.ms.SetGauge(w, req)
			result := w.Result()
			result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestMetrics_SetCounter(t *testing.T) {
	store := data.NewMetStore()
	memSt := NewMetrics(store)

	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		ms      *Metrics
		request string
		want    want
	}{
		{
			ms:      &memSt,
			name:    "positive test ",
			request: "/update/counter/somemetric/3500",
			want: want{
				statusCode: 200,
			},
		},
		{
			ms:      &memSt,
			name:    "negative test ",
			request: "/update/counter/somemetric/k",
			want: want{
				statusCode: 400,
			},
		},
		{
			ms:      &memSt,
			name:    "zero test ",
			request: "/update/counter/somemetric/0",
			want: want{
				statusCode: 200,
			},
		},
		{
			ms:      &memSt,
			name:    "big number ",
			request: "/update/counter/somemetric/98456498816513",
			want: want{
				statusCode: 200,
			},
		},
		{
			ms:      &memSt,
			name:    "big negative number ",
			request: "/update/counter/somemetric/-98456498816513",
			want: want{
				statusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			tt.ms.SetCounter(w, req)
			result := w.Result()
			result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
