package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration20250310_0003 = &gormigrate.Migration{
	ID: "20250310_0003",
	Migrate: func(tx *gorm.DB) error {
		return tx.Exec(`
			ALTER TABLE tickets_usuarios
			ADD COLUMN IF NOT EXISTS tipo_documento VARCHAR(50) DEFAULT '',
			ADD COLUMN IF NOT EXISTS numero_documento VARCHAR(50) DEFAULT '',
			ADD COLUMN IF NOT EXISTS telefono VARCHAR(30) DEFAULT '',
			ADD COLUMN IF NOT EXISTS regional VARCHAR(100) DEFAULT '',
			ADD COLUMN IF NOT EXISTS municipio VARCHAR(100) DEFAULT '';
		`).Error
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Exec(`
			ALTER TABLE tickets_usuarios
			DROP COLUMN IF EXISTS tipo_documento,
			DROP COLUMN IF EXISTS numero_documento,
			DROP COLUMN IF EXISTS telefono,
			DROP COLUMN IF EXISTS regional,
			DROP COLUMN IF EXISTS municipio;
		`).Error
	},
}
