package handlers

import "github.com/gofiber/fiber/v2"

// GetTiposDocumentosIdentificacion godoc
// @Summary Obtener tipos de documentos de identificación
// @Description Devuelve el listado de tipos de documentos de identificación (hardcodeado temporalmente)
// @Tags Catalogos
// @Produce json
// @Success 200 {object} map[string]interface{} "Listado de tipos de documento"
// @Router /catalogos/tipos-documento [get]
// @Security CookieAuth
// @Security TwoFactorAuth
func GetTiposDocumentosIdentificacion(c *fiber.Ctx) error {
	tiposDocumento := []fiber.Map{
		{"id": 1, "nombre": "Cédula de Ciudadanía"},
		{"id": 2, "nombre": "Cédula de Extranjería"},
		{"id": 3, "nombre": "Pasaporte"},
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"tipos_documento": tiposDocumento,
		},
	})
}

// GetRegionales godoc
// @Summary Obtener regionales
// @Description Devuelve el listado de regionales (hardcodeado temporalmente)
// @Tags Catalogos
// @Produce json
// @Success 200 {object} map[string]interface{} "Listado de regionales"
// @Router /catalogos/regionales [get]
// @Security CookieAuth
// @Security TwoFactorAuth
func GetRegionales(c *fiber.Ctx) error {
	regionales := []fiber.Map{
		{"nombre": "CENTRAL", "identificador": 100},
		{"nombre": "OCCIDENTE", "identificador": 200},
		{"nombre": "NORTE", "identificador": 300},
		{"nombre": "ORIENTE", "identificador": 400},
		{"nombre": "NOROESTE", "identificador": 500},
		{"nombre": "VIEJO CALDAS", "identificador": 600},
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"regionales": regionales,
		},
	})
}
