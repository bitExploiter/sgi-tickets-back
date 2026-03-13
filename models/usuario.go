package models

import (
	"time"

	"gorm.io/gorm"
)

type TicketUsuario struct {
	gorm.Model

	Id                  uint                 `json:"id" gorm:"primaryKey"`
	Nombres             string               `json:"nombres" validate:"required"`
	Apellidos           string               `json:"apellidos" validate:"required"`
	TipoDocumentoID     *uint                `json:"tipo_documento_id" validate:"-"`
	TipoDocumentoRef    *TicketTipoDocumento `json:"tipo_documento_ref,omitempty" validate:"-" gorm:"foreignKey:TipoDocumentoID;references:Id"`
	NumeroDocumento     string               `json:"numero_documento" validate:"-"`
	Email               string               `json:"email" gorm:"unique" validate:"required"`
	Telefono            string               `json:"telefono" validate:"-"`
	RegionalID          *uint                `json:"regional_id" validate:"-"`
	RegionalRef         *TicketRegional      `json:"regional_ref,omitempty" validate:"-" gorm:"foreignKey:RegionalID;references:Id"`
	DepartamentoID      *uint                `json:"departamento_id" validate:"-"`
	DepartamentoRef     *TicketDepartamento  `json:"departamento_ref,omitempty" validate:"-" gorm:"foreignKey:DepartamentoID;references:Id"`
	MunicipioID         *uint                `json:"municipio_id" validate:"-"`
	MunicipioRef        *TicketMunicipio     `json:"municipio_ref,omitempty" validate:"-" gorm:"foreignKey:MunicipioID;references:Id"`
	Password            string               `json:"-"`
	Rol                 string               `json:"rol" validate:"required"`
	Origen              string               `json:"origen" validate:"-"`
	SgiUsuarioID        *uint                `json:"sgi_usuario_id" validate:"-"`
	DependenciaID       *uint                `json:"dependencia_id" validate:"-"`
	Dependencia         *TicketDependencia   `json:"dependencia,omitempty" validate:"-" gorm:"foreignKey:DependenciaID"`
	TotpToken           string               `json:"-"`
	ResetToken          string               `json:"-"`
	ResetTokenExpiry    *time.Time           `json:"-"`
	UltimaFechaConexion *time.Time           `json:"ultima_fecha_conexion" validate:"-"`
	Activo              bool                 `json:"activo" gorm:"default:true"`
}

func (TicketUsuario) TableName() string {
	return "tickets_usuarios"
}
