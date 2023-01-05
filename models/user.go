package models

type User struct {
	ID       int    `gorm:"primaryKey,autoIncrement"`
	Username string `gorm:"not null" form:"username"`
	Password string `gorm:"not null" form:"password"`
}
