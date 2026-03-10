package models

import "time"

type TicketEstadoLog struct {
	Id               uint         `json:"id" gorm:"primaryKey"`
	TicketID         uint         `json:"ticket_id" validate:"required"`
	Ticket           Ticket       `json:"ticket" validate:"-"`
	EstadoAnteriorID uint         `json:"estado_anterior_id" validate:"-"`
	EstadoAnterior   TicketEstado `json:"estado_anterior" gorm:"foreignKey:EstadoAnteriorID" validate:"-"`
	EstadoNuevoID    uint         `json:"estado_nuevo_id" validate:"required"`
	EstadoNuevo      TicketEstado `json:"estado_nuevo" gorm:"foreignKey:EstadoNuevoID" validate:"-"`
	UsuarioID        uint         `json:"usuario_id" validate:"required"`
	Usuario          TicketUsuario `json:"usuario" validate:"-"`
	Observaciones    string       `json:"observaciones" validate:"-"`
	CreatedAt        time.Time    `json:"created_at"`
}

func (TicketEstadoLog) TableName() string {
	return "tickets_estados_log"
}
