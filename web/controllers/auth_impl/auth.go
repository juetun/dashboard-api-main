/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-09
 * Time: 21:34
 */
package auth_impl

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/controllers"
	"github.com/juetun/app-dashboard/web/pojos"
	"github.com/juetun/app-dashboard/web/services"
	// "github.com/juetun/dashboard/jwt"
	"github.com/mojocn/base64Captcha"
	"golang.org/x/crypto/bcrypt"
)

type ControllerAuth struct {
	base.ControllerBase
}

func NewControllerAuth() controllers.ConsoleAuth {
	controller := &ControllerAuth{}
	controller.ControllerBase.Init()
	return controller
}
func (c *ControllerAuth) Register(ctx *gin.Context) {
	srv := services.NewAuthService()
	cnt, err := srv.GetUserCnt()
	if err != nil {
		c.Log.Error(map[string]string{
			"message": "auth.Register",
			"error":   err.Error(),
		})
		c.Response(ctx, 400001004, nil)
		return
	}
	if cnt >= int64(common.Conf.UserCnt) {
		c.Log.Error(map[string]string{
			"message": "auth.Register",
			"error":   "User cnt beyond expectation",
		})
		c.Response(ctx, 407000015, nil)
		return
	}
	c.Response(ctx, 0, nil)
	return
}
func (c *ControllerAuth) AuthRegister(ctx *gin.Context) {
	srv := services.NewAuthService()
	requestJson, exists := ctx.Get("json")
	if !exists {
		c.Log.Error(map[string]string{
			"message": "auth.AuthRegister",
			"error":   "get request_params from context fail",
		})
		c.Response(ctx, 401000004, nil)
		return
	}
	ar, ok := requestJson.(pojos.AuthRegister)
	if !ok {
		c.Log.Error(map[string]string{
			"message": "auth.AuthRegister",
			"error":   "request_params turn to error",
		})
		c.Response(ctx, 400001001, nil)
		return
	}
	cnt, err := srv.GetUserCnt()
	if err != nil {
		c.Log.Error(map[string]string{
			"message": "auth.AuthRegister",
			"error":   err.Error(),
		})
		c.Response(ctx, 400001004, nil)
		return
	}
	if cnt >= int64(common.Conf.UserCnt) {
		c.Log.Error(map[string]string{
			"message": "auth.AuthRegister",
			"error":   "User cnt beyond expectation",
		})
		c.Response(ctx, 400001004, nil)
		return
	}
	srv.UserStore(ar)
	c.Response(ctx, 0, nil)
	return
}
func (c *ControllerAuth) Login(ctx *gin.Context) {
	srv := services.NewAuthService()
	data, err := srv.Login()
	if err != nil {
		c.Response(ctx, 407000115, nil)
		return
	}
	c.Response(ctx, 0, data)
	return
}
func (c *ControllerAuth) AuthLogin(ctx *gin.Context) {
	srv := services.NewAuthService()
	requestJson, exists := ctx.Get("json")
	if !exists {
		c.Log.Error(map[string]string{
			"message": "auth.AuthLogin",
			"error":   "get request_params from context fail",
		})
		c.Response(ctx, 401000004, nil)
		return
	}
	al, ok := requestJson.(pojos.AuthLogin)
	if !ok {
		c.Log.Error(map[string]string{
			"message": "auth.AuthLogin",
			"error":   "request_params turn to error",
		})
		c.Response(ctx, 400001001, nil)
		return
	}
	verifyResult := base64Captcha.VerifyCaptcha(al.CaptchaKey, al.Captcha)
	if !verifyResult {
		c.Log.Error(map[string]string{
			"message": "auth.AuthLogin",
			"error":   "captcha is error",
		})
		c.Response(ctx, 407000008, nil)
		return
	}

	user, err := srv.GetUserByEmail(al.Email)
	if err != nil {
		c.Log.Error(map[string]string{
			"message": "auth.AuthLogin",
			"error":   "captcha is error",
		})
		c.Response(ctx, 407000010, nil)
		return
	}
	if user.Id <= 0 {
		c.Log.Error(map[string]string{
			"message": "auth.AuthLogin",
			"error":   "Can get user",
		})
		c.Response(ctx, 407000010, nil)
		return
	}

	password := []byte(al.Password)
	hashedPassword := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {

		c.Log.Error(map[string]string{
			"message": "auth.AuthLogin",
			"error":   err.Error(),
		})
		c.Response(ctx, 407000010, nil)
		return
	}

	userIdStr := strconv.Itoa(user.Id)
	token, err := common.CreateToken(userIdStr)
	if err != nil {
		c.Log.Error(map[string]string{
			"message": "auth.AuthLogin",
			"error":   err.Error(),
		})
		c.Response(ctx, 407000011, nil)
		return
	}
	c.Response(ctx, 0, token)
	return
}
func (c *ControllerAuth) Logout(ctx *gin.Context) {
	token, exist := ctx.Get("token")
	if !exist || token == "" {
		c.Log.Error(map[string]string{
			"message": "auth.Logout",
			"error":   "Can not get token",
		})
		c.Response(ctx, 400001004, nil)
		return
	}
	_, err := common.UnsetToken(token.(string))
	if err != nil {
		c.Log.Error(map[string]string{
			"message": "auth.Logout",
			"error":   "Can not get token",
		})
		c.Response(ctx, 407000014, nil)
		return
	}
	c.Response(ctx, 0, token)
	return
}
func (c *ControllerAuth) DelCache(ctx *gin.Context) {
	srv := services.NewAuthService()
	srv.DelAllCache()
	c.Response(ctx, 0, nil)
	return
}
