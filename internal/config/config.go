package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Directories []string `yaml:"Directories"`
}

func LoadConfig() Config {
	cfg := Config{}
	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		panic("unable to load config.yml")
	}
	log.Println(cfg)
	return cfg
}
