package models

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	UserLogin
	Email string `json:"email" db:"email" binding:"required"`
	Id    int    `json:"-" db:"id"`
}
