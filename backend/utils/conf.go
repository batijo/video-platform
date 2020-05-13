package utils

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	SD         string
	TD         string
	DD         string
	TND        string
	Debug      bool
	DebugStart string
	DebugEnd   string
	TempJson   string
	TempTxt    string
	VBW        int
	ABW        int
	DataGen    string
	LogP       string
	TNTS       string
	TNNum      int
	Presets    bool
	FileTypes  []string
}

// Load config file
func GetConf() (Config, error) {
	var conf Config
	if _, err := toml.DecodeFile("utils/conf.toml", &conf); err != nil {
		log.Println("error geting conf.toml")
		return conf, err
	}
	return conf, nil
}
