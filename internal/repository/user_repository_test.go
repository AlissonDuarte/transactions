package repository_test

import (
	"context"
	"testing"

	"github.com/AlissonDuarte/transactions/internal/models"
	"github.com/AlissonDuarte/transactions/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	assert.NoError(t, err)

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	user := &models.User{
		Name:     "Test User",
		Cpf:      "123.456.789-00",
		Email:    "test@example.com",
		Password: "securepassword",
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}
