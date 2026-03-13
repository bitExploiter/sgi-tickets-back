package migrations

import (
	"fmt"

	"sgi-tickets-back/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migration20260313_0006 = &gormigrate.Migration{
	ID: "20260313_0006",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(
			&models.TicketDepartamento{},
			&models.TicketMunicipio{},
			&models.TicketTipoDocumento{},
			&models.TicketRegional{},
		); err != nil {
			return err
		}

		if err := tx.Exec(`
			ALTER TABLE tickets_usuarios
			ADD COLUMN IF NOT EXISTS tipo_documento_id BIGINT,
			ADD COLUMN IF NOT EXISTS regional_id BIGINT,
			ADD COLUMN IF NOT EXISTS departamento_id BIGINT,
			ADD COLUMN IF NOT EXISTS municipio_id BIGINT;
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			CREATE INDEX IF NOT EXISTS idx_municipios_departamento_id ON municipios(departamento_id);
			CREATE UNIQUE INDEX IF NOT EXISTS idx_departamentos_codigo_dane ON departamentos(codigo_dane);
			CREATE UNIQUE INDEX IF NOT EXISTS idx_municipios_codigo_dane ON municipios(codigo_dane);
			CREATE UNIQUE INDEX IF NOT EXISTS idx_tipo_documentos_nombre ON tipo_documentos(nombre);
			CREATE UNIQUE INDEX IF NOT EXISTS idx_regionales_identificador ON regionales(identificador);
			CREATE INDEX IF NOT EXISTS idx_tickets_usuarios_tipo_documento_id ON tickets_usuarios(tipo_documento_id);
			CREATE INDEX IF NOT EXISTS idx_tickets_usuarios_regional_id ON tickets_usuarios(regional_id);
			CREATE INDEX IF NOT EXISTS idx_tickets_usuarios_departamento_id ON tickets_usuarios(departamento_id);
			CREATE INDEX IF NOT EXISTS idx_tickets_usuarios_municipio_id ON tickets_usuarios(municipio_id);
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			DO $$
			BEGIN
				IF NOT EXISTS (
					SELECT 1 FROM pg_constraint WHERE conname = 'fk_tickets_usuarios_tipo_documento'
				) THEN
					ALTER TABLE tickets_usuarios
					ADD CONSTRAINT fk_tickets_usuarios_tipo_documento
					FOREIGN KEY (tipo_documento_id) REFERENCES tipo_documentos(id);
				END IF;
			END $$;
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			DO $$
			BEGIN
				IF NOT EXISTS (
					SELECT 1 FROM pg_constraint WHERE conname = 'fk_tickets_usuarios_regional'
				) THEN
					ALTER TABLE tickets_usuarios
					ADD CONSTRAINT fk_tickets_usuarios_regional
					FOREIGN KEY (regional_id) REFERENCES regionales(id);
				END IF;
			END $$;
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			DO $$
			BEGIN
				IF NOT EXISTS (
					SELECT 1 FROM pg_constraint WHERE conname = 'fk_tickets_usuarios_departamento'
				) THEN
					ALTER TABLE tickets_usuarios
					ADD CONSTRAINT fk_tickets_usuarios_departamento
					FOREIGN KEY (departamento_id) REFERENCES departamentos(id);
				END IF;
			END $$;
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			DO $$
			BEGIN
				IF NOT EXISTS (
					SELECT 1 FROM pg_constraint WHERE conname = 'fk_tickets_usuarios_municipio'
				) THEN
					ALTER TABLE tickets_usuarios
					ADD CONSTRAINT fk_tickets_usuarios_municipio
					FOREIGN KEY (municipio_id) REFERENCES municipios(id);
				END IF;
			END $$;
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			DO $$
			BEGIN
				IF NOT EXISTS (
					SELECT 1 FROM pg_constraint WHERE conname = 'fk_municipios_departamento'
				) THEN
					ALTER TABLE municipios
					ADD CONSTRAINT fk_municipios_departamento
					FOREIGN KEY (departamento_id) REFERENCES departamentos(id);
				END IF;
			END $$;
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE tickets_usuarios u
			SET tipo_documento_id = td.id
			FROM tipo_documentos td
			WHERE u.tipo_documento_id IS NULL
			AND btrim(u.tipo_documento) <> ''
			AND (
				LOWER(btrim(u.tipo_documento)) = LOWER(btrim(td.nombre))
				OR (LOWER(btrim(u.tipo_documento)) = 'cc' AND LOWER(td.nombre) = LOWER('Cédula de Ciudadanía'))
				OR (LOWER(btrim(u.tipo_documento)) = 'ce' AND LOWER(td.nombre) = LOWER('Cédula de Extranjería'))
			);
		`).Error; err != nil {
			return fmt.Errorf("error mapeando tipo_documento_id: %w", err)
		}

		if err := tx.Exec(`
			UPDATE tickets_usuarios u
			SET regional_id = r.id
			FROM regionales r
			WHERE u.regional_id IS NULL
			AND btrim(u.regional) <> ''
			AND LOWER(btrim(u.regional)) = LOWER(btrim(r.nombre));
		`).Error; err != nil {
			return fmt.Errorf("error mapeando regional_id: %w", err)
		}

		if err := tx.Exec(`
			UPDATE tickets_usuarios u
			SET municipio_id = m.id,
				departamento_id = m.departamento_id
			FROM municipios m
			WHERE u.municipio_id IS NULL
			AND btrim(u.municipio) <> ''
			AND LOWER(btrim(u.municipio)) = LOWER(btrim(m.nombre));
		`).Error; err != nil {
			return fmt.Errorf("error mapeando municipio_id/departamento_id: %w", err)
		}

		return nil
	},
	Rollback: func(tx *gorm.DB) error {
		if err := tx.Exec(`
			ALTER TABLE tickets_usuarios
			DROP CONSTRAINT IF EXISTS fk_tickets_usuarios_tipo_documento,
			DROP CONSTRAINT IF EXISTS fk_tickets_usuarios_regional,
			DROP CONSTRAINT IF EXISTS fk_tickets_usuarios_departamento,
			DROP CONSTRAINT IF EXISTS fk_tickets_usuarios_municipio;
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			ALTER TABLE municipios
			DROP CONSTRAINT IF EXISTS fk_municipios_departamento;
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			DROP INDEX IF EXISTS idx_tickets_usuarios_tipo_documento_id;
			DROP INDEX IF EXISTS idx_tickets_usuarios_regional_id;
			DROP INDEX IF EXISTS idx_tickets_usuarios_departamento_id;
			DROP INDEX IF EXISTS idx_tickets_usuarios_municipio_id;
			DROP INDEX IF EXISTS idx_municipios_departamento_id;
			DROP INDEX IF EXISTS idx_departamentos_codigo_dane;
			DROP INDEX IF EXISTS idx_municipios_codigo_dane;
			DROP INDEX IF EXISTS idx_tipo_documentos_nombre;
			DROP INDEX IF EXISTS idx_regionales_identificador;
		`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			ALTER TABLE tickets_usuarios
			DROP COLUMN IF EXISTS tipo_documento_id,
			DROP COLUMN IF EXISTS regional_id,
			DROP COLUMN IF EXISTS departamento_id,
			DROP COLUMN IF EXISTS municipio_id;
		`).Error; err != nil {
			return err
		}

		return tx.Migrator().DropTable("regionales", "tipo_documentos", "municipios", "departamentos")
	},
}
