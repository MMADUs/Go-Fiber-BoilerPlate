package model

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID       uint16    `gorm:"primaryKey;autoIncrement"`
	Name     string    `json:"name" gorm:"type:varchar(50);not null"`
	Products []Product `gorm:"foreignKey:CategoryID"`
}
