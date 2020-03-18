package app_log

import (
	"github.com/sirupsen/logrus"
)

var LogApp *AppLog

type AppLog struct {
	*logrus.Logger
}

func NewAppLog() *AppLog {
	return &AppLog{}
}
func GetLog() *AppLog {
	return LogApp
}
func (r *AppLog) SetLog(log *logrus.Logger) *AppLog {
	r.Logger = log
	return r
}

func (r *AppLog) Error(data map[string]string, message ...string) {
	fields := logrus.Fields{}
	if len(data) > 0 {
		for key, value := range data {
			fields[key] = value
		}
	}
	r.WithFields(fields).Error(message)
}
func (r *AppLog) Info(data map[string]string, message ...string) {
	fields := logrus.Fields{}
	if len(data) > 0 {
		for key, value := range data {
			fields[key] = value
		}
	}
	r.WithFields(fields).Error(message)
}
