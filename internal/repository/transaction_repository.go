package repository

import (
	"context"
	"log"

	"github.com/AlissonDuarte/transactions/internal/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *models.Transaction) error
	Update(ctx context.Context, tx *models.Transaction) error
	Transaction(ctx context.Context, fn func(txRepo TransactionRepository) error) error
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

func (r *transactionRepository) Transaction(ctx context.Context, fn func(txRepo TransactionRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &transactionRepository{db: tx}
		return fn(txRepo)
	})
}
