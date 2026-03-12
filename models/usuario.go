package models

import (
	"time"

	"gorm.io/gorm"
)

type TicketUsuario struct {
	gorm.Model

	Id                  uint               `json:"id" gorm:"primaryKey"`
	Nombres             string             `json:"nombres" validate:"required"`
	Apellidos           string             `json:"apellidos" validate:"required"`
	TipoDocumento       string             `json:"tipo_documento" validate:"-"`
	NumeroDocumento     string             `json:"numero_documento" validate:"-"`
	Email               string             `json:"email" gorm:"unique" validate:"required"`
	Telefono            string             `json:"telefono" validate:"-"`
	Regional            string             `json:"regional" validate:"-"`
	Municipio           string             `json:"municipio" validate:"-"`
	Password            string             `json:"-"`
	Rol                 string             `json:"rol" validate:"required"`
	Origen              string             `json:"origen" validate:"-"`
	SgiUsuarioID        *uint              `json:"sgi_usuario_id" validate:"-"`
	DependenciaID       *uint              `json:"dependencia_id" validate:"-"`
	Dependencia         *TicketDependencia `json:"dependencia,omitempty" validate:"-" gorm:"foreignKey:DependenciaID"`
	TotpToken           string             `json:"-"`
	ResetToken          string             `json:"-"`
	ResetTokenExpiry    *time.Time         `json:"-"`
	UltimaFechaConexion *time.Time         `json:"ultima_fecha_conexion" validate:"-"`
	Activo              bool               `json:"activo" gorm:"default:true"`
}

func (TicketUsuario) TableName() string {
	return "tickets_usuarios"
}
