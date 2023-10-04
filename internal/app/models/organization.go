package models

import "time"

type Organization struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	Name      *string   `json:"name" validate:"required,min=2,max=100" gorm:"not null;unique"`
}

type OrganizationCreate struct {
	Name     *string `json:"name" validate:"required,min=2,max=100"`
	Username *string `json:"username" validate:"required,min=2,max=100"`
	Password *string `json:"password" validate:"required,min=6"`
}
