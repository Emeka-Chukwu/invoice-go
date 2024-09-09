package auth_repository

import (
	"go-invoice/domain"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	authRepo := &authRepository{Db: db}
	req := domain.CreateUserRequestDto{
		Email:     "john.doe@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "hashedpassword",
	}
	query := `insert into users \(email, first_name, last_name, password_hash\) values \(\$1, \$2, \$3, \$4\)
	returning id, email, first_name, last_name, created_at, updated_at`
	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "created_at", "updated_at"}).
		AddRow(1, "john.doe@example.com", "John", "Doe", time.Now(), time.Now())

	mock.ExpectQuery(query).
		WithArgs(req.Email, req.FirstName, req.LastName, req.Password).
		WillReturnRows(rows)

	userResponse, err := authRepo.CreateUser(req)
	require.NoError(t, err)
	require.Equal(t, req.Email, userResponse.Email)
	require.Equal(t, req.FirstName, userResponse.FirstName)
	require.Equal(t, req.LastName, userResponse.LastName)
	require.NotZero(t, userResponse.Id)
	require.NotZero(t, userResponse.CreatedAt)
	require.NotZero(t, userResponse.UpdatedAt)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFetchUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	authRepo := &authRepository{Db: db}
	email := "john.doe@example.com"

	query := `select id, email, first_name, last_name, password_hash, created_at, updated_at from users where email = \$1`
	rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "password_hash", "created_at", "updated_at"}).
		AddRow(1, "john.doe@example.com", "John", "Doe", "hashedpassword", time.Now(), time.Now())

	mock.ExpectQuery(query).
		WithArgs(email).
		WillReturnRows(rows)

	userResponse, err := authRepo.FetchUserByEmail(email)

	require.NoError(t, err)
	require.Equal(t, email, userResponse.Email)
	require.Equal(t, "John", userResponse.FirstName)
	require.Equal(t, "Doe", userResponse.LastName)
	require.NotZero(t, userResponse.Id)
	require.NotZero(t, userResponse.CreatedAt)
	require.NotZero(t, userResponse.UpdatedAt)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
