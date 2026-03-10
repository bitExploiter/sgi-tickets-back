package models

import "time"

type TicketNotificacion struct {
	Id         uint          `json:"id" gorm:"primaryKey"`
	TicketID   uint          `json:"ticket_id" validate:"required"`
	Ticket     Ticket        `json:"ticket" validate:"-"`
	UsuarioID  uint          `json:"usuario_id" validate:"required"`
	Usuario    TicketUsuario `json:"usuario" validate:"-"`
	Tipo       string        `json:"tipo" validate:"required"`
	Asunto     string        `json:"asunto" validate:"-"`
	Cuerpo     string        `json:"cuerpo" validate:"-"`
	Enviado    bool          `json:"enviado" gorm:"default:false"`
	FechaEnvio *time.Time    `json:"fecha_envio" validate:"-"`
	Error      string        `json:"error" validate:"-"`
	CreatedAt  time.Time     `json:"created_at"`
}

func (TicketNotificacion) TableName() string {
	return "tickets_notificaciones"
}
