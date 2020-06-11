/**
* @Author:changjiang
* @Description:
* @File:export_data
* @Version: 1.0.0
* @Date 2020/6/10 10:19 上午
 */
package models

import (
	"github.com/juetun/base-wrapper/lib/base"
	utils2 "github.com/juetun/dashboard-api-main/basic/utils"
	"github.com/juetun/dashboard-api-main/basic/utils/hashid"
)

type ZExportData struct {
	base.Model
	Hid           string `gorm:"column:hid;" json:"hid"`
	Name          string `json:"name"`
	Progress      int    `json:"progress"`
	Type          string `json:"type"`
	Arguments     string `json:"arguments"`
	DownloadLink  string `json:"download_link"`
	CreateUserHid string `json:"create_user_hid"`
}

func (r *ZExportData) TableName() string {
	return "z_export_data"
}

func (r *ZExportData) SaltForHID() string {
	return r.TableName()
}

func (r *ZExportData) GetID() int {
	return r.Id
}

func (r *ZExportData) HidUpdate() (err error) {
	if r.Hid == "" {
		// 根据id 设置hid
		r.Hid, err = hashid.Encode(r.SaltForHID(), r.GetID())
		if err != nil {
			return
		}
	}
	return
}

func (r *ZExportData) StartHidInit() {
	r.Hid = utils2.Guid(r.TableName())
}
