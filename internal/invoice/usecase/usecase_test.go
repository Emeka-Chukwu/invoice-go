package invoice_usecase

import (
	"database/sql"
	"go-invoice/domain"
	mockInvoice "go-invoice/internal/invoice/repository/mock"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func randomCreateInvoiceDTO() domain.CreateInvoiceRequestDTO {
	return domain.CreateInvoiceRequestDTO{
		InvoiceDTO: domain.InvoiceDTO{
			UserID:        1,
			CustomerID:    101,
			InvoiceNumber: "INV-2024-001",
			Status:        "unpaid",
			IssueDate:     time.Now(),
			DueDate:       time.Now().AddDate(0, 1, 0),
			TotalAmount:   1500.75,
		},
		CreateInvoiceItem: domain.CreateInvoiceItem{
			{
				Description: "Product 1",
				Quantity:    2,
				UnitPrice:   500.25,
			},
			{
				Description: "Product 2",
				Quantity:    1,
				UnitPrice:   500.50,
			},
		},
	}
}

func TestCreateUser(t *testing.T) {

	req := randomCreateInvoiceDTO()

	tests := []struct {
		name          string
		requestBody   domain.CreateInvoiceRequestDTO
		buildStubs    func(store *mockInvoice.MockInvoiceRepository)
		checkResponse func(status, expectedStatus int, err error)
		StatusCode    int
	}{
		{
			name:        "Ok",
			requestBody: req,
			buildStubs: func(store *mockInvoice.MockInvoiceRepository) {
				store.EXPECT().CreateInvoiceWithItems(gomock.Any()).
					Times(1).
					Return(2, nil)
			},
			StatusCode: http.StatusCreated,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.NoError(t, err)
			},
		},
		{
			name:        "InternalServerError",
			requestBody: req,
			StatusCode:  http.StatusInternalServerError,
			buildStubs: func(store *mockInvoice.MockInvoiceRepository) {
				store.EXPECT().CreateInvoiceWithItems(gomock.Any()).
					Times(1).
					Return(0, sql.ErrConnDone)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status, expectedStatus)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockInvoice.NewMockInvoiceRepository(ctrl)
			handler := NewInvoiceUsecase(store)
			tt.buildStubs(store)
			status, _, err := handler.CreateInvoiceWithItems(tt.requestBody)
			tt.checkResponse(status, tt.StatusCode, err)
		})

	}
}

func TestFetchInvoices(t *testing.T) {
	tests := []struct {
		name          string
		buildStubs    func(store *mockInvoice.MockInvoiceRepository)
		checkResponse func(status, expectedStatus int, err error)
		StatusCode    int
		userId        int
		page          int
		limit         int
	}{
		{
			name:   "Ok",
			userId: 1,
			page:   1,
			limit:  10,
			buildStubs: func(store *mockInvoice.MockInvoiceRepository) {
				store.EXPECT().FetchInvoices(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return([]domain.InvoiceResponse{}, nil)
			},
			StatusCode: http.StatusOK,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.NoError(t, err)
			},
		},
		{
			name:       "InternalServerError",
			userId:     1,
			page:       1,
			limit:      10,
			StatusCode: http.StatusInternalServerError,
			buildStubs: func(store *mockInvoice.MockInvoiceRepository) {
				store.EXPECT().FetchInvoices(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return([]domain.InvoiceResponse{}, sql.ErrConnDone)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status, expectedStatus)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockInvoice.NewMockInvoiceRepository(ctrl)
			handler := NewInvoiceUsecase(store)
			tt.buildStubs(store)
			status, _, err := handler.FetchInvoices(tt.userId, int64(tt.page), int64(tt.limit))
			tt.checkResponse(status, tt.StatusCode, err)
		})

	}
}

func TestFetchInvoiceStats(t *testing.T) {
	tests := []struct {
		name          string
		buildStubs    func(store *mockInvoice.MockInvoiceRepository)
		checkResponse func(status, expectedStatus int, err error)
		StatusCode    int
		userId        int
	}{
		{
			name:   "Ok",
			userId: 1,

			buildStubs: func(store *mockInvoice.MockInvoiceRepository) {
				store.EXPECT().FetchInvoiceStats(gomock.Any()).
					Times(1).
					Return(make(map[string]domain.InvoiceStats), nil)
			},
			StatusCode: http.StatusOK,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.NoError(t, err)
			},
		},
		{
			name:   "InternalServerError",
			userId: 1,

			StatusCode: http.StatusInternalServerError,
			buildStubs: func(store *mockInvoice.MockInvoiceRepository) {
				store.EXPECT().FetchInvoiceStats(gomock.Any()).
					Times(1).
					Return(make(map[string]domain.InvoiceStats), sql.ErrConnDone)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status, expectedStatus)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockInvoice.NewMockInvoiceRepository(ctrl)
			handler := NewInvoiceUsecase(store)
			tt.buildStubs(store)
			status, _, err := handler.FetchInvoiceStats(tt.userId)
			tt.checkResponse(status, tt.StatusCode, err)
		})

	}
}

func TestFetchInvoiceWithItems(t *testing.T) {
	tests := []struct {
		name          string
		buildStubs    func(store *mockInvoice.MockInvoiceRepository)
		checkResponse func(status, expectedStatus int, err error)
		StatusCode    int
		userId        int
		invoiceId     int
	}{
		{
			name:      "Ok",
			userId:    1,
			invoiceId: 1,
			buildStubs: func(store *mockInvoice.MockInvoiceRepository) {
				resp := domain.InvoiceResponse{}
				resp.UserID = 1
				store.EXPECT().FetchInvoiceWithItems(gomock.Any()).
					Times(1).
					Return(resp, nil)
			},
			StatusCode: http.StatusOK,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.NoError(t, err)
			},
		},
		{
			name:      "ForbiddenResource",
			userId:    1,
			invoiceId: 1,
			buildStubs: func(store *mockInvoice.MockInvoiceRepository) {

				store.EXPECT().FetchInvoiceWithItems(gomock.Any()).
					Times(1).
					Return(domain.InvoiceResponse{}, nil)
			},
			StatusCode: http.StatusForbidden,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
		},
		{
			name:       "InternalServerError",
			userId:     1,
			invoiceId:  1,
			StatusCode: http.StatusInternalServerError,
			buildStubs: func(store *mockInvoice.MockInvoiceRepository) {
				store.EXPECT().FetchInvoiceWithItems(gomock.Any()).
					Times(1).
					Return(domain.InvoiceResponse{}, sql.ErrConnDone)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status, expectedStatus)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockInvoice.NewMockInvoiceRepository(ctrl)
			handler := NewInvoiceUsecase(store)
			tt.buildStubs(store)
			status, _, err := handler.FetchInvoiceWithItems(tt.userId, tt.invoiceId)
			tt.checkResponse(status, tt.StatusCode, err)
		})

	}
}

func TestUpdateInvoiceStatus(t *testing.T) {
	tests := []struct {
		name          string
		buildStubs    func(store *mockInvoice.MockInvoiceRepository)
		checkResponse func(status, expectedStatus int, err error)
		StatusCode    int
		userId        int
		invoiceId     int
		status        string
	}{
		{
			name:      "Ok",
			userId:    1,
			invoiceId: 1,
			status:    "paid",
			buildStubs: func(store *mockInvoice.MockInvoiceRepository) {

				store.EXPECT().UpdateInvoiceStatus(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			StatusCode: http.StatusOK,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.NoError(t, err)
			},
		},

		{
			name:       "InternalServerError",
			userId:     1,
			invoiceId:  1,
			status:     "paid",
			StatusCode: http.StatusInternalServerError,
			buildStubs: func(store *mockInvoice.MockInvoiceRepository) {
				store.EXPECT().UpdateInvoiceStatus(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status, expectedStatus)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockInvoice.NewMockInvoiceRepository(ctrl)
			handler := NewInvoiceUsecase(store)
			tt.buildStubs(store)
			status, err := handler.UpdateInvoiceStatus(tt.status, tt.invoiceId, tt.userId)
			tt.checkResponse(status, tt.StatusCode, err)
		})

	}
}
