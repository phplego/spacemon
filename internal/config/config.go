package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"path/filepath"
)

type Config struct {
	// Directories is a slice of directory paths that will be scanned.
	Directories []string `yaml:"directories"`

	// Title is a string that displays at the top of the report, typically used to identify the host machine.
	// If left empty, the hostname will be used as the title.
	Title          string `yaml:"title"`
	MaxHistorySize int    `yaml:"max-history-size"`
	DaemonPort     int    `yaml:"daemon-port"`
}

func LoadConfig() Config {
	cfg := Config{}
	err := cleanenv.ReadConfig(GetAppDir()+"/config.yml", &cfg)
	if err != nil {
		println("unable to load config.yml:", err.Error())
		os.Exit(1)
	}
	return cfg
}

func GetAppDir() string {
	path, _ := os.Executable()
	path, _ = filepath.EvalSymlinks(path)
	return filepath.Dir(path)
}
