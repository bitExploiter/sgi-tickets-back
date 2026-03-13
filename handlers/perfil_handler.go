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

// GetPerfil godoc
// @Summary Obtener perfil del usuario autenticado
// @Description Devuelve la información completa del perfil del usuario autenticado, incluyendo sus datos personales, rol y dependencia
// @Tags Perfil
// @Produce json
// @Success 200 {object} map[string]interface{} "Datos del perfil del usuario"
// @Failure 401 {object} map[string]interface{} "No autenticado"
// @Router /perfil [get]
// @Security CookieAuth
// @Security TwoFactorAuth
func GetPerfil(c *fiber.Ctx) error {
	usuario := c.Locals("CurrentUser").(models.TicketUsuario)

	// Recargar usuario con relación Dependencia
	storage.DB.
		Preload("Dependencia").
		Preload("TipoDocumentoRef").
		Preload("RegionalRef").
		Preload("DepartamentoRef").
		Preload("MunicipioRef").
		First(&usuario, usuario.Id)

	regionalNombre := ""
	if usuario.RegionalRef != nil {
		regionalNombre = usuario.RegionalRef.Nombre
	}

	municipioNombre := ""
	if usuario.MunicipioRef != nil {
		municipioNombre = usuario.MunicipioRef.Nombre
	}

	tipoDocumentoNombre := ""
	if usuario.TipoDocumentoRef != nil {
		tipoDocumentoNombre = usuario.TipoDocumentoRef.Nombre
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":                    usuario.Id,
			"nombres":               usuario.Nombres,
			"apellidos":             usuario.Apellidos,
			"tipo_documento":        tipoDocumentoNombre,
			"tipo_documento_id":     usuario.TipoDocumentoID,
			"numero_documento":      usuario.NumeroDocumento,
			"email":                 usuario.Email,
			"telefono":              usuario.Telefono,
			"rol":                   usuario.Rol,
			"regional":              regionalNombre,
			"regional_id":           usuario.RegionalID,
			"departamento_id":       usuario.DepartamentoID,
			"municipio":             municipioNombre,
			"municipio_id":          usuario.MunicipioID,
			"origen":                usuario.Origen,
			"dependencia_id":        usuario.DependenciaID,
			"dependencia":           usuario.Dependencia,
			"activo":                usuario.Activo,
			"totp_enabled":          usuario.TotpToken != "",
			"ultima_fecha_conexion": usuario.UltimaFechaConexion,
			"created_at":            usuario.CreatedAt,
		},
	})
}

// UpdatePerfil godoc
// @Summary Actualizar información personal del usuario
// @Description Permite al usuario autenticado actualizar su información personal (nombres, apellidos, documento, teléfono)
// @Tags Perfil
// @Accept json
// @Produce json
// @Param body body UpdatePerfilRequest true "Datos a actualizar"
// @Success 200 {object} map[string]interface{} "Perfil actualizado exitosamente"
// @Failure 400 {object} map[string]interface{} "Datos inválidos o errores de validación"
// @Failure 401 {object} map[string]interface{} "No autenticado"
// @Router /perfil [put]
// @Security CookieAuth
// @Security TwoFactorAuth
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
		"numero_documento": request.NumeroDocumento,
		"telefono":         request.Telefono,
	})

	// 4. Log de acción
	toolbox.SaveLoggerAction(usuario, "Perfil", "perfil_actualizado", c.IP())

	// 5. Recargar usuario actualizado
	storage.DB.
		Preload("Dependencia").
		Preload("TipoDocumentoRef").
		Preload("RegionalRef").
		Preload("DepartamentoRef").
		Preload("MunicipioRef").
		First(&usuario, usuario.Id)

	regionalNombre := ""
	if usuario.RegionalRef != nil {
		regionalNombre = usuario.RegionalRef.Nombre
	}

	municipioNombre := ""
	if usuario.MunicipioRef != nil {
		municipioNombre = usuario.MunicipioRef.Nombre
	}

	tipoDocumentoNombre := ""
	if usuario.TipoDocumentoRef != nil {
		tipoDocumentoNombre = usuario.TipoDocumentoRef.Nombre
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":                    usuario.Id,
			"nombres":               usuario.Nombres,
			"apellidos":             usuario.Apellidos,
			"tipo_documento":        tipoDocumentoNombre,
			"tipo_documento_id":     usuario.TipoDocumentoID,
			"numero_documento":      usuario.NumeroDocumento,
			"email":                 usuario.Email,
			"telefono":              usuario.Telefono,
			"rol":                   usuario.Rol,
			"regional":              regionalNombre,
			"regional_id":           usuario.RegionalID,
			"departamento_id":       usuario.DepartamentoID,
			"municipio":             municipioNombre,
			"municipio_id":          usuario.MunicipioID,
			"dependencia_id":        usuario.DependenciaID,
			"dependencia":           usuario.Dependencia,
			"activo":                usuario.Activo,
			"totp_enabled":          usuario.TotpToken != "",
			"ultima_fecha_conexion": usuario.UltimaFechaConexion,
			"created_at":            usuario.CreatedAt,
		},
	})
}

// ChangePassword godoc
// @Summary Cambiar contraseña del usuario autenticado
// @Description Permite al usuario cambiar su contraseña verificando primero la contraseña actual. La nueva contraseña debe tener al menos 8 caracteres
// @Tags Perfil
// @Accept json
// @Produce json
// @Param body body ChangePasswordRequest true "Contraseña actual y nueva contraseña"
// @Success 200 {object} map[string]interface{} "Contraseña actualizada exitosamente"
// @Failure 400 {object} map[string]interface{} "Datos inválidos o errores de validación"
// @Failure 401 {object} map[string]interface{} "Contraseña actual incorrecta"
// @Router /perfil/password [put]
// @Security CookieAuth
// @Security TwoFactorAuth
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

// Activar2FA godoc
// @Summary Activar 2FA en el perfil
// @Description GET: Genera un código QR para activar 2FA. POST: Verifica el código y activa 2FA en el perfil del usuario
// @Tags Perfil - Seguridad
// @Accept json
// @Produce json
// @Param body body Activar2FARequest false "Código de verificación (solo para POST)"
// @Success 200 {object} map[string]interface{} "QR generado o 2FA activado"
// @Failure 400 {object} map[string]interface{} "Código incorrecto o no hay configuración pendiente"
// @Router /perfil/2fa/activar [get]
// @Router /perfil/2fa/activar [post]
// @Security CookieAuth
// @Security TwoFactorAuth
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

// Desactivar2FA godoc
// @Summary Desactivar 2FA en el perfil
// @Description Desactiva la autenticación de dos factores del usuario. Requiere contraseña para confirmar
// @Tags Perfil - Seguridad
// @Accept json
// @Produce json
// @Param body body Desactivar2FARequest true "Contraseña del usuario"
// @Success 200 {object} map[string]interface{} "2FA desactivado exitosamente"
// @Failure 400 {object} map[string]interface{} "2FA no está activo"
// @Failure 401 {object} map[string]interface{} "Contraseña incorrecta"
// @Router /perfil/2fa/desactivar [post]
// @Security CookieAuth
// @Security TwoFactorAuth
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
