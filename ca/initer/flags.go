package initer

import (
	"flag"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"strings"

	cfssl_config "github.com/flowshield/cfssl/config"
	"github.com/spf13/viper"

	"github.com/flowshield/flowshield/ca/core"
	"github.com/flowshield/flowshield/ca/core/config"
)

const (
	G_       = "IS"
	ConfName = "conf"
)

func parseConfigs(c *cli.Context) (core.Config, error) {
	// Cmdline flags
	flag.Parse()
	// Default config
	viper.SetConfigName(fmt.Sprintf("%v", ConfName))
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return core.Config{}, err
	}

	// Merge config frm ENV
	viper.SetEnvPrefix(G_)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	hs, err := os.Hostname()
	if err != nil {
		return core.Config{}, err
	}
	hostname := hs
	if v := os.Getenv("HOSTNAME"); v != "" {
		hostname = v
	}

	conf := core.Config{IConfig: config.IConfig{
		Log: config.Log{
			LogProxy: config.LogProxy{
				Host: viper.GetString("log.log-proxy.host"),
				Port: viper.GetInt("log.log-proxy.port"),
				Key:  viper.GetString("log.log-proxy.key"),
			},
		},
		Keymanager: config.Keymanager{
			SelfSign: viper.GetBool("keymanager.self-sign"),
			CsrTemplates: config.CsrTemplates{
				RootCa: config.RootCa{
					O:      viper.GetString("keymanager.csr-templates.root-ca.o"),
					Expiry: viper.GetString("keymanager.csr-templates.root-ca.expiry"),
				},
				IntermediateCa: config.IntermediateCa{
					O:      viper.GetString("keymanager.csr-templates.intermediate-ca.o"),
					Ou:     viper.GetString("keymanager.csr-templates.intermediate-ca.ou"),
					Expiry: viper.GetString("keymanager.csr-templates.intermediate-ca.expiry"),
				},
			},
		},
		Singleca: config.Singleca{
			ConfigPath: viper.GetString("singleca.config-path"),
		},
		OCSPHost: viper.GetString("ocsp-host"),
		HTTP: config.HTTP{
			OcspListen: viper.GetString("http.ocsp-listen"),
			CaListen:   viper.GetString("http.ca-listen"),
			Listen:     viper.GetString("http.listen"),
		},
		Debug:    viper.GetBool("debug"),
		Version:  viper.GetString("version"),
		Hostname: hostname,
	}}

	cfg, err := cfssl_config.LoadFile(conf.Singleca.ConfigPath)
	if err != nil {
		return conf, fmt.Errorf("cfssl configuration file %s Error: %s", conf.Singleca.ConfigPath, err)
	}

	conf.Singleca.CfsslConfig = cfg

	return conf, nil
}
