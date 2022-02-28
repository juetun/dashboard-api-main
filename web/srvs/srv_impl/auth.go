// Package srv_impl
/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-11
 * Time: 00:17
 */
package srv_impl

// import (
// 	"errors"
// 	"fmt"
// 	"time"
//
// 	"gorm.io/gorm"
// 	"github.com/juetun/base-wrapper/lib/app/app_obj"
// 	"github.com/juetun/base-wrapper/lib/base"
// 	"github.com/juetun/base-wrapper/lib/common"
// 	"github.com/juetun/dashboard-api-main/basic/utils"
// 	"github.com/juetun/dashboard-api-main/web/models"
// 	"github.com/juetun/dashboard-api-main/web/pojos"
// 	"github.com/mojocn/base64Captcha"
// 	"golang.org/x/crypto/bcrypt"
// )
//
// const MaxPermitRegisterUserCount = 100000000000
//
// type AuthService struct {
// 	base.ServiceBase
// }
//
// func NewAuthService(context ...*base.Context) (p *AuthService) {
// 	p = &AuthService{}
// 	p.SetContext(context)
// 	return
// }
//
// func (r *AuthService) AuthLogin(logArg *pojos.AuthLogin) (user *models.ZUsers, token string, err error) {
// 	verifyResult := base64Captcha.VerifyCaptcha(logArg.CaptchaKey, logArg.Captcha)
// 	if !verifyResult {
// 		r.Context.Log.Error(map[string]string{
// 			"message": "service.AuthLogin",
// 			"error":   "captcha is error",
// 		})
// 		err = fmt.Errorf("您输入的验证码不正确")
// 		return
// 	}
// 	user, err = r.GetUserByEmail(logArg.Email)
// 	if err != nil {
// 		if gorm.IsRecordNotFoundError(err) {
// 			err = fmt.Errorf("您要登录的用户信息不存在")
// 		}
// 		r.Context.Log.Error(map[string]string{
// 			"message": "service.AuthLogin",
// 			"error":   "Can get user",
// 		})
// 		return
// 	}
//
// 	password := []byte(logArg.Password)
// 	hashedPassword := []byte(user.Password)
// 	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
// 	if err != nil {
// 		r.Context.Log.Error(map[string]string{
// 			"message": "auth.AuthLogin",
// 			"error":   err.Error(),
// 		})
// 		err = fmt.Errorf("您要登录的账号或密码错误")
// 		return
// 	}
//
// 	token, err = common.CreateToken(app_obj.JwtUser{
// 		UserId: user.UserHid,
// 		Name:   user.Name,
// 		Status: user.Status,
// 	})
// 	if err != nil {
// 		r.Context.Log.Error(map[string]string{
// 			"message": "auth.AuthLogin",
// 			"error":   err.Error(),
// 		})
// 		err = fmt.Errorf("系统异常，请刷新或稍后重试")
// 		return
// 	}
// 	return
// }
// func (r *AuthService) Login() (res *map[string]string, err error) {
// 	// srv := services.NewAuthService()
// 	idKeyD, base64stringD := utils.NewAuthCaptcha().SetContext(&utils.CustomizeRdsStore{
// 		RedisClient: r.Context.CacheClient,
// 		Log:         r.Context.Log,
// 	}).InitGet()
// 	data := make(map[string]string)
// 	data["key"] = idKeyD
// 	data["png"] = base64stringD
// 	return &data, err
// }
// func (r *AuthService) GetUserByEmail(email string) (user *models.ZUsers, err error) {
// 	user, err = r.GetByAccount(email)
// 	if err != nil && gorm.IsRecordNotFoundError(err) {
// 		err = errors.New("账号不存在")
// 	}
// 	return
// }
// func (r *AuthService) GetByAccount(account string) (user *models.ZUsers, err error) {
// 	user = &models.ZUsers{}
// 	err = r.Context.Db.Table((&models.ZUsers{}).TableName()).
// 		Where("email = ? OR mobile=?", account, account).
// 		Find(user).Error
// 	return
// }
//
// func (r *AuthService) GetRegisterCappcha() (data map[string]interface{}, err error) {
// 	data = make(map[string]interface{})
// 	var cnt float64
// 	cnt, err = r.GetUserCnt()
// 	if err != nil {
// 		return
// 	}
// 	if cnt >= MaxPermitRegisterUserCount {
// 		r.Context.Log.Error(map[string]string{
// 			"message": "auth.Register",
// 			"error":   "User cnt beyond expectation",
// 		})
// 		err = fmt.Errorf("当前注册用户的数超过了系统规定的数,请尝试联系客服")
// 		return
// 	}
// 	idKeyD, base64stringD := utils.NewAuthCaptcha().SetContext(&utils.CustomizeRdsStore{
// 		RedisClient: r.Context.CacheClient,
// 		Log:         r.Context.Log,
// 	}).InitGet()
// 	data["key"] = idKeyD
// 	data["png"] = base64stringD
// 	data["cnt"] = cnt
// 	return
// }
// func (r *AuthService) GetUserCnt() (cnt float64, err error) {
// 	err = r.Context.Db.Table((&models.ZUsers{}).TableName()).
// 		Unscoped().
// 		Where("deleted_at IS NULL").
// 		Count(&cnt).Error
// 	return
// }
//
// func (r *AuthService) AuthRegister(regArg *pojos.AuthRegister) (res bool, err error) {
//
// 	verifyResult := base64Captcha.VerifyCaptcha(regArg.CaptchaKey, regArg.Captcha)
// 	if !verifyResult {
// 		r.Context.Log.Error(map[string]string{
// 			"message": "service.AuthLogin",
// 			"error":   "captcha is error",
// 		})
// 		err = fmt.Errorf("您输入的验证码不正确")
// 		return
// 	}
// 	_, err = r.UserStore(regArg)
// 	if err != nil {
// 		r.Context.Log.Error(map[string]string{
// 			"message": "auth.AuthRegister",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}
// 	res = true
// 	return
// }
// func defaultRegister() (t1 time.Time) {
// 	time.LoadLocation("")
// 	loc, _ := time.LoadLocation("Local")
// 	t1, _ = time.ParseInLocation("2006-01-02 15:04:05", "2000-01-01 00:00:00", loc)
// 	return
// }
// func (r *AuthService) UserStore(ar *pojos.AuthRegister) (user *models.ZUsers, err error) {
// 	user, err = r.GetByAccount(ar.Email)
// 	if err != nil && !gorm.IsRecordNotFoundError(err) { // 如果出异常了且异常不为没查到数据。
// 		return
// 	}
// 	if user.UserHid != "" {
// 		err = errors.New("您输入的手机号或邮箱已注册")
// 		return
// 	}
//
// 	password := []byte(ar.Password)
// 	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
// 	if err != nil {
// 		r.Context.Log.Error(map[string]string{
// 			"message": "service.UserStore", "error": err.Error(),
// 		})
// 		return
// 	}
//
// 	userInsert := models.ZUsers{
// 		Name:            "昵称",
// 		Email:           ar.Email,
// 		EmailVerifiedAt: defaultRegister(),
// 		Password:        string(hashedPassword),
// 		Status:          1,
// 	}
// 	err = r.Context.Db.Create(&userInsert).Error
// 	if err != nil {
// 		r.Context.Log.Error(map[string]string{
// 			"message": "service.UserStore", "error": err.Error(),
// 		})
// 		return
// 	}
// 	return
// }
//
// func (r *AuthService) DelAllCache() {
// 	conf := common.Conf
// 	r.Context.CacheClient.Del(
// 		conf.TagListKey,
// 		conf.CateListKey,
// 		conf.ArchivesKey,
// 		conf.LinkIndexKey,
// 		conf.PostIndexKey,
// 		conf.SystemIndexKey,
// 		conf.TagPostIndexKey,
// 		conf.CatePostIndexKey,
// 		conf.PostDetailIndexKey,
// 	)
// }
