package base

import (
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	"github.com/juetun/app-dashboard/lib/app_log"
)

type ServiceBase struct {
	Log         *app_log.AppLog
	Db          *xorm.Engine
	CacheClient *redis.Client
}
