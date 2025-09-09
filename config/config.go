package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// ServiceConfig holds service-specific configuration
type ServiceConfig struct {
	Name                   string     `mapstructure:"name"`
	Env                    string     `mapstructure:"env"`
	GracefulShutdownSecond int        `mapstructure:"gracefulShutdownSecond"`
	HTTP                   HTTPConfig `mapstructure:"http"`
}

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Port string `mapstructure:"port"`
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DatabaseName string `mapstructure:"databaseName"`
	SSLMode      string `mapstructure:"sslMode"`
	TimeZone     string `mapstructure:"timeZone"`
}

// RedisConfig holds redis-specific configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig holds JWT-specific configuration
type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	ExpirationHours int    `mapstructure:"expirationHours"`
}

// RateLimitConfig holds rate limit configuration for OTP
type RateLimitConfig struct {
	Count int `mapstructure:"count"`
	Time  int `mapstructure:"time"` // in minutes
}

// OTPConfig holds OTP-specific configuration
type OTPConfig struct {
	Expiration int             `mapstructure:"expiration"` // in seconds
	Length     int             `mapstructure:"length"`
	RateLimit  RateLimitConfig `mapstructure:"rateLimit"`
}

// Config holds all configuration for the application
type Config struct {
	Service  ServiceConfig  `mapstructure:"service"`
	Postgres DatabaseConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	OTP      OTPConfig      `mapstructure:"otp"`
}

// ConfigSetup holds the configuration setup
type ConfigSetup struct {
	path   string
	config Config
}

// NewConfigSetup creates a new config setup
func NewConfigSetup(path string) *ConfigSetup {
	return &ConfigSetup{
		path: path,
	}
}

// SetUp reads and sets up the configuration
func (cs *ConfigSetup) SetUp() *Config {
	viper.SetConfigFile(cs.path)

	if err := viper.ReadInConfig(); err != nil {
		log.Panic("Error reading config file: ", err)
	}

	if err := viper.Unmarshal(&cs.config); err != nil {
		log.Panic("Error unmarshalling config: ", err)
	}

	return &cs.config
}

// LoadConfig loads configuration from the YAML file
func LoadConfig() *Config {
	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Panic("Failed to get current directory: ", err)
	}

	// Fall back to environment variables if config file not found
	viper.AutomaticEnv()

	// Set up default config path
	configPath := filepath.Join(dir, "config.local.yaml")

	// Check if config path provided as environment variable
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	// Try to load the config
	cs := NewConfigSetup(configPath)
	config := cs.SetUp()

	// Convert config values to the expected format
	return &Config{
		Service:  config.Service,
		Postgres: config.Postgres,
		Redis:    config.Redis,
		JWT:      config.JWT,
		OTP:      config.OTP,
	}
}

// GetOTPExpiration GetExpiration returns the OTP expiration as time.Duration
func (c *Config) GetOTPExpiration() time.Duration {
	return time.Duration(c.OTP.Expiration) * time.Second
}

// GetRateLimitDuration returns the rate limit duration as time.Duration
func (c *Config) GetRateLimitDuration() time.Duration {
	return time.Duration(c.OTP.RateLimit.Time) * time.Minute
}

// GetGracefulShutdownDuration returns the graceful shutdown duration
func (c *Config) GetGracefulShutdownDuration() time.Duration {
	return time.Duration(c.Service.GracefulShutdownSecond) * time.Second
}

// GetDSN returns the PostgreSQL DSN
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.DatabaseName,
		c.Postgres.SSLMode,
		c.Postgres.TimeZone,
	)
}

// GetRedisAddr returns the full Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}
