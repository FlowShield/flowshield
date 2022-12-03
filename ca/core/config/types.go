package config

import (
	cfssl_config "github.com/flowshield/cfssl/config"
)

const (
	MetricsTablePrefix = "ca_"
)

type IConfig struct {
	Log            Log        `yaml:"log"`
	Keymanager     Keymanager `yaml:"keymanager"`
	Singleca       Singleca   `yaml:"singleca"`
	OCSPHost       string     `yaml:"ocsp-host"`
	HTTP           HTTP       `yaml:"http"`
	SwaggerEnabled bool       `yaml:"swagger-enabled"`
	Debug          bool       `yaml:"debug"`
	Version        string     `yaml:"version"`
	Hostname       string     `yaml:"hostname"`
	Ocsp           Ocsp       `yaml:"ocsp"`
}

type Registry struct {
	SelfName string `yaml:"self-name"`
	Command  string `yaml:"command"`
}

// ocsp
type Ocsp struct {
	CacheTime int `yaml:"cache-time"`
}

type LogProxy struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Key  string `yaml:"key"`
}
type Log struct {
	LogProxy LogProxy `yaml:"log-proxy"`
}
type Singleca struct {
	ConfigPath string `yaml:"config-path"`

	// Raw
	CfsslConfig *cfssl_config.Config
}
type HTTP struct {
	OcspListen string `yaml:"ocsp-listen"`
	CaListen   string `yaml:"ca-listen"`
	Listen     string `yaml:"listen"`
}
type Mysql struct {
	Dsn string `yaml:"dsn"`
}
type RootCa struct {
	O      string `yaml:"o"`
	Expiry string `yaml:"expiry"`
}
type IntermediateCa struct {
	O      string `yaml:"o"`
	Ou     string `yaml:"ou"`
	Expiry string `yaml:"expiry"`
}
type CsrTemplates struct {
	RootCa         RootCa         `yaml:"root-ca"`
	IntermediateCa IntermediateCa `yaml:"intermediate-ca"`
}
type Keymanager struct {
	UpperCa      []string     `yaml:"upper-ca"`
	SelfSign     bool         `yaml:"self-sign"`
	CsrTemplates CsrTemplates `yaml:"csr-templates"`
}
type Vault struct {
	Enabled bool   `yaml:"enabled"`
	Addr    string `yaml:"addr"`
	Token   string `yaml:"token"`
	Prefix  string `yaml:"prefix"`
}
