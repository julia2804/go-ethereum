package ethereum

import (
	"github.com/Unknwon/goconfig"
	"github.com/ethereum/go-ethereum/log"
)

var cfg *goconfig.ConfigFile

func GetConfig(path string) *goconfig.ConfigFile {
	if cfg != nil {
		return cfg
	} else {
		cfg, _ = goconfig.LoadConfigFile(path)
		return cfg
	}
}

func GetConfigDefault() *goconfig.ConfigFile {
	return cfg
}

func GetValue(path string, section string, key string) string {
	cfg := GetConfig(path)
	if cfg == nil {
		log.Info("no cfg")
		return ""
	}
	value, err := cfg.GetValue(section, key)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return value
}

func GetValueFromDefaultPath(section string, key string) string {
	path1 := "/root/config.ini"
	path2 := "/home/julia/config.ini"
	path3 := "/home/mimota/config.ini"

	if GetConfig(path1) == nil && GetConfig(path2) == nil && GetConfig(path3) == nil {
		log.Info("no cfg in these path", path1, path2, path3)
		return ""
	}
	cfg := GetConfigDefault()
	value, err := cfg.GetValue(section, key)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	return value
}
