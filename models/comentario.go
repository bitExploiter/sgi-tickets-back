package models

import "gorm.io/gorm"

type TicketComentario struct {
	gorm.Model

	Id        uint          `json:"id" gorm:"primaryKey"`
	TicketID  uint          `json:"ticket_id" validate:"required"`
	Ticket    Ticket        `json:"ticket" validate:"-"`
	UsuarioID uint          `json:"usuario_id" validate:"required"`
	Usuario   TicketUsuario `json:"usuario" validate:"-"`
	Contenido string        `json:"contenido" validate:"required"`
	Interno   bool          `json:"interno" gorm:"default:false"`
}

func (TicketComentario) TableName() string {
	return "tickets_comentarios"
}
