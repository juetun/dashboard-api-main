/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-11
 * Time: 00:43
 */
package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/pojos"
)

type AuthRegisterV struct {
}

func (r *AuthRegisterV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := common.NewGin(c)
		var json pojos.AuthRegister
		if err := c.ShouldBindJSON(&json); err != nil {
			// o.Response(http.StatusOK, 400001000, nil)
			appG.Response(400001000, nil)
			return
		}

		reqValidate := &AuthRegister{
			Email:    json.Email,
			Password: json.Password,
			UserName: json.UserName,
		}
 		if b := appG.Validate(reqValidate); !b {
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type AuthRegister struct {
	UserName string `valid:"Required;MaxSize(30)"`
	Email    string `valid:"Required;Email"`
	Password string `valid:"Required;MaxSize(30)"`
}

func (r *AuthRegister) Message() map[string]int {
	return map[string]int{
		"Email.Required":    407000000,
		"Email.Email":       407000001,
		"Password.Required": 407000002,
		"Password.MaxSize":  407000003,
		"UserName.Required": 407000012,
		"UserName.MaxSize":  407000013,
	}
}
