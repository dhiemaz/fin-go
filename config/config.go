package config

import (
	"github.com/dhiemaz/fin-go/infrastructure/database/postgres"
	"github.com/dhiemaz/fin-go/infrastructure/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

const (
	APP_PREFIX = "fin_go"
)

type Config struct {
	Port                  string `envconfig:"PORT"`
	DBHost                string `envconfig:"DB_HOST"`
	DBUsername            string `envconfig:"DB_USERNAME"`
	DBPort                string `envconfig:"DB_PORT"`
	DBPassword            string `envconfig:"DB_PASSWORD"`
	DBName                string `envconfig:"DB_NAME"`
	DBMaxConn             int    `envconfig:"DB_MAX_CONN"`
	DBMaxIdle             int    `envconfig:"DB_MAX_IDLE"`
	HTTPMaxConnPerIP      int    `envconfig:"HTTP_MAX_CONN_PER_IP"`
	HTTPMaxRequestPerConn int    `envconfig:"HTTP_MAX_REQUEST_PER_CONN"`
	HTTPMaxConcurrency    int    `envconfig:"HTTP_MAX_CONCURRENCY"`
	HTTPMaxKeepAlive      int    `envconfig:"HTTP_MAX_KEEP_ALIVE_DURATION"`
	JWT                   string `envconfig:"JWT_SECRET"`
	DBPool                *pgxpool.Pool
}

var cfg Config

// Initiate configuration
func InitConfig() {
	err := LoadConfigs()
	if err != nil {
		log.Fatalf("failed load config, error : %v", err)
		os.Exit(0)
	}

	InitLogger() // initialize logger instance

	cfg.DBPool, err = postgres.InitDBConnection()
	if err != nil {
		log.Fatalf("failed connect to database, error : %v", err)
		os.Exit(0)
	}
}

// Loads general configs
func LoadConfigs() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	err = envconfig.Process(APP_PREFIX, &cfg)
	return err
}

// InitLogger : initialize logger instance
func InitLogger() {
	logConfig := logger.Configuration{
		EnableConsole:     true,    // next, get from configuration
		ConsoleJSONFormat: true,    // next, get from configuration
		ConsoleLevel:      "debug", // next, get from configuration
	}

	if err := logger.NewLogger(logConfig, logger.InstanceZapLogger); err != nil {
		log.Fatalf("Could not instantiate log %v", err)
	}
}

// GetConfig : get configuration stored
func GetConfig() *Config {
	if &cfg == nil {
		cfg := new(Config)
		return cfg
	}

	return &cfg
}
