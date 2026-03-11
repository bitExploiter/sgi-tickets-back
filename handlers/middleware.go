package handlers

import (
	"strings"

	"sgi-tickets-back/models"
	"sgi-tickets-back/toolbox"

	"github.com/gofiber/fiber/v2"
)

func CookieMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("sgi_tickets_user_email")
		if cookie == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "No tienes permisos para acceder a este recurso",
			})
		}
		if !toolbox.CheckCookie(cookie) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "No tienes permisos para acceder a este recurso",
			})
		}
		c.Locals("CurrentUser", toolbox.GetUserByCookie(cookie))
		return c.Next()
	}
}

func TwoFaMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("sgi_tickets_identity")
		if cookie == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "No tienes permisos para acceder a este recurso",
			})
		}
		if !toolbox.CheckCookie(cookie) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "No tienes permisos para acceder a este recurso",
			})
		}
		c.Locals("CurrentUser", toolbox.GetUserByCookie(cookie))
		return c.Next()
	}
}

func PermisosMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("CurrentUser").(models.TicketUsuario)
		route := strings.TrimPrefix(c.Route().Path, "/api/v1")
		method := c.Route().Method
		if !toolbox.HasPermissionRoute(currentUser.Rol, route, method) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "No tienes permisos suficientes para acceder a este recurso",
			})
		}
		return c.Next()
	}
}
