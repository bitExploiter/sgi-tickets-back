package models

import "gorm.io/gorm"

type TicketSubtipo struct {
	gorm.Model

	Id     uint       `json:"id" gorm:"primaryKey"`
	TipoID uint       `json:"tipo_id" validate:"required"`
	Tipo   TicketTipo `json:"tipo" validate:"-"`
	Nombre string     `json:"nombre" validate:"required"`
	Activo bool       `json:"activo" gorm:"default:true"`
}

func (TicketSubtipo) TableName() string {
	return "tickets_subtipos"
}
