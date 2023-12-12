package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/BurntSushi/toml"
	"go.etcd.io/etcd/client/pkg/v3/tlsutil"
	etcd "go.etcd.io/etcd/client/v3"
	"time"
)

type Config struct {
	EtcdClient EtcdClientConfig

	EtcdConfigKey string
}

type EtcdClientConfig struct {
	Endpoints        []string
	AutoSyncInterval string
	DialTimeout      string

	Username string
	Password string

	InsecureTransport     bool
	InsecureSkipTLSVerify bool
	CertFile              string
	KeyFile               string
	TrustedCAFile         string
}

func LoadConfig(data string) (*Config, error) {
	config := &Config{}
	if _, err := toml.Decode(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func ParseEtcdClientConfig(src EtcdClientConfig) (*etcd.Config, error) {
	if len(src.Endpoints) == 0 {
		return nil, errors.New("config with empty endpoints")
	}
	ret := &etcd.Config{}
	ret.Endpoints = src.Endpoints
	ret.AutoSyncInterval = time.Minute
	if autoSyncInterval, err := time.ParseDuration(
		src.AutoSyncInterval); err == nil {
		ret.AutoSyncInterval = autoSyncInterval
	}
	ret.DialTimeout = time.Second * 10
	if dialTimeout, err := time.ParseDuration(
		src.DialTimeout); err == nil {
		ret.DialTimeout = dialTimeout
	}

	ret.Username = src.Username
	ret.Password = src.Password

	if src.InsecureTransport {
		return ret, nil
	}

	var (
		cert     *tls.Certificate
		certPool *x509.CertPool
	)
	if src.CertFile != "" && src.KeyFile != "" {
		if nc, err := tlsutil.NewCert(src.CertFile,
			src.KeyFile, nil); err != nil {
			return nil, err
		} else {
			cert = nc
		}
	}
	if src.TrustedCAFile != "" {
		if ncp, err := tlsutil.NewCertPool(
			[]string{src.TrustedCAFile}); err != nil {
			return nil, err
		} else {
			certPool = ncp
		}
	}
	tlsCfg := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: src.InsecureSkipTLSVerify,
		RootCAs:            certPool,
	}
	if cert != nil {
		tlsCfg.Certificates = []tls.Certificate{*cert}
	}
	ret.TLS = tlsCfg
	return ret, nil
}
