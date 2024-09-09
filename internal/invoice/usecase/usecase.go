package invoice_usecase

import (
	"database/sql"
	"errors"
	"go-invoice/domain"
	invoice_repository "go-invoice/internal/invoice/repository"
	"net/http"
)

type InvoiceUsecase interface {
	CreateInvoiceWithItems(req domain.CreateInvoiceRequestDTO) (int, int, error)
	FetchInvoicesWithItems(userId int, limit, offset int64) (int, []domain.InvoiceResponse, error)
	FetchInvoices(userId int, page, limit int64) (int, []domain.InvoiceResponse, error)
	FetchInvoiceWithItems(invoiceId, userId int) (int, domain.InvoiceResponse, error)
	DeleteInvoiceItems(itemIds []int) (int, error)
	FetchInvoiceStats(userId int) (int, map[string]domain.InvoiceStats, error)
	GetPagination(payload domain.PaginationDTO) domain.PaginationDTO
	GenerateInvoicePDF() ([]byte, error)
	UpdateInvoiceStatus(status string, id, userId int) (int, error)
}

type invoiceUsecase struct {
	Repo invoice_repository.InvoiceRepository
}

// GenerateInvoicePDF implements InvoiceUsecase.
func (i *invoiceUsecase) GenerateInvoicePDF() ([]byte, error) {
	panic("unimplemented")
}

// CreateInvoiceWithItems implements InvoiceUsecase.
func (i invoiceUsecase) CreateInvoiceWithItems(req domain.CreateInvoiceRequestDTO) (int, int, error) {
	var totalAmount float64 = 0
	for _, item := range req.CreateInvoiceItem {
		totalAmount += item.UnitPrice * float64(item.Quantity)
	}
	req.TotalAmount = totalAmount
	resp, err := i.Repo.CreateInvoiceWithItems(req)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}
	return http.StatusOK, resp, nil
}

// DeleteInvoiceItems implements InvoiceUsecase.
func (i invoiceUsecase) DeleteInvoiceItems(itemIds []int) (int, error) {
	panic("unimplemented")
}

// FetchInvoiceStats implements InvoiceUsecase.
func (i invoiceUsecase) FetchInvoiceStats(userId int) (int, map[string]domain.InvoiceStats, error) {
	resp, err := i.Repo.FetchInvoiceStats(userId)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}
	return http.StatusOK, resp, nil
}

// FetchInvoiceWithItems implements InvoiceUsecase.
func (i invoiceUsecase) FetchInvoiceWithItems(invoiceId, userId int) (int, domain.InvoiceResponse, error) {
	resp, err := i.Repo.FetchInvoiceWithItems(invoiceId)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, resp, err
	} else if err != nil {
		return http.StatusInternalServerError, resp, err
	}
	if resp.UserID != userId {
		return http.StatusForbidden, domain.InvoiceResponse{}, errors.New("forbidden resource")
	}
	return http.StatusOK, resp, nil
}

// FetchInvoices implements InvoiceUsecase.
func (i invoiceUsecase) FetchInvoices(userId int, page int64, limit int64) (int, []domain.InvoiceResponse, error) {
	resp, err := i.Repo.FetchInvoices(userId, page, limit)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}
	return http.StatusOK, resp, nil

}

// FetchInvoicesWithItems implements InvoiceUsecase.
func (i invoiceUsecase) FetchInvoicesWithItems(userId int, limit int64, offset int64) (int, []domain.InvoiceResponse, error) {
	resp, err := i.Repo.FetchInvoicesWithItems(userId, limit, offset)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}
	return http.StatusOK, resp, nil
}

func NewInvoiceUsecase(Repo invoice_repository.InvoiceRepository) InvoiceUsecase {
	return &invoiceUsecase{Repo: Repo}
}

type QueryParams struct {
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
	SortBy  string `form:"sort_by"`
	OrderBy string `form:"order_by"`
}

func (i invoiceUsecase) GetPagination(payload domain.PaginationDTO) domain.PaginationDTO {
	page, limit := payload.Page, payload.Limit
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	page = (page - 1) * limit
	return domain.PaginationDTO{Limit: limit, Page: page}
}

func (i invoiceUsecase) UpdateInvoiceStatus(status string, id, userId int) (int, error) {
	err := i.Repo.UpdateInvoiceStatus(id, userId, status)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
