package domain

import "time"

type PaymentDto struct {
	InvoiceID   int       `json:"invoice_id"`
	Amount      float64   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
}

type PaymentReponse struct {
	ID          int       `json:"id"`
	InvoiceID   int       `json:"invoice_id"`
	Amount      float64   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
