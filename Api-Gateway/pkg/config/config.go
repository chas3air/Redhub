package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Env              string    `yaml:"env" env-default:"local"`
	UsersStorageHost string    `yaml:"usersStorageHost" env-default:"user_service"`
	UsersStoragePort int       `yaml:"usersStoragePort" env-default:"50051"`
	AuthHost         string    `yaml:"authHost" env-default:"auth"`
	AuthPort         int       `yaml:"authPost" env-default:"50051"`
	API              APIConfig `yaml:"api"`
}

type APIConfig struct {
	Port    int           `yaml:"port" env-default:"50051"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	// Используем cleanenv для чтения конфигурации из YAML
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		// Если cleanenv не справляется, попробуем YAML
		file, err := os.ReadFile(configPath)
		if err != nil {
			panic("cannot read config file: " + err.Error())
		}
		if err := yaml.Unmarshal(file, &cfg); err != nil {
			panic("cannot unmarshal config: " + err.Error())
		}
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	// --config=./config/local.yaml
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
