package customer_repository

import (
	"context"
	"database/sql"
	"go-invoice/domain"
	"go-invoice/util"
)

type customerRepository struct {
	Db *sql.DB
}

type CustomerRepository interface {
	CreateCustomer(req domain.CustomerDTO) (int64, error)
	FetchCustomersByUserId(userId int) ([]domain.CustomerResponse, error)
}

func (cr *customerRepository) CreateCustomer(req domain.CustomerDTO) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), util.DbTimeout)
	defer cancel()
	query := `
    INSERT INTO customers (user_id, name, email, phone, created_at, updated_at) 
    VALUES ($1, $2, $3, $4,  CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
    RETURNING id`

	var id int64
	err := cr.Db.QueryRowContext(ctx, query, req.UserID, req.Name, req.Email, req.Phone).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (cr *customerRepository) FetchCustomersByUserId(userId int) ([]domain.CustomerResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), util.DbTimeout)
	defer cancel()
	query := `SELECT id, user_id, name, email, phone, created_at, updated_at 
              FROM customers WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := cr.Db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var customers []domain.CustomerResponse
	for rows.Next() {
		var customer domain.CustomerResponse
		err := rows.Scan(
			&customer.ID, &customer.UserID, &customer.Name, &customer.Email, &customer.Phone,
			&customer.CreatedAt, &customer.UpdatedAt)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{Db: db}
}
