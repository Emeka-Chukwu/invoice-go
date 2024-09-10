package invoice_usecase

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"go-invoice/domain"
	bankinfo_repository "go-invoice/internal/bank_info/repository"
	customer_repository "go-invoice/internal/customer/repository"
	"go-invoice/internal/invoice/helper"
	invoice_repository "go-invoice/internal/invoice/repository"
	"net/http"

	"github.com/jung-kurt/gofpdf"
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
	DownloadSingleInvoice(id, customerId, userId int) (int, []byte, error)
}

type invoiceUsecase struct {
	Repo         invoice_repository.InvoiceRepository
	RepoCustomer customer_repository.CustomerRepository
	RepoBank     bankinfo_repository.BankInfoRepository
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
	return http.StatusCreated, resp, nil
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

func NewInvoiceUsecase(Repo invoice_repository.InvoiceRepository,
	RepoCustomer customer_repository.CustomerRepository,
	RepoBank bankinfo_repository.BankInfoRepository) InvoiceUsecase {
	return &invoiceUsecase{Repo: Repo, RepoCustomer: RepoCustomer, RepoBank: RepoBank}
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

func (i invoiceUsecase) DownloadSingleInvoice(id, customerId, userId int) (int, []byte, error) {
	invoiceResp, err := i.Repo.FetchInvoiceWithItems(id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	customers, err := i.RepoCustomer.FetchCustomersByUserId(userId)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	var customer domain.CustomerResponse
	if len(customers) == 0 {
		return http.StatusInternalServerError, nil, errors.New("no customer selected")
	}
	customer = customers[0]
	bankInfo, err := i.RepoBank.FetchBankInformation(int64(userId))
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFillColor(252, 242, 244)
	pdf.Rect(0, 0, 210, 297, "F")
	pdf.SetFillColor(255, 255, 255)
	pdf.RoundedRect(10, 10, 190, 277, 5, "F", "")
	helper.AddSenderCustomerInfo(pdf, customer)
	helper.AddInvoiceDetails(pdf, invoiceResp)
	helper.AddItemsTable(pdf, invoiceResp.Items)
	subTotal := 0.0
	for _, item := range invoiceResp.Items {
		subTotal += float64(item.UnitPrice) * float64(item.Quantity)
	}
	discountedAmount := subTotal * 0.01
	totalAmount := subTotal - discountedAmount

	helper.AddTotals(pdf, fmt.Sprintf("%f", subTotal), "DISCOUNT(1%)", fmt.Sprintf("%f", discountedAmount), fmt.Sprintf("%f", totalAmount))
	helper.AddPaymentInfo(pdf, bankInfo)
	helper.AddNote(pdf, "Thank you for your patronage")
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, buf.Bytes(), nil

}
