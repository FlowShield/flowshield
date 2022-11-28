package config

import (
	"bufio"
	"fmt"
	"github.com/cloudslit/cloudslit/client/configs"
	"github.com/cloudslit/cloudslit/client/pkg/util"
	"github.com/cloudslit/cloudslit/client/pkg/util/json"
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
	Is   = new(I)
)

// I ...
type I struct {
	HttpClient *http.Client
}

// 读取内置配置文件并创建
func createConfigFile(path string) error {
	ok, err := util.PathExists(path)
	if err != nil {
		return err
	}
	if !ok {
		f, err := os.Create(path) //创建文件
		if err != nil {
			return err
		}
		defer f.Close()
		w := bufio.NewWriter(f)
		_, err = w.WriteString(string(configs.ConfigFileData))
		if err != nil {
			return err
		}
		err = w.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

// MustLoad load config
func MustLoad(fpaths ...string) error {
	if len(fpaths) <= 0 || fpaths[0] == "" {
		fpaths[0] = "config.toml"
		// 无配置文件，读取默认配置文件
		err := createConfigFile(fpaths[0])
		if err != nil {
			return err
		}
	}
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
	if v := os.Getenv("CLI_APP_LOCAL_ADDR"); v != "" {
		C.App.LocalAddr = v
	}
	if v := os.Getenv("CLI_APP_LOCAL_PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("environment variable [%s] parsing error:%v", "CLI_APP_LOCAL_PORT", err)
		}
		C.App.LocalPort = p
	}
	if v := os.Getenv("CLI_APP_CONTROL_HOST"); v != "" {
		C.App.ControlHost = v
	}

	// w3s
	if v := os.Getenv("CLI_WEB3_W3S_TOKEN"); v != "" {
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

	// Log
	if v := os.Getenv("CLI_LOG_HOOK_ENABLED"); v == "true" {
		C.Log.EnableHook = true
	}
	if v := os.Getenv("CLI_LOG_REDIS_ADDR"); v != "" {
		C.LogRedisHook.Addr = v
	}
	if v := os.Getenv("CLI_LOG_REDIS_KEY"); v != "" {
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
	Machine      Machine
	Log          Log
	LogRedisHook LogRedisHook
	Web3         Web3
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

// Common Configuration parameters
type App struct {
	LocalAddr   string
	LocalPort   int
	ControlHost string
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

// Machine
type Machine struct {
	MachineId string
	Cookie    string
}

const lockPath = "./machine.lock"

func (a *Machine) SetMachineId(macid string) {
	C.Machine.MachineId = macid
	a.MachineId = macid
}

func (a *Machine) SetCookie(cookie string) {
	C.Machine.Cookie = cookie
	a.Cookie = cookie
}

func (a *Machine) Write() error {
	b, err := json.Marshal(a)
	if err != nil {
		return err
	}
	// Write lock to filesystem to indicate an existing running daemon.
	err = os.WriteFile(lockPath, b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (a *Machine) Read() (*Machine, error) {
	in, err := os.ReadFile(lockPath)
	if err != nil {
		return nil, err
	}
	var result Machine
	err = json.Unmarshal(in, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
