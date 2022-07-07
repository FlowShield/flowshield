package config

import (
	"fmt"
	"github.com/cloudslit/cloudslit/provider/pkg/util/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/koding/multiconfig"
)

var (
	// C Global configuration (Must Load first, otherwise the configuration will not be available)
	C    = new(Config)
	once sync.Once
)

// MustLoad load config
func MustLoad(fpaths ...string) error {
	once.Do(func() {
		loaders := []multiconfig.Loader{
			&multiconfig.TagLoader{},
			&multiconfig.EnvironmentLoader{},
		}

		for _, fpath := range fpaths {
			if strings.HasSuffix(fpath, "toml") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "json") {
				loaders = append(loaders, &multiconfig.JSONLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "yaml") {
				loaders = append(loaders, &multiconfig.YAMLLoader{Path: fpath})
			}
		}
		m := multiconfig.DefaultLoader{
			Loader:    multiconfig.MultiLoader(loaders...),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}
		m.MustLoad(C)
	})
	ParseConfigByEnv()
	return nil
}

func ParseConfigByEnv() error {
	if v := os.Getenv("LOCAL_ADDR"); v != "" {
		C.App.LocalAddr = v
	}
	if v := os.Getenv("LOCAL_PORT"); v != "" {
		p, _ := strconv.Atoi(v)
		C.App.LocalPort = p
	}
	if v := os.Getenv("LOG_HOOK_ENABLED"); v == "true" {
		C.Log.EnableHook = true
	}
	if v := os.Getenv("LOG_REDIS_ADDR"); v != "" {
		C.LogRedisHook.Addr = v
	}
	if v := os.Getenv("LOG_REDIS_KEY"); v != "" {
		C.LogRedisHook.Key = v
	}
	return nil
}

func PrintWithJSON() {
	if C.PrintConfig {
		b, err := json.MarshalIndent(C, "", " ")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
			return
		}
		os.Stdout.WriteString(string(b) + "\n")
	}
}

type Config struct {
	HttpClient   *http.Client
	RunMode      string
	PrintConfig  bool
	App          App
	P2p          P2p
	Web3         Web3
	Log          Log
	LogRedisHook LogRedisHook
	Certificate  Certificate
}

func (c *Config) IsDebugMode() bool {
	return c.RunMode == "debug"
}

func (c *Config) IsReleaseMode() bool {
	return c.RunMode == "release"
}

type LogHook string

func (h LogHook) IsRedis() bool {
	return h == "redis"
}

type Log struct {
	Level         int
	Format        string
	Output        string
	OutputFile    string
	EnableHook    bool
	HookLevels    []string
	Hook          LogHook
	HookMaxThread int
	HookMaxBuffer int
}

type LogGormHook struct {
	DBType       string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	Table        string
}

type LogRedisHook struct {
	Addr string
	Key  string
}

// App Configuration parameters
type App struct {
	LocalAddr  string
	LocalPort  int
	ControHost string
}

// Certificate certificate
type Certificate struct {
	CertPem string
	CaPem   string
	KeyPem  string

	CertPemPath string
	CaPemPath   string
	KeyPemPath  string
}

// P2p Configuration parameters
type P2p struct {
	Enable               bool
	ServiceDiscoveryID   string
	ServiceDiscoveryMode string
	ServiceMetadataTopic string
}

type Web3 struct {
	Account    string
	Price      int
	PrivateKey string
	Contract   Contract
	W3S        W3S
	ETH        ETH
}

type Contract struct {
	Token string
}

type W3S struct {
	Token string
}

type ETH struct {
	URL       string
	ProjectID string
}

func (w *Web3) EthAddress() string {
	return fmt.Sprintf("%s/%s", w.ETH.URL, w.ETH.ProjectID)
}
