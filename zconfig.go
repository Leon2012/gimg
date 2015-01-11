package gimg

import (
	"code.google.com/p/gcfg"
	_ "errors"
	_ "strings"
)

type AppConfig struct {
	System struct {
		IsDaemon    int
		Host        string
		Port        int
		Headers     string
		Etag        int
		LogOutput   string
		LogLevel    int
		LogName     string
		DisableArgs int
		Format      string
		Quality     int
	}

	Cache struct {
		Cache        int
		MemcacheHost string
		MemcachePort int
	}

	Storage struct {
		Mode         int
		SaveNew      int
		MaxSize      int
		AllowedTypes string
		ImgPath      string

		BeansdbHost string
		BeansdbPort int

		SsdbHost string
		SsdbPort int
	}
}

// func (a *AppConfig) Types() []string {
//  return strings.Split(a.Storage.AllowedTypes, ",")
// }

func LoadConfig(cfgFile string) (AppConfig, error) {
	var err error
	var cfg AppConfig

	// isExists := util.IsExist(cfgFile)
	// if !isExists {
	//  return cfg, errors.New("Config file is not exists!")
	// }

	err = gcfg.ReadFileInto(&cfg, cfgFile)
	if err != nil {
		//panic(err)
		return cfg, err
	}

	return cfg, nil
}
