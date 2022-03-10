package models

import "fmt"

type ZMigrations struct {
	ID        int    `gorm:"column:id;primary_key" json:"id"`
	Migration string  `gorm:"column:migration;" json:"migration"`
	Batch     int     `gorm:"column:batch;" json:"batch"`
}
func (r *ZMigrations) TableName() string {
	return fmt.Sprintf("%smigrations", TablePrefix)
}