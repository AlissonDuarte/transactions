package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	ID           int64
	SenderID     int64
	SenderType   string
	ReceiverID   int64
	ReceiverType string
	Amount       float64
	Status       string
	Message      string
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.Status == "" {
		t.Status = "Pending"
	}
	return nil
}
