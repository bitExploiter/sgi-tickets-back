package toolbox

import (
	"sgi-tickets-back/models"
	"sgi-tickets-back/storage"
)

func GetPermisosByRol(rol string) []models.PermisoRol {
	permisos := []models.PermisoRol{}
	storage.DB.Preload("Permiso").Where("rol = ?", rol).Find(&permisos)
	return permisos
}

func HasPermissionRoute(rol string, route string, method string) bool {
	permisos := GetPermisosByRol(rol)
	for _, permiso := range permisos {
		if permiso.Permiso.Ruta == route && permiso.Permiso.Metodo == method {
			return true
		}
	}
	return false
}
