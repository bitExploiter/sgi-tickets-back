package models

import "gorm.io/gorm"

type TicketTipoDocumento struct {
	gorm.Model

	Id     uint   `json:"id" gorm:"primaryKey"`
	Nombre string `json:"nombre" validate:"required"`
}

func (TicketTipoDocumento) TableName() string {
	return "tipo_documentos"
}
