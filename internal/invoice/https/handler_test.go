package invoice_https_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-invoice/domain"
	mockInvoiceUse "go-invoice/internal/invoice/usecase/mock"
	"go-invoice/security"
	"net/http"
	"net/http/httptest"
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
		name           string
		requestBody    domain.CreateInvoiceRequestDTO
		buildStubs     func(store *mockInvoiceUse.MockInvoiceUsecase)
		checkResponse  func(status, expectedStatus int, err error)
		expectedStatus int
		setupAuth      func(t *testing.T, request *http.Request, tokenMaker security.Maker)
	}{
		{
			name:        "Ok",
			requestBody: req,
			buildStubs: func(store *mockInvoiceUse.MockInvoiceUsecase) {
				store.EXPECT().CreateInvoiceWithItems(gomock.Any()).
					Times(1).
					Return(201, 1, nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.NoError(t, err)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
		},
		{
			name:           "InternalServerError",
			requestBody:    req,
			expectedStatus: http.StatusInternalServerError,
			buildStubs: func(store *mockInvoiceUse.MockInvoiceUsecase) {
				store.EXPECT().CreateInvoiceWithItems(gomock.Any()).
					Times(1).
					Return(http.StatusInternalServerError, 0, sql.ErrConnDone)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
		},
		{
			name:           "Unauthorized",
			requestBody:    req,
			expectedStatus: http.StatusUnauthorized,
			buildStubs: func(store *mockInvoiceUse.MockInvoiceUsecase) {
				store.EXPECT().CreateInvoiceWithItems(gomock.Any()).
					Times(0)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockInvoiceUse.NewMockInvoiceUsecase(ctrl)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)
			ctx := context.Background()
			tt.buildStubs(store)
			request, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/invoice/create", bytes.NewReader(data))
			require.NoError(t, err)
			tt.setupAuth(t, request, server.TokenMaker)
			server.Router.ServeHTTP(recorder, request)
			require.Equal(t, tt.expectedStatus, recorder.Code)
		})

	}
}

func TestFetchInvoices(t *testing.T) {
	tests := []struct {
		name           string
		buildStubs     func(store *mockInvoiceUse.MockInvoiceUsecase)
		checkResponse  func(status, expectedStatus int, err error)
		expectedStatus int
		page           int
		limit          int
		setupAuth      func(t *testing.T, request *http.Request, tokenMaker security.Maker)
	}{
		{
			name: "Ok",
			buildStubs: func(store *mockInvoiceUse.MockInvoiceUsecase) {
				store.EXPECT().GetPagination(gomock.Any()).Times(1)
				store.EXPECT().FetchInvoices(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(200, []domain.InvoiceResponse{}, nil)

			},
			expectedStatus: http.StatusOK,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.NoError(t, err)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
		},
		{
			name:           "InternalServerError",
			expectedStatus: http.StatusInternalServerError,
			buildStubs: func(store *mockInvoiceUse.MockInvoiceUsecase) {
				store.EXPECT().GetPagination(gomock.Any()).Times(1)
				store.EXPECT().FetchInvoices(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(http.StatusInternalServerError, []domain.InvoiceResponse{}, sql.ErrConnDone)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
		},
		{
			name:           "Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			buildStubs: func(store *mockInvoiceUse.MockInvoiceUsecase) {
				store.EXPECT().GetPagination(gomock.Any()).Times(0)
				store.EXPECT().FetchInvoices(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0).
					Return(http.StatusInternalServerError, []domain.InvoiceResponse{}, sql.ErrConnDone)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockInvoiceUse.NewMockInvoiceUsecase(ctrl)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			ctx := context.Background()
			tt.buildStubs(store)
			request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/invoice/invoices", nil)
			require.NoError(t, err)
			tt.setupAuth(t, request, server.TokenMaker)
			server.Router.ServeHTTP(recorder, request)
			require.Equal(t, tt.expectedStatus, recorder.Code)
		})

	}
}

func TestUpdateInvoiceStatus(t *testing.T) {
	tests := []struct {
		name           string
		buildStubs     func(store *mockInvoiceUse.MockInvoiceUsecase)
		checkResponse  func(status, expectedStatus int, err error)
		expectedStatus int
		Invoice        int
		limit          int
		setupAuth      func(t *testing.T, request *http.Request, tokenMaker security.Maker)
	}{
		{
			name: "Ok",
			buildStubs: func(store *mockInvoiceUse.MockInvoiceUsecase) {

				store.EXPECT().UpdateInvoiceStatus(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(200, nil)

			},
			expectedStatus: http.StatusOK,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.NoError(t, err)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			Invoice: 2,
		},
		{
			name:           "InternalServerError",
			expectedStatus: http.StatusInternalServerError,
			buildStubs: func(store *mockInvoiceUse.MockInvoiceUsecase) {
				store.EXPECT().UpdateInvoiceStatus(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(500, sql.ErrConnDone)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockInvoiceUse.NewMockInvoiceUsecase(ctrl)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			ctx := context.Background()
			tt.buildStubs(store)
			request, err := http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("/api/v1/invoice/update/%d", tt.Invoice), nil)
			require.NoError(t, err)
			tt.setupAuth(t, request, server.TokenMaker)
			server.Router.ServeHTTP(recorder, request)
			require.Equal(t, tt.expectedStatus, recorder.Code)
		})

	}
}

func TestFetchInvoiceStats(t *testing.T) {
	tests := []struct {
		name           string
		buildStubs     func(store *mockInvoiceUse.MockInvoiceUsecase)
		checkResponse  func(status, expectedStatus int, err error)
		expectedStatus int
		Invoice        int
		limit          int
		setupAuth      func(t *testing.T, request *http.Request, tokenMaker security.Maker)
	}{
		{
			name: "Ok",
			buildStubs: func(store *mockInvoiceUse.MockInvoiceUsecase) {
				store.EXPECT().FetchInvoiceStats(gomock.Any()).
					Times(1).
					Return(200, make(map[string]domain.InvoiceStats), nil)

			},
			expectedStatus: http.StatusOK,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.NoError(t, err)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			Invoice: 2,
		},
		{
			name:           "InternalServerError",
			expectedStatus: http.StatusInternalServerError,
			buildStubs: func(store *mockInvoiceUse.MockInvoiceUsecase) {
				store.EXPECT().FetchInvoiceStats(gomock.Any()).
					Times(1).
					Return(http.StatusInternalServerError, make(map[string]domain.InvoiceStats), sql.ErrConnDone)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
		},
		{
			name:           "Unauthorized",
			expectedStatus: http.StatusUnauthorized,
			buildStubs: func(store *mockInvoiceUse.MockInvoiceUsecase) {
				store.EXPECT().GetPagination(gomock.Any()).Times(0)
				store.EXPECT().FetchInvoiceStats(gomock.Any()).
					Times(0).
					Return(http.StatusInternalServerError, make(map[string]domain.InvoiceStats), sql.ErrConnDone)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker security.Maker) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockInvoiceUse.NewMockInvoiceUsecase(ctrl)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			ctx := context.Background()
			tt.buildStubs(store)
			request, err := http.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/invoice/invoices/stats", nil)
			require.NoError(t, err)
			tt.setupAuth(t, request, server.TokenMaker)
			server.Router.ServeHTTP(recorder, request)
			require.Equal(t, tt.expectedStatus, recorder.Code)
		})

	}
}
