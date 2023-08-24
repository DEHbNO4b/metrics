package middleware

import (
	"bytes"
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
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// b, err := io.ReadAll(req.Body)
		// if err != nil {
		// 	logger.Log.Error("UNABLE TO READALL from r.Body()", zap.String("err:", err.Error()))
		// 	http.Error(w, "", http.StatusBadRequest)
		// 	return
		// }
		// bodyBytes, _ := io.ReadAll(req.Body)
		// req.Body.Close() //  must close
		// req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// b := make([]byte, 0, 1000)
		// n, err := req.Body.Read(b)
		// if err != nil {
		// 	logger.Log.Error("UNABLE TO READALL from r.Body()", zap.String("err:", err.Error()))
		// }
		// c := b[0:n]

		if req.Header.Get("HashSHA256") != "" {

			bodyBytes, _ := io.ReadAll(req.Body)
			req.Body.Close() //  must close

			logger.Log.Info("body: ", zap.String("body:", string(bodyBytes)))
			hasher := hmac.New(sha256.New, h.Key)
			// b, err := io.ReadAll(req.Body)
			logger.Log.Info("body: ", zap.String("body:", string(bodyBytes)))

			hasher.Write(bodyBytes)
			dst := hasher.Sum(nil)
			logger.Log.Sugar().Infof("%x", req.Header.Get("HashSHA256"))
			logger.Log.Sugar().Infof("%x", dst)
			if !hmac.Equal(dst, []byte(req.Header.Get("HashSHA256"))) {
				logger.Log.Sugar().Error("BAD REQUEST HASH")
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		next.ServeHTTP(w, req)
	})
}
