package models

import "gorm.io/gorm"

type TicketDependencia struct {
	gorm.Model

	Id          uint   `json:"id" gorm:"primaryKey"`
	Nombre      string `json:"nombre" validate:"required"`
	Codigo      string `json:"codigo" validate:"required"`
	Descripcion string `json:"descripcion" validate:"-"`
	Activo      bool   `json:"activo" gorm:"default:true"`
}

func (TicketDependencia) TableName() string {
	return "tickets_dependencias"
}
