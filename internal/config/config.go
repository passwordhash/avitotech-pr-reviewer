package config

import (
	"flag"
	"fmt"
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
	App  AppConfig  `yaml:"app" env-required:"true"`
	HTTP HTTPConfig `yaml:"http" env-required:"true"`
	PG   PGConfig   `yaml:"postgres" env-required:"true"`
}

type AppConfig struct {
	Env               string `yaml:"env" env:"APP_ENV" env-required:"true"`
	AdminToken        string `yaml:"admin_token" env:"ADMIN_TOKEN" env-required:"true"`
	MaxReviewersPerPR int    `yaml:"max_reviewers_per_pr" env:"REVIEWERS_PER_PR" env-required:"true"`
}

type HTTPConfig struct {
	Port           int           `yaml:"port" env:"HTTP_PORT" env-required:"true"`
	ReadTimeout    time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-required:"true"`
	WriteTimeout   time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-required:"true"`
	GatewayTimeout time.Duration `yaml:"gateway_timeout" env:"HTTP_GATEWAY_TIMEOUT" env-required:"true"`
}

type PGConfig struct {
	Host     string `env:"POSTGRES_HOST" yaml:"host" env-required:"true"`
	Port     int    `env:"POSTGRES_PORT" yaml:"port" env-required:"true"`
	Username string `env:"POSTGRES_USER" yaml:"user" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" yaml:"password" env-required:"true"`
	Database string `env:"POSTGRES_DB" yaml:"database" env-required:"true"`
	SSLMode  string `env:"POSTGRES_SSLMODE" yaml:"sslmode" env-default:"disable"`

	MaxConns int32 `env:"POSTGRES_MAX_CONNS" yaml:"max_conns"`
}

func (p PGConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		p.Username,
		p.Password,
		p.Host,
		p.Port,
		p.Database,
		p.SSLMode,
	)
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

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
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
