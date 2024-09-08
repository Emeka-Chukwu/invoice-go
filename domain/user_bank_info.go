package domain

import "time"

type BankInformation struct {
	ID            int64     `json:"id"`
	AccountName   string    `json:"account_name"`
	AccountNumber string    `json:"account_number"`
	ACHRoutingNo  string    `json:"ach_routing_no"`
	BankName      string    `json:"bank_name"`
	BankAddress   string    `json:"bank_address"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserID        int       `json:"user_id"`
}

type CreateBankInformationDTO struct {
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	ACHRoutingNo  string `json:"ach_routing_no"`
	BankName      string `json:"bank_name"`
	BankAddress   string `json:"bank_address"`
	UserID        int    `json:"user_id"`
}
