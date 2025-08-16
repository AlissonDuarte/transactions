package dto

type TransactionDTO struct {
	SenderID     int64   `json:"sender_id"`
	SenderType   string  `json:"sender_type"`
	ReceiverID   int64   `json:"receiver_id"`
	ReceiverType string  `json:"receiver_type"`
	Amount       float64 `json:"amount"`
	Status       string  `json:"status"`
	Message      string  `json:"message"`
}

func (t TransactionDTO) ValidateStatus() bool {
	switch t.Status {
	case "Pending", "Processing", "Success", "Failed":
		return true
	default:
		return false
	}
}
