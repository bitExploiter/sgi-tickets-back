package models

import "gorm.io/gorm"

type TicketMunicipio struct {
	gorm.Model

	Id             uint               `json:"id" gorm:"primaryKey"`
	Nombre         string             `json:"nombre" validate:"required"`
	CodigoDane     string             `json:"codigo_dane" validate:"required"`
	DepartamentoID uint               `json:"departamento_id" validate:"required"`
	Departamento   TicketDepartamento `json:"departamento" validate:"-" gorm:"foreignKey:DepartamentoID;references:Id"`
}

func (TicketMunicipio) TableName() string {
	return "municipios"
}
