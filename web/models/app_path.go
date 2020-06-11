/**
* @Author:changjiang
* @Description:
* @File:app_path
* @Version: 1.0.0
* @Date 2020/6/11 10:18 上午
 */
package models

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/juetun/base-wrapper/lib/base"
)

type ZAppPath struct {
	base.Model
	Key     string `json:"key"`
	Name    string `json:"name"`
	Port    int    `json:"port"`
	Domains string `json:"domains"`
}

// 随机获取一个域名
func (r *ZAppPath) GetRandomDomain() (res string, err error) {
	if r.Domains == "" {
		err = fmt.Errorf("您没有在数据库中配置%s(%s)的domain信息", r.Name, r.Key)
		return
	}
	var domain []string
	err = json.Unmarshal([]byte(r.Domains), &domain)
	if err != nil {
		err = fmt.Errorf("您在数据库中配置%s(%s)的domain信息不正确 可能不为json格式", r.Name, r.Key)
		return
	}
	rNum := len(domain) - 1
	if rNum > 0 {
		rNum = rand.Intn(rNum)
	}
	res = domain[rNum]
	return
}
func (r *ZAppPath) TableName() string {
	return "z_app_path"
}
