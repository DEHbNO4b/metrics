package middleware

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/DEHbNO4b/metrics/internal/config"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
)

func CryptoHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		mes := bytes.NewReader(r.Body)

		k, err := decrypt()
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
		c, err := gzip.NewReader(r.Body)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		// defer gzr.Close()
		r.Body = c

		next.ServeHTTP(w, r)
	})
}
func decrypt(b *bytes.Reader) (*bytes.Reader, error) {
	k, err := config.GetPr()
	if err != nil {
		return b, err
	}
	rng := rand.Reader
	pub, ok := k.(rsa.PrivateKey)
	if !ok {
		return b, errors.New("wrong crypto key")
	}
	mes := make([]byte, 0)
	b.Read(mes)
	text, err := rsa.DecryptOAEP(sha256.New(), rng, &pub, mes, []byte("metrics"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		return b, err
	}
	r := bytes.NewReader(text)
	return r, nil
}
