package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Directories []string `yaml:"directories"`
	Title       string   `yaml:"title"`
}

func LoadConfig() Config {
	cfg := Config{}
	err := cleanenv.ReadConfig(GetAppDir()+"/config.yml", &cfg)
	if err != nil {
		panic("unable to load config.yml")
	}
	log.Println(cfg)
	return cfg
}

func GetAppDir() string {
	path, _ := os.Executable()
	path, _ = filepath.EvalSymlinks(path)
	return filepath.Dir(path)
}
