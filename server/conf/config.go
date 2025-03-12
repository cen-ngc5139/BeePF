package conf

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/toolkits/file"
)

type HttpConfig struct {
	Enabled  bool   `json:"enabled"`
	Listen   string `json:"listen"`
	Backdoor bool   `json:"backdoor"`
}

type GlobalConfig struct {
	Http     *HttpConfig    `json:"http"`
	Database *Database      `json:"database"`
	LogMode  string         `json:"logMode"`
	Env      string         `json:"env"`
	Metrics  *MetricsConfig `json:"metrics"`
}

type MetricsConfig struct {
	PrometheusHost string `json:"prometheusHost"`
}

type Database struct {
	Type     string `json:"type"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Name     string `json:"name"`
	MaxIdle  int    `json:"maxIdle"`
	MaxOpen  int    `json:"maxOpen"`
	LogMode  string `json:"logMode"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string, reload bool) {
	if cfg == "" {
		if reload {
			logrus.Error("configuration file is nil")
			return
		}
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		if reload {
			logrus.Error("config file:", cfg, "is not existent")
			return
		}
		log.Fatalln("config file:", cfg, "is not existent")
	}
	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		if reload {
			logrus.Error("read config file:", cfg, "fail:", err)
			return
		}
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		if reload {
			logrus.Error("parse config file:", cfg, "fail:", err)
			return
		}
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	config = &c

	if !reload {
		log.Println("read config file:", cfg, "successfully")
	}
}
