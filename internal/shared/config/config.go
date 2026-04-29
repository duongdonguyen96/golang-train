package config

import (
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string

	PublicSchema       string
    TenantSchemaPrefix string
}

type JWTConfig struct {
	Secret string
}

type S3Config struct {
	Bucket    string
	Region    string
	Endpoint  string
	AccessKey string
	SecretKey string
}

type Config struct {
	AppAddr string
	DB      DBConfig
	JWT     JWTConfig
	S3      S3Config
}

func LoadFromEnv() Config {
	return Config{
		AppAddr: getenv("APP_ADDR", ":8080"),
		DB: DBConfig{
			Host:     getenv("DB_HOST", "localhost"),
			Port:     getenv("DB_PORT", "3306"),
			User:     getenv("DB_USER", "root"),
			Password: getenv("DB_PASSWORD", "root"),
		},
		JWT: JWTConfig{Secret: getenv("JWT_SECRET", "dev-secret")},
		S3: S3Config{
			Bucket:    getenv("S3_BUCKET", "dev-bucket"),
			Region:    getenv("S3_REGION", ""),
			Endpoint:  getenv("S3_ENDPOINT", ""),
			AccessKey: getenv("S3_ACCESS_KEY", ""),
			SecretKey: getenv("S3_SECRET_KEY", ""),
		},
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
