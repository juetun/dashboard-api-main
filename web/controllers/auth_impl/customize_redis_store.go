package auth_impl

import (
	"github.com/go-redis/redis"
	"github.com/juetun/app-dashboard/lib/app_log"

	"time"
)

// customizeRdsStore An object implementing Store interface
type customizeRdsStore struct {
	redisClient *redis.Client
	Log         *app_log.AppLog
}

// customizeRdsStore implementing Set method of  Store interface
func (r *customizeRdsStore) Set(id string, value string) {
	err := r.redisClient.Set(id, value, time.Minute*10).Err()
	if err != nil {
		r.Log.Errorln("message", "auth.AuthLogin", "error", err.Error())
	}
}

// customizeRdsStore implementing Get method of  Store interface
func (r *customizeRdsStore) Get(id string, clear bool) (value string) {
	val, err := r.redisClient.Get(id).Result()
	if err != nil {
		r.Log.Errorln("message", "auth.AuthLogin", "error", err.Error())
		return
	}
	if clear {
		err := r.redisClient.Del(id).Err()
		if err != nil {
			r.Log.Errorln("message", "auth.AuthLogin", "error", err.Error())
			return
		}
	}
	return val
}
