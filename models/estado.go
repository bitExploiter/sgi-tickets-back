package models

import "gorm.io/gorm"

type TicketEstado struct {
	gorm.Model

	Id                  uint   `json:"id" gorm:"primaryKey"`
	Nombre              string `json:"nombre" validate:"required"`
	Codigo              string `json:"codigo" validate:"required"`
	Descripcion         string `json:"descripcion" validate:"-"`
	Color               string `json:"color" validate:"-"`
	EsInicial           bool   `json:"es_inicial" gorm:"default:false"`
	EsFinal             bool   `json:"es_final" gorm:"default:false"`
	RequiereAprobacion  bool   `json:"requiere_aprobacion" gorm:"default:false"`
	Orden               int    `json:"orden" validate:"-"`
	Activo              bool   `json:"activo" gorm:"default:true"`
}

func (TicketEstado) TableName() string {
	return "tickets_estados"
}
