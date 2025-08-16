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

func setupTestDBTransaction(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Transaction{})
	assert.NoError(t, err)

	return db
}

func TestTransactionRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	db := setupTestDBTransaction(t)
	repo := repository.NewTransactionRepository(db)

	tx := &models.Transaction{
		SenderID:   1,
		ReceiverID: 2,
		Amount:     100.0,
		Status:     "Pending",
		Message:    "Test transaction",
	}

	err := repo.Create(ctx, tx)
	assert.NoError(t, err)
	assert.NotZero(t, tx.ID)

	tx.Status = "Processing"
	tx.Message = "Updated message"
	err = repo.Update(ctx, tx)
	assert.NoError(t, err)

	var updatedTx models.Transaction
	err = db.WithContext(ctx).First(&updatedTx, tx.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Processing", updatedTx.Status)
	assert.Equal(t, "Updated message", updatedTx.Message)
}

func TestTransactionRepository_CreateMultiple(t *testing.T) {
	ctx := context.Background()
	db := setupTestDBTransaction(t)
	repo := repository.NewTransactionRepository(db)

	tx1 := &models.Transaction{SenderID: 1, ReceiverID: 2, Amount: 50, Status: "Pending"}
	tx2 := &models.Transaction{SenderID: 2, ReceiverID: 1, Amount: 30, Status: "Pending"}

	assert.NoError(t, repo.Create(ctx, tx1))
	assert.NoError(t, repo.Create(ctx, tx2))

	var count int64
	db.Model(&models.Transaction{}).Count(&count)
	assert.Equal(t, int64(3), count)
}
