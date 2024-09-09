package auth_usecase

import (
	"database/sql"
	"go-invoice/domain"
	mockAuth "go-invoice/internal/auths/usecase/mock"
	"go-invoice/security"
	"go-invoice/util"
	"net/http"
	"testing"
	"time"

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
		name          string
		requestBody   domain.CreateUserRequestDto
		buildStubs    func(store *mockAuth.MockAuthRepository)
		checkResponse func(status, expectedStatus int, err error)
		StatusCode    int
	}{
		{
			name:        "Ok",
			requestBody: req,
			buildStubs: func(store *mockAuth.MockAuthRepository) {
				store.EXPECT().CreateUser(gomock.Any()).
					Times(1).
					Return(domain.UserReponse{Id: 2}, nil)
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
			buildStubs: func(store *mockAuth.MockAuthRepository) {
				store.EXPECT().CreateUser(gomock.Any()).
					Times(1).
					Return(domain.UserReponse{}, sql.ErrConnDone)
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
			store := mockAuth.NewMockAuthRepository(ctrl)
			config := util.Config{
				TokenSymmetricKey:   util.RandomString(32),
				AccessTokenDuration: time.Minute,
			}
			tokenMaker, err := security.NewJWTMaker(config.TokenSymmetricKey)
			require.NoError(t, err)
			handler := NewAuthUsecase(store, tokenMaker, util.Config{AccessTokenDuration: time.Hour})
			tt.buildStubs(store)
			status, _, err := handler.CreateUser(tt.requestBody)
			tt.checkResponse(status, tt.StatusCode, err)
		})

	}
}

func TestLoginUser(t *testing.T) {

	req := domain.LoginRequestDto{
		Email:    "email@yopmail.com",
		Password: "Password",
	}

	tests := []struct {
		name          string
		requestBody   domain.LoginRequestDto
		buildStubs    func(store *mockAuth.MockAuthRepository)
		checkResponse func(status, expectedStatus int, err error)
		StatusCode    int
	}{
		{
			name:        "Ok",
			requestBody: req,
			buildStubs: func(store *mockAuth.MockAuthRepository) {
				hashedPassword, err := util.HashPassword(req.Password)
				require.NoError(t, err)
				store.EXPECT().FetchUserByEmail(gomock.Any()).
					Times(1).
					Return(domain.UserReponse{Id: 2, Password: hashedPassword, Email: req.Email}, nil)
			},
			StatusCode: http.StatusOK,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.NoError(t, err)
			},
		},
		{
			name:        "Unauthorized",
			requestBody: req,
			StatusCode:  http.StatusUnauthorized,
			buildStubs: func(store *mockAuth.MockAuthRepository) {
				store.EXPECT().FetchUserByEmail(gomock.Any()).
					Times(1).
					Return(domain.UserReponse{}, sql.ErrConnDone)
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
			store := mockAuth.NewMockAuthRepository(ctrl)
			config := util.Config{
				TokenSymmetricKey:   util.RandomString(32),
				AccessTokenDuration: time.Minute,
			}
			tokenMaker, err := security.NewJWTMaker(config.TokenSymmetricKey)
			require.NoError(t, err)
			handler := NewAuthUsecase(store, tokenMaker, util.Config{AccessTokenDuration: time.Hour})
			tt.buildStubs(store)
			status, _, err := handler.LoginUser(tt.requestBody)
			tt.checkResponse(status, tt.StatusCode, err)
		})

	}
}

func TestFetchUser(t *testing.T) {

	req := domain.EmailPayload{
		Email: "email@yopmail.com",
	}

	tests := []struct {
		name          string
		requestBody   domain.EmailPayload
		buildStubs    func(store *mockAuth.MockAuthRepository)
		checkResponse func(status, expectedStatus int, err error)
		StatusCode    int
	}{
		{
			name:        "Ok",
			requestBody: req,
			buildStubs: func(store *mockAuth.MockAuthRepository) {
				store.EXPECT().FetchUserByEmail(gomock.Any()).
					Times(1).
					Return(domain.UserReponse{Id: 2, Email: req.Email}, nil)
			},
			StatusCode: http.StatusOK,
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.NoError(t, err)
			},
		},
		{
			name:        "RecordNotFound",
			requestBody: req,
			StatusCode:  http.StatusNotFound,
			buildStubs: func(store *mockAuth.MockAuthRepository) {
				store.EXPECT().FetchUserByEmail(gomock.Any()).
					Times(1).
					Return(domain.UserReponse{}, sql.ErrNoRows)
			},
			checkResponse: func(status, expectedStatus int, err error) {
				require.NotEmpty(t, status)
				require.Equal(t, status, expectedStatus)
				require.Error(t, err)
			},
		},
		{
			name:        "InternalServerError",
			requestBody: req,
			StatusCode:  http.StatusInternalServerError,
			buildStubs: func(store *mockAuth.MockAuthRepository) {
				store.EXPECT().FetchUserByEmail(gomock.Any()).
					Times(1).
					Return(domain.UserReponse{}, sql.ErrConnDone)
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
			store := mockAuth.NewMockAuthRepository(ctrl)
			config := util.Config{
				TokenSymmetricKey:   util.RandomString(32),
				AccessTokenDuration: time.Minute,
			}
			tokenMaker, err := security.NewJWTMaker(config.TokenSymmetricKey)
			require.NoError(t, err)
			handler := NewAuthUsecase(store, tokenMaker, util.Config{AccessTokenDuration: time.Hour})
			tt.buildStubs(store)
			status, _, err := handler.FetchUser(tt.requestBody.Email)
			tt.checkResponse(status, tt.StatusCode, err)
		})

	}
}
