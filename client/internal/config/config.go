package config

import (
	"fmt"
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

// MustLoad load config
func MustLoad(fpaths ...string) {
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
}

func ParseConfigByEnv() error {
	if v := os.Getenv("LOCAL_ADDR"); v != "" {
		C.App.LocalAddr = v
	}
	if v := os.Getenv("LOCAL_PORT"); v != "" {
		p, _ := strconv.Atoi(v)
		C.App.LocalPort = p
	}
	if v := os.Getenv("CONTRO_HOST"); v != "" {
		C.App.ControlHost = v
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
	Token string
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
