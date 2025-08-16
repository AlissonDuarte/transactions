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

func TestUserRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)

	user := &models.User{
		Name:     "Test User",
		Cpf:      "123.456.789-00",
		Email:    "test@example.com",
		Password: "securepassword",
	}
	err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	got, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Name, got.Name)
	assert.Equal(t, user.Email, got.Email)

	got, err = repo.GetByEmail(ctx, user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)

	user.Name = "Updated Name"
	err = repo.Update(ctx, user)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)

	users, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, user.ID, users[0].ID)

	err = repo.Delete(ctx, user.ID)
	assert.NoError(t, err)

	deleted, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Nil(t, deleted)
}

func TestUserRepository_GetNonExistent(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)

	user, err := repo.GetByID(ctx, 999)
	assert.NoError(t, err)
	assert.Nil(t, user)

	user, err = repo.GetByEmail(ctx, "notfound@example.com")
	assert.NoError(t, err)
	assert.Nil(t, user)
}
