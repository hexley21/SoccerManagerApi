package config

import (
	"math"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		IsProd     bool
		Server     Server     `yaml:"server"`
		HTTP       HTTP       `yaml:"http"`
		Postgres   Postgres   `yaml:"postgres"`
		Metrics    Metrics    `yaml:"metrics"`
		JWT        JWT        `yaml:"jwt"`
		Pagination Pagination `yaml:"pagination"`
		Globe      Globe      `yaml:"globe"`
		Argon2     Argon2     `yaml:"argon2"`
		Logging    Logging    `yaml:"logging"`
	}

	Server struct {
		ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
		InstanceId      int64         `yaml:"instance_id"`
		MaxFileSize     int64         `yaml:"max_file_size"`
	}

	HTTP struct {
		Port         int           `yaml:"port"`
		CorsOrigins  string        `yaml:"cors_origins"`
		IdleTimeout  time.Duration `yaml:"idle_timeout"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
	}

	Globe struct {
		TTL time.Duration `yaml:"ttl"`
	}

	Postgres struct {
		Port              int    `yaml:"port"`
		Host              string `yaml:"host"`
		DBName            string `yaml:"db_name"`
		User              string
		Password          string
		SslMode           string
		MaxConns          int32         `yaml:"max-connections"`
		MinConns          int32         `yaml:"min-connections"`
		HealthCheckPeriod time.Duration `yaml:"healthcheck-period"`
		MaxConnLifetime   time.Duration `yaml:"max-conn-lifetime"`
		MaxConnIdleTime   time.Duration `yaml:"max-conn-idle-time"`
		Version           uint          `yaml:"version"`
	}

	Argon2 struct {
		SaltLen    uint32 `yaml:"salt_len"`
		KeyLen     uint32 `yaml:"key_len"`
		Time       uint32 `yaml:"time"`
		Memory     uint32 `yaml:"memory"`
		Threads    uint8  `yaml:"threads"`
		Breakpoint int    // indicates where the salt starts
	}

	Metrics struct {
		Port int `yaml:"port"`
	}

	JWT struct {
		Access  TokenParams `yaml:"access"`
		Refresh TokenParams `yaml:"refresh"`
	}

	TokenParams struct {
		Secret string
		TTL    time.Duration `yaml:"ttl"`
	}

	Pagination struct {
		S   int32 `yaml:"s"`
		M   int32 `yaml:"m"`
		L   int32 `yaml:"l"`
		XL  int32 `yaml:"xl"`
		XXL int32 `yaml:"2xl"`
	}

	Logging struct {
		DevLogLevel   string `yaml:"dev_level"`
		ProdLogLevel  string `yaml:"prod_level"`
		CallerEnabled bool   `yaml:"caller_enabled"`
	}
)

func LoadConfig(
	generalConfigDir string,
	serviceConfigDir string,
	envPath string,
) (*Config, error) {
	cfg := new(Config)

	if serviceConfigDir != "" {
		if err := cfg.parseYaml(serviceConfigDir); err != nil {
			return cfg, err
		}
	}

	if err := cfg.parseYaml(generalConfigDir); err != nil {
		return cfg, err
	}

	if err := cfg.parseEnv(envPath); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (cfg *Config) parseYaml(configDir string) error {
	yamlFile, err := os.ReadFile(configDir)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(yamlFile, cfg); err != nil {
		return err
	}

	if (cfg.Argon2 != Argon2{}) {
		cfg.Argon2.Breakpoint = Argon2KeylenBreakpoint(cfg.Argon2.KeyLen)
	}

	return nil
}

func (cfg *Config) parseEnv(envPath string) error {
	err := godotenv.Load(envPath)
	if err != nil {
		return err
	}

	cfg.IsProd, err = strconv.ParseBool(os.Getenv("IS_PROD"))
	if err != nil {
		return err
	}

	cfg.JWT.Access.Secret = os.Getenv("JWT_ACCESS_SECRET")
	cfg.JWT.Access.Secret = os.Getenv("JWT_REFRESH_SECRET")

	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.SslMode = os.Getenv("POSTGRES_SSL_MODE")

	return nil
}

func Argon2KeylenBreakpoint(keylen uint32) int {
	return int(math.Ceil(float64(keylen) * 4.0 / 3.0))
}
