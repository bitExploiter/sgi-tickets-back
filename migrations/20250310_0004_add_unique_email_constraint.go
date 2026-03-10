package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration20250310_0004 = &gormigrate.Migration{
	ID: "20250310_0004",
	Migrate: func(tx *gorm.DB) error {
		return tx.Exec(`
			ALTER TABLE tickets_usuarios
			ADD CONSTRAINT tickets_usuarios_email_unique UNIQUE (email);
		`).Error
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Exec(`
			ALTER TABLE tickets_usuarios
			DROP CONSTRAINT IF EXISTS tickets_usuarios_email_unique;
		`).Error
	},
}
