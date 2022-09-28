package wrapper_admin

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	"strconv"
	"strings"
)

type (
	ArgAdminUserUpdateWithColumn struct {
		UserHids    string  `json:"user_hids" form:"user_hids"`
		UserHidVals []int64 `json:"-" form:"-"`
		Column      string  `json:"column" form:"column"`
		Value       string  `json:"value"  form:"value"`
	}
	ResultAdminUserUpdateWithColumn struct {
		Result bool `json:"result"`
	}
)

func (r *ArgAdminUserUpdateWithColumn) Default(c *base.Context) (err error) {
	_ = c
	if r.UserHids == "" {
		err = fmt.Errorf("请选择要编辑的管理员")
		return
	}
	r.UserHidVals = []int64{}
	v := strings.Split(r.UserHids, ",")
	var uid int64
	for _, s := range v {
		if s == "" {
			continue
		}
		if uid, err = strconv.ParseInt(s, 10, 64); err != nil {
			return
		}
		r.UserHidVals = append(r.UserHidVals, uid)
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
