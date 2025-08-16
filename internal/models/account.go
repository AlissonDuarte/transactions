package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	ID         int64
	OwnerID    int64
	OwnerType  string
	Balance    float64
	CanSend    bool
	CanReceive bool
}
