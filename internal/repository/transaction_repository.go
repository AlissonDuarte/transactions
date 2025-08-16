package repository

import (
	"context"
	"log"

	"github.com/AlissonDuarte/transactions/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	Update(ctx context.Context, transaction *models.Transaction) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	result := r.db.WithContext(ctx).Create(transaction)
	if result.Error != nil {
		log.Printf("Error creating transaction: %v", result.Error)
		return result.Error
	}
	return nil
}

func (r *transactionRepository) Update(ctx context.Context, transaction *models.Transaction) error {
	result := r.db.WithContext(ctx).Save(transaction)
	if result.Error != nil {
		log.Printf("Error updating transaction: %v", result.Error)
		return result.Error
	}
	return nil
}
