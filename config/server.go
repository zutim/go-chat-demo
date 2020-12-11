package config

import (
	"chat/pkg/utils"
	"github.com/spf13/viper"
)


// ServerConfig
type ServerConfig struct {
	Port int
	LogPath string
	JwtSign []byte
}


// Server return server config
func Server() (options *ServerConfig) {
	if err := container.Invoke(func(o *ServerConfig) {
		options = o
	}); err != nil {
		options = &ServerConfig{
			Port:               viper.GetInt(httpPortKey),
			LogPath:            viper.GetString(logPathKey),
			JwtSign:         utils.Str2Byte(viper.GetString(jwtSignKey)),
		}

		_ = container.Provide(func() *ServerConfig {
			return options
		})
	}
	return
}
