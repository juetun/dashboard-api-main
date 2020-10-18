package models

import (
	"time"

	"github.com/juetun/base-wrapper/lib/base"
)

type Users struct {
	UserHid         string           `gorm:"column:user_hid;" json:"user_hid"`
	Portrait        string           `gorm:"column:portrait;" json:"portrait"`
	Name            string           `gorm:"column:name;" json:"name"`
	Email           string           `gorm:"column:email;" json:"email"`
	Mobile          string           `gorm:"column:mobile;" json:"mobile"`
	Gender          int              `gorm:"column:gender;" json:"gender"`
	Status          int              `gorm:"column:status;" json:"status"`
	EmailVerifiedAt time.Time        `gorm:"column:email_verified_at;" json:"email_verified_at"`
	RememberToken   string           `gorm:"column:remember_token;" json:"remember_token"`
	CreatedAt       base.TimeNormal  `json:"created_at"`
	UpdatedAt       base.TimeNormal  `json:"updated_at"`
	DeletedAt       *base.TimeNormal `sql:"index" json:"-"`
}
