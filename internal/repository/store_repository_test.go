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

func setupTestDBStore(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Store{})
	assert.NoError(t, err)

	return db
}

func TestStoreRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	db := setupTestDBStore(t)
	repo := repository.NewStoreRepository(db)

	store := &models.Store{
		Name:     "Test Store",
		CNPJ:     "0001/564123-00",
		Email:    "test@example.com",
		Password: "securepassword",
	}

	err := repo.Create(ctx, store)
	assert.NoError(t, err)
	assert.NotZero(t, store.ID)

	got, err := repo.GetByID(ctx, int64(store.ID))
	assert.NoError(t, err)
	assert.Equal(t, store.Name, got.Name)
	assert.Equal(t, store.Email, got.Email)

	store.Name = "Updated Name"
	err = repo.Update(ctx, store)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ctx, int64(store.ID))
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
}
