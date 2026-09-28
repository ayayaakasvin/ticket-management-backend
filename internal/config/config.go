package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	configPathEnvKey  = "CONFIG_PATH"
	postgresURLEnvKey = "POSTGRES_URL"
	valkeyURLEnvKey   = "VALKEY_URL"
	smtpPasswordKey   = "SMTP_PASSWORD"
)

// Config represents the configuration structure
type Config struct {
	HTTP     HTTPServerConfig      `yaml:"http_server"`
	CORS     CorsConfig            `yaml:"cors"`
	Database SQLiteConfig          `yaml:"sqlite"`
	LFS      LocalFileSystemConfig `yaml:"lfs"`
	// Cache RedisConfig
	SMTP   SMTPConfig   `yaml:"smtp"`
	Logger LoggerConfig `yaml:"logger"`
	JWT    JWTConfig    `yaml:"jwt"`
}

type CorsConfig struct {
	AllowedOrigins     []string `yaml:"allowed_origins"`
	AllowedMethods     []string `yaml:"allowed_methods"`
	AllowedHeaders     []string `yaml:"allowed_headers"`
	AllowedCredentials bool     `yaml:"allow_credentials"`
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration `yaml:"access_ttl" env-default:"24h"`
	RefreshTokenTTL time.Duration `yaml:"refresh_ttl" env-default:"168h"`
}

type HTTPServerConfig struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

type SQLiteConfig struct {
	FilePath string `yaml:"db_name"`
}

type PostgresConfig struct {
	URL string `yaml:"url" env-required:"true"`
}

type RedisConfig struct {
	URL string `yaml:"url" env-required:"true"`
}

type SMTPConfig struct {
	Username string `yaml:"username" env-required:"true"`
	Password string
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
}

type LoggerConfig struct {
	Env     string `yaml:"env" env-required:"true"`
	Service string `yaml:"service" env-required:"true"`
	// 	const (
	// 	LevelDebug Level = -4
	// 	LevelInfo  Level = 0
	// 	LevelWarn  Level = 4
	// 	LevelError Level = 8
	// )
	Level int  `yaml:"level" env-default:"0"`
	JSON  bool `yaml:"json" env-default:"false"`
}

type LocalFileSystemConfig struct {
	BasePath string `yaml:"path"`
}

// MustLoadConfig loads the configuration from the specified path
func MustLoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found, falling back to environment only")
	}

	configPath := os.Getenv(configPathEnvKey)
	if configPath == "" {
		log.Fatalf("%s is not set up", configPathEnvKey)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist: %s", configPath, err.Error())
	}

	log.Println(configPath)

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Println(cfg)
		log.Fatalf("failed to read config file: %s", err.Error())
	}

	err := loadFromEnv(&cfg)
	if err != nil {
		log.Fatalf("failed to read env values:%s", err.Error())
	}

	return &cfg
}

func loadFromEnv(cfg *Config) error {
	// postgresURL := os.Getenv(postgresURLEnvKey)
	// valkeyURL := os.Getenv(valkeyURLEnvKey)
	smtpPassword := os.Getenv(smtpPasswordKey)

	if smtpPassword == "" {
		return errors.New("SMTP password is not provided")
	}

	// if postgresURL == "" || valkeyURL == "" {
	// 	log.Fatalf("failed to read URLs")
	// }

	// cfg.Database.URL = postgresURL
	// cfg.Cache.URL = valkeyURL
	cfg.SMTP.Password = smtpPassword

	return nil
}
