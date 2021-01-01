package model

type User struct {
	ID       int    `gorm:"primaryKey"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
