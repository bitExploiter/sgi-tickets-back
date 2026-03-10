package models

import "gorm.io/gorm"

type TicketTipo struct {
	gorm.Model

	Id          uint   `json:"id" gorm:"primaryKey"`
	Nombre      string `json:"nombre" validate:"required"`
	Descripcion string `json:"descripcion" validate:"-"`
	Activo      bool   `json:"activo" gorm:"default:true"`
}

func (TicketTipo) TableName() string {
	return "tickets_tipos"
}
