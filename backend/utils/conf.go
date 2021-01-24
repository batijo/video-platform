package utils

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	JWTExp     int
	JWTSecret  string
	SD         string
	TD         string
	DD         string
	TND        string
	APTGD      string
	Debug      bool
	DebugStart string
	DebugEnd   string
	TempJson   string
	TempTxt    string
	VBW        int
	DataGen    string
	LogP       string
	TNTS       string
	Presets    bool
	FileTypes  []string
}

// Conf configuration file
var Conf Config

// GetConf load and return config file
func GetConf() (Config, error) {
	var conf Config
	if _, err := toml.DecodeFile("./utils/conf.toml", &conf); err != nil {
		log.Println("error geting conf.toml")
		return conf, err
	}
	return conf, nil
}
