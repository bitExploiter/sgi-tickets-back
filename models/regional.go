package models

import "gorm.io/gorm"

type TicketRegional struct {
	gorm.Model

	Id            uint   `json:"id" gorm:"primaryKey"`
	Nombre        string `json:"nombre" validate:"required"`
	Identificador int    `json:"identificador" validate:"required"`
}

func (TicketRegional) TableName() string {
	return "regionales"
}
