package models

import "gorm.io/gorm"

type TicketSubdependencia struct {
	gorm.Model

	Id             uint               `json:"id" gorm:"primaryKey"`
	DependenciaID  uint               `json:"dependencia_id" validate:"required"`
	Dependencia    TicketDependencia   `json:"dependencia" validate:"-"`
	Nombre         string             `json:"nombre" validate:"required"`
	Activo         bool               `json:"activo" gorm:"default:true"`
}

func (TicketSubdependencia) TableName() string {
	return "tickets_subdependencias"
}
