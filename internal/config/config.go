package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port         string
	Env          string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

// LoadConfig loads configuration from environment variables.
// In development (APP_ENV empty or "development") it loads .env.dev first.
func LoadConfig() (*Config, error) {
	appEnv := strings.TrimSpace(os.Getenv("APP_ENV"))
	if appEnv == "" {
		appEnv = "development"
	}

	// Load .env.dev only in development
	if appEnv == "development" {
		_ = godotenv.Load(".env.dev")
	}

	// Re-read APP_ENV after potential .env.dev load
	appEnv = strings.TrimSpace(os.Getenv("APP_ENV"))
	if appEnv == "" {
		appEnv = "development"
	}

	serverCfg, err := loadServerConfig(appEnv)
	if err != nil {
		return nil, err
	}

	dbCfg, err := loadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	redisCfg, err := loadRedisConfig()
	if err != nil {
		return nil, err
	}

	jwtCfg, err := loadJWTConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Server:   serverCfg,
		Database: dbCfg,
		Redis:    redisCfg,
		JWT:      jwtCfg,
	}, nil
}

func loadServerConfig(appEnv string) (ServerConfig, error) {
	port, err := getRequiredEnv("SERVER_PORT")
	if err != nil {
		return ServerConfig{}, err
	}

	readTimeout, err := parseDurationEnv("SERVER_READ_TIMEOUT")
	if err != nil {
		return ServerConfig{}, err
	}

	writeTimeout, err := parseDurationEnv("SERVER_WRITE_TIMEOUT")
	if err != nil {
		return ServerConfig{}, err
	}

	return ServerConfig{
		Port:         port,
		Env:          appEnv,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}, nil
}

func loadDatabaseConfig() (DatabaseConfig, error) {
	dsn, err := getRequiredEnv("DB_DSN")
	if err != nil {
		return DatabaseConfig{}, err
	}

	maxOpenConns, err := parseIntEnv("DB_MAX_OPEN_CONNS")
	if err != nil {
		return DatabaseConfig{}, err
	}

	maxIdleConns, err := parseIntEnv("DB_MAX_IDLE_CONNS")
	if err != nil {
		return DatabaseConfig{}, err
	}

	connMaxLifetime, err := parseDurationEnv("DB_CONN_MAX_LIFETIME")
	if err != nil {
		return DatabaseConfig{}, err
	}

	return DatabaseConfig{
		DSN:             dsn,
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
	}, nil
}

func loadRedisConfig() (RedisConfig, error) {
	addr, err := getRequiredEnv("REDIS_ADDR")
	if err != nil {
		return RedisConfig{}, err
	}

	// Redis password MAY be empty in some setups — do NOT force non-empty
	password := strings.TrimSpace(os.Getenv("REDIS_PASSWORD"))

	db, err := parseIntEnv("REDIS_DB")
	if err != nil {
		return RedisConfig{}, err
	}

	return RedisConfig{
		Addr:     addr,
		Password: password,
		DB:       db,
	}, nil
}

func loadJWTConfig() (JWTConfig, error) {
	secret, err := getRequiredEnv("JWT_SECRET")
	if err != nil {
		return JWTConfig{}, err
	}

	accessTTL, err := parseDurationEnv("JWT_ACCESS_TTL")
	if err != nil {
		return JWTConfig{}, err
	}

	refreshTTL, err := parseDurationEnv("JWT_REFRESH_TTL")
	if err != nil {
		return JWTConfig{}, err
	}

	issuer, err := getRequiredEnv("JWT_ISSUER")
	if err != nil {
		return JWTConfig{}, err
	}

	return JWTConfig{
		Secret:          secret,
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
		Issuer:          issuer,
	}, nil
}

// getRequiredEnv returns the trimmed env value or an error if it is empty.
func getRequiredEnv(key string) (string, error) {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return "", fmt.Errorf("missing required environment variable: %s", key)
	}
	return val, nil
}

func parseIntEnv(key string) (int, error) {
	raw, err := getRequiredEnv(key)
	if err != nil {
		return 0, err
	}

	n, convErr := strconv.Atoi(raw)
	if convErr != nil {
		return 0, fmt.Errorf("invalid integer value for %s: %q", key, raw)
	}
	return n, nil
}

func parseDurationEnv(key string) (time.Duration, error) {
	raw, err := getRequiredEnv(key)
	if err != nil {
		return 0, err
	}

	d, convErr := time.ParseDuration(raw)
	if convErr != nil {
		return 0, fmt.Errorf("invalid duration value for %s: %q", key, raw)
	}
	return d, nil
}
