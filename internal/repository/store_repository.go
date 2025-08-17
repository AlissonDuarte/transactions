package repository

import (
	"context"
	"log"

	"github.com/AlissonDuarte/transactions/internal/models"
	"gorm.io/gorm"
)

type StoreRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Store, error)
	Create(ctx context.Context, store *models.Store) error
	Update(ctx context.Context, store *models.Store) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*models.Store, error)
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) GetByID(ctx context.Context, id int64) (*models.Store, error) {
	var store models.Store
	result := r.db.WithContext(ctx).First(&store, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &store, nil
}

func (r *storeRepository) Create(ctx context.Context, store *models.Store) error {
	result := r.db.WithContext(ctx).Create(store)
	if result.Error != nil {
		log.Printf("Error creating store: %v", result.Error)
		return result.Error
	}
	return nil
}

func (r *storeRepository) Update(ctx context.Context, store *models.Store) error {
	result := r.db.WithContext(ctx).Save(store)
	if result.Error != nil {
		log.Printf("Error updating store: %v", result.Error)
		return result.Error
	}
	return nil
}

func (r *storeRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&models.Store{}, id)
	if result.Error != nil {
		log.Printf("Error deleting store: %v", result.Error)
		return result.Error
	}
	return nil
}

func (r *storeRepository) List(ctx context.Context) ([]*models.Store, error) {
	var stores []*models.Store
	result := r.db.WithContext(ctx).Find(&stores)
	if result.Error != nil {
		log.Printf("Error listing stores: %v", result.Error)
		return nil, result.Error
	}
	return stores, nil
}
