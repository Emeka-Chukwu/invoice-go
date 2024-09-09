package auth_http

import (
	"fmt"
	"go-invoice/domain"
	auth_usecase "go-invoice/internal/auths/usecase"
	"go-invoice/util"

	"github.com/gin-gonic/gin"
)

type AuthsHandler interface {
	RegisterUser(*gin.Context)
	LoginUser(*gin.Context)
	FetchUser(*gin.Context)
}

type authsHandler struct {
	authusecase auth_usecase.AuthUsecase
}

// FetchUser implements AuthsHandler.
func (ah authsHandler) FetchUser(ctx *gin.Context) {
	payload := util.GetUrlParams[domain.EmailPayload](ctx)
	status, user, err := ah.authusecase.FetchUser(payload.Email)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, gin.H{"data": user})
}

// LoginUser implements AuthsHandler.
func (ah authsHandler) LoginUser(ctx *gin.Context) {
	payload := util.GetBody[domain.LoginRequestDto](ctx)
	status, resp, err := ah.authusecase.LoginUser(payload)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, resp)
}

// RegisterUser implements AuthsHandler.
func (ah authsHandler) RegisterUser(ctx *gin.Context) {
	payload := util.GetBody[domain.CreateUserRequestDto](ctx)
	fmt.Println(payload)
	status, resp, err := ah.authusecase.CreateUser(payload)
	fmt.Println(resp)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, resp)
}

func NewAuthsHandlers(authusecase auth_usecase.AuthUsecase) AuthsHandler {
	return authsHandler{authusecase: authusecase}
}
