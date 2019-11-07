package ethereum

import (
	"fmt"
	"github.com/Unknwon/goconfig"
)

var cfg *goconfig.ConfigFile

func GetConfig(path string) *goconfig.ConfigFile {
	if cfg != nil {
		return cfg
	} else {
		var err error
		cfg, err = goconfig.LoadConfigFile(path)
		if err != nil {
			fmt.Println("error in getConfig")
			return nil
		}
		return cfg
	}
}

func GetConfigDefault() *goconfig.ConfigFile {
	return cfg
}

func GetValue(path string, section string, key string) string {
	cfg := GetConfig(path)
	if cfg == nil {
		fmt.Println("no cfg")
		return ""
	}
	value, err := cfg.GetValue(section, key)
	if err != nil {
		fmt.Println("error in getValue")
		return ""
	}
	return value
}

func GetValueFromDefaultPath(section string, key string) string {
	path1 := "/home/mimota/go/src/github.com/ethereum/go-ethereum/config.ini"
	path2 := "/root/go/src/github.com/ethereum/go-ethereum/config.ini"

	if GetConfig(path1) == nil && GetConfig(path2) == nil {
		fmt.Println("no cfg in two path")
		return ""
	}
	cfg := GetConfigDefault()
	value, err := cfg.GetValue(section, key)
	if err != nil {
		fmt.Println("error in GetValueFromDefaultPath")
		return ""
	}
	return value
}
