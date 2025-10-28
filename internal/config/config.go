package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	configPathEnvKey 	= "CONFIG_PATH"
	postgresURLEnvKey 	= "POSTGRES_URL"
	valkeyURLEnvKey 	= "VALKEY_URL"
	smtpPasswordKey 	= "SMTP_PASSWORD"
)

// Config represents the configuration structure
type Config struct {
					HTTPServer 		`yaml:"http_server"															env-required:"true"`
	Database   		StorageConfig
	Valkey     		RedisConfig
	SMTP			SMTPConfig		`yaml:"smtp"																env-required:"true"`
}

type HTTPServer struct {
	Address     	string        	`yaml:"address"					env-default:"localhost:8080`
	Timeout     	time.Duration 	`yaml:"timeout" 															env-required:"true"`
	IdleTimeout 	time.Duration 	`yaml:"iddle_timeout" 														env-required:"true"`
	TLS         	TLSConfig     	`yaml:"tls"																	env-required:"true"`
}

type StorageConfig struct {
	URL 			string			`yaml:url																	env-required:"true"`
}

type TLSConfig struct {
	CertFile 		string 			`yaml:"certfile"`
	KeyFile  		string 			`yaml:"keyfile"	`
}

type RedisConfig struct {
	URL				string			`yaml:"url"																	env-required:"true"`
}

type SMTPConfig struct {
	Username 		string			`yaml:"username"															env-required:"true"`
	Password		string
	Host			string			`yaml:"host"																env-required:"true"`
	Port			int				`yaml:"port"																env-required:"true"`
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

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config file: %s", err.Error())
	}

	postgresURL := os.Getenv(postgresURLEnvKey)
	valkeyURL := os.Getenv(valkeyURLEnvKey)
	smtpPassword := os.Getenv(smtpPasswordKey)

	if postgresURL == "" || valkeyURL == "" {
		log.Fatalf("failed to read URLs")
	}

	cfg.Database.URL = postgresURL
	cfg.Valkey.URL = valkeyURL
	cfg.SMTP.Password = smtpPassword

	return &cfg
}
