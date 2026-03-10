package models

import "time"

type TicketAprobacion struct {
	Id            uint          `json:"id" gorm:"primaryKey"`
	TicketID      uint          `json:"ticket_id" validate:"required"`
	Ticket        Ticket        `json:"ticket" validate:"-"`
	UsuarioID     uint          `json:"usuario_id" validate:"required"`
	Usuario       TicketUsuario `json:"usuario" validate:"-"`
	TipoAprobador string        `json:"tipo_aprobador" validate:"required"`
	Aprobado      bool          `json:"aprobado" gorm:"default:false"`
	Observaciones string        `json:"observaciones" validate:"-"`
	CreatedAt     time.Time     `json:"created_at"`
}

func (TicketAprobacion) TableName() string {
	return "tickets_aprobaciones"
}
