package migrations

import (
	"fmt"

	"sgi-tickets-back/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func getMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		migration20250309_0001,
		migration20250309_0002,
		migration20250310_0003,
		migration20250310_0004,
		migration20260312_0005,
		migration20260313_0006,
		migration20260313_0007,
		// Futuras migraciones se agregan aqui
	}
}

// RunMigrations ejecuta todas las migraciones pendientes
func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, getMigrations())

	// InitSchema se ejecuta solo si la tabla migrations no existe
	// y la BD esta vacia — crea todo el esquema de una vez
	m.InitSchema(func(tx *gorm.DB) error {
		return tx.AutoMigrate(
			&models.TicketDepartamento{},
			&models.TicketMunicipio{},
			&models.TicketTipoDocumento{},
			&models.TicketRegional{},
			// Catalogos
			&models.TicketDependencia{},
			&models.TicketSubdependencia{},
			&models.TicketTipo{},
			&models.TicketSubtipo{},
			&models.TicketPrioridad{},
			&models.TicketEstado{},
			// Usuarios
			&models.TicketUsuario{},
			// Email
			&models.TicketEmailBandeja{},
			// Ticket
			&models.Ticket{},
			// Historial y seguimiento
			&models.TicketEstadoLog{},
			&models.TicketComentario{},
			&models.TicketArchivo{},
			&models.TicketNotificacion{},
			&models.TicketAprobacion{},
			// Auxiliares
			&models.Cookie{},
			&models.Logger{},
			&models.Permiso{},
			&models.PermisoRol{},
		)
	})

	if err := m.Migrate(); err != nil {
		return fmt.Errorf("error ejecutando migraciones: %w", err)
	}

	fmt.Println("Migraciones ejecutadas correctamente")
	return nil
}

// RollbackMigration revierte la ultima migracion aplicada
func RollbackMigration(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, getMigrations())

	if err := m.RollbackLast(); err != nil {
		return fmt.Errorf("error revirtiendo migracion: %w", err)
	}

	fmt.Println("Rollback ejecutado correctamente")
	return nil
}
