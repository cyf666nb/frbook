package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Platform PlatformConfig
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Mode         string
}

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	MaxOpenConns int
	MaxIdleConns int
	ConnLifetime time.Duration
}

type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type JWTConfig struct {
	Secret     string
	ExpireTime time.Duration
	Issuer     string
}

type PlatformConfig struct {
	FeeRate       float64
	MaxAge        int
	PayTimeout    time.Duration
	AutoSettleDay int
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func LoadConfigFromEnv() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			Mode:         getEnvAsString("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:         getEnvAsString("DB_HOST", "localhost"),
			Port:         getEnvAsInt("DB_PORT", 3306),
			User:         getEnvAsString("DB_USER", "root"),
			Password:     getEnvAsString("DB_PASSWORD", "root"),
			DBName:       getEnvAsString("DB_NAME", "bookshare"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			ConnLifetime: getEnvAsDuration("DB_CONN_LIFETIME", 1*time.Hour),
		},
		Redis: RedisConfig{
			Host:         getEnvAsString("REDIS_HOST", "localhost"),
			Port:         getEnvAsInt("REDIS_PORT", 6379),
			Password:     getEnvAsString("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 100),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 10),
			MaxRetries:   getEnvAsInt("REDIS_MAX_RETRIES", 3),
			DialTimeout:  getEnvAsDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
			ReadTimeout:  getEnvAsDuration("REDIS_READ_TIMEOUT", 3*time.Second),
			WriteTimeout: getEnvAsDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
		},
		JWT: JWTConfig{
			Secret:     getEnvAsString("JWT_SECRET", "bookshare-secret-key"),
			ExpireTime: getEnvAsDuration("JWT_EXPIRE_TIME", 168*time.Hour),
			Issuer:     getEnvAsString("JWT_ISSUER", "bookshare"),
		},
		Platform: PlatformConfig{
			FeeRate:       getEnvAsFloat("PLATFORM_FEE_RATE", 0.05),
			MaxAge:        getEnvAsInt("PLATFORM_MAX_AGE", 30),
			PayTimeout:    getEnvAsDuration("PLATFORM_PAY_TIMEOUT", 15*time.Minute),
			AutoSettleDay: getEnvAsInt("PLATFORM_AUTO_SETTLE_DAY", 7),
		},
	}
}

func getEnvAsString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if value := viper.GetString(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	if value := viper.GetInt(key); value != 0 {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if value := viper.GetDuration(key); value != 0 {
		return value
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	if value := viper.GetFloat64(key); value != 0 {
		return value
	}
	return defaultValue
}