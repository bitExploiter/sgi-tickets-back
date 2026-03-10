package models

import "gorm.io/gorm"

type Permiso struct {
	gorm.Model

	Id     uint   `json:"id" gorm:"primaryKey"`
	Ruta   string `json:"ruta" validate:"required"`
	Metodo string `json:"metodo" validate:"required"`
	Nombre string `json:"nombre" validate:"-"`
}

func (Permiso) TableName() string {
	return "tickets_permisos"
}

type PermisoRol struct {
	gorm.Model

	Id        uint    `json:"id" gorm:"primaryKey"`
	Rol       string  `json:"rol" validate:"required"`
	PermisoID uint    `json:"permiso_id" validate:"required"`
	Permiso   Permiso `json:"permiso" validate:"-"`
}

func (PermisoRol) TableName() string {
	return "tickets_permisos_rol"
}
