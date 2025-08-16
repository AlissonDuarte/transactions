package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	ID         int64
	SenderID   uint
	ReceiverID uint
	Amount     float64
	Status     string
	Message    string
}
