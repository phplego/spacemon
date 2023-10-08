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
	DaemonBindAddr string `yaml:"daemon-bind-addr"`
	DaemonPort     int    `yaml:"daemon-port"`
	DaemonUsername string `yaml:"daemon-username"`
	DaemonPassword string `yaml:"daemon-password"`
}

func LoadConfig() Config {
	// Default config values
	cfg := Config{
		MaxHistorySize: 20,
		DaemonBindAddr: "0.0.0.0",
		DaemonPort:     18080,
		DaemonUsername: "user",
		DaemonPassword: "<default-password>",
	}
	err := cleanenv.ReadConfig(GetAppDir()+"/config.yml", &cfg)
	if err != nil {
		println("unable to load config.yml:", err.Error())
		os.Exit(1)
	}

	// If title is empty, use hostname.
	if cfg.Title == "" {
		cfg.Title, _ = os.Hostname()
	}

	return cfg
}

func GetAppDir() string {
	path, _ := os.Executable()
	path, _ = filepath.EvalSymlinks(path)
	return filepath.Dir(path)
}
