package config

import "gopkg.in/ini.v1"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Storage  StorageConfig
}

type ServerConfig struct {
	Port      string `ini:"PORT"`
	LogLevel  string `ini:"LOG_LEVEL"`
	PrettyLog bool   `ini:"PRETTY_LOG"`
}

type DatabaseConfig struct {
	Host     string `ini:"HOST"`
	Port     string `ini:"PORT"`
	User     string `ini:"USER"`
	Password string `ini:"PASSWORD"`
	DbName   string `ini:"DB_NAME"`
	SSLMode  string `ini:"SSL_MODE"`
}

type JWTConfig struct {
	Secret          string `ini:"SECRET"`
	ExpirationHours int    `ini:"EXPIRATION_HOURS"`
}

type StorageConfig struct {
	UploadDir string `ini:"UPLOAD_DIR"`
	MaxSizeMB int64  `ini:"MAX_SIZE_MB"`
}

func Load(path string) (*Config, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	cfg.Section("server").MapTo(&config.Server)
	cfg.Section("database").MapTo(&config.Database)
	cfg.Section("jwt").MapTo(&config.JWT)
	cfg.Section("storage").MapTo(&config.Storage)
	return config, nil
}
