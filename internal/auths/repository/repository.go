package auth_repository

import (
	"context"
	"database/sql"
	"go-invoice/domain"
	"go-invoice/util"
)

type authRepository struct {
	Db *sql.DB
}

func (a *authRepository) CreateUser(req domain.CreateUserRequestDto) (domain.UserReponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), util.DbTimeout)
	defer cancel()
	stmt := `insert into users (email, first_name, last_name, password_hash) values ($1, $2, $3, $4)
	returning id, email, first_name, last_name, created_at, updated_at`
	var response domain.UserReponse
	err := a.Db.QueryRowContext(ctx, stmt, req.Email, req.FirstName, req.LastName, req.Password).
		Scan(&response.Id, &response.Email, &response.FirstName, &response.LastName, &response.CreatedAt, &response.UpdatedAt)
	return response, err
}

func (a *authRepository) FetchUserByEmail(email string) (domain.UserReponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), util.DbTimeout)
	defer cancel()
	stmt := `select  id, email, first_name, last_name, password_hash, created_at, updated_at from users where email = $1`
	var response domain.UserReponse
	err := a.Db.QueryRowContext(ctx, stmt, email).
		Scan(&response.Id, &response.Email, &response.FirstName, &response.LastName, &response.Password, &response.CreatedAt, &response.UpdatedAt)
	return response, err
}

type AuthRepository interface {
	CreateUser(req domain.CreateUserRequestDto) (domain.UserReponse, error)
	FetchUserByEmail(email string) (domain.UserReponse, error)
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{Db: db}
}
