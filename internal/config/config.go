package config

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	// s     sync.Once
	Pub, Pr any
	RestPub []byte
	RestPr  []byte
)

func DecPub(key string) (any, error) {
	c, b := pem.Decode([]byte(key))
	if c == nil || c.Type != "PUBLIC KEY" {
		return nil, errors.New("no public key in data")
	}
	pub, err := x509.ParsePKIXPublicKey(c.Bytes)
	if err != nil {
		return nil, err
	}
	Pub = pub
	RestPub = b
	return Pub, nil
}

func GetPub() (any, error) {
	if Pub == nil {
		return nil, errors.New("there is no crypto key in config")
	}
	return Pub, nil
}
func DecPr(key string) (any, error) {
	c, b := pem.Decode([]byte(key))
	if c == nil || c.Type != "PRIVAT KEY" {
		return nil, errors.New("no privat key in data")
	}
	pr, err := x509.ParsePKCS8PrivateKey(c.Bytes)
	if err != nil {
		return nil, err
	}
	Pr = pr
	RestPr = b
	return Pr, nil
}

func GetPr() (any, error) {
	if Pr == nil {
		return nil, errors.New("there is no crypto key in config")
	}
	return Pr, nil
}
