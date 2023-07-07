package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DEHbNO4b/metrics/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestIsPostReq(t *testing.T) {
	memSt := handlers.NewMemStorage()

	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		ms   *handlers.MemStorage
		req  *http.Request
		want want
	}{
		{
			ms:   memSt,
			name: "positive test ",
			req:  httptest.NewRequest(http.MethodPost, "/update/counter/somemetric/100", nil),
			want: want{
				statusCode: 200,
			},
		},
		{
			ms:   memSt,
			name: "negative test ",
			req:  httptest.NewRequest(http.MethodGet, "/update/counter/somemetric/100", nil),
			want: want{
				statusCode: 405,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			IsPostReq(testHend{})
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}

	// type args struct {
	// 	next http.Handler
	// }
	// tests := []struct {
	// 	name string
	// 	args args
	// 	want http.Handler
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		if got := IsPostReq(tt.args.next); !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("IsPostReq() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}

type testHend struct {
	req *http.Request
	w   http.ResponseWriter
}

func (th testHend) ServeHTTP(http.ResponseWriter, *http.Request) {

}
