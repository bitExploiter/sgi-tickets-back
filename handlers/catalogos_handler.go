package handlers

import (
	"sgi-tickets-back/models"
	"sgi-tickets-back/storage"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetTiposDocumentosIdentificacion godoc
// @Summary Obtener tipos de documentos de identificación
// @Description Devuelve el listado de tipos de documentos de identificación
// @Tags Catalogos
// @Produce json
// @Success 200 {object} map[string]interface{} "Listado de tipos de documento"
// @Router /catalogos/tipos-documento [get]
// @Security CookieAuth
// @Security TwoFactorAuth
func GetTiposDocumentosIdentificacion(c *fiber.Ctx) error {
	var tiposDocumento []models.TicketTipoDocumento

	if err := storage.DB.Order("id ASC").Find(&tiposDocumento).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error obteniendo tipos de documento",
		})
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
// @Description Devuelve el listado de regionales
// @Tags Catalogos
// @Produce json
// @Success 200 {object} map[string]interface{} "Listado de regionales"
// @Router /catalogos/regionales [get]
// @Security CookieAuth
// @Security TwoFactorAuth
func GetRegionales(c *fiber.Ctx) error {
	var regionales []models.TicketRegional

	if err := storage.DB.Order("identificador ASC").Find(&regionales).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error obteniendo regionales",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"regionales": regionales,
		},
	})
}

// GetDepartamentos godoc
// @Summary Obtener departamentos
// @Description Devuelve el listado de departamentos
// @Tags Catalogos
// @Produce json
// @Success 200 {object} map[string]interface{} "Listado de departamentos"
// @Router /catalogos/departamentos [get]
// @Security CookieAuth
// @Security TwoFactorAuth
func GetDepartamentos(c *fiber.Ctx) error {
	var departamentos []models.TicketDepartamento

	if err := storage.DB.Order("nombre ASC").Find(&departamentos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error obteniendo departamentos",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"departamentos": departamentos,
		},
	})
}

// GetMunicipios godoc
// @Summary Obtener municipios
// @Description Devuelve el listado de municipios, opcionalmente filtrados por departamento_id
// @Tags Catalogos
// @Produce json
// @Param departamento_id query int false "Filtrar por departamento"
// @Success 200 {object} map[string]interface{} "Listado de municipios"
// @Router /catalogos/municipios [get]
// @Security CookieAuth
// @Security TwoFactorAuth
func GetMunicipios(c *fiber.Ctx) error {
	departamentoID := c.Query("departamento_id", "")

	query := storage.DB.Model(&models.TicketMunicipio{}).Order("nombre ASC")

	if departamentoID != "" {
		id, err := strconv.Atoi(departamentoID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "departamento_id invalido",
			})
		}

		query = query.Where("departamento_id = ?", id)
	}

	var municipios []models.TicketMunicipio
	if err := query.Find(&municipios).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Error obteniendo municipios",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"municipios": municipios,
		},
	})
}
