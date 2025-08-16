package models

type Account struct {
	ID         int64
	OwnerID    int64
	OwnerType  string
	Balance    float64
	CanSend    bool
	CanReceive bool
}
