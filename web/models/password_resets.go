package models

import (
	"github.com/juetun/base-wrapper/lib/base"
)

type ZPasswordResets struct {
	base.Model
	Email string `gorm:"column:email;" json:"email"`
	Token string `gorm:"column:token;" json:"token"`
}
