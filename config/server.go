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

var SrvCfg SrvConfig

func SrvCfgInit() {
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
		SrvCfg.TlsConf = &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cer},
			//NextProtos:   []string{"h2"},
		}
	}

	consulCfg := consulapi.DefaultConfig()
	consulCfg.Address = cfgSec.Key("consul_addrs").String()
	if cfgSec.Key("consul_acl_token").String() != "" {
		consulCfg.Token = cfgSec.Key("consul_acl_token").String()
	}

	SrvCfg.ConsulConf = consulCfg
	SrvCfg.ServerName = cfgSec.Key("srv_name").String()
	SrvCfg.LogTrack, _ = cfgSec.Key("log_track_enable").Bool()
}
