package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envTest  = "test"
)

type Config struct {
	Env  string     `yaml:"env" env:"APP_ENV" env-required:"true"`
	HTTP HTTPConfig `yaml:"http" env-required:"true"`
}

type HTTPConfig struct {
	Port           int           `yaml:"port" env:"HTTP_PORT" env-required:"true"`
	ReadTimeout    time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-required:"true"`
	WriteTimeout   time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-required:"true"`
	GatewayTimeout time.Duration `yaml:"gateway_timeout" env:"HTTP_GATEWAY_TIMEOUT" env-required:"true"`
}

// MustLoad загружает конфигурацию из файла и переменных окружения.
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is not set")
	}

	return MustLoadByPath(path)
}

// MustLoadByPath загружает конфигурацию из указанного файла.
// Если файл не существует или нет прав доступа, вызывает панику.
func MustLoadByPath(configPath string) *Config {
	_, err := os.Stat(configPath)
	if err != nil && os.IsPermission(err) {
		panic("no permission to config file: " + configPath)
	}
	if err != nil && os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to load config: " + err.Error())
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to load environment variables: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
