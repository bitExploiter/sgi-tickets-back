package models

import "gorm.io/gorm"

type TicketArchivo struct {
	gorm.Model

	Id             uint             `json:"id" gorm:"primaryKey"`
	TicketID       uint             `json:"ticket_id" validate:"required"`
	Ticket         Ticket           `json:"ticket" validate:"-"`
	ComentarioID   *uint            `json:"comentario_id" validate:"-"`
	Comentario     TicketComentario `json:"comentario" validate:"-"`
	UsuarioID      uint             `json:"usuario_id" validate:"required"`
	Usuario        TicketUsuario    `json:"usuario" validate:"-"`
	Nombre         string           `json:"nombre" validate:"required"`
	NombreStorage  string           `json:"nombre_storage" validate:"-"`
	Path           string           `json:"path" validate:"-"`
	Tipo           string           `json:"tipo" validate:"-"`
	Tamanio        int64            `json:"tamanio" validate:"-"`
}

func (TicketArchivo) TableName() string {
	return "tickets_archivos"
}
