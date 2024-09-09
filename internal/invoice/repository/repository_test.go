package invoice_repository

import (
	"go-invoice/domain"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateInvoiceWithItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	invoiceRepo := NewInvoiceWithItems(db)
	req := domain.CreateInvoiceRequestDTO{}
	req.UserID = 1
	req.CustomerID = 1
	req.InvoiceNumber = "INV-001"
	req.Status = "unpaid"
	req.IssueDate = time.Now()
	req.DueDate = time.Now()
	req.TotalAmount = 1000
	req.CreateInvoiceItem = []domain.InvoiceItemDTO{
		{Description: "Item 1", Quantity: 1, UnitPrice: 500},
		{Description: "Item 2", Quantity: 1, UnitPrice: 500},
	}
	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO invoices`).
		WithArgs(req.UserID, req.CustomerID, req.InvoiceNumber, req.Status, req.IssueDate, req.DueDate, req.TotalAmount).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	for _, item := range req.CreateInvoiceItem {
		mock.ExpectExec(`INSERT INTO invoice_items`).
			WithArgs(1, item.Description, item.Quantity, item.UnitPrice).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
	mock.ExpectCommit()
	invoiceID, err := invoiceRepo.CreateInvoiceWithItems(req)
	require.NoError(t, err)
	require.Equal(t, 1, invoiceID)
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestFetchInvoicesWithItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	invoiceRepo := &invoiceRepository{Db: db}

	userId := 1
	limit := int64(10)
	offset := int64(0)

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
               ) FILTER \(WHERE ii.id IS NOT NULL\), '[]'
           ) AS items
    FROM invoices i
    LEFT JOIN invoice_items ii ON i.id = ii.invoice_id
    WHERE i.user_id = \$1
    GROUP BY i.id
    ORDER BY i.created_at DESC
    LIMIT \$2 OFFSET \$3`

	items := `[{"id":1,"description":"Item 1","quantity":1,"unit_price":500,"created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"}]`
	rows := sqlmock.NewRows([]string{
		"invoice_id", "user_id", "customer_id", "invoice_number", "status", "issue_date", "due_date", "total_amount", "invoice_created_at", "invoice_updated_at", "items",
	}).
		AddRow(1, userId, 1, "INV-001", "unpaid", time.Now(), time.Now(), 1000, time.Now(), time.Now(), items)

	mock.ExpectQuery(query).
		WithArgs(userId, limit, offset).
		WillReturnRows(rows)

	invoices, err := invoiceRepo.FetchInvoicesWithItems(userId, limit, offset)

	require.NoError(t, err)
	require.Len(t, invoices, 1)
	require.Equal(t, "INV-001", invoices[0].InvoiceNumber)
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestFetchInvoices(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	invoiceRepo := &invoiceRepository{Db: db}

	userId := 1
	page := int64(1)
	limit := int64(10)
	offset := (page - 1) * limit

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "customer_id", "invoice_number", "status", "issue_date", "due_date", "total_amount", "created_at", "updated_at",
	}).
		AddRow(1, userId, 1, "INV-001", "unpaid", time.Now(), time.Now(), 1000.0, time.Now(), time.Now()).
		AddRow(2, userId, 2, "INV-002", "paid", time.Now(), time.Now(), 2000.0, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT id, user_id, customer_id, invoice_number, status, issue_date, due_date, total_amount, created_at, updated_at FROM invoices WHERE user_id = \$1 ORDER BY created_at DESC LIMIT \$2 OFFSET \$3`).
		WithArgs(userId, limit, offset).
		WillReturnRows(rows)

	invoices, err := invoiceRepo.FetchInvoices(userId, page, limit)

	require.NoError(t, err)
	require.Len(t, invoices, 2)

	require.Equal(t, "INV-001", invoices[0].InvoiceNumber)
	require.Equal(t, "unpaid", invoices[0].Status)
	require.Equal(t, 1000.0, invoices[0].TotalAmount)

	require.Equal(t, "INV-002", invoices[1].InvoiceNumber)
	require.Equal(t, "paid", invoices[1].Status)
	require.Equal(t, 2000.0, invoices[1].TotalAmount)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestUpdateInvoiceStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	invoiceRepo := &invoiceRepository{Db: db}

	invoiceID := 1
	userID := 1
	newStatus := "paid"

	mock.ExpectExec(`
        UPDATE invoices
        SET status = \$1, updated_at = CURRENT_TIMESTAMP
        WHERE id = \$2 AND user_id = \$3;
    `).
		WithArgs(newStatus, invoiceID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = invoiceRepo.UpdateInvoiceStatus(invoiceID, userID, newStatus)

	require.NoError(t, err)
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestFetchInvoiceWithItems(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	invoiceRepo := NewInvoiceWithItems(db)

	invoiceID := 1
	itemsJson := `[{
		"id": 1,
		"description": "Item 1",
		"quantity": 2,
		"unit_price": 100,
		"created_at": "2023-01-01T00:00:00Z",
		"updated_at": "2023-01-01T00:00:00Z"
	}]`
	createdAt := time.Now()
	updatedAt := time.Now()

	mock.ExpectQuery(`
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
		WHERE i.id = \$1 
		GROUP BY i.id
	`).
		WithArgs(invoiceID).
		WillReturnRows(sqlmock.NewRows([]string{
			"invoice_id", "user_id", "customer_id", "invoice_number", "status",
			"issue_date", "due_date", "total_amount", "invoice_created_at",
			"invoice_updated_at", "items",
		}).AddRow(
			invoiceID, 1, 1, "INV-001", "paid",
			createdAt, createdAt.AddDate(0, 0, 30), 200.00, createdAt,
			updatedAt, itemsJson,
		))

	invoice, err := invoiceRepo.FetchInvoiceWithItems(invoiceID)

	require.NoError(t, err)
	require.Equal(t, invoiceID, invoice.ID)
	require.Equal(t, "INV-001", invoice.InvoiceNumber)
	require.Equal(t, "paid", invoice.Status)
	require.Equal(t, 1, len(invoice.Items))
	require.Equal(t, "Item 1", invoice.Items[0].Description)

	// Ensure that all expectations were met
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
