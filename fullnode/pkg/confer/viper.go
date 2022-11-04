package confer

import (
	"fmt"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

var globalConfig *ServerConfig
var mutex sync.RWMutex

func Init(configURL string) (err error) {
	v := viper.New()
	v.SetConfigFile(configURL)
	err = v.ReadInConfig()
	if err != nil {
		err = fmt.Errorf("fatal error config file: %w", err)
		return
	}
	if err = v.Unmarshal(&globalConfig); err != nil {
		return
	}
	return handleConfig(globalConfig)
}

func handleConfig(config *ServerConfig) error {
	config.replaceByEnv(&config.Redis.Addr)
	config.replaceByEnv(&config.Mysql.Write.Host)
	config.replaceByEnv(&config.Mysql.Write.User)
	config.replaceByEnv(&config.Mysql.Write.Password)
	config.replaceByEnv(&config.CA.SignURL)
	config.replaceByEnv(&config.CA.AuthKey)
	config.replaceByEnv(&config.Oauth2.ClientID)
	config.replaceByEnv(&config.Oauth2.ClientSecret)
	config.replaceByEnv(&config.Web3.Register)
	config.replaceByEnv(&config.Web3.PrivateKey)
	config.replaceByEnv(&config.P2P.ServiceDiscoveryID)
	config.replaceByEnv(&config.P2P.ServiceMetadataTopic)
	privateKey, err := crypto.HexToECDSA(config.Web3.PrivateKey)
	if err != nil {
		return err
	}
	// Account 处理
	config.P2P.Account = crypto.PubkeyToAddress(privateKey.PublicKey).String()
	config.replaceByEnv(&config.Web3.Contract.Token)
	config.replaceByEnv(&config.Web3.W3S.Token)
	config.replaceByEnv(&config.Web3.ETH.URL)
	config.Mysql.Write.DBName = globalConfig.Mysql.DBName
	config.Mysql.Write.Prefix = globalConfig.Mysql.Prefix
	return nil
}

func (*ServerConfig) replaceByEnv(conf *string) {
	if s := os.Getenv(*conf); len(s) > 0 {
		*conf = s
	}
}

func GlobalConfig() *ServerConfig {
	mutex.RLock()
	defer mutex.RUnlock()
	return globalConfig
}
