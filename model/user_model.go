package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          string `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username    string `json:"username" gorm:"type:varchar(30);not null"`
	Email       string `json:"email" gorm:"type:varchar(50);not null"`
	PhoneNumber string `json:"phone_number" gorm:"type:varchar(20);not null"`
	Password    string `json:"password" gorm:"type:varchar;not null"`
	Role        int    `json:"role" gorm:"type:int;default:0"`
}