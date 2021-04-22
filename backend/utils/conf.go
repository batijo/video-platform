package utils

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	JWTExp      int
	JWTRef      int
	JWTSecret   string
	SD          string
	TD          string
	DD          string
	TND         string
	APTGD       string
	Debug       bool
	ShowStructs bool
	DebugStart  string
	DebugEnd    string
	VBW         int
	LogP        string
	FileTypes   []string
}

// Conf configuration file
var Conf Config

// GetConf load and return config file
func getConf() (Config, error) {
	var conf Config
	if _, err := toml.DecodeFile("./utils/conf.toml", &conf); err != nil {
		log.Println("error geting conf.toml")
		return conf, err
	}
	return conf, nil
}

func (c *Config) Load() error {
	conf, err := getConf()
	if err != nil {
		return err
	}
	(*c) = conf
	return nil
}
