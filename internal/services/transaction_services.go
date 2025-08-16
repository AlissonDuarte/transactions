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
	"github.com/streadway/amqp"
)

type TransactionService interface {
	GetTransactionByID(ctx context.Context, id int64) (*models.Transaction, error)
	EnqueueTransaction(ctx context.Context, tx *models.Transaction) error
	StartTransactionWorker(ctx context.Context)
}

type transactionService struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	httpClient      *http.Client
	amqpChannel     *amqp.Channel
	queueName       string
}

func NewTransactionService(
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
	amqpChannel *amqp.Channel,
	queueName string,
) TransactionService {
	return &transactionService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		httpClient:      &http.Client{Timeout: 5 * time.Second},
		amqpChannel:     amqpChannel,
		queueName:       queueName,
	}
}

func (s *transactionService) GetTransactionByID(ctx context.Context, id int64) (*models.Transaction, error) {
	tx, err := s.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
func (s *transactionService) EnqueueTransaction(ctx context.Context, tx *models.Transaction) error {
	tx.Status = "Pending"
	if err := s.transactionRepo.Create(ctx, tx); err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	body, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	return s.amqpChannel.Publish(
		"",
		s.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (s *transactionService) StartTransactionWorker(ctx context.Context) {
	msgs, err := s.amqpChannel.Consume(
		s.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	go func() {
		for msg := range msgs {
			var tx models.Transaction
			if err := json.Unmarshal(msg.Body, &tx); err != nil {
				msg.Nack(false, false)
				continue
			}

			err := s.transactionRepo.Transaction(ctx, func(txRepo repository.TransactionRepository) error {
				senderAcc, err := s.accountRepo.GetByOwnerID(ctx, int64(tx.SenderID), "")
				if err != nil || senderAcc == nil {
					tx.Status = "Failed"
					return s.transactionRepo.Update(ctx, &tx)
				}

				receiverAcc, err := s.accountRepo.GetByOwnerID(ctx, int64(tx.ReceiverID), "")
				if err != nil || receiverAcc == nil {
					tx.Status = "Failed"
					return s.transactionRepo.Update(ctx, &tx)
				}

				if !senderAcc.CanSend || senderAcc.Balance < tx.Amount {
					tx.Status = "Failed"
					return s.transactionRepo.Update(ctx, &tx)
				}

				authOK, err := s.checkAuthorization()
				if err != nil || !authOK {
					tx.Status = "Failed"
					return s.transactionRepo.Update(ctx, &tx)
				}

				senderAcc.Balance -= tx.Amount
				receiverAcc.Balance += tx.Amount

				if err := s.accountRepo.Update(ctx, senderAcc); err != nil {
					return err
				}

				if err := s.accountRepo.Update(ctx, receiverAcc); err != nil {
					return err
				}

				if err := s.notificationExternalService(); err != nil {
					return err
				}

				tx.Status = "Success"
				return s.transactionRepo.Update(ctx, &tx)
			})

			if err != nil {
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}()
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
