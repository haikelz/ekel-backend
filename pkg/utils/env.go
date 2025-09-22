package utils

import (
	"guestbook-backend/pkg/configs"
	"guestbook-backend/pkg/entities"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Env() entities.EnvVariables {
	v := configs.NewViper()
	logger := configs.NewZap()

	_ = v.BindEnv("TURSO_DATABASE_URL")
	_ = v.BindEnv("TURSO_AUTH_TOKEN")
	_ = v.BindEnv("ADMIN_EMAIL")
	_ = v.BindEnv("JWT_SECRET")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Error("Failed to read config", zap.Error(err))
		}
	}

	return entities.EnvVariables{
		TURSO_DATABASE_URL: v.GetString("TURSO_DATABASE_URL"),
		TURSO_AUTH_TOKEN:   v.GetString("TURSO_AUTH_TOKEN"),
		ADMIN_EMAIL:        v.GetString("ADMIN_EMAIL"),
		ADMIN_PASSWORD:     v.GetString("ADMIN_PASSWORD"),
		JWT_SECRET:         v.GetString("JWT_SECRET"),
	}
}
