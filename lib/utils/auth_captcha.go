/**
* @Author:changjiang 验证码操作方法
* @Description:
* @File:auth_captcha
* @Version: 1.0.0
* @Date 2020/4/5 3:51 下午
 */
package utils

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/juetun/app-dashboard/lib/app_log"
	"github.com/mojocn/base64Captcha"
)

type AuthCaptcha struct {
	Context *CustomizeRdsStore
}

func NewAuthCaptcha() *AuthCaptcha {
	return &AuthCaptcha{}
}
func (r *AuthCaptcha) SetContext(context *CustomizeRdsStore) *AuthCaptcha {
	r.Context = context
	return r
}
func (r *AuthCaptcha) InitGet() (idKeyD string, base64stringD string) {
	base64Captcha.SetCustomStore(r.Context)
	var configD = base64Captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 5,
	}
	idKeyD, capD := base64Captcha.GenerateCaptcha("", configD)
	base64stringD = base64Captcha.CaptchaWriteToBase64Encoding(capD)
	return idKeyD, base64stringD
}

// customizeRdsStore An object implementing Store interface
type CustomizeRdsStore struct {
	RedisClient *redis.Client
	Log         *app_log.AppLog
}

// customizeRdsStore implementing Set method of  Store interface
func (r *CustomizeRdsStore) Set(id string, value string) {
	err := r.RedisClient.Set(id, value, time.Minute*10).Err()
	if err != nil {
		r.Log.Errorln("message", "auth.AuthLogin", "error", err.Error())
	}
}

// customizeRdsStore implementing Get method of  Store interface
func (r *CustomizeRdsStore) Get(id string, clear bool) (value string) {
	val, err := r.RedisClient.Get(id).Result()
	if err != nil {
		r.Log.Errorln("message", "auth.AuthLogin", "error", err.Error())
		return
	}
	if !clear {
		return val
	}
	err = r.RedisClient.Del(id).Err()
	if err != nil {
		r.Log.Errorln("message", "auth.AuthLogin", "error", err.Error())
		return
	}
	return val

}
