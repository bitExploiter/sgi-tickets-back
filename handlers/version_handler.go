package handlers

import "github.com/gofiber/fiber/v2"

// GetVersion godoc
// @Summary Obtener versión de la API
// @Description Devuelve la versión actual de la API y el nombre del sistema
// @Tags Sistema
// @Produce json
// @Success 200 {object} map[string]interface{} "Información de versión"
// @Router /api/version [get]
func GetVersion(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"version": "1.0.0",
			"nombre":  "SGI Tickets API",
		},
	})
}
