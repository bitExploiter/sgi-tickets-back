package toolbox

import (
	"sgi-tickets-back/models"
	"sgi-tickets-back/storage"
)

func SaveLoggerAction(usuario models.TicketUsuario, tipo string, accion string, ip string) {
	logger := models.Logger{}
	logger.UsuarioID = usuario.Id
	logger.Tipo = tipo
	logger.Accion = accion
	logger.Descripcion = usuario.Email
	logger.Ip = ip
	storage.DB.Create(&logger)
}
