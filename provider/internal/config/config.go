package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/cloudslit/cloudslit/provider/pkg/util/json"
	"github.com/ethereum/go-ethereum/crypto"

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
	return ParseConfigByEnv()
}

func ParseConfigByEnv() error {
	// APP
	if v := os.Getenv("PR_APP_LOCAL_ADDR"); v != "" {
		C.App.LocalAddr = v
	}
	if v := os.Getenv("PR_APP_LOCAL_PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("environment variable [%s] parsing error:%v", "PR_APP_LOCAL_PORT", err)
		}
		C.App.LocalPort = p
	}

	// Web3
	if v := os.Getenv("PR_WEB3_PRIVATE_KEY"); v != "" {
		C.Web3.PrivateKey = v
	}
	// Account 处理
	privateKey, err := crypto.HexToECDSA(C.Web3.PrivateKey)
	if err != nil {
		return fmt.Errorf("wallet account private key resolution failed:%v", err)
	}
	C.Web3.Account = crypto.PubkeyToAddress(privateKey.PublicKey).String()
	if v := os.Getenv("PR_WEB3_PRICE"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("environment variable [%s] parsing error:%v", "PR_WEB3_PRICE", err)
		}
		C.Web3.Price = p
	}
	if v := os.Getenv("PR_WEB3_CONTRACT_TOKEN"); v != "" {
		C.Web3.Contract.Token = v
	}
	if v := os.Getenv("PR_WEB3_W3S_TOKEN"); v != "" {
		C.Web3.W3S.Token = v
	}
	if v := os.Getenv("PR_WEB3_W3S_TIMEOUT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("environment variable [%s] parsing error:%v", "PR_WEB3_W3S_TIMEOUT", err)
		}
		C.Web3.W3S.Timeout = p
	}
	if v := os.Getenv("PR_WEB3_W3S_RETRY_COUNT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("environment variable [%s] parsing error:%v", "PR_WEB3_W3S_RETRY_COUNT", err)
		}
		C.Web3.W3S.RetryCount = p
	}
	if v := os.Getenv("PR_WEB3_ETH_URL"); v != "" {
		C.Web3.ETH.URL = v
	}
	if v := os.Getenv("PR_WEB3_ETH_PROJECT_ID"); v != "" {
		C.Web3.ETH.ProjectID = v
	}

	// Mysql
	if v := os.Getenv("PR_MYSQL_HOST"); v != "" {
		C.Mysql.Host = v
	}
	if v := os.Getenv("PR_MYSQL_PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("environment variable [%s] parsing error:%v", "PR_MYSQL_PORT", err)
		}
		C.Mysql.Port = p
	}
	if v := os.Getenv("PR_MYSQL_USER"); v != "" {
		C.Mysql.User = v
	}
	if v := os.Getenv("PR_MYSQL_PASSWORD"); v != "" {
		C.Mysql.Password = v
	}

	// Log
	if v := os.Getenv("PR_LOG_LEVEL"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("environment variable [%s] parsing error:%v", "PR_LOG_LEVEL", err)
		}
		C.Log.Level = p
	}
	if v := os.Getenv("PR_LOG_HOOK_ENABLED"); v == "true" {
		C.Log.EnableHook = true
	}
	if v := os.Getenv("PR_LOG_REDIS_ADDR"); v != "" {
		C.LogRedisHook.Addr = v
	}
	if v := os.Getenv("PR_LOG_REDIS_KEY"); v != "" {
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
	RunMode      string
	PrintConfig  bool
	App          App
	P2p          P2p
	Web3         Web3
	Log          Log
	LogRedisHook LogRedisHook
	Mysql        Mysql
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
	LocalAddr      string
	LocalPort      int
	CertFile       string
	KeyFile        string
	HttpListenAddr string
}

type Mysql struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	Prefix          string
	Parameters      string
	PoolMinCap      int
	PoolExCap       int
	PoolMaxCap      int
	PoolIdleTimeout int
	PoolWaitCount   int
	PoolWaiTimeout  int
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
	Token      string
	Timeout    int
	RetryCount int
}

type ETH struct {
	URL       string
	ProjectID string
}

func (w *Web3) EthAddress() string {
	return fmt.Sprintf("%s/%s", w.ETH.URL, w.ETH.ProjectID)
}
