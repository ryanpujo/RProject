package models

import "time"

type User struct {
	Id        int       `json:"id"`
	Fname     string    `json:"fname"`
	Lname     string    `json:"lname"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserPayload struct {
	Id        int       `json:"id"`
	Fname     string    `json:"fname" binding:"required,max=10"`
	Lname     string    `json:"lname" binding:"required,max=10"`
	Password  string    `json:"password" binding:"required,min=12"`
	Email     string    `json:"email" binding:"required,email"`
	Username  string    `json:"username" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
