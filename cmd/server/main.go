package main

import (
	"log"
	"net/http"

	"github.com/AlissonDuarte/transactions/internal/models"
	"github.com/AlissonDuarte/transactions/internal/repository"
	"github.com/AlissonDuarte/transactions/internal/routers"
	"github.com/AlissonDuarte/transactions/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.Open(
		"host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"),
		&gorm.Config{},
	)
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Transaction{},
	); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	amqpConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer amqpConn.Close()

	amqpChannel, err := amqpConn.Channel()
	if err != nil {
		panic(err)
	}
	defer amqpChannel.Close()

	queueName := "transactions"

	_, err = amqpChannel.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		panic(err)
	}

	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	txService := services.NewTransactionService(accountRepo, transactionRepo, amqpChannel, queueName)
	txService.StartTransactionWorker(nil)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	txHandler := routers.NewTransactionHandler(txService)
	txHandler.RegisterRoutes(r)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
