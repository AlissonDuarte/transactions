package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	SenderID   uint
	ReceiverID uint
	Amount     float64
	Status     string
	Message    string
}
