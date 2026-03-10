package models

import "gorm.io/gorm"

type Cookie struct {
	gorm.Model

	Id         uint   `json:"id" gorm:"primaryKey"`
	Token      string `json:"token" validate:"required"`
	Habilitado bool   `json:"habilitado" gorm:"default:true"`
	Level      string `json:"level" validate:"-"`
}

func (Cookie) TableName() string {
	return "tickets_cookies"
}
