package config

import (
	"github.com/restore/user/repository"
	"github.com/restore/user/service"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const EmailHeader = "X-Consumer-Username"

type Configuration struct {
	Kong  service.KongConfig `yaml:"kong"`
	Mysql repository.Config  `yaml:"mysql"`
}

var config Configuration

func Init() {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
}

func NewKongConfig() *service.KongConfig {
	return &config.Kong
}

func NewDBConfig() *repository.Config {
	return &config.Mysql
}
