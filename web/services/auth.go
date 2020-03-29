/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-11
 * Time: 00:17
 */
package services

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/juetun/app-dashboard/lib/app_log"
	"github.com/juetun/app-dashboard/lib/base"
	"github.com/juetun/app-dashboard/lib/common"
	"github.com/juetun/app-dashboard/web/models"
	"github.com/juetun/app-dashboard/web/pojos"
	"github.com/mojocn/base64Captcha"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	base.ServiceBase
}

func NewAuthService(context ...*base.Context) (p *AuthService) {
	p = &AuthService{}
	p.SetContext(context)
	return
}

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
	if !clear {
		return val
	}
	err = r.redisClient.Del(id).Err()
	if err != nil {
		r.Log.Errorln("message", "auth.AuthLogin", "error", err.Error())
		return
	}
	return val

}

func (r *AuthService) Login() (res *map[string]string, err error) {
	// srv := services.NewAuthService()
	customStore := customizeRdsStore{
		redisClient: r.Context.CacheClient,
		Log:         r.Context.Log,
	}
	base64Captcha.SetCustomStore(&customStore)
	var configD = base64Captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 5,
	}
	idKeyD, capD := base64Captcha.GenerateCaptcha("", configD)
	base64stringD := base64Captcha.CaptchaWriteToBase64Encoding(capD)
	data := make(map[string]string)
	data["key"] = idKeyD
	data["png"] = base64stringD
	return &data, err
}
func (r *AuthService) GetUserByEmail(email string) (user *models.ZUsers, err error) {
	user, err = r.GetByAccount(email)
	if err != nil && gorm.IsRecordNotFoundError(err) {
		err = errors.New("账号不存在")
	}
	return
}
func (r *AuthService) GetByAccount(account string) (user *models.ZUsers, err error) {
	user = &models.ZUsers{}
	err = r.Context.Db.Table((&models.ZUsers{}).TableName()).
		Where("email = ? OR mobile=?", account, account).
		Find(user).Error
	return
}

func (r *AuthService) GetUserCnt() (cnt int64, err error) {
	err = r.Context.Db.Table((&models.ZUsers{}).TableName()).Count(&cnt).Error
	return
}

func defaultRegister() (t1 time.Time) {
	time.LoadLocation("")
	loc, _ := time.LoadLocation("Local")
	t1, _ = time.ParseInLocation("2006-01-02 15:04:05", "2000-01-01 00:00:00", loc)
	return
}
func (r *AuthService) UserStore(ar pojos.AuthRegister) (user *models.ZUsers, err error) {
	user, err = r.GetByAccount(ar.Email)

	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {

			return
		}
	}
	if user.Id != 0 {
		err = errors.New("您输入的手机号或邮箱已注册")
		return
	}

	password := []byte(ar.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.UserStore", "error": err.Error(),
		})
		return
	}

	userInsert := models.ZUsers{
		Name:            ar.UserName,
		Email:           ar.Email,
		EmailVerifiedAt: defaultRegister(),
		Password:        string(hashedPassword),
		Status:          1,
	}
	err = r.Context.Db.Create(&userInsert).Error
	if err != nil {
		r.Context.Log.Error(map[string]string{
			"message": "service.UserStore", "error": err.Error(),
		})
		return
	}
	return
}

func (r *AuthService) DelAllCache() {
	conf := common.Conf
	r.Context.CacheClient.Del(
		conf.TagListKey,
		conf.CateListKey,
		conf.ArchivesKey,
		conf.LinkIndexKey,
		conf.PostIndexKey,
		conf.SystemIndexKey,
		conf.TagPostIndexKey,
		conf.CatePostIndexKey,
		conf.PostDetailIndexKey,
	)
}
