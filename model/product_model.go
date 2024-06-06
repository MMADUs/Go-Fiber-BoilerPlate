package model

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID         uint64   `gorm:"primaryKey;autoIncrement"`
	Name       string   `json:"name" gorm:"type:varchar(100);not null"`
	Price      float64  `json:"price" gorm:"type:decimal(10,2);not null"`
	Stock      int      `json:"stock" gorm:"type:int;default:0"`
	CategoryID uint     `json:"category_id" gorm:"index;not null"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`
}
