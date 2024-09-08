package invoice_repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"go-invoice/domain"
	"go-invoice/util"
	"strings"
)

type invoiceRepository struct {
	Db *sql.DB
}
type InvoiceRepository interface {
	CreateInvoiceWithItems(req domain.CreateInvoiceRequestDTO) (int, error)
	FetchInvoicesWithItems(userId int, limit, offset int64) ([]domain.InvoiceResponse, error)
	FetchInvoices(userId int, page, limit int64) ([]domain.InvoiceResponse, error)
	FetchInvoiceWithItems(invoiceId int) (domain.InvoiceResponse, error)
	DeleteInvoiceItems(itemIds []int) error
	FetchInvoiceStats(userId int) (map[string]domain.InvoiceStats, error)
}

func NewInvoiceWithItems(Db *sql.DB) InvoiceRepository {
	return &invoiceRepository{Db: Db}
}

func (iR *invoiceRepository) CreateInvoiceWithItems(req domain.CreateInvoiceRequestDTO) (int, error) {
	tx, err := iR.Db.Begin()
	if err != nil {
		return 0, err
	}
	var invoiceID int
	err = tx.QueryRow(`
        INSERT INTO invoices (user_id, customer_id, invoice_number, status, issue_date, due_date, total_amount) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) 
        RETURNING id`,
		req.UserID, req.CustomerID, req.InvoiceNumber, req.Status, req.IssueDate, req.DueDate, req.TotalAmount).Scan(&invoiceID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	for _, item := range req.CreateInvoiceItem {
		_, err = tx.Exec(`
            INSERT INTO invoice_items (invoice_id, description, quantity, unit_price) 
            VALUES ($1, $2, $3, $4)`,
			invoiceID, item.Description, item.Quantity, item.UnitPrice)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (ir invoiceRepository) FetchInvoicesWithItems(userId int, limit, offset int64) ([]domain.InvoiceResponse, error) {
	query := `
    SELECT i.id AS invoice_id, i.user_id, i.customer_id, i.invoice_number, i.status, 
           i.issue_date, i.due_date, i.total_amount, i.created_at AS invoice_created_at, 
           i.updated_at AS invoice_updated_at, 
           COALESCE(
               json_agg(
                   json_build_object(
                       'id', ii.id, 'description', ii.description, 
                       'quantity', ii.quantity, 'unit_price', ii.unit_price, 
                       'created_at', ii.created_at, 'updated_at', ii.updated_at
                   )
               ) FILTER (WHERE ii.id IS NOT NULL), '[]'
           ) AS items
    FROM invoices i
    LEFT JOIN invoice_items ii ON i.id = ii.invoice_id
    GROUP BY i.id
    ORDER BY i.created_at DESC
    LIMIT $2 OFFSET $3
	WHERE user_id = $1`

	rows, err := ir.Db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var invoices []domain.InvoiceResponse
	for rows.Next() {
		var invoice domain.InvoiceResponse
		var itemsJson json.RawMessage
		err := rows.Scan(
			&invoice.ID, &invoice.UserID, &invoice.CustomerID, &invoice.InvoiceNumber,
			&invoice.Status, &invoice.IssueDate, &invoice.DueDate,
			&invoice.TotalAmount, &invoice.CreatedAt, &invoice.UpdatedAt,
			&itemsJson,
		)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(itemsJson, &invoice.Items)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

func (ir *invoiceRepository) FetchInvoices(userId int, page, limit int64) ([]domain.InvoiceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), util.DbTimeout)
	defer cancel()

	offset := (page - 1) * limit
	stmt := `
		SELECT id, user_id, customer_id, invoice_number, status, issue_date, due_date, total_amount, created_at, updated_at
		FROM invoices
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := ir.Db.QueryContext(ctx, stmt, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var invoices []domain.InvoiceResponse
	for rows.Next() {
		var invoice domain.InvoiceResponse
		err := rows.Scan(
			&invoice.ID,
			&invoice.UserID,
			&invoice.CustomerID,
			&invoice.InvoiceNumber,
			&invoice.Status,
			&invoice.IssueDate,
			&invoice.DueDate,
			&invoice.TotalAmount,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return invoices, nil
}

func (ir invoiceRepository) FetchInvoiceWithItems(invoiceId int) (domain.InvoiceResponse, error) {
	query := `
    SELECT i.id AS invoice_id, i.user_id, i.customer_id, i.invoice_number, i.status, 
           i.issue_date, i.due_date, i.total_amount, i.created_at AS invoice_created_at, 
           i.updated_at AS invoice_updated_at, 
           COALESCE(
               json_agg(
                   json_build_object(
                       'id', ii.id, 'description', ii.description, 
                       'quantity', ii.quantity, 'unit_price', ii.unit_price, 
                       'created_at', ii.created_at, 'updated_at', ii.updated_at
                   )
               ) FILTER (WHERE ii.id IS NOT NULL), '[]'
           ) AS items
    FROM invoices i
    LEFT JOIN invoice_items ii ON i.id = ii.invoice_id
    WHERE i.id = $1 
    GROUP BY i.id
    `
	row := ir.Db.QueryRow(query, invoiceId)
	var invoice domain.InvoiceResponse
	var itemsJson json.RawMessage
	err := row.Scan(
		&invoice.ID, &invoice.UserID, &invoice.CustomerID, &invoice.InvoiceNumber,
		&invoice.Status, &invoice.IssueDate, &invoice.DueDate,
		&invoice.TotalAmount, &invoice.CreatedAt, &invoice.UpdatedAt,
		&itemsJson,
	)
	if err != nil {
		return invoice, err
	}
	err = json.Unmarshal(itemsJson, &invoice.Items)
	if err != nil {
		return invoice, err
	}
	return invoice, nil
}

func (ir *invoiceRepository) DeleteInvoiceItems(itemIds []int) error {
	tx, err := ir.Db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("DELETE FROM invoice_items WHERE id = $1")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	for _, id := range itemIds {
		_, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *invoiceRepository) FetchInvoiceStats(userId int) (map[string]domain.InvoiceStats, error) {
	var defaultStatuses = []string{"paid", "overdue", "draft", "unpaid"}
	ctx, cancel := context.WithTimeout(context.Background(), util.DbTimeout)
	defer cancel()
	stmt := `SELECT status,  COUNT(*) AS count, SUM(total_amount) AS total_amount FROM invoices WHERE user_id=$1 GROUP BY status;`
	rows, err := a.Db.QueryContext(ctx, stmt, userId)
	statsMap := make(map[string]domain.InvoiceStats)
	for _, status := range defaultStatuses {
		statsMap[strings.ToLower(status)] = domain.InvoiceStats{
			Status:      status,
			Count:       0,
			TotalAmount: 0.0,
		}
	}
	for rows.Next() {
		var stat domain.InvoiceStats
		err := rows.Scan(&stat.Status, &stat.Count, &stat.TotalAmount)
		if err != nil {
			return nil, err
		}
		statsMap[strings.ToLower(stat.Status)] = stat
	}
	return statsMap, err
}
