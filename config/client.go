package config

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	"github.com/go-ini/ini"
	consulapi "github.com/hashicorp/consul/api"
)

// CliConfig 服务配置
type CliConfig struct {
	ServerName string
	ConsulConf *consulapi.Config
	TlsConf    *tls.Config
}

// CliInit 服务配置初始化
// section 服务实例节点名称
func CliInit(section string) *CliConfig {
	cli := &CliConfig{}
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		log.Fatalln(err)
	}
	if len(section) == 0 {
		section = "client"
	}
	cfgSec := cfg.Section(section)

	certFile := cfgSec.Key("grpc_tls_certfile").String()
	if certFile != "" {
		b, err := ioutil.ReadFile(certFile)
		if err != nil {
			log.Fatalln(err)
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			log.Fatalln("credentials: failed to append certificates")
		}

		cli.TlsConf = &tls.Config{
			ServerName: cfgSec.Key("srv_name").String(),
			RootCAs:    cp,
		}
	}

	consulCfg := consulapi.DefaultConfig()
	consulCfg.Address = cfgSec.Key("consul_addrs").String()
	if cfgSec.Key("consul_acl_token").String() != "" {
		consulCfg.Token = cfgSec.Key("consul_acl_token").String()
	}

	cli.ConsulConf = consulCfg
	cli.ServerName = cfgSec.Key("srv_name").String()

	return cli
}
