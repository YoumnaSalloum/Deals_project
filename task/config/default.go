package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	DBUri    string `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri string `mapstructure:"REDIS_URL"`
	Port     string `mapstructure:"PORT"`
}

/*
LoadConfig reads app.env from `path` if it exists, and lets real environment
variables override it. In Docker/Kubernetes there is no app.env file, so all
values come from environment variables (set by the Deployment manifest).
*/
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()
	// AutomaticEnv only works for keys viper already knows, so register them.
	for _, k := range []string{"MONGODB_LOCAL_URI", "REDIS_URL", "PORT"} {
		_ = viper.BindEnv(k)
	}
	viper.SetDefault("PORT", "8000")

	if err = viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return
		}
		err = nil // no app.env file: fine, use env vars only
	}

	err = viper.Unmarshal(&config)
	return
}
