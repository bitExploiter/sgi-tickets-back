package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration20260312_0005 = &gormigrate.Migration{
	ID: "20260312_0005",
	Migrate: func(tx *gorm.DB) error {
		return tx.Exec(`
			ALTER TABLE tickets_usuarios
			ADD COLUMN IF NOT EXISTS ultima_fecha_conexion TIMESTAMP;
		`).Error
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Exec(`
			ALTER TABLE tickets_usuarios
			DROP COLUMN IF EXISTS ultima_fecha_conexion;
		`).Error
	},
}
