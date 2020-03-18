package plugins

import (
	"github.com/juetun/app-dashboard/lib/app_log"
	"github.com/sirupsen/logrus"
	"os"
)

func PluginLog() (err error) {
	app_log.LogApp = NewAppLog()
	return
}

//初始化日志操作对象
func NewAppLog() *app_log.AppLog {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.WarnLevel)
	return app_log.NewAppLog().SetLog(log)
}
