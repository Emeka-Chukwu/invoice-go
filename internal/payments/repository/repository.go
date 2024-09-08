package payment_repository

import (
	"database/sql"
	"go-invoice/domain"
)

type paymentRepo struct {
	db *sql.DB
}

type PaymentRepo interface {
	CreatePayment(payment domain.PaymentDto) (domain.PaymentReponse, error)
}

func NewPaymentRepo(db *sql.DB) PaymentRepo {
	return &paymentRepo{db: db}
}

func (pr *paymentRepo) CreatePayment(payment domain.PaymentDto) (domain.PaymentReponse, error) {
	query := `
		INSERT INTO payments (invoice_id, amount, payment_date) 
		VALUES ($1, $2, $3) 
		RETURNING id, invoice_id, amount, payment_date, created_at, updated_at`
	var paymentResp domain.PaymentReponse
	err := pr.db.QueryRow(query, payment.InvoiceID, payment.Amount, payment.PaymentDate).
		Scan(&paymentResp.ID, &paymentResp.InvoiceID, &paymentResp.Amount, &paymentResp.PaymentDate, &paymentResp.CreatedAt, &paymentResp.UpdatedAt)
	if err != nil {
		return domain.PaymentReponse{}, err
	}
	return paymentResp, nil
}
