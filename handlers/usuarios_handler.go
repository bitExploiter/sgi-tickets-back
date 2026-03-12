package handlers

import (
	"fmt"
	"os"
	"sgi-tickets-back/models"
	"sgi-tickets-back/storage"
	"sgi-tickets-back/toolbox"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Structs para requests
type CreateUsuarioRequest struct {
	Nombres         string `json:"nombres" validate:"required"`
	Apellidos       string `json:"apellidos" validate:"required"`
	TipoDocumento   string `json:"tipo_documento" validate:"-"`
	NumeroDocumento string `json:"numero_documento" validate:"-"`
	Email           string `json:"email" validate:"required,email"`
	Telefono        string `json:"telefono" validate:"-"`
	Regional        string `json:"regional" validate:"-"`
	Municipio       string `json:"municipio" validate:"-"`
	Rol             string `json:"rol" validate:"required"`
	DependenciaID   *uint  `json:"dependencia_id" validate:"-"`
}

type UpdateUsuarioRequest struct {
	Nombres         string `json:"nombres" validate:"required"`
	Apellidos       string `json:"apellidos" validate:"required"`
	TipoDocumento   string `json:"tipo_documento" validate:"-"`
	NumeroDocumento string `json:"numero_documento" validate:"-"`
	Telefono        string `json:"telefono" validate:"-"`
	Regional        string `json:"regional" validate:"-"`
	Municipio       string `json:"municipio" validate:"-"`
	Rol             string `json:"rol" validate:"required"`
	DependenciaID   *uint  `json:"dependencia_id" validate:"-"`
	Activo          *bool  `json:"activo" validate:"-"`
}

// ListarUsuarios godoc
// @Summary Listar usuarios con paginación y filtros
// @Description Obtiene un listado paginado de usuarios del sistema con opciones de filtrado por búsqueda, rol, estado y regional
// @Tags Usuarios
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param page_size query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por nombre, apellido o email"
// @Param rol query string false "Filtrar por rol (admin, supervisor, agente, entidad, contratista)"
// @Param estado query string false "Filtrar por estado (activo, inactivo)"
// @Param regional query string false "Filtrar por regional"
// @Success 200 {object} map[string]interface{} "Lista paginada de usuarios"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /usuarios [get]
// @Security CookieAuth
// @Security TwoFactorAuth
func ListarUsuarios(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	search := c.Query("search", "")
	rol := c.Query("rol", "")
	estado := c.Query("estado", "")
	regional := c.Query("regional", "")

	offset := (page - 1) * pageSize

	// Query base
	query := storage.DB.Model(&models.TicketUsuario{}).Preload("Dependencia")

	// Aplicar filtro de búsqueda (nombre, apellido, email, ID)
	if search != "" {
		// Intentar convertir a ID si es numérico
		if id, err := strconv.Atoi(search); err == nil {
			query = query.Where("id = ?", id)
		} else {
			query = query.Where(
				"LOWER(nombres) LIKE ? OR LOWER(apellidos) LIKE ? OR LOWER(email) LIKE ?",
				"%"+strings.ToLower(search)+"%",
				"%"+strings.ToLower(search)+"%",
				"%"+strings.ToLower(search)+"%",
			)
		}
	}

	// Aplicar filtro de rol
	if rol != "" && rol != "todos" {
		query = query.Where("rol = ?", rol)
	}

	// Aplicar filtro de estado (activo/inactivo)
	if estado != "" && estado != "todos" {
		if estado == "activo" {
			query = query.Where("activo = ?", true)
		} else if estado == "inactivo" {
			query = query.Where("activo = ?", false)
		}
	}

	// Aplicar filtro de regional
	if regional != "" && regional != "todas" {
		query = query.Where("regional = ?", regional)
	}

	// Contar total de registros
	var totalRows int64
	query.Count(&totalRows)

	// Obtener usuarios paginados
	var usuarios []models.TicketUsuario
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&usuarios)

	// Formatear respuesta
	usuariosResponse := make([]fiber.Map, len(usuarios))
	for i, usuario := range usuarios {
		usuariosResponse[i] = fiber.Map{
			"id":                    usuario.Id,
			"nombres":               usuario.Nombres,
			"apellidos":             usuario.Apellidos,
			"tipo_documento":        usuario.TipoDocumento,
			"numero_documento":      usuario.NumeroDocumento,
			"email":                 usuario.Email,
			"telefono":              usuario.Telefono,
			"rol":                   usuario.Rol,
			"regional":              usuario.Regional,
			"municipio":             usuario.Municipio,
			"origen":                usuario.Origen,
			"dependencia_id":        usuario.DependenciaID,
			"dependencia":           usuario.Dependencia,
			"activo":                usuario.Activo,
			"totp_enabled":          usuario.TotpToken != "",
			"ultima_fecha_conexion": usuario.UltimaFechaConexion,
			"created_at":            usuario.CreatedAt,
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"usuarios":    usuariosResponse,
			"page":        page,
			"page_size":   pageSize,
			"total_rows":  totalRows,
			"total_pages": (totalRows + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// ObtenerUsuario godoc
// @Summary Obtener un usuario por ID
// @Description Obtiene los datos completos de un usuario específico por su ID
// @Tags Usuarios
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} map[string]interface{} "Datos del usuario"
// @Failure 404 {object} map[string]interface{} "Usuario no encontrado"
// @Router /usuarios/{id} [get]
// @Security CookieAuth
// @Security TwoFactorAuth
func ObtenerUsuario(c *fiber.Ctx) error {
	id := c.Params("id")

	var usuario models.TicketUsuario
	result := storage.DB.Preload("Dependencia").First(&usuario, id)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Usuario no encontrado",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":                    usuario.Id,
			"nombres":               usuario.Nombres,
			"apellidos":             usuario.Apellidos,
			"tipo_documento":        usuario.TipoDocumento,
			"numero_documento":      usuario.NumeroDocumento,
			"email":                 usuario.Email,
			"telefono":              usuario.Telefono,
			"rol":                   usuario.Rol,
			"regional":              usuario.Regional,
			"municipio":             usuario.Municipio,
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

// CrearUsuario godoc
// @Summary Crear un nuevo usuario
// @Description Crea un nuevo usuario en el sistema. Se genera una contraseña aleatoria y se envía por email
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param body body CreateUsuarioRequest true "Datos del nuevo usuario"
// @Success 201 {object} map[string]interface{} "Usuario creado exitosamente"
// @Failure 400 {object} map[string]interface{} "Datos inválidos o email duplicado"
// @Router /usuarios [post]
// @Security CookieAuth
// @Security TwoFactorAuth
func CrearUsuario(c *fiber.Ctx) error {
	currentUser := c.Locals("CurrentUser").(models.TicketUsuario)

	var request CreateUsuarioRequest

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

	// 3. Validar que el email no exista
	var existingUser models.TicketUsuario
	if storage.DB.Where("email = ?", request.Email).First(&existingUser).Error == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "El email ya esta registrado",
		})
	}

	// 4. Generar contraseña aleatoria
	password := toolbox.GenerateRandomPassword(12)
	hashedPassword, err := toolbox.HashPassword(password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error generando contraseña",
		})
	}

	// 5. Crear usuario
	usuario := models.TicketUsuario{
		Nombres:         request.Nombres,
		Apellidos:       request.Apellidos,
		TipoDocumento:   request.TipoDocumento,
		NumeroDocumento: request.NumeroDocumento,
		Email:           request.Email,
		Telefono:        request.Telefono,
		Regional:        request.Regional,
		Municipio:       request.Municipio,
		Password:        hashedPassword,
		Rol:             request.Rol,
		Origen:          "local",
		DependenciaID:   request.DependenciaID,
		Activo:          true,
	}

	if err := storage.DB.Create(&usuario).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error creando usuario",
		})
	}

	// 6. Enviar email con credenciales
	go toolbox.SendEmailAsync(
		usuario.Email,
		"Bienvenido a SGI Tickets",
		fmt.Sprintf(`
			<h2>Bienvenido a SGI Tickets</h2>
			<p>Hola %s %s,</p>
			<p>Tu cuenta ha sido creada exitosamente.</p>
			<p><strong>Email:</strong> %s</p>
			<p><strong>Contraseña temporal:</strong> %s</p>
			<p>Por favor, inicia sesión y cambia tu contraseña en tu perfil.</p>
			<p><a href="%s">Ir a SGI Tickets</a></p>
		`,
			usuario.Nombres,
			usuario.Apellidos,
			usuario.Email,
			password,
			os.Getenv("FRONTEND_URL"),
		),
	)

	// 7. Log de acción
	toolbox.SaveLoggerAction(currentUser, "Usuario", fmt.Sprintf("crear_usuario_%d", usuario.Id), c.IP())

	// 8. Recargar usuario con dependencia
	storage.DB.Preload("Dependencia").First(&usuario, usuario.Id)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":                    usuario.Id,
			"nombres":               usuario.Nombres,
			"apellidos":             usuario.Apellidos,
			"tipo_documento":        usuario.TipoDocumento,
			"numero_documento":      usuario.NumeroDocumento,
			"email":                 usuario.Email,
			"telefono":              usuario.Telefono,
			"rol":                   usuario.Rol,
			"regional":              usuario.Regional,
			"municipio":             usuario.Municipio,
			"dependencia_id":        usuario.DependenciaID,
			"dependencia":           usuario.Dependencia,
			"activo":                usuario.Activo,
			"totp_enabled":          false,
			"ultima_fecha_conexion": usuario.UltimaFechaConexion,
			"created_at":            usuario.CreatedAt,
		},
	})
}

// ActualizarUsuario godoc
// @Summary Actualizar un usuario
// @Description Actualiza los datos de un usuario existente. No permite actualizar el email ni la contraseña
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param body body UpdateUsuarioRequest true "Datos a actualizar"
// @Success 200 {object} map[string]interface{} "Usuario actualizado exitosamente"
// @Failure 400 {object} map[string]interface{} "Datos inválidos"
// @Failure 404 {object} map[string]interface{} "Usuario no encontrado"
// @Router /usuarios/{id} [put]
// @Security CookieAuth
// @Security TwoFactorAuth
func ActualizarUsuario(c *fiber.Ctx) error {
	currentUser := c.Locals("CurrentUser").(models.TicketUsuario)
	id := c.Params("id")

	var request UpdateUsuarioRequest

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

	// 3. Buscar usuario
	var usuario models.TicketUsuario
	if err := storage.DB.First(&usuario, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Usuario no encontrado",
		})
	}

	// 4. Actualizar campos permitidos
	updates := map[string]interface{}{
		"nombres":          request.Nombres,
		"apellidos":        request.Apellidos,
		"tipo_documento":   request.TipoDocumento,
		"numero_documento": request.NumeroDocumento,
		"telefono":         request.Telefono,
		"regional":         request.Regional,
		"municipio":        request.Municipio,
		"rol":              request.Rol,
		"dependencia_id":   request.DependenciaID,
	}

	// Actualizar estado si se proporciona
	if request.Activo != nil {
		updates["activo"] = *request.Activo
	}

	storage.DB.Model(&usuario).Updates(updates)

	// 5. Log de acción
	toolbox.SaveLoggerAction(currentUser, "Usuario", fmt.Sprintf("actualizar_usuario_%s", id), c.IP())

	// 6. Recargar usuario actualizado
	storage.DB.Preload("Dependencia").First(&usuario, id)

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id":                    usuario.Id,
			"nombres":               usuario.Nombres,
			"apellidos":             usuario.Apellidos,
			"tipo_documento":        usuario.TipoDocumento,
			"numero_documento":      usuario.NumeroDocumento,
			"email":                 usuario.Email,
			"telefono":              usuario.Telefono,
			"rol":                   usuario.Rol,
			"regional":              usuario.Regional,
			"municipio":             usuario.Municipio,
			"dependencia_id":        usuario.DependenciaID,
			"dependencia":           usuario.Dependencia,
			"activo":                usuario.Activo,
			"totp_enabled":          usuario.TotpToken != "",
			"ultima_fecha_conexion": usuario.UltimaFechaConexion,
			"created_at":            usuario.CreatedAt,
		},
	})
}

// EliminarUsuario godoc
// @Summary Eliminar un usuario (soft delete)
// @Description Elimina lógicamente un usuario del sistema. El usuario no se borra físicamente, solo se marca como eliminado
// @Tags Usuarios
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} map[string]interface{} "Usuario eliminado exitosamente"
// @Failure 404 {object} map[string]interface{} "Usuario no encontrado"
// @Router /usuarios/{id} [delete]
// @Security CookieAuth
// @Security TwoFactorAuth
func EliminarUsuario(c *fiber.Ctx) error {
	currentUser := c.Locals("CurrentUser").(models.TicketUsuario)
	id := c.Params("id")

	var usuario models.TicketUsuario
	if err := storage.DB.First(&usuario, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Usuario no encontrado",
		})
	}

	// Soft delete
	storage.DB.Delete(&usuario)

	// Log de acción
	toolbox.SaveLoggerAction(currentUser, "Usuario", fmt.Sprintf("eliminar_usuario_%s", id), c.IP())

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Usuario eliminado exitosamente",
	})
}

// ResetearPasswordUsuario godoc
// @Summary Resetear contraseña de un usuario
// @Description Genera un nuevo token de recuperación y envía un email al usuario con el enlace para restablecer su contraseña
// @Tags Usuarios
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} map[string]interface{} "Email de recuperación enviado"
// @Failure 404 {object} map[string]interface{} "Usuario no encontrado"
// @Router /usuarios/{id}/reset-password [post]
// @Security CookieAuth
// @Security TwoFactorAuth
func ResetearPasswordUsuario(c *fiber.Ctx) error {
	currentUser := c.Locals("CurrentUser").(models.TicketUsuario)
	id := c.Params("id")

	var usuario models.TicketUsuario
	if err := storage.DB.First(&usuario, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Usuario no encontrado",
		})
	}

	// Generar token de reset (igual que en RecoverPassword)
	token, hashToken, err := toolbox.GenerateResetToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error generando token de recuperacion",
		})
	}

	// Guardar hash del token en DB con expiración de 1 hora
	expiry := toolbox.AddHours(1)
	storage.DB.Model(&usuario).Updates(map[string]interface{}{
		"reset_token":        hashToken,
		"reset_token_expiry": expiry,
	})

	// Enviar email con enlace de recuperación
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), token)

	go toolbox.SendEmailAsync(
		usuario.Email,
		"Restablecer contraseña - SGI Tickets",
		fmt.Sprintf(`
			<h2>Restablecer contraseña</h2>
			<p>Hola %s,</p>
			<p>El administrador ha solicitado restablecer tu contraseña.</p>
			<p>Haz clic en el siguiente enlace para crear una nueva contraseña:</p>
			<p><a href="%s">Restablecer contraseña</a></p>
			<p>Este enlace expirará en 1 hora.</p>
			<p>Si no solicitaste este cambio, ignora este mensaje.</p>
		`,
			usuario.Nombres,
			resetURL,
		),
	)

	// Log de acción
	toolbox.SaveLoggerAction(currentUser, "Usuario", fmt.Sprintf("resetear_password_usuario_%s", id), c.IP())

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Email de recuperacion enviado exitosamente",
	})
}
