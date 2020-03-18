package auth_impl

import (
	"github.com/go-redis/redis"
	"time"
)

// customizeRdsStore An object implementing Store interface
type customizeRdsStore struct {
	redisClient *redis.Client
}

// customizeRdsStore implementing Set method of  Store interface
func (s *customizeRdsStore) Set(id string, value string) {
	err := s.redisClient.Set(id, value, time.Minute*10).Err()
	if err != nil {
		zgh.ZLog().Error("message","auth.AuthLogin","error",err.Error())
	}
}

// customizeRdsStore implementing Get method of  Store interface
func (s *customizeRdsStore) Get(id string, clear bool) (value string) {
	val, err := s.redisClient.Get(id).Result()
	if err != nil {
		zgh.ZLog().Error("message","auth.AuthLogin","error",err.Error())
		return
	}
	if clear {
		err := s.redisClient.Del(id).Err()
		if err != nil {
			zgh.ZLog().Error("message","auth.AuthLogin","error",err.Error())
			return
		}
	}
	return val
}