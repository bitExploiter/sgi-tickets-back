package models

import "gorm.io/gorm"

type TicketPrioridad struct {
	gorm.Model

	Id     uint   `json:"id" gorm:"primaryKey"`
	Nombre string `json:"nombre" validate:"required"`
	Nivel  int    `json:"nivel" validate:"required"`
	Color  string `json:"color" validate:"-"`
	Activo bool   `json:"activo" gorm:"default:true"`
}

func (TicketPrioridad) TableName() string {
	return "tickets_prioridades"
}
