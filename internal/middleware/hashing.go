package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"net/http"

	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

type Hash struct {
	Key []byte
}

func (h *Hash) WithHash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HashSHA256") != "" {
			h := hmac.New(sha256.New, h.Key)
			b, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Log.Error("UNABLE TO READALL from r.Body()", zap.String("err:", err.Error()))
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			h.Write(b)
			dst := h.Sum(nil)
			logger.Log.Sugar().Infof("%x", r.Header.Get("HashSHA256"))
			logger.Log.Sugar().Infof("%x", dst)
			if !hmac.Equal(dst, []byte(r.Header.Get("HashSHA256"))) {
				logger.Log.Sugar().Error("BAD REQUEST HASH")
				http.Error(w, "", http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
