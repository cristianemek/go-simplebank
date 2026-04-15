package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config es una estructura que contiene la configuración de la aplicación
// Los valores son leidos por viper desde un archivo de configuración o variables de entorno

type Config struct {
	ENVIROMENT           string        `mapstructure:"ENVIROMENT"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationUrl         string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccesTokenDuration   time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
}

// Load config lee configuracion de una variable de entorno
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
