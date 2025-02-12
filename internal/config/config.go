package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var _ Config = (*AppConfig)(nil)

// Config interface provides methods to retrieve configuration values.
type Config interface {
	Address() string
	DbConfig() PostgresConfig
	RDConfig() RedisConfig
	ShutdownTimeout() time.Duration
	SecretKeyByte() string
}

// AppConfig holds all configuration settings for the application.
type AppConfig struct {
	Server         ServerConfig   `json:"Server" validate:"required"`
	Postgres       PostgresConfig `json:"Postgres" validate:"required"`
	ShutdownTime   time.Duration  `json:"ShutdownTime"`
	Redis          RedisConfig    `json:"Redis" validate:"required"`
	SecretKeyValue string         `json:"SecretKeyByte" validate:"required"`
}

func (c *AppConfig) SecretKeyByte() string {
	return c.SecretKeyValue
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	HTTPPort string `json:"HTTPPort" validate:"required,numeric,min=1,max=65535"`
}

// PostgresConfig holds PostgreSQL database configuration.
type PostgresConfig struct {
	Host     string `json:"Host" validate:"required"`
	Port     string `json:"Port" validate:"required,numeric,min=1,max=65535"`
	User     string `json:"User" validate:"required"`
	Password string `json:"Password" validate:"required"`
	DBName   string `json:"DBName" validate:"required"`
}

type RedisConfig struct {
	Host     string `json:"Host" validate:"required"`
	Port     string `json:"Port" validate:"required,numeric,min=1,max=65535"`
	Password string `json:"Password" validate:"required,min=1,max=65535"`
}

// LoadEnv loads environment variables from a file.
func LoadEnv(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	return nil
}

// NewConfig creates a new AppConfig instance by reading environment variables and validates the configuration.
func NewConfig() (Config, error) {
	appConfig := &AppConfig{}

	// Load server configuration.
	appConfig.Server.HTTPPort = os.Getenv("HTTP_PORT")

	// Load Postgres configuration.
	appConfig.Postgres = PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
	appConfig.Redis = RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Port:     os.Getenv("REDIS_PORT"),
	}

	appConfig.SecretKeyValue = os.Getenv("SECRET_KEY")

	// Load shutdown timeout.
	shutdownTime, err := strconv.Atoi(os.Getenv("SHUTDOWN_TIMEOUT"))
	if err != nil || shutdownTime < 0 {
		shutdownTime = 15 // Default value in seconds.
	}
	appConfig.ShutdownTime = time.Duration(shutdownTime) * time.Second

	// Validate the loaded configuration.
	if err := validateConfig(appConfig); err != nil {
		return nil, err
	}

	return appConfig, nil
}

// validateConfig validates the configuration using go-playground/validator.
func validateConfig(cfg *AppConfig) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}

// Address returns the server address in the format "host:port".
func (c *AppConfig) Address() string {
	return net.JoinHostPort("", c.Server.HTTPPort)
}

// DbConfig returns the PostgreSQL configuration.
func (c *AppConfig) DbConfig() PostgresConfig {
	return c.Postgres
}

func (c *AppConfig) RDConfig() RedisConfig {
	return c.Redis
}

// ShutdownTimeout returns the application shutdown timeout.
func (c *AppConfig) ShutdownTimeout() time.Duration {
	return c.ShutdownTime
}

func (p *PostgresConfig) GenerateDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.DBName)
}

func (r *RedisConfig) GenerateDSN() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port) // Redis DSN format
}
