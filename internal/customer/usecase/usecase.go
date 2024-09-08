package customer_usecase

import (
	"go-invoice/domain"
	customer_repository "go-invoice/internal/customer/repository"
	"net/http"
)

type CustomerUsecase interface {
	CreateCustomer(req domain.CustomerDTO) (int, int64, error)
	FetchCustomers(user int) (int, []domain.CustomerResponse, error)
}

type customerUsecase struct {
	Repo customer_repository.CustomerRepository
}

func NewCustomerUsecase(Repo customer_repository.CustomerRepository) CustomerUsecase {
	return customerUsecase{Repo: Repo}
}

// CreateCustomer implements CustomerUsecase.
func (c customerUsecase) CreateCustomer(req domain.CustomerDTO) (int, int64, error) {
	resp, err := c.Repo.CreateCustomer(req)
	if err != nil {
		return http.StatusInternalServerError, 0, err
	}
	return http.StatusCreated, resp, nil
}

// FetchCustomers implements CustomerUsecase.
func (c customerUsecase) FetchCustomers(user int) (int, []domain.CustomerResponse, error) {
	resp, err := c.Repo.FetchCustomersByUserId(user)
	if err != nil {
		return http.StatusInternalServerError, []domain.CustomerResponse{}, err
	}
	return http.StatusOK, resp, nil
}
