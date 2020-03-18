/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2018-12-14
 * Time: 23:48
 */
package common

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	conf2 "github.com/juetun/study/app-dashboard/conf"
	"github.com/juetun/study/app-dashboard/lib/app_log"
	utils2 "github.com/juetun/study/app-dashboard/lib/utils"
)

type JwtParam struct {
	DefaultIss      string
	DefaultAudience string
	DefaultJti      string
	SecretKey       string
	TokenKey        string
	TokenLife       time.Duration
	RedisCache      *redis.Client
}

func (jp *JwtParam) SetTokenKey(tk string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.TokenKey
		jp.TokenKey = tk
		return i
	}
}

func (jp *JwtParam) SetTokenLife(tl time.Duration) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.TokenLife
		jp.TokenLife = tl
		return i
	}
}

func (jp *JwtParam) SetDefaultIss(iss string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.DefaultIss
		jp.DefaultIss = iss
		return i
	}
}

func (jp *JwtParam) SetDefaultAudience(ad string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.DefaultAudience
		jp.DefaultAudience = ad
		return i
	}
}

func (jp *JwtParam) SetDefaultJti(jti string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.DefaultJti
		jp.DefaultJti = jti
		return i
	}
}

func (jp *JwtParam) SetDefaultSecretKey(sk string) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.SecretKey
		jp.SecretKey = sk
		return i
	}
}

func (jp *JwtParam) SetRedisCache(rc *redis.Client) func(jp *JwtParam) interface{} {
	return func(jp *JwtParam) interface{} {
		i := jp.RedisCache
		jp.RedisCache = rc
		return i
	}
}

var jwtParam *JwtParam

func (jp *JwtParam) JwtInit(options ...func(jp *JwtParam) interface{}) error {
	q := &JwtParam{
		DefaultJti:      conf2.JWTJTI,
		DefaultAudience: conf2.JWTAUDIENCE,
		DefaultIss:      conf2.JWTISS,
		SecretKey:       conf2.JWTSECRETKEY,
		TokenLife:       conf2.JWTTOKENLIFE,
		TokenKey:        conf2.JWTTOKENKEY,
	}
	for _, option := range options {
		option(q)
	}
	jwtParam = q
	return nil
}

func CreateToken(userIdString string) (token string, err error) {
	//	iss: jwt签发者
	//	sub: jwt所面向的用户
	//	aud: 接收jwt的一方
	//	exp: jwt的过期时间，这个过期时间必须要大于签发时间
	//	nbf: 定义在什么时间之前，该jwt都是不可用的.
	//	iat: jwt的签发时间
	//	jti: jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击。

	tk := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	// claims["exp"] = time.Now().Add(time.Hour * time.Duration(72)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["iss"] = jwtParam.DefaultIss
	claims["sub"] = userIdString
	claims["aud"] = jwtParam.DefaultAudience
	claims["jti"] = utils2.Md5(jwtParam.DefaultJti + jwtParam.DefaultIss)
	tk.Claims = claims

	SecretKey := jwtParam.SecretKey
	tokenString, err := tk.SignedString([]byte(SecretKey))
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "token create error",
			"error":   err.Error(),
		})
		return "", err
	}

	err = jwtParam.RedisCache.Set(jwtParam.TokenKey+userIdString, tokenString, jwtParam.TokenLife).Err()
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "token create error",
			"error":   err.Error(),
		})
		return "", err
	}

	return tokenString, nil
}

func ParseToken(myToken string) (userId string, err error) {

	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtParam.SecretKey), nil
	})
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "parse token has error",
			"error":   err.Error(),
		})
		return "", err
	}

	if !token.Valid {
		app_log.GetLog().Error(map[string]string{
			"content": fmt.Sprintf("token is invalid(%s)", myToken),
			"error":   err.Error(),
		})
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)

	sub, ok := claims["sub"].(string)
	if !ok {
		app_log.GetLog().Error(map[string]string{
			"content": "claims duan yan is error",
			"error":   err.Error(),
		})
		return "", errors.New("claims duan yan is error")
	}

	res, err := jwtParam.RedisCache.Get(jwtParam.TokenKey + sub).Result()

	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "get token from redis error",
			"error":   err.Error(),
		})
		return "", err
	}

	if res == "" || res != myToken {
		app_log.GetLog().Error(map[string]string{
			"content": "token is invalid",
			"error":   err.Error(),
		})
		return "", errors.New("token is invalid")
	}

	// refresh the token life time
	err = jwtParam.RedisCache.Set(jwtParam.TokenKey+sub, myToken, jwtParam.TokenLife).Err()
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "token create error",
			"error":   err.Error(),
		})
		return "", err
	}

	return sub, nil
}

func UnsetToken(myToken string) (bool, error) {
	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtParam.SecretKey), nil
	})
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "parse token has error",
			"error":   err.Error(),
		})
		return false, err
	}
	claims := token.Claims.(jwt.MapClaims)

	sub, ok := claims["sub"].(string)
	if !ok {
		app_log.GetLog().Error(map[string]string{
			"content": "claims duan yan is error",
			"error":   err.Error(),
		})
		return false, errors.New("claims duan yan is error")
	}
	err = jwtParam.RedisCache.Del(jwtParam.TokenKey + sub).Err()
	if err != nil {
		app_log.GetLog().Error(map[string]string{
			"content": "unset token has error",
			"error":   err.Error(),
		})
		return false, err
	}
	return true, nil
}
