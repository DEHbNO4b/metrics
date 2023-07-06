package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_SetMetrics(t *testing.T) {
	memSt := NewMemStorage()

	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		ms      *MemStorage
		request string
		want    want
	}{
		{
			ms:      memSt,
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
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestMemStorage_SetGauge(t *testing.T) {

	memSt := NewMemStorage()

	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		ms      *MemStorage
		request string
		want    want
	}{
		{
			ms:      memSt,
			name:    "positive test #1",
			request: "/update/gauge/somemetric/300",
			want: want{
				statusCode: 200,
			},
		},
		{
			ms:      memSt,
			name:    "negative test #2",
			request: "/update/gauge/somemetric/k",
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
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestMemStorage_SetCounter(t *testing.T) {
	memSt := NewMemStorage()

	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name    string
		ms      *MemStorage
		request string
		want    want
	}{
		{
			ms:      memSt,
			name:    "positive test #1",
			request: "/update/counter/somemetric/3500",
			want: want{
				statusCode: 200,
			},
		},
		{
			ms:      memSt,
			name:    "negative test #2",
			request: "/update/counter/somemetric/k",
			want: want{
				statusCode: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			tt.ms.SetCounter(w, req)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
