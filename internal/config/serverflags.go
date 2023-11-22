package config

import (
	"bufio"
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"strings"
	"sync"

	logger "github.com/DEHbNO4b/metrics/internal/loger"
)

type ServerConfig struct {
	Adress        string `json:"address"`
	Restore       bool   `json:"restore"`
	StoreInterval int    `json:"store_intervAL"`
	StoreFile     string `json:"store_file"`
	Dsn           string `json:"database_dsn"`
	CryptoKey     string `json:"crypto_key"`
	HashKey       string
	ConfPath      string
}

var (
	Cfg  ServerConfig
	once sync.Once
)

func GetServCfg() ServerConfig {
	once.Do(func() {
		parseFlag()
		parseEnv()
		if Cfg.ConfPath != "" {
			c, err := readServConfFile(Cfg.ConfPath)
			if err != nil {
				logger.Log.Error(err.Error())
			}
			if Cfg.Adress == "" {
				Cfg.Adress = c.Adress
			}
			if !Cfg.Restore {
				Cfg.Restore = c.Restore
			}
			if Cfg.StoreInterval == 0 {
				Cfg.StoreInterval = c.StoreInterval
			}
			if Cfg.StoreFile == "" {
				Cfg.StoreFile = c.StoreFile
			}
			if Cfg.Dsn == "" {
				Cfg.Dsn = c.Dsn
			}
			if Cfg.CryptoKey == "" {
				Cfg.CryptoKey = c.CryptoKey
			}
		}
	})
	return Cfg
}
func parseFlag() {
	Cfg = ServerConfig{}
	flag.StringVar(&Cfg.Adress, "a", "localhost:8080", "adress and port for running")
	flag.StringVar(&Cfg.HashKey, "k", "", "hash key")
	flag.StringVar(&Cfg.CryptoKey, "crypto-key", "", "crypto config file path")
	flag.StringVar(&Cfg.StoreFile, "f", "/tmp/metrics-db.json", "file storage path")
	flag.StringVar(&Cfg.Dsn, "d", "", "dsn for postgres")
	flag.StringVar(&Cfg.ConfPath, "c", "", "path for conf file")
	flag.IntVar(&Cfg.StoreInterval, "i", 300, "data store interval")
	flag.BoolVar(&Cfg.Restore, "r", true, "restore_flag")
	flag.Parse()

}
func parseEnv() {
	if ep := os.Getenv("ADDRESS"); ep != "" {
		Cfg.Adress = ep
	}
	if si := os.Getenv("STORE_INTERVAL"); si != "" {
		sInt, err := strconv.Atoi(si)
		if err != nil {
			logger.Log.Sugar().Error("unable to convert STORE_INTERVAL to int", err.Error())
			return
		}
		Cfg.StoreInterval = sInt
	}
	if ck := os.Getenv("CRYPTO_KEY"); ck != "" {
		Cfg.CryptoKey = ck
	}
	if fp := os.Getenv("FILE_STORAGE_PATH"); fp != "" {
		Cfg.StoreFile = fp
	}
	if cnf := os.Getenv("CONFIG"); cnf != "" {
		Cfg.ConfPath = cnf
	}
	if k := os.Getenv("KEY"); k != "" {
		Cfg.HashKey = k
	}
	if dbdsn := os.Getenv("DATABASE_DSN"); dbdsn != "" {
		Cfg.Dsn = dbdsn
	}
	if r := os.Getenv("RESTORE"); r != "" {
		re, err := strconv.ParseBool(r)
		if err != nil {
			logger.Log.Sugar().Error("unable to convert STORE_INTERVAL to int", err.Error())
			return
		}
		Cfg.Restore = re
	}
	if Cfg.CryptoKey != "" {
		_, err := DecPr(Cfg.CryptoKey)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}
}
func readServConfFile(path string) (ServerConfig, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return ServerConfig{}, err
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
		return ServerConfig{}, err
	}

	cfg := ServerConfig{}
	err = json.Unmarshal([]byte(text), &cfg)
	if err != nil {
		return ServerConfig{}, err
	}
	return cfg, nil
}
