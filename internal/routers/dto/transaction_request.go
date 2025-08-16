package dto

type TransactionDTO struct {
	SenderID   uint    `json:"sender_id"`
	ReceiverID uint    `json:"receiver_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	Message    string  `json:"message"`
}

func (t TransactionDTO) ValidateStatus() bool {
	switch t.Status {
	case "Pending", "Processing", "Success", "Failed":
		return true
	default:
		return false
	}
}
