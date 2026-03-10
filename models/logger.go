package models

import "gorm.io/gorm"

type Logger struct {
	gorm.Model

	Id          uint          `json:"id" gorm:"primaryKey"`
	UsuarioID   uint          `json:"usuario_id" validate:"-"`
	Usuario     TicketUsuario `json:"usuario" validate:"-"`
	Tipo        string        `json:"tipo" validate:"-"`
	Accion      string        `json:"accion" validate:"-"`
	Descripcion string        `json:"descripcion" validate:"-"`
	Ip          string        `json:"ip" validate:"-"`
}

func (Logger) TableName() string {
	return "tickets_logger"
}
