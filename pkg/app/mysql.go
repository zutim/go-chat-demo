package app

import (
	"chat/config"
	"github.com/ebar-go/egu"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"sync"
)

var dbPool = new(Database)

type Database struct {
	once sync.Once
	instance *gorm.DB
}

// connect
func (database *Database) connect(dsn string) (err error) {
	// dsn的格式为 用户名:密码/tcp(主机地址)/数据库名称?charset=字符格式
	database.instance, err = gorm.Open("mysql", dsn)
	return
}

// DB return gorm instance
func DB() *gorm.DB {
	dbPool.once.Do(func() {
		egu.SecurePanic(dbPool.connect(config.Mysql().Dsn()))
	})
	return dbPool.instance
}
