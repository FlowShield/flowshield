package confer

import "github.com/flowshield/flowshield/verifier/pkg/util"

func ConfigAppGetString(key string, defaultConfig string) string {
	config := ConfigAppGet(key)
	if config == nil {
		return defaultConfig
	} else {
		configStr := config.(string)
		if util.UtilIsEmpty(configStr) {
			configStr = defaultConfig
		}
		return configStr
	}
}

func ConfigAppGetInt(key string, defaultConfig int) int {
	config := ConfigAppGet(key)
	if config == nil {
		return defaultConfig
	} else {
		configInt := config.(int)
		if configInt == 0 {
			configInt = defaultConfig
		}
		return configInt
	}
}

func ConfigAppGet(key string) interface{} {
	globalConfig.RLock()
	defer globalConfig.RUnlock()
	//将配置文件中的app字段转为map
	config, exists := globalConfig.App[key]
	if !exists {
		return nil
	}
	return config
}

func ConfigEnvGet() string {
	strEnv := ConfigAppGet("env")
	return strEnv.(string)
}

func ConfigEnvIsDev() bool {
	env := ConfigEnvGet()
	if env == "dev" || env == "debug" {
		return true
	}
	return false
}

func ConfigEnvIsBeta() bool {
	env := ConfigEnvGet()
	if env == "beta" {
		return true
	}
	return false
}
