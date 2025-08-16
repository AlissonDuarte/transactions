package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AlissonDuarte/transactions/internal/models"
	"github.com/AlissonDuarte/transactions/internal/repository"
)

type TransactionService interface {
	SendTransaction(
		ctx context.Context,
		senderID int64,
		senderType string,
		receiverID int64,
		receiverType string,
		amount float64,
		message string) (*models.Transaction, error)
}

type transactionService struct {
	userRepo        repository.UserRepository
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	httpClient      *http.Client
}

func NewTransactionService(userRepo repository.UserRepository, accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository) TransactionService {
	return &transactionService{
		userRepo:        userRepo,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		httpClient:      &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *transactionService) checkAuthorization() (bool, error) {
	resp, err := s.httpClient.Get("https://util.devi.tools/api/v2/authorize")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var body struct {
		Status string `json:"status"`
		Data   struct {
			Authorization bool `json:"authorization"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return false, err
	}

	return body.Status == "success" && body.Data.Authorization, nil
}

func (s *transactionService) notificationExternalService() error {
	reqBody := bytes.NewBuffer([]byte(`{}`))
	resp, err := s.httpClient.Post("https://util.devi.tools/api/v2/notify", "application/json", reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	var body struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}

	return errors.New(body.Message)
}

func (s *transactionService) SendTransaction(
	ctx context.Context,
	senderID int64,
	senderType string,
	receiverID int64,
	receiverType string,
	amount float64,
	message string) (*models.Transaction, error) {
	senderAccount, err := s.accountRepo.GetByOwnerID(ctx, senderID, senderType)
	if err != nil || senderAccount == nil {
		return nil, errors.New("sender account not found")
	}
	receiverAccount, err := s.accountRepo.GetByOwnerID(ctx, receiverID, receiverType)
	if err != nil || receiverAccount == nil {
		return nil, errors.New("receiver account not found")
	}

	if !senderAccount.CanSend {
		return nil, errors.New("sender account cannot send transactions")
	}

	if senderAccount.Balance < amount {
		return nil, errors.New("sender account balance is not enough")
	}

	authOK, err := s.checkAuthorization()

	if err != nil {
		return nil, fmt.Errorf("error checking authorization: %v", err)
	}

	if !authOK {
		return nil, errors.New("authorization failed")
	}

	tx := &models.Transaction{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     amount,
		Status:     "Pending",
		Message:    message,
	}

	if err := s.transactionRepo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("error creating transaction: %v", err)
	}

	senderAccount.Balance -= amount
	receiverAccount.Balance += amount

	if err := s.accountRepo.Update(ctx, senderAccount); err != nil {
		return nil, err
	}
	if err := s.accountRepo.Update(ctx, receiverAccount); err != nil {
		return nil, err
	}

	if err := s.notificationExternalService(); err != nil {
		return nil, err
	}

	tx.Status = "Success"
	if err := s.transactionRepo.Update(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}
