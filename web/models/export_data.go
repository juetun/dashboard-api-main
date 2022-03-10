// Package models
/**
* @Author:changjiang
* @Description:
* @File:export_data
* @Version: 1.0.0
* @Date 2020/6/10 10:19 上午
 */
package models

import (
	"fmt"
	"github.com/juetun/base-wrapper/lib/base"
	utils2 "github.com/juetun/base-wrapper/lib/utils"
	"github.com/juetun/base-wrapper/lib/utils/hashid"
	"github.com/juetun/dashboard-api-main/web"
)

type ZExportData struct {
	base.Model
	Hid           string `gorm:"column:hid;" json:"hid"`
	Name          string `gorm:"column:name;" json:"name"`
	Progress      int    `gorm:"column:progress;" json:"progress"`
	Status        int    `gorm:"column:status;" json:"status"`
	Type          string `gorm:"column:type;" json:"type"`
	Arguments     string `gorm:"column:arguments;" json:"arguments"`
	DownloadLink  string `gorm:"column:download_link;" json:"download_link"`
	Domain        string `gorm:"column:domain;" json:"domain"`
	FilePath      string `gorm:"column:file_path;" json:"file_path"`
	CreateUserHid int64  `gorm:"column:create_user_hid;" json:"create_user_hid"`
}

func (r *ZExportData) TableName() string {
	return fmt.Sprintf("%sexport_data", TablePrefix)
}

func (r *ZExportData) SaltForHID() string {
	return r.TableName()
}

func (r *ZExportData) GetCacheKey() string {
	return web.RedisCacheKeyPrefixExport + r.Hid
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
