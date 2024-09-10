package domain

import "time"

type Invoice struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	CustomerID    int       `json:"customer_id"`
	InvoiceNumber string    `json:"invoice_number"`
	Status        string    `json:"status"`
	IssueDate     time.Time `json:"issue_date"`
	DueDate       time.Time `json:"due_date"`
	TotalAmount   float64   `json:"total_amount"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type InvoiceDTO struct {
	UserID        int       `json:"user_id" validate:"required"`
	CustomerID    int       `json:"customer_id" validate:"required"`
	InvoiceNumber string    `json:"invoice_number" validate:"required"`
	Status        string    `json:"status" validate:"required"`
	IssueDate     time.Time `json:"issue_date" validate:"required"`
	DueDate       time.Time `json:"due_date" validate:"required"`
	TotalAmount   float64   `json:"total_amount" validate:"required,gt=0"`
}

type InvoiceItem struct {
	ID          int       `json:"id"`
	InvoiceID   int       `json:"invoice_id"`
	Description string    `json:"description"`
	Title       string    `json:"title"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type InvoiceItemDTO struct {
	InvoiceID   int     `json:"invoice_id"`
	Description string  `json:"description" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
	UnitPrice   float64 `json:"unit_price" validate:"required,gt=0"`
}

type CreateInvoiceItem []InvoiceItemDTO
type Items []InvoiceItem

type CreateInvoiceRequestDTO struct {
	InvoiceDTO
	CreateInvoiceItem
}

type InvoiceResponse struct {
	Invoice
	Items
}

type InvoiceStats struct {
	Count       int     `json:"count"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}
