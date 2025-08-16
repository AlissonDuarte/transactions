package repository

import (
	"context"
	"errors"
	"log"

	"github.com/AlissonDuarte/transactions/internal/models"
	"gorm.io/gorm"
)

type AccountRepository interface {
	GetByOwnerID(ctx context.Context, ownerID int64, ownerType string) (*models.Account, error)
	Create(ctx context.Context, account *models.Account) error
	Update(ctx context.Context, account *models.Account) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) GetByOwnerID(ctx context.Context, ownerID int64, ownerType string) (*models.Account, error) {
	var account models.Account
	result := r.db.WithContext(ctx).Where("owner_id = ? AND owner_type = ?", ownerID).First(&account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("Error getting account by owner ID: %v", result.Error)
		return nil, result.Error
	}
	return &account, nil
}

func (r *accountRepository) Create(ctx context.Context, account *models.Account) error {
	result := r.db.WithContext(ctx).Create(account)
	if result.Error != nil {
		log.Printf("Error creating account: %v", result.Error)
		return result.Error
	}
	return nil
}

func (r *accountRepository) Update(ctx context.Context, account *models.Account) error {
	var existingAccount models.Account
	if err := r.db.WithContext(ctx).First(&existingAccount, account.ID).Error; err != nil {
		log.Printf("Error finding account for update: %v", err)
		return errors.New("account not found")
	}

	if account.OwnerID != existingAccount.OwnerID || account.OwnerType != existingAccount.OwnerType {
		return errors.New("cannot change account owner information (OwnerID or OwnerType)")
	}

	result := r.db.WithContext(ctx).Model(&models.Account{}).
		Where("id = ?", account.ID).
		Updates(map[string]interface{}{
			"balance":     account.Balance,
			"can_send":    account.CanSend,
			"can_receive": account.CanReceive,
			"active":      account.Active,
			"updated_at":  gorm.Expr("CURRENT_TIMESTAMP"),
		})

	if result.Error != nil {
		log.Printf("Error updating account: %v", result.Error)
		return result.Error
	}

	return nil
}
