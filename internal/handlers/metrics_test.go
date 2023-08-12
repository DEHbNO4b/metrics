package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DEHbNO4b/metrics/internal/maindb"
	"github.com/stretchr/testify/assert"
)

func TestMetrics_SetMetricsJSON(t *testing.T) {
	store := maindb.NewRamStore(maindb.StoreConfig{})
	metrics := NewMetrics(store)
	b := []byte("")
	type want struct {
		code int
		// response    string

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
		{
			name: "positiv test",
			ms:   &metrics,
			args: args{
				body: bytes.NewReader([]byte(`{"id":"some","type":"gauge","value":100}`)),
			},
			want: want{
				code: 200,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/update/", test.args.body)
			w := httptest.NewRecorder()
			test.ms.SetMetricsJSON(w, request)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			res.Body.Close()
		})
	}
}

func TestMetrics_GetMetricJSON(t *testing.T) {
	store := maindb.NewRamStore(maindb.StoreConfig{})
	metrics := NewMetrics(store)
	type want struct {
		code int
		// response    string

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
