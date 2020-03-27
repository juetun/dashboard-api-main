/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-10
 * Time: 23:23
 */
package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/pojos"
)

type AuthLoginV struct {
}

func (av *AuthLoginV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG:=common.NewGin(c)
		var json pojos.AuthLogin
		if err := c.ShouldBindJSON(&json); err != nil {
			appG.Response(400001000, nil)
			return
		}

		reqValidate := &AuthLogin{
			Email:      json.Email,
			Password:   json.Password,
			Captcha:    json.Captcha,
			CaptchaKey: json.CaptchaKey,
		}
		if b := appG.Validate(reqValidate); !b {
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type AuthLogin struct {
	Email      string `valid:"Required;Email"`
	Password   string `valid:"Required;MaxSize(30)"`
	Captcha    string `valid:"Required;MaxSize(5)"`
	CaptchaKey string `valid:"Required;MaxSize(30)"`
}

func (av *AuthLogin) Message() map[string]int {
	return map[string]int{
		"Email.Required":      407000000,
		"Email.Email":         407000001,
		"Password.Required":   407000002,
		"Password.MaxSize":    407000003,
		"Captcha.Required":    407000004,
		"Captcha.MaxSize":     407000005,
		"CaptchaKey.Required": 407000006,
		"CaptchaKey.MaxSize":  407000007,
	}
}
