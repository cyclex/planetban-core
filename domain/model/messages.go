package model

import "github.com/jinzhu/gorm"

type Messages struct {
	gorm.Model
	Message string `gorm:"message" json:"message"`
}
