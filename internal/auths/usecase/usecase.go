package auth_usecase

import (
	"database/sql"
	"fmt"
	"go-invoice/domain"
	auth_repository "go-invoice/internal/auths/repository"
	"go-invoice/security"
	"go-invoice/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthUsecase interface {
	CreateUser(req domain.CreateUserRequestDto) (int, gin.H, error)
	LoginUser(req domain.LoginRequestDto) (int, gin.H, error)
	FetchUser(email string) (int, domain.UserReponse, error)
}

type authUsecase struct {
	Repo     auth_repository.AuthRepository
	Security security.Maker
	config   util.Config
}

// CreateUser implements AuthUsecase.
func (au authUsecase) CreateUser(req domain.CreateUserRequestDto) (int, gin.H, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return http.StatusBadRequest, gin.H{}, err
	}
	req.Password = hashedPassword
	user, err := au.Repo.CreateUser(req)
	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, gin.H{}, err
	}
	token, payload, err := au.Security.CreateToken(user.Id, au.config.AccessTokenDuration, true)
	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, gin.H{}, err
	}
	return http.StatusCreated, gin.H{"data": user, "token": token, "token_payload": payload}, nil
}

// FetchUser implements AuthUsecase.
func (au authUsecase) FetchUser(email string) (int, domain.UserReponse, error) {
	user, err := au.Repo.FetchUserByEmail(email)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, domain.UserReponse{}, err
	} else if err != nil {
		return http.StatusInternalServerError, domain.UserReponse{}, err
	}
	return http.StatusOK, user, nil
}

// LoginUser implements AuthUsecase.
func (au authUsecase) LoginUser(req domain.LoginRequestDto) (int, gin.H, error) {
	user, err := au.Repo.FetchUserByEmail(req.Email)
	if err != nil {
		return http.StatusUnauthorized, gin.H{"error": "invalid email or password"}, err
	}
	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		return http.StatusUnauthorized, gin.H{"error": "invalid email or password"}, err
	}
	token, payload, err := au.Security.CreateToken(user.Id, au.config.AccessTokenDuration, true)
	if err != nil {
		return http.StatusInternalServerError, gin.H{}, err
	}
	return http.StatusOK, gin.H{"data": user, "token": token, "token_payload": payload}, nil
}

func NewAuthUsecase(Repo auth_repository.AuthRepository, Security security.Maker, config util.Config) AuthUsecase {
	return authUsecase{Repo: Repo, Security: Security, config: config}
}
