package customer_repository

import (
	"go-invoice/domain"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateCustomer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	customerRepo := NewCustomerRepository(db)
	req := domain.CustomerDTO{
		UserID: 101,
		Name:   "Jane Doe",
		Email:  "jane.doe@example.com",
		Phone:  "123-456-7890",
	}

	query := `
    INSERT INTO customers \(user_id, name, email, phone, created_at, updated_at\) 
    VALUES \(\$1, \$2, \$3, \$4,  CURRENT_TIMESTAMP, CURRENT_TIMESTAMP\) 
    RETURNING id`
	mock.ExpectQuery(query).
		WithArgs(req.UserID, req.Name, req.Email, req.Phone).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := customerRepo.CreateCustomer(req)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFetchCustomersByUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	customerRepo := NewCustomerRepository(db)
	query := `SELECT id, user_id, name, email, phone, created_at, updated_at 
              FROM customers WHERE user_id = \$1 ORDER BY created_at DESC`
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "email", "phone", "created_at", "updated_at"}).
		AddRow(1, 101, "Jane Doe", "jane.doe@example.com", "123-456-7890", time.Now(), time.Now()).
		AddRow(2, 101, "John Doe", "john.doe@example.com", "098-765-4321", time.Now(), time.Now())

	mock.ExpectQuery(query).
		WithArgs(101).
		WillReturnRows(rows)

	customers, err := customerRepo.FetchCustomersByUserId(101)

	assert.NoError(t, err)
	assert.Len(t, customers, 2)
	assert.Equal(t, "Jane Doe", customers[0].Name)
	assert.Equal(t, "John Doe", customers[1].Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
