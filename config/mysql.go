package config

import (
	"fmt"
	"github.com/spf13/viper"
	"net"
	"strconv"
)

type MysqlConfig struct {
	// database name
	Name string

	// host
	Host string

	// port, default 3306
	Port int

	// user, default root
	User string

	// password
	Password string

	// log mode
	LogMode bool

	// MaxIdleConnections, default 10
	MaxIdleConnections int

	// MaxOpenConnections, default 40
	MaxOpenConnections int

	// max life time, default 10
	MaxLifeTime int
}

// Dsn return mysql dsn
func (options MysqlConfig) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		options.User,
		options.Password,
		net.JoinHostPort(options.Host, strconv.Itoa(options.Port)),
		options.Name)
}


// Server return server config
func Mysql() (options *MysqlConfig) {
	if err := container.Invoke(func(o *MysqlConfig) {
		options = o
	}); err != nil {
		options = &MysqlConfig{
			Name:               viper.GetString(dbNameKey),
			Host:               viper.GetString(dbHostKey),
			Port:               viper.GetInt(dbPortKey),
			User:               viper.GetString(dbUserKey),
			Password:           viper.GetString(dbPassKey),
			MaxOpenConnections: viper.GetInt(dbMaxOpenConnectionsKey),
			MaxIdleConnections: viper.GetInt(dbMaxIdleConnectionsKey),
			MaxLifeTime:        viper.GetInt(dbMaxLifeTimeKey),
		}

		_ = container.Provide(func() *MysqlConfig {
			return options
		})
	}
	return
}