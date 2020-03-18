/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2019-05-11
 * Time: 00:17
 */
package services

import (
	"fmt"

	"github.com/juetun/study/app-dashboard/lib/base"
	"github.com/juetun/study/app-dashboard/lib/common"
	"github.com/juetun/study/app-dashboard/web/models"
	"github.com/juetun/study/app-dashboard/web/pojos"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	base.ServiceBase
}

func NewAuthService() *AuthService {
	return &AuthService{}
}
func (r *AuthService) GetUserByEmail(email string) (user *models.ZUsers, err error) {
	user = new(models.ZUsers)
	_, err = r.Db.Where("email = ?", email).Get(user)
	return
}

func (r *AuthService) GetUserCnt() (cnt int64, err error) {
	user := new(models.ZUsers)
	cnt, err = r.Db.Count(user)
	return
}

func (r *AuthService) UserStore(ar pojos.AuthRegister) (user *models.ZUsers, err error) {
	password := []byte(ar.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.UserStore", "error": err.Error(),
		})
		return
	}
	userInsert := models.ZUsers{
		Name:     ar.UserName,
		Email:    ar.Email,
		Password: string(hashedPassword),
		Status:   1,
	}
	_, err = r.Db.Insert(&userInsert)
	if err != nil {
		r.Log.Error(map[string]string{
			"message": "service.UserStore", "error": err.Error(),
		})
		return
	}
	fmt.Println(userInsert.Id)
	return
}

func (r *AuthService) DelAllCache() {
	conf := common.Conf
	r.CacheClient.Del(
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
