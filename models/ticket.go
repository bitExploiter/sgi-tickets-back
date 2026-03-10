package models

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	gorm.Model

	Id                             uint                `json:"id" gorm:"primaryKey"`
	Numero                         string              `json:"numero" validate:"-"`
	Asunto                         string              `json:"asunto" validate:"required"`
	Descripcion                    string              `json:"descripcion" validate:"-"`
	DependenciaID                  uint                `json:"dependencia_id" validate:"required"`
	Dependencia                    TicketDependencia    `json:"dependencia" validate:"-"`
	SubdependenciaID               uint                `json:"subdependencia_id" validate:"-"`
	Subdependencia                 TicketSubdependencia `json:"subdependencia" validate:"-"`
	TipoID                         uint                `json:"tipo_id" validate:"required"`
	Tipo                           TicketTipo           `json:"tipo" validate:"-"`
	SubtipoID                      uint                `json:"subtipo_id" validate:"-"`
	Subtipo                        TicketSubtipo        `json:"subtipo" validate:"-"`
	PrioridadID                    uint                `json:"prioridad_id" validate:"required"`
	Prioridad                      TicketPrioridad      `json:"prioridad" validate:"-"`
	EstadoID                       uint                `json:"estado_id" validate:"required"`
	Estado                         TicketEstado         `json:"estado" validate:"-"`
	SolicitanteID                  uint                `json:"solicitante_id" validate:"-"`
	Solicitante                    TicketUsuario        `json:"solicitante" gorm:"foreignKey:SolicitanteID" validate:"-"`
	AgenteID                       uint                `json:"agente_id" validate:"-"`
	Agente                         TicketUsuario        `json:"agente" gorm:"foreignKey:AgenteID" validate:"-"`
	ContratistaID                  uint                `json:"contratista_id" validate:"-"`
	Contratista                    TicketUsuario        `json:"contratista" gorm:"foreignKey:ContratistaID" validate:"-"`
	PplNombre                      string              `json:"ppl_nombre" validate:"-"`
	PplNui                         string              `json:"ppl_nui" validate:"-"`
	LinkCorreo                     string              `json:"link_correo" validate:"-"`
	VisitaID                       uint                `json:"visita_id" validate:"-"`
	FechaLimite                    *time.Time          `json:"fecha_limite" validate:"-"`
	AprobacionEntidad              bool                `json:"aprobacion_entidad" gorm:"default:false"`
	AprobacionInterventoria        bool                `json:"aprobacion_interventoria" gorm:"default:false"`
	FechaAprobacionEntidad         *time.Time          `json:"fecha_aprobacion_entidad" validate:"-"`
	FechaAprobacionInterventoria   *time.Time          `json:"fecha_aprobacion_interventoria" validate:"-"`
}

func (Ticket) TableName() string {
	return "tickets"
}
