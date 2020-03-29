package models

type ZMigrations struct {
	ID        int    `gorm:"primary_key" json:"id"`
	Migration string  `gorm:"column:migration;" json:"migration"`
	Batch     int     `gorm:"column:batch;" json:"batch"`
}
