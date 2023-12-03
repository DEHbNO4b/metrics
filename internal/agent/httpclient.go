package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DEHbNO4b/metrics/internal/config"
	"github.com/DEHbNO4b/metrics/internal/data"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

type HttpClient struct {
	client http.Client
	addr   string
}

func NewHttpCLient(addr string) *HttpClient {
	cl := http.Client{Timeout: 1000 * time.Millisecond}
	return &HttpClient{
		client: cl,
		addr:   addr,
	}
}
func (h *HttpClient) SendMetric(ctx context.Context, m data.Metrics, key string) {
	var req *http.Request
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(&m)
	if err != nil {
		logger.Log.Info("unable to encode metric", zap.String("err: ", err.Error()))
		return
	}
	compressed := bytes.Buffer{}
	compressor, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		return
	}
	compressor.Write(buf.Bytes())
	compressor.Close()
	req, err = http.NewRequest(http.MethodPost, h.addr+"/update/", &compressed) // (1)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-encoding", "gzip")
	req.Header.Add("Accept-encoding", "gzip")
	if key != "" {
		b := signature(key, buf.Bytes())
		req.Header.Add("HashSHA256", string(b))
	}
	resp, err := h.client.Do(req)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		return
	}
	resp.Body.Close()
}
func signature(key string, b []byte) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(b)
	dst := h.Sum(nil)
	logger.Log.Sugar().Infof("%x", dst)
	return dst
}
func (h *HttpClient) SendMetrics(ctx context.Context, metrics []data.Metrics) {
	fmt.Println("http send metrics")
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(&metrics)
	if err != nil {
		logger.Log.Info("unable to encode metric", zap.String("err: ", err.Error()))
		return
	}
	compressed := bytes.Buffer{}
	compressor, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		return
	}
	compressor.Write(buf.Bytes())
	compressor.Close()
	mes, err := encrypt(compressed)
	if err != nil {
		logger.Log.Error(err.Error())
	} else {
		compressed = mes
	}
	req, err := http.NewRequest(http.MethodPost, h.addr+"/updates/", &compressed) // (1)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-encoding", "gzip")
	req.Header.Add("Accept-encoding", "gzip")
	resp, err := h.client.Do(req)
	if err != nil {
		logger.Log.Error("request returned err ", zap.Error(err))
		return
	}
	resp.Body.Close()
}
func encrypt(b bytes.Buffer) (bytes.Buffer, error) {
	k, err := config.GetPub()
	if err != nil {
		return b, err
	}
	rng := rand.Reader
	pub, ok := k.(rsa.PublicKey)
	if !ok {
		return b, errors.New("wrong crypto key")
	}
	text, err := rsa.EncryptOAEP(sha256.New(), rng, &pub, b.Bytes(), []byte("metrics"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		return b, err
	}
	buf := bytes.Buffer{}
	buf.Write(text)
	return buf, nil
}
