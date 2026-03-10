package models

import "time"

type TicketEmailBandeja struct {
	Id               uint       `json:"id" gorm:"primaryKey"`
	MessageID        string     `json:"message_id" validate:"-"`
	Asunto           string     `json:"asunto" validate:"required"`
	Cuerpo           string     `json:"cuerpo" validate:"-"`
	RemitenteEmail   string     `json:"remitente_email" validate:"required"`
	RemitenteNombre  string     `json:"remitente_nombre" validate:"-"`
	LinkCorreo       string     `json:"link_correo" validate:"-"`
	FechaRecibido    *time.Time `json:"fecha_recibido" validate:"-"`
	Procesado        bool       `json:"procesado" gorm:"default:false"`
	TicketID         *uint      `json:"ticket_id" validate:"-"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (TicketEmailBandeja) TableName() string {
	return "tickets_email_bandeja"
}
