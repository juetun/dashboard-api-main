package models

type UserMain struct {
	UserHid  string `gorm:"column:user_hid;" json:"user_hid"`
	Portrait string `gorm:"column:portrait;" json:"portrait"`
	Name     string `gorm:"column:name;" json:"name"`
	Email    string `gorm:"column:email;" json:"email"`
	Mobile   string `gorm:"column:mobile;" json:"mobile"`
	Gender   int    `gorm:"column:gender;" json:"gender"`
	Status   int    `gorm:"column:status;" json:"status"`
}
