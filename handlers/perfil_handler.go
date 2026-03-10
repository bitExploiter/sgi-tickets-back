package handlers

import (
	"sgi-tickets-back/models"
	"sgi-tickets-back/storage"
	"sgi-tickets-back/toolbox"

	"github.com/gofiber/fiber/v2"
)

// Struct para actualizar información personal
type UpdatePerfilRequest struct {
	Nombres         string `json:"nombres" validate:"required"`
	Apellidos       string `json:"apellidos" validate:"required"`
	TipoDocumento   string `json:"tipo_documento" validate:"-"`
	NumeroDocumento string `json:"numero_documento" validate:"-"`
	Telefono        string `json:"telefono" validate:"-"`
}

// Struct para cambiar contraseña
type ChangePasswordRequest struct {
	PasswordActual string `json:"password_actual" validate:"required"`
	NuevaPassword  string `json:"nueva_password" validate:"required,min=8"`
}

// Struct para confirmar código 2FA al activar
type Activar2FARequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// Struct para desactivar 2FA (requiere contraseña)
type Desactivar2FARequest struct {
	Password string `json:"password" validate:"required"`
}

// ==========================================
// GET /api/v1/perfil
// Obtener perfil del usuario autenticado
// ==========================================
func GetPerfil(c *fiber.Ctx) error {
	usuario := c.Locals("CurrentUser").(models.TicketUsuario)

	// Recargar usuario con relación Dependencia
	storage.DB.Preload("Dependencia").First(&usuario, usuario.Id)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":                usuario.Id,
			"nombres":           usuario.Nombres,
			"apellidos":         usuario.Apellidos,
			"tipo_documento":    usuario.TipoDocumento,
			"numero_documento":  usuario.NumeroDocumento,
			"email":             usuario.Email,
			"telefono":          usuario.Telefono,
			"rol":               usuario.Rol,
			"regional":          usuario.Regional,
			"municipio":         usuario.Municipio,
			"origen":            usuario.Origen,
			"dependencia_id":    usuario.DependenciaID,
			"dependencia":       usuario.Dependencia,
			"activo":            usuario.Activo,
			"totp_enabled":      usuario.TotpToken != "",
			"created_at":        usuario.CreatedAt,
		},
	})
}

// ==========================================
// PUT /api/v1/perfil
// Actualizar información personal del usuario
// ==========================================
func UpdatePerfil(c *fiber.Ctx) error {
	usuario := c.Locals("CurrentUser").(models.TicketUsuario)

	var request UpdatePerfilRequest

	// 1. Parse body
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Datos invalidos",
		})
	}

	// 2. Validar campos
	if errors, err := toolbox.FormatValidationErrors(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errors,
		})
	}

	// 3. Actualizar campos permitidos
	storage.DB.Model(&usuario).Updates(map[string]interface{}{
		"nombres":          request.Nombres,
		"apellidos":        request.Apellidos,
		"tipo_documento":   request.TipoDocumento,
		"numero_documento": request.NumeroDocumento,
		"telefono":         request.Telefono,
	})

	// 4. Log de acción
	toolbox.SaveLoggerAction(usuario, "Perfil", "perfil_actualizado", c.IP())

	// 5. Recargar usuario actualizado
	storage.DB.Preload("Dependencia").First(&usuario, usuario.Id)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":                usuario.Id,
			"nombres":           usuario.Nombres,
			"apellidos":         usuario.Apellidos,
			"tipo_documento":    usuario.TipoDocumento,
			"numero_documento":  usuario.NumeroDocumento,
			"email":             usuario.Email,
			"telefono":          usuario.Telefono,
			"rol":               usuario.Rol,
			"regional":          usuario.Regional,
			"municipio":         usuario.Municipio,
			"dependencia_id":    usuario.DependenciaID,
			"dependencia":       usuario.Dependencia,
			"activo":            usuario.Activo,
			"totp_enabled":      usuario.TotpToken != "",
			"created_at":        usuario.CreatedAt,
		},
	})
}

// ==========================================
// PUT /api/v1/perfil/password
// Cambiar contraseña del usuario autenticado
// ==========================================
func ChangePassword(c *fiber.Ctx) error {
	usuario := c.Locals("CurrentUser").(models.TicketUsuario)

	var request ChangePasswordRequest

	// 1. Parse body
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Datos invalidos",
		})
	}

	// 2. Validar campos
	if errors, err := toolbox.FormatValidationErrors(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errors,
		})
	}

	// 3. Recargar usuario para obtener password hash actual
	storage.DB.First(&usuario, usuario.Id)

	// 4. Verificar contraseña actual
	if !toolbox.CheckPasswordHash(request.PasswordActual, usuario.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Contraseña actual incorrecta",
		})
	}

	// 5. Hash de la nueva contraseña
	hashedPassword, err := toolbox.HashPassword(request.NuevaPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error procesando contraseña",
		})
	}

	// 6. Actualizar contraseña
	storage.DB.Model(&usuario).Update("password", hashedPassword)

	// 7. Log de acción
	toolbox.SaveLoggerAction(usuario, "Perfil", "password_cambiado", c.IP())

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Contraseña actualizada exitosamente",
	})
}

// ==========================================
// GET /api/v1/perfil/2fa/activar  (generar QR)
// POST /api/v1/perfil/2fa/activar (confirmar código)
// Activar autenticación de dos factores
// ==========================================
func Activar2FA(c *fiber.Ctx) error {
	usuario := c.Locals("CurrentUser").(models.TicketUsuario)

	// GET: Generar secreto y devolver QR
	if c.Method() == "GET" {
		// Generar secreto TOTP
		secret, err := toolbox.GenerateTOTPSecret(usuario.Email, "SGI Tickets")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Error generando codigo 2FA",
			})
		}

		// Generar URL y QR
		totpURL := toolbox.GetTOTPURL(secret, usuario.Email, "SGI Tickets")
		qrBase64, err := toolbox.GenerateQRCodeBase64(totpURL)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Error generando codigo QR",
			})
		}

		// Guardar secreto temporalmente en DB
		storage.DB.Model(&usuario).Update("totp_token", secret)

		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"secret":  secret,
				"qr_code": qrBase64,
				"issuer":  "SGI Tickets",
				"account": usuario.Email,
			},
		})
	}

	// POST: Verificar código y confirmar activación
	if c.Method() == "POST" {
		var request Activar2FARequest

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Datos invalidos",
			})
		}

		if errors, err := toolbox.FormatValidationErrors(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"errors":  errors,
			})
		}

		// Recargar usuario para obtener el secreto guardado
		storage.DB.First(&usuario, usuario.Id)

		if usuario.TotpToken == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "No hay configuracion 2FA pendiente",
			})
		}

		// Validar código
		if !toolbox.ValidateTOTPCode(request.Code, usuario.TotpToken) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Codigo incorrecto",
			})
		}

		// Log de activación exitosa
		toolbox.SaveLoggerAction(usuario, "Perfil", "2fa_activado", c.IP())

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Autenticacion de dos factores activada exitosamente",
		})
	}

	return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
		"success": false,
		"error":   "Metodo no permitido",
	})
}

// ==========================================
// POST /api/v1/perfil/2fa/desactivar
// Desactivar autenticación de dos factores
// ==========================================
func Desactivar2FA(c *fiber.Ctx) error {
	usuario := c.Locals("CurrentUser").(models.TicketUsuario)

	var request Desactivar2FARequest

	// 1. Parse body
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Datos invalidos",
		})
	}

	// 2. Validar campos
	if errors, err := toolbox.FormatValidationErrors(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  errors,
		})
	}

	// 3. Recargar usuario para obtener password hash
	storage.DB.First(&usuario, usuario.Id)

	// 4. Verificar contraseña
	if !toolbox.CheckPasswordHash(request.Password, usuario.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Contraseña incorrecta",
		})
	}

	// 5. Verificar que tenga 2FA activo
	if usuario.TotpToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "La autenticacion de dos factores no esta activa",
		})
	}

	// 6. Limpiar token TOTP
	storage.DB.Model(&usuario).Update("totp_token", "")

	// 7. Log de acción
	toolbox.SaveLoggerAction(usuario, "Perfil", "2fa_desactivado", c.IP())

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Autenticacion de dos factores desactivada exitosamente",
	})
}
