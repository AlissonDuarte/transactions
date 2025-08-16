package models

import "gorm.io/gorm"

type Store struct {
	gorm.Model
	Name     string
	CNPJ     string
	Email    string
	Password string
}
