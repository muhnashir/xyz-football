package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort string
	AppEnv  string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret     string
	JWTAccessTTL  int
	JWTRefreshTTL int

	SeasonStartDate time.Time

	UploadDir string
	BaseURL   string

	RateLimit string
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = v.ReadInConfig() // .env is optional; real env vars still apply

	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "xyz_football")
	v.SetDefault("JWT_SECRET", "ganti-dengan-secret-kuat")
	v.SetDefault("JWT_ACCESS_TTL", 3600)
	v.SetDefault("JWT_REFRESH_TTL", 604800)
	v.SetDefault("SEASON_START_DATE", time.Now().Format("2006")+"-01-01")
	v.SetDefault("UPLOAD_DIR", "./uploads")
	v.SetDefault("BASE_URL", "http://localhost:8080")
	v.SetDefault("RATE_LIMIT", "5-M")

	seasonStart, err := time.Parse("2006-01-02", v.GetString("SEASON_START_DATE"))
	if err != nil {
		seasonStart = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	}

	return &Config{
		AppPort:         v.GetString("APP_PORT"),
		AppEnv:          v.GetString("APP_ENV"),
		DBHost:          v.GetString("DB_HOST"),
		DBPort:          v.GetString("DB_PORT"),
		DBUser:          v.GetString("DB_USER"),
		DBPassword:      v.GetString("DB_PASSWORD"),
		DBName:          v.GetString("DB_NAME"),
		JWTSecret:       v.GetString("JWT_SECRET"),
		JWTAccessTTL:    v.GetInt("JWT_ACCESS_TTL"),
		JWTRefreshTTL:   v.GetInt("JWT_REFRESH_TTL"),
		SeasonStartDate: seasonStart,
		UploadDir:       v.GetString("UPLOAD_DIR"),
		BaseURL:         v.GetString("BASE_URL"),
		RateLimit:       v.GetString("RATE_LIMIT"),
	}, nil
}
