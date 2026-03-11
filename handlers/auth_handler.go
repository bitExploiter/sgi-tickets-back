package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"sgi-tickets-back/models"
	"sgi-tickets-back/storage"
	"sgi-tickets-back/toolbox"

	"github.com/gofiber/fiber/v2"
)

// Struct para parsear body de login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Struct para parsear body de 2FA
type Verify2FARequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// Struct para parsear body de setup 2FA
type Setup2FARequest struct {
	Code string `json:"code" validate:"required,len=6"`
}

// Struct para parsear body de recover
type RecoverRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Struct para parsear body de reset
type ResetRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// Login godoc
// @Summary Iniciar sesión
// @Description Autentica un usuario con email y contraseña. Si el rol requiere 2FA, devuelve el estado de configuración 2FA
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Credenciales de acceso"
// @Success 200 {object} map[string]interface{} "Login exitoso - puede requerir 2FA"
// @Failure 400 {object} map[string]interface{} "Datos inválidos"
// @Failure 401 {object} map[string]interface{} "Credenciales incorrectas"
// @Failure 403 {object} map[string]interface{} "Usuario deshabilitado"
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	var request LoginRequest

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

	// 3. Buscar usuario por email
	var usuario models.TicketUsuario
	storage.DB.Where("email = ?", request.Email).First(&usuario)

	if usuario.Id == 0 {
		// Usuario no encontrado (mensaje genérico por seguridad)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Credenciales incorrectas",
		})
	}

	// 4. Verificar que el usuario esté activo
	if !usuario.Activo {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "Usuario deshabilitado",
		})
	}

	// 5. Verificar password
	if !toolbox.CheckPasswordHash(request.Password, usuario.Password) {
		// Log de intento fallido
		toolbox.SaveLoggerAction(usuario, "Auth", "login_fallido", c.IP())

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Credenciales incorrectas",
		})
	}

	// 6. Generar cookie de sesión básica (sgi_tickets_user_email)
	sessionToken, err := toolbox.GenerateSessionToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error generando sesion",
		})
	}

	// Guardar cookie en DB
	toolbox.SaveCookieToStorage(sessionToken, usuario.Email)

	// Setear cookie HTTP (encriptada por middleware)
	c.Cookie(&fiber.Cookie{
		Name:     "sgi_tickets_user_email",
		Value:    sessionToken,
		MaxAge:   86400 * 7, // 7 días
		HTTPOnly: true,
		Secure:   false, // true en producción
		SameSite: "Strict",
	})

	// 7. Verificar si requiere 2FA
	require2FA := toolbox.Require2FA(usuario.Rol)
	hasTOTP := usuario.TotpToken != ""

	// Log de login exitoso (nivel 1)
	toolbox.SaveLoggerAction(usuario, "Auth", "login_exitoso", c.IP())

	// 8. Devolver respuesta según estado de 2FA
	if !require2FA {
		// No requiere 2FA - generar cookie completa directamente
		identityToken, _ := toolbox.GenerateSessionToken()
		toolbox.SaveCookieToStorage(identityToken, usuario.Email)

		c.Cookie(&fiber.Cookie{
			Name:     "sgi_tickets_identity",
			Value:    identityToken,
			MaxAge:   86400 * 7,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Strict",
		})

		return c.JSON(fiber.Map{
			"success":     true,
			"require_2fa": false,
			"data": fiber.Map{
				"id":      usuario.Id,
				"email":   usuario.Email,
				"nombres": usuario.Nombres,
				"rol":     usuario.Rol,
			},
		})
	}

	if !hasTOTP {
		// Requiere 2FA pero no tiene TOTP configurado
		return c.JSON(fiber.Map{
			"success":      true,
			"require_2fa":  true,
			"totp_enabled": false,
			"message":      "Debes configurar autenticacion de dos factores",
		})
	}

	// Requiere 2FA y tiene TOTP configurado
	return c.JSON(fiber.Map{
		"success":      true,
		"require_2fa":  true,
		"totp_enabled": true,
		"message":      "Ingresa el codigo de tu aplicacion de autenticacion",
	})
}

// Setup2FA godoc
// @Summary Configurar autenticación de dos factores
// @Description GET: Genera un código QR para configurar 2FA. POST: Verifica el código y activa 2FA
// @Tags Autenticación - 2FA
// @Accept json
// @Produce json
// @Param body body Setup2FARequest false "Código de verificación (solo para POST)"
// @Success 200 {object} map[string]interface{} "QR generado o 2FA activado exitosamente"
// @Failure 400 {object} map[string]interface{} "Código incorrecto o no hay configuración pendiente"
// @Failure 401 {object} map[string]interface{} "No autenticado"
// @Router /auth/2fa/setup [get]
// @Router /auth/2fa/setup [post]
// @Security CookieAuth
func Setup2FA(c *fiber.Ctx) error {
	// Obtener usuario de la sesión (via CookieMiddleware)
	usuario := c.Locals("CurrentUser").(models.TicketUsuario)

	// Si el método es GET, generar nuevo secreto y devolver QR
	if c.Method() == "GET" {
		// Generar secreto TOTP
		secret, err := toolbox.GenerateTOTPSecret(usuario.Email, "SGI Tickets")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Error generando codigo 2FA",
			})
		}

		// Generar URL TOTP
		totpURL := toolbox.GetTOTPURL(secret, usuario.Email, "SGI Tickets")

		// Generar QR code en base64
		qrBase64, err := toolbox.GenerateQRCodeBase64(totpURL)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Error generando codigo QR",
			})
		}

		// Guardar secreto en DB
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

	// Si el método es POST, verificar código y confirmar setup
	if c.Method() == "POST" {
		var request Setup2FARequest

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
				"error":   "No hay configuracion pendiente",
			})
		}

		// Validar código
		if !toolbox.ValidateTOTPCode(request.Code, usuario.TotpToken) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Codigo incorrecto",
			})
		}

		// Código correcto - 2FA configurado exitosamente
		// Generar cookie de identidad completa
		identityToken, _ := toolbox.GenerateSessionToken()
		toolbox.SaveCookieToStorage(identityToken, usuario.Email)

		c.Cookie(&fiber.Cookie{
			Name:     "sgi_tickets_identity",
			Value:    identityToken,
			MaxAge:   86400 * 7,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Strict",
		})

		toolbox.SaveLoggerAction(usuario, "Auth", "2fa_configurado", c.IP())

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Autenticacion de dos factores configurada exitosamente",
			"data": fiber.Map{
				"id":      usuario.Id,
				"email":   usuario.Email,
				"nombres": usuario.Nombres,
				"rol":     usuario.Rol,
			},
		})
	}

	return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
		"success": false,
		"error":   "Metodo no permitido",
	})
}

// Verify2FA godoc
// @Summary Verificar código 2FA
// @Description Verifica el código TOTP generado por la aplicación de autenticación y completa el inicio de sesión
// @Tags Autenticación - 2FA
// @Accept json
// @Produce json
// @Param body body Verify2FARequest true "Código de verificación TOTP"
// @Success 200 {object} map[string]interface{} "Código correcto - sesión completa establecida"
// @Failure 400 {object} map[string]interface{} "Datos inválidos o 2FA no configurado"
// @Failure 401 {object} map[string]interface{} "Código incorrecto"
// @Router /auth/2fa/verify [post]
// @Security CookieAuth
func Verify2FA(c *fiber.Ctx) error {
	var request Verify2FARequest

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

	// 3. Obtener usuario de la sesión (via CookieMiddleware)
	usuario := c.Locals("CurrentUser").(models.TicketUsuario)

	// 4. Verificar que tenga TOTP configurado
	if usuario.TotpToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "2FA no configurado",
		})
	}

	// 5. Validar código TOTP
	if !toolbox.ValidateTOTPCode(request.Code, usuario.TotpToken) {
		toolbox.SaveLoggerAction(usuario, "Auth", "2fa_fallido", c.IP())

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Codigo incorrecto",
		})
	}

	// 6. Código correcto - generar cookie de identidad completa
	identityToken, err := toolbox.GenerateSessionToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error generando sesion",
		})
	}

	toolbox.SaveCookieToStorage(identityToken, usuario.Email)

	c.Cookie(&fiber.Cookie{
		Name:     "sgi_tickets_identity",
		Value:    identityToken,
		MaxAge:   86400 * 7, // 7 días
		HTTPOnly: true,
		Secure:   false, // true en producción
		SameSite: "Strict",
	})

	// Log de 2FA exitoso
	toolbox.SaveLoggerAction(usuario, "Auth", "2fa_exitoso", c.IP())

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":      usuario.Id,
			"email":   usuario.Email,
			"nombres": usuario.Nombres,
			"rol":     usuario.Rol,
		},
	})
}

// Logout godoc
// @Summary Cerrar sesión
// @Description Invalida las cookies de sesión y cierra la sesión del usuario
// @Tags Autenticación
// @Produce json
// @Success 200 {object} map[string]interface{} "Sesión cerrada exitosamente"
// @Router /auth/logout [post]
// @Security CookieAuth
func Logout(c *fiber.Ctx) error {
	// Obtener usuario de la sesión
	usuario := c.Locals("CurrentUser").(models.TicketUsuario)

	// Obtener tokens de las cookies
	emailToken := c.Cookies("sgi_tickets_user_email")
	identityToken := c.Cookies("sgi_tickets_identity")

	// Deshabilitar cookies en BD
	if emailToken != "" {
		toolbox.DisableCookie(emailToken)
	}
	if identityToken != "" {
		toolbox.DisableCookie(identityToken)
	}

	// Limpiar cookies del navegador
	c.Cookie(&fiber.Cookie{
		Name:     "sgi_tickets_user_email",
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "sgi_tickets_identity",
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
	})

	// Log de logout
	toolbox.SaveLoggerAction(usuario, "Auth", "logout", c.IP())

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Sesion cerrada exitosamente",
	})
}

// RecoverPassword godoc
// @Summary Solicitar recuperación de contraseña
// @Description Envía un email con un enlace para restablecer la contraseña si el usuario existe
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param body body RecoverRequest true "Email del usuario"
// @Success 200 {object} map[string]interface{} "Email enviado (si el usuario existe)"
// @Failure 400 {object} map[string]interface{} "Datos inválidos"
// @Router /auth/recover [post]
func RecoverPassword(c *fiber.Ctx) error {
	var request RecoverRequest

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

	// 3. Buscar usuario por email
	var usuario models.TicketUsuario
	storage.DB.Where("email = ?", request.Email).First(&usuario)

	// SEGURIDAD: Siempre devolver éxito (no revelar si el email existe)
	// Pero solo enviar email si el usuario existe

	if usuario.Id != 0 && usuario.Activo {
		// 4. Generar token de recuperación
		tokenPlain, tokenHash, err := toolbox.GenerateResetToken()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Error generando token",
			})
		}

		// 5. Guardar hash del token y expiración (1 hora) en DB
		expiry := time.Now().Add(1 * time.Hour)
		storage.DB.Model(&usuario).Updates(map[string]interface{}{
			"reset_token":        tokenHash,
			"reset_token_expiry": expiry,
		})

		// 6. Enviar email con el token
		resetURL := "http://localhost:5173/reset-password?token=" + tokenPlain

		emailData := fiber.Map{
			"nombre":    usuario.Nombres + " " + usuario.Apellidos,
			"reset_url": resetURL,
		}

		toolbox.SendNotificacionEmail(
			usuario.Nombres,
			usuario.Email,
			"Recuperar contraseña - SGI Tickets",
			"./templates/emails/recover.html",
			emailData,
		)

		// Log de recuperación solicitada
		toolbox.SaveLoggerAction(usuario, "Auth", "recuperacion_solicitada", c.IP())
	}

	// Siempre devolver éxito (no revelar si el email existe)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Si el email existe, recibiras instrucciones de recuperacion",
	})
}

// ResetPassword godoc
// @Summary Restablecer contraseña
// @Description Cambia la contraseña del usuario usando el token recibido por email
// @Tags Autenticación
// @Accept json
// @Produce json
// @Param body body ResetRequest true "Token y nueva contraseña"
// @Success 200 {object} map[string]interface{} "Contraseña actualizada exitosamente"
// @Failure 400 {object} map[string]interface{} "Token inválido, expirado o datos inválidos"
// @Router /auth/reset [post]
func ResetPassword(c *fiber.Ctx) error {
	var request ResetRequest

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

	// 3. Buscar usuario con el token (hash)
	// Calcular hash del token recibido
	hash := sha256.Sum256([]byte(request.Token))
	hashString := hex.EncodeToString(hash[:])

	var usuario models.TicketUsuario
	storage.DB.Where("reset_token = ?", hashString).First(&usuario)

	if usuario.Id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Token invalido o expirado",
		})
	}

	// 4. Validar token (verificar expiración)
	// Nota: reset_token_expiry puede ser null si la migración aún no se ha ejecutado
	var expiry time.Time
	storage.DB.Model(&models.TicketUsuario{}).Select("reset_token_expiry").Where("id = ?", usuario.Id).Scan(&expiry)

	if !expiry.IsZero() && !toolbox.ValidateResetToken(request.Token, usuario.ResetToken, expiry) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Token invalido o expirado",
		})
	}

	// 5. Hash de la nueva contraseña
	hashedPassword, err := toolbox.HashPassword(request.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error procesando contraseña",
		})
	}

	// 6. Actualizar contraseña y limpiar token
	storage.DB.Model(&usuario).Updates(map[string]interface{}{
		"password":           hashedPassword,
		"reset_token":        "",
		"reset_token_expiry": nil,
	})

	// 7. Deshabilitar todas las sesiones activas del usuario (seguridad)
	toolbox.DisableAllCookies(usuario.Email)

	// Log de contraseña cambiada
	toolbox.SaveLoggerAction(usuario, "Auth", "password_cambiado", c.IP())

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Contraseña actualizada exitosamente",
	})
}
