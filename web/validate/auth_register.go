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
			appG.Response(400001000, err.Error())
			c.Abort()
			return
		}

		reqValidate := &AuthRegister{
			Email:    json.Email,
			Password: json.Password,
			UserName: json.UserName,
		}
		if b := appG.Validate(reqValidate); !b {
			c.Abort()
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

func (r *AuthRegister) Message() map[string]common.ValidationMessage {
	return map[string]common.ValidationMessage{
		"Email.Required.":    {Code: 407000000, Message: "请输入邮箱"},
		"Email.Email.":       {Code: 407000001, Message: "您输入的邮箱格式不正确"},
		"Password.Required.": {Code: 407000002, Message: "请输入密码"},
		"Password.MaxSize.":  {Code: 407000003, Message: "密码不超过30个字符"},
		"UserName.Required.": {Code: 407000012, Message: "请输入用户名"},
		"UserName.MaxSize.":  {Code: 407000013, Message: "用户名不超过30个字符"},
	}

}
