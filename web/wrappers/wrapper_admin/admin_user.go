package wrapper_admin

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	ArgAdminUserUpdateWithColumn struct {
		UserHids    string   `json:"user_hids" form:"user_hids"`
		UserHidVals []string `json:"-" form:"-"`
		Column      string   `json:"column" form:"column"`
		Value       string   `json:"value"  form:"value"`
	}
	ResultAdminUserUpdateWithColumn struct {
		Result bool `json:"result"`
	}
)

func (r *ArgAdminUserUpdateWithColumn) Default(c *gin.Context) (err error) {
	if r.UserHids == "" {
		err = fmt.Errorf("请选择要编辑的管理员")
		return
	}
	r.UserHidVals = []string{}
	v := strings.Split(r.UserHids, ",")

	for _, s := range v {
		if s == "" {
			continue
		}
		r.UserHidVals = append(r.UserHidVals, s)
	}
	if r.Column == "" {
		err = fmt.Errorf("请选择要修改的数据字段")
		return
	}
	if err = r.columnArea(); err != nil {
		return
	}
	return
}

func (r *ArgAdminUserUpdateWithColumn) columnArea() (err error) {
	permitColumn := map[string][]string{
		"can_not_use": {"0", "1"},
	}
	if _, ok := permitColumn[r.Column]; !ok {
		err = fmt.Errorf("当前不支持你选择的字段(%s)", r.Column)
		return
	}
	vArea := permitColumn[r.Column]
	for _, s := range vArea {
		if s == r.Value {
			return
		}
	}
	err = fmt.Errorf("您输入的字段(%s)值(%s)不正确", r.Column, r.Value)
	return
}
