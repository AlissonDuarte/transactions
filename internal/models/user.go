package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       int64
	Name     string
	Cpf      string
	Email    string
	Password string
}
