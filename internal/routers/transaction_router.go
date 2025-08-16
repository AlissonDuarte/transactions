package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AlissonDuarte/transactions/internal/models"
	"github.com/AlissonDuarte/transactions/internal/routers/dto"
	"github.com/AlissonDuarte/transactions/internal/services"
	"github.com/go-chi/chi/v5"
)

type TransactionHandler struct {
	Service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: service}
}

// RegisterRoutes registra as rotas no router Chi
func (h *TransactionHandler) RegisterRoutes(r chi.Router) {
	r.Route("/transactions", func(r chi.Router) {
		r.Post("/", h.CreateTransaction)
		r.Get("/{id}", h.GetTransactionByID)
	})
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var txDTO dto.TransactionDTO

	if err := json.NewDecoder(r.Body).Decode(&txDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx := &models.Transaction{
		SenderID:     txDTO.SenderID,
		SenderType:   txDTO.SenderType,
		ReceiverID:   txDTO.ReceiverID,
		ReceiverType: txDTO.ReceiverType,
		Amount:       txDTO.Amount,
		Status:       "Pending",
		Message:      txDTO.Message,
	}

	if err := h.Service.EnqueueTransaction(r.Context(), tx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}

func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := h.Service.GetTransactionByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tx)
}
