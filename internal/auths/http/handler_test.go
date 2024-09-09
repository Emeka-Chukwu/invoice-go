package auth_http_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"go-invoice/domain"
	mockAuthUse "go-invoice/internal/auths/http/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {

	req := domain.CreateUserRequestDto{
		Email:     "email@yopmail.com",
		FirstName: "FirstName",
		LastName:  "LastName",
		Password:  "Password",
	}

	tests := []struct {
		name           string
		requestBody    domain.CreateUserRequestDto
		buildStubs     func(store *mockAuthUse.MockAuthUsecase)
		checkResponse  func(status, expectedStatus int, err error)
		expectedStatus int
	}{
		{
			name:        "Ok",
			requestBody: req,
			buildStubs: func(store *mockAuthUse.MockAuthUsecase) {
				store.EXPECT().CreateUser(gomock.Any()).
					Times(1).
					Return(201, gin.H{}, nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.NoError(t, err)
			},
		},
		{
			name:           "InternalServerError",
			requestBody:    req,
			expectedStatus: http.StatusInternalServerError,
			buildStubs: func(store *mockAuthUse.MockAuthUsecase) {
				store.EXPECT().CreateUser(gomock.Any()).
					Times(1).
					Return(http.StatusInternalServerError, gin.H{}, sql.ErrConnDone)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockAuthUse.NewMockAuthUsecase(ctrl)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)
			ctx := context.Background()
			tt.buildStubs(store)
			request, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/auths/register", bytes.NewReader(data))
			require.NoError(t, err)
			server.Router.ServeHTTP(recorder, request)
			require.Equal(t, tt.expectedStatus, recorder.Code)
		})

	}
}

func TestLoginUser(t *testing.T) {
	req := domain.LoginRequestDto{
		Email:    "email@yopmail.com",
		Password: "Password",
	}

	tests := []struct {
		name           string
		requestBody    domain.LoginRequestDto
		buildStubs     func(store *mockAuthUse.MockAuthUsecase)
		checkResponse  func(status, expectedStatus int, err error)
		expectedStatus int
	}{
		{
			name:        "Ok",
			requestBody: req,
			buildStubs: func(store *mockAuthUse.MockAuthUsecase) {

				store.EXPECT().LoginUser(gomock.Any()).
					Times(1).
					Return(200, gin.H{}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.NoError(t, err)
			},
		},
		{
			name:           "InternalServerError",
			requestBody:    req,
			expectedStatus: http.StatusInternalServerError,
			buildStubs: func(store *mockAuthUse.MockAuthUsecase) {
				store.EXPECT().LoginUser(gomock.Any()).
					Times(1).
					Return(http.StatusInternalServerError, gin.H{}, sql.ErrConnDone)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockAuthUse.NewMockAuthUsecase(ctrl)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)
			ctx := context.Background()
			tt.buildStubs(store)
			request, err := http.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/auths/login", bytes.NewReader(data))
			require.NoError(t, err)
			server.Router.ServeHTTP(recorder, request)
			require.Equal(t, tt.expectedStatus, recorder.Code)
		})

	}
}
