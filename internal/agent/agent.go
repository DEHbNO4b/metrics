// packag agent provides automatic collection of metrics and sending them to the server.
package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/DEHbNO4b/metrics/internal/config"
	"github.com/DEHbNO4b/metrics/internal/data"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

// Agent struct collects runtume metrics and send them to server
type Agent struct {
	m      *runtime.MemStats
	client http.Client
	url    string
	gauges []data.Metrics
}

// NewAgent return Agent structure with endpoint path inside.
func NewAgent(endpoint string) Agent {
	var m runtime.MemStats
	cl := http.Client{Timeout: 1000 * time.Millisecond}
	u := "http://" + endpoint
	a := Agent{m: &m, client: cl, url: u, gauges: data.NewGauges()}
	return a
}

// ReadRuntimeMetrics reads runtume metrics.
func (a Agent) ReadRuntimeMetrics(interval int) {
	var pollInterval = time.Duration(interval) * time.Second
	for {
		runtime.ReadMemStats(a.m)
		time.Sleep(pollInterval)
	}

}

// PullMetrics sends metrics to server.
func (a Agent) PullMetrics(interval int, key, crypto string) {

	var reportInterval = time.Duration(interval) * time.Second
	metrics := make([]data.Metrics, 0, 30)
	for {
		for _, el := range a.gauges {
			el.ReadValue(a.m)
			metrics = append(metrics, el)
		}
		var d int64 = 1
		metrics = append(metrics, data.Metrics{ID: "PollCount", MType: "counter", Delta: &d})
		go a.sendMetrics(metrics)
		time.Sleep(reportInterval)
	}

}

func (a Agent) sendMetrics(metrics []data.Metrics) {
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
	req, err := http.NewRequest(http.MethodPost, a.url+"/updates/", &compressed) // (1)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-encoding", "gzip")
	req.Header.Add("Accept-encoding", "gzip")
	resp, err := a.client.Do(req)
	if err != nil {
		logger.Log.Error("request returned err ", zap.Error(err))
		return
	}
	resp.Body.Close()
}
func (a Agent) sendMetric(m data.Metrics, key string) {
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
	req, err = http.NewRequest(http.MethodPost, a.url+"/update/", &compressed) // (1)
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
	resp, err := a.client.Do(req)
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
