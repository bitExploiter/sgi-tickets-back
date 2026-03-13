package models

import "gorm.io/gorm"

type TicketDepartamento struct {
	gorm.Model

	Id         uint              `json:"id" gorm:"primaryKey"`
	Nombre     string            `json:"nombre" validate:"required"`
	CodigoDane string            `json:"codigo_dane" validate:"required"`
	Municipios []TicketMunicipio `json:"municipios,omitempty" validate:"-" gorm:"foreignKey:DepartamentoID;references:Id"`
}

func (TicketDepartamento) TableName() string {
	return "departamentos"
}
