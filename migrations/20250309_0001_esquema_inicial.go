package migrations

import (
	"sgi-tickets-back/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration20250309_0001 = &gormigrate.Migration{
	ID: "20250309_0001",
	Migrate: func(tx *gorm.DB) error {
		// Crear tablas en orden de dependencias
		// Sin dependencia circular, GORM puede manejar el orden automáticamente
		return tx.AutoMigrate(
			&models.TicketDepartamento{},
			&models.TicketMunicipio{},
			&models.TicketTipoDocumento{},
			&models.TicketRegional{},
			// Catalogos (sin dependencias)
			&models.TicketDependencia{},
			&models.TicketSubdependencia{},
			&models.TicketTipo{},
			&models.TicketSubtipo{},
			&models.TicketPrioridad{},
			&models.TicketEstado{},
			// Usuarios
			&models.TicketUsuario{},
			// Auxiliares (sin dependencias a tickets)
			&models.Cookie{},
			&models.Logger{},
			&models.Permiso{},
			&models.PermisoRol{},
			// Ticket (depende de catalogos y usuario)
			&models.Ticket{},
			// Email bandeja (depende de Ticket - tiene FK ticket_id)
			&models.TicketEmailBandeja{},
			// Historial y seguimiento (dependen de Ticket)
			&models.TicketEstadoLog{},
			&models.TicketComentario{},
			&models.TicketArchivo{},
			&models.TicketNotificacion{},
			&models.TicketAprobacion{},
		)
	},
	Rollback: func(tx *gorm.DB) error {
		// Eliminar tablas en orden inverso de dependencias
		return tx.Migrator().DropTable(
			"regionales",
			"tipo_documentos",
			"municipios",
			"departamentos",
			// Primero las tablas que dependen de tickets
			"tickets_aprobaciones",
			"tickets_notificaciones",
			"tickets_archivos",
			"tickets_comentarios",
			"tickets_estados_log",
			"tickets_email_bandeja",
			// Luego tickets
			"tickets",
			// Auxiliares
			"tickets_permisos_rol",
			"tickets_permisos",
			"tickets_logger",
			"tickets_cookies",
			// Usuarios
			"tickets_usuarios",
			// Catalogos al final
			"tickets_estados",
			"tickets_prioridades",
			"tickets_subtipos",
			"tickets_tipos",
			"tickets_subdependencias",
			"tickets_dependencias",
		)
	},
}
