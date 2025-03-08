package entity

import (
	"encoding/json"
	"time"
)

type Transaction struct {
	Id        string
	UserId    string
	PhotoId   string
	SnapToken string
	// ExternalStatus           enum.MidtransPaymentStatus
	ExternalCallbackResponse json.RawMessage
	PaidAt                   time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
}
