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
		statusCode int
		// response    string
		// contentType string
	}

	tests := []struct {
		name    string
		ms      *Metrics
		request string
		want    want
	}{
		{
			ms:      &memSt,
			name:    "negative test #1",
			request: "/update/",
			want: want{
				statusCode: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			tt.ms.SetMetricsJSON(w, req)
			result := w.Result()
			result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
