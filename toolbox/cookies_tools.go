package toolbox

import (
	"sgi-tickets-back/models"
	"sgi-tickets-back/storage"
)

func SaveCookieToStorage(token string, level string) {
	cookie := models.Cookie{Token: token, Habilitado: true, Level: level}
	storage.DB.Create(&cookie)
}

func CheckCookie(token string) bool {
	cookie := models.Cookie{}
	storage.DB.Where("token = ? AND habilitado = ?", token, true).First(&cookie)
	return cookie.Id != 0
}

func GetUserByCookie(token string) models.TicketUsuario {
	cookie := models.Cookie{}
	storage.DB.Where("token = ?", token).First(&cookie)
	usuario := models.TicketUsuario{}
	storage.DB.Where("email = ?", cookie.Level).First(&usuario)
	return usuario
}

func DisableCookie(token string) {
	cookie := models.Cookie{}
	storage.DB.Where("token = ?", token).First(&cookie)
	storage.DB.Model(&cookie).Update("habilitado", false)
}

func DisableAllCookies(email string) {
	storage.DB.Model(&models.Cookie{}).Where("level = ?", email).Update("habilitado", false)
}
