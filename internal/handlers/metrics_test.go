package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestMetrics_SetMetricsJSON(t *testing.T) {
	store := data.NewMetStore(data.StoreConfig{})
	metrics := NewMetrics(store)
	b := []byte("")
	type want struct {
		code int
		// response    string
		// contentType string
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
