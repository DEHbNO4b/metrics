package config

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	logger "github.com/DEHbNO4b/metrics/internal/loger"
)

// var (
// 	endpoint       string
// 	key            string
// 	crypto         string
// 	reportInterval int
// 	pollInterval   int
// )

var (
	AgentCfg  AgentConfig
	agentOnce sync.Once
)

type AgentConfig struct {
	Adress         string `json:"adress"`          //"address": "localhost:8080", // аналог переменной окружения ADDRESS или флага -a
	ReportInterval int    `json:"report_interval"` //"report_interval": "1s", // аналог переменной окружения REPORT_INTERVAL или флага -r
	PollInterval   int    `json:"poll_interval"`   //"poll_interval": "1s", // аналог переменной окружения POLL_INTERVAL или флага -p
	CryptoKey      string `json:"crypto_key"`      //"crypto_key": "/path/to/key.pem" // аналог переменной окружения CRYPTO_KEY или флага -crypto-key
	HashKey        string
	ConfPath       string
}

func GetAgentCfg() AgentConfig {
	AgentCfg = AgentConfig{}
	agentOnce.Do(func() {
		parseAgentFlag()
		parseAgentEnv()
		if AgentCfg.ConfPath != "" {
			c, err := readAgentConfFile(AgentCfg.ConfPath)
			if err != nil {
				logger.Log.Error(err.Error())
			}
			if AgentCfg.Adress == "" {
				AgentCfg.Adress = c.Adress
			}
			if AgentCfg.ReportInterval == 0 {
				AgentCfg.ReportInterval = c.ReportInterval
			}
			if AgentCfg.PollInterval == 0 {
				AgentCfg.PollInterval = c.PollInterval
			}

			if AgentCfg.CryptoKey == "" {
				AgentCfg.CryptoKey = c.CryptoKey
			}
		}
	})

	return AgentCfg
}

func parseAgentFlag() {
	flag.StringVar(&AgentCfg.Adress, "a", "localhost:8080", "endpoint adress")
	flag.StringVar(&AgentCfg.HashKey, "k", "", "hash key")
	flag.IntVar(&AgentCfg.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&AgentCfg.PollInterval, "p", 2, "poll interval")
	flag.StringVar(&AgentCfg.CryptoKey, "crypto-key", "", "crypto config file path")
	flag.Parse()

}
func parseAgentEnv() {
	if ep := os.Getenv("ADDRESS"); ep != "" {
		AgentCfg.Adress = ep
	}
	if ck := os.Getenv("CRYPTO_KEY"); ck != "" {
		AgentCfg.CryptoKey = ck
	}
	if k := os.Getenv("KEY"); k != "" {
		AgentCfg.HashKey = k
	}
	if ri := os.Getenv("REPORT_INTERVAL"); ri != "" {
		rInt, err := strconv.Atoi(ri)
		if err != nil {
			AgentCfg.ReportInterval = rInt
		}

	}
	if pi := os.Getenv("POLL_INTERVAL"); pi != "" {
		pInt, err := strconv.Atoi(pi)
		if err != nil {
			AgentCfg.PollInterval = pInt
		}
	}
	if AgentCfg.CryptoKey != "" {
		_, err := DecPub(AgentCfg.CryptoKey)
		if err != nil {
			fmt.Println(err)
		}
	}
}
func readAgentConfFile(path string) (AgentConfig, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return AgentConfig{}, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	var text = ""
	for sc.Scan() {
		line := sc.Text()
		if strings.Contains(line, "//") {
			str := strings.Split(line, "//")
			line = str[0] + "\n"
		}
		text = text + line
	}
	if err := sc.Err(); err != nil {
		return AgentConfig{}, err
	}

	cfg := AgentConfig{}
	err = json.Unmarshal([]byte(text), &cfg)
	if err != nil {
		return AgentConfig{}, err
	}
	return cfg, nil
}
