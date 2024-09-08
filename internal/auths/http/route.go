package auth_http

import (
	"go-invoice/domain"
	auth_usecase "go-invoice/internal/auths/usecase"
	"go-invoice/middleware"

	"github.com/gin-gonic/gin"
)

func NewAuthRoutes(router *gin.RouterGroup, usecase auth_usecase.AuthUsecase) {
	authsHandler := NewAuthsHandlers(usecase)
	route := router.Group("/auths")
	route.POST("/register", middleware.ValidatorMiddleware[domain.CreateUserRequestDto], authsHandler.RegisterUser)
	route.POST("/login", middleware.ValidatorMiddleware[domain.LoginRequestDto], authsHandler.LoginUser)
}
