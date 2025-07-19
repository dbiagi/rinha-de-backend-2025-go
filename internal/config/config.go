package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	AppName        = "rinha2025"
	DevelopmentEnv = "dev"
	ProductionEnv  = "prod"
)

type Configuration struct {
	WebConfig
	AppConfig
	DatabaseConfig
	ProcessorConfig
}

type AppConfig struct {
	Name        string
	Version     string
	Environment string
}

type WebConfig struct {
	Port                     int
	IdleTimeout              time.Duration
	ReadTimeout              time.Duration
	WriteTimeout             time.Duration
	ShutdownTimeout          time.Duration
	GracefulShutdownDisabled bool
}

type DatabaseConfig struct {
	Host               string
	Port               int
	User               string
	Password           string
	DatabaseName       string
	MaxIdleConnections int
	MaxOpenConnections int
}

type ProcessorConfig struct {
	DefaultHost            string
	FallbackHost           string
	MaxAllowedResponseTime int
}

func LoadConfig(env string) Configuration {
	filePath := os.Getenv("CONFIG_FILE_PATH")

	loadFromFile(filePath)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	mrt, _ := strconv.Atoi(os.Getenv("PROCESSOR_MAX_RESPONSE_TIME"))
	idleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	openConn, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTIONS"))

	return Configuration{
		WebConfig: WebConfig{
			Port:            port,
			IdleTimeout:     time.Second * 10,
			ReadTimeout:     time.Second * 10,
			WriteTimeout:    time.Second * 10,
			ShutdownTimeout: time.Second * 20,
		},
		AppConfig: AppConfig{
			Name:        AppName,
			Version:     "1.0.0",
			Environment: env,
		},
		DatabaseConfig: DatabaseConfig{
			Host:               os.Getenv("DB_HOST"),
			Port:               dbPort,
			User:               os.Getenv("DB_USER"),
			Password:           os.Getenv("DB_PASSWORD"),
			DatabaseName:       os.Getenv("DB_NAME"),
			MaxIdleConnections: idleConn,
			MaxOpenConnections: openConn,
		},
		ProcessorConfig: ProcessorConfig{
			DefaultHost:            os.Getenv("PROCESSOR_DEFAULT_HOST"),
			FallbackHost:           os.Getenv("PROCESSOR_FALLBACK_HOST"),
			MaxAllowedResponseTime: mrt,
		},
	}
}

func loadFromFile(configFilePath string) {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		slog.Warn("Config file does not exist.")
		return
	}

	err := godotenv.Load(configFilePath)

	if err != nil {
		slog.Error("Error loading .env file.", slog.String("error", err.Error()))
	}
}
