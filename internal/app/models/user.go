package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  *string   `json:"username" validate:"required,min=2,max=100" gorm:"not null;unique"`
	Password  *string   `json:"password" validate:"required,min=6" gorm:"not null"`
	UserType  *string   `json:"user_type" validate:"required,eq=ADMIN|eq=USER" gorm:"not null"`
}

type UserCreate struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	UserType *string `json:"user_type"`
}

type UserLogin struct {
	Username     *string `json:"username"`
	Password     *string `json:"password"`
	Organization *string `json:"organization"`
}

type UserResponse struct {
	ID       uint    `json:"id"`
	Username *string `json:"username"`
	UserType *string `json:"user_type"`
}

func EncodeToUserResponse(user *User) *UserResponse {
	userRes := &UserResponse{}

	userRes.ID = user.ID
	userRes.Username = user.Username
	userRes.UserType = user.UserType

	return userRes
}
