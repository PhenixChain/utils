package config

import (
	"crypto/tls"
	"log"
	"net"
	"os"

	"github.com/go-ini/ini"
	consulapi "github.com/hashicorp/consul/api"
)

// SrvConfig 服务配置
type SrvConfig struct {
	ServerName string
	ConsulConf *consulapi.Config
	TlsConf    *tls.Config
	LogTrack   bool
}

// SrvInit 服务配置初始化
func SrvInit() *SrvConfig {
	srv := &SrvConfig{}
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		log.Fatalln(err)
	}
	cfgSec := cfg.Section("server")

	addr := cfgSec.Key("srv_addrs").String()
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatalln(err)
	}

	// Bind address for the server
	err = os.Setenv("MICRO_SERVER_ADDRESS", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalln(err)
	}

	// Used instead of the server address when registering with discovery
	err = os.Setenv("MICRO_SERVER_ADVERTISE", addr)
	if err != nil {
		log.Fatalln(err)
	}

	certFile := cfgSec.Key("grpc_tls_certfile").String()
	keyFile := cfgSec.Key("grpc_tls_keyfile").String()
	if certFile != "" && keyFile != "" {
		cer, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalln(err)
		}
		srv.TlsConf = &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cer},
		}
	}

	consulCfg := consulapi.DefaultConfig()
	consulCfg.Address = cfgSec.Key("consul_addrs").String()
	if cfgSec.Key("consul_acl_token").String() != "" {
		consulCfg.Token = cfgSec.Key("consul_acl_token").String()
	}

	srv.ConsulConf = consulCfg
	srv.ServerName = cfgSec.Key("srv_name").String()
	srv.LogTrack, _ = cfgSec.Key("log_track_enable").Bool()

	return srv
}
