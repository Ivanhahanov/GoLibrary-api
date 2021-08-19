package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Elastic struct {
		Address string `yaml:"address"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"elastic"`
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func (cfg *Config) LoadConfig(configName string) {
	f, err := os.Open(configName)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		processError(err)
	}

}
