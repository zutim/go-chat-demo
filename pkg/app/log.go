package app

import (
	"chat/config"
	"chat/pkg/constant"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"sync"
)

type LoggerInstance struct {
	once     sync.Once
	instance *logrus.Logger
}

var logger *LoggerInstance

func init() {
	logger = new(LoggerInstance)
}

// getInstance
func (l *LoggerInstance) getInstance(writer io.Writer) {
	// 实例化,实际项目中一般用全局变量来初始化一个日志管理器
	l.instance = logrus.New()

	// 设置日志内容为json格式
	l.instance.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   constant.TimeFormat, // 时间格式
		DisableTimestamp:  false,               //是否禁用日期
		//DisableHTMLEscape: false,               // 是否禁用html转义
		DataKey:           "",
		//CallerPrettyfier:  nil,
		PrettyPrint:       false, // 是否需要格式化
	})

	l.instance.Level = logrus.DebugLevel

	l.instance.Out = writer

}

// Logger return logrus instance
func Logger() *logrus.Logger {
	// do once
	logger.once.Do(func() {
		// 指定日志的输出为文件，默认是os.Stdout
		f, err := os.OpenFile(path.Join(config.Server().LogPath, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Failed to open:", err.Error())
		}
		logger.getInstance(f)
	})

	return logger.instance
}
