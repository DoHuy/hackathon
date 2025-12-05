package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Storage  StorageConfig
}

type ServerConfig struct {
	Port      string
	LogLevel  string
	PrettyLog bool
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

type StorageConfig struct {
	UploadDir    string
	MaxSizeMB    int64
	AllowedTypes string
}

func Load() (*Config, error) {
	config := new(Config)

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// Load all environment variables into a map
	envMap := loadEnvToMap()

	// Populate config from the map
	populateConfig(config, envMap)

	return config, nil
}

func loadEnvToMap() map[string]string {
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envMap[pair[0]] = pair[1]
		}
	}
	return envMap
}

func populateConfig(config *Config, envMap map[string]string) {
	// Server
	config.Server.Port = getString(envMap, "SERVER_PORT", "8080")
	config.Server.LogLevel = getString(envMap, "SERVER_LOG_LEVEL", "info")
	config.Server.PrettyLog = getBool(envMap, "SERVER_PRETTY_LOG", false)

	// Database
	config.Database.Host = getString(envMap, "DATABASE_HOST", "localhost")
	config.Database.Port = getString(envMap, "DATABASE_PORT", "5432")
	config.Database.User = getString(envMap, "DATABASE_USER", "user")
	config.Database.Password = getString(envMap, "DATABASE_PASSWORD", "password")
	config.Database.DbName = getString(envMap, "DATABASE_DB_NAME", "dbname")
	config.Database.SSLMode = getString(envMap, "DATABASE_SSL_MODE", "disable")

	// JWT
	config.JWT.Secret = getString(envMap, "JWT_SECRET", "your-secret-key")
	config.JWT.ExpirationHours = getInt(envMap, "JWT_EXPIRATION_HOURS", 24)

	// Storage
	config.Storage.UploadDir = getString(envMap, "STORAGE_UPLOAD_DIR", "uploads")
	config.Storage.MaxSizeMB = getInt64(envMap, "STORAGE_MAX_SIZE_MB", 10)
	config.Storage.AllowedTypes = getString(envMap, "STORAGE_ALLOWED_TYPES", "image/jpeg,image/png")
}

func getString(envMap map[string]string, key string, defaultValue string) string {
	if value, ok := envMap[key]; ok {
		return value
	}
	return defaultValue
}

func getInt(envMap map[string]string, key string, defaultValue int) int {
	if valueStr, ok := envMap[key]; ok {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getBool(envMap map[string]string, key string, defaultValue bool) bool {
	if valueStr, ok := envMap[key]; ok {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getInt64(envMap map[string]string, key string, defaultValue int64) int64 {
	if valueStr, ok := envMap[key]; ok {
		if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			return value
		}
	}
	return defaultValue
}
