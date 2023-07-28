package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var contentToCompress = []string{"application/json", "text/html"}

func isNeedToCompress(s string) bool {
	for _, el := range contentToCompress {
		if s == el {
			return true
		}
	}
	return false
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandle(next http.Handler) http.Handler {
	//&& isNeedToCompress(r.Header.Get("Accept"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("headers: ", r.Header)
		var nextWriter = w
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && isNeedToCompress(r.Header.Get("Accept")) {
			fmt.Println("i find accept-encoding:gzip")
			gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer gz.Close()
			w.Header().Set("Content-Encoding", "gzip")
			nextWriter = gzipWriter{ResponseWriter: w, Writer: gz}
		}
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gzr, err := gzip.NewReader(r.Body)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer gzr.Close()

			// b := bytes.Buffer{}
			// buf := make([]byte, 30)
			// gzr.Read(buf)
			// b.ReadFrom(gzr)
			// fmt.Println("buf: ", buf)
			r.Body = gzr

		}

		next.ServeHTTP(nextWriter, r)
	})
}
