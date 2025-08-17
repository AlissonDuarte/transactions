package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Printf("Connected to database: %v\n", db)
			break
		}
		fmt.Println("Database not ready, retrying in 2s...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Transaction{},
		&models.Store{},
	); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitPort := os.Getenv("RABBITMQ_PORT")
	rabbitUser := os.Getenv("RABBITMQ_USER")
	rabbitPass := os.Getenv("RABBITMQ_PASS")

	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		rabbitUser, rabbitPass, rabbitHost, rabbitPort,
	)

	// exemplo usando streadway/amqp
	var amqpConn *amqp.Connection
	for i := 0; i < 10; i++ {
		amqpConn, err = amqp.Dial(rabbitURL)
		if err == nil {
			fmt.Println("Connected to RabbitMQ")
			break
		}
		fmt.Println("RabbitMQ not ready, retrying in 2s...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		panic(fmt.Sprintf("failed to connect to RabbitMQ: %v", err))
	}
	defer amqpConn.Close()
	fmt.Println("Connected to RabbitMQ")

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

	numConsumers := 5
	for i := 0; i < numConsumers; i++ {
		ch, err := amqpConn.Channel()
		if err != nil {
			panic(err)
		}
		go txService.StartTransactionWorker(context.Background(), ch)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	txHandler := routers.NewTransactionHandler(txService)
	txHandler.RegisterRoutes(r)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
