// internal/repository/user.go
package repository

import (
	"context"
	"errors"
	"log"

	"github.com/AlissonDuarte/transactions/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		log.Printf("Error creating user: %v", result.Error)
		return result.Error
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("Error getting user by ID: %v", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("Error getting user by email: %v", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		log.Printf("Error updating user: %v", result.Error)
		return result.Error
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		log.Printf("Error deleting user: %v", result.Error)
		return result.Error
	}
	return nil
}

func (r *userRepository) List(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	result := r.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		log.Printf("Error listing users: %v", result.Error)
		return nil, result.Error
	}
	return users, nil
}
