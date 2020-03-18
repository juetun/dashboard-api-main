/**
 * Created by GoLand.
 * User: xzghua@gmail.com
 * Date: 2018-12-15
 * Time: 16:44
 */
package qq_captcha

import (
	conf2 "github.com/juetun/app-dashboard/conf"
	request2 "github.com/juetun/app-dashboard/request"
	"net/http"
	"time"
)

type QQCaptcha struct {
	Aid string
	AppSecretKey string
	Ticket string
	Randstr string
	UserIP  string
	Url string
}


type qct func(qc *QQCaptcha) interface{}

func (qc *QQCaptcha) SetAid(aid string) qct {
	return func(qc *QQCaptcha) interface{} {
		a := qc.Aid
		qc.Aid = aid
		return a
	}
}

func (qc *QQCaptcha) SetSecretKey(sk string) qct {
	return func(qc *QQCaptcha) interface{} {
		a := qc.AppSecretKey
		qc.AppSecretKey = sk
		return a
	}
}

var qqCaptcha *QQCaptcha


func (qc *QQCaptcha)QQCaptchaInit(options ...qct) error {
	q := &QQCaptcha{
	}
	for _,option := range options {
		option(q)
	}
	qqCaptcha = q
	return nil
}

type QqCaptchaResponse struct {
	Response int `json:"response"`
	EvilLevel int `json:"evil_level"`
	errMsg string `json:"err_msg"`
}

func QQCaptchaVerify(ticket string,randStr string,userIP string) (*http.Response,[]error) {
	resp := new(QqCaptchaResponse)
	res, _,err := request2.New().Get(conf2.QCapUrl).
		Param("aid",qqCaptcha.Aid).
		Param("AppSecretKey",qqCaptcha.AppSecretKey).
		Param("Ticket",ticket).
		Param("Randstr",randStr).
		Param("UserIP",userIP).
		Timeout(time.Minute * 1).Type(request2.TypeUrlencoded).EndStruct(resp)
	return res,err
}