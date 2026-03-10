package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration20250309_0002 = &gormigrate.Migration{
	ID: "20250309_0002",
	Migrate: func(tx *gorm.DB) error {
		// Agregar campo reset_token_expiry a tickets_usuarios
		return tx.Exec(`
			ALTER TABLE tickets_usuarios
			ADD COLUMN IF NOT EXISTS reset_token_expiry TIMESTAMP;
		`).Error
	},
	Rollback: func(tx *gorm.DB) error {
		// Eliminar campo reset_token_expiry
		return tx.Exec(`
			ALTER TABLE tickets_usuarios
			DROP COLUMN IF EXISTS reset_token_expiry;
		`).Error
	},
}
