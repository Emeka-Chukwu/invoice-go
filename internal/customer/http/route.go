package customer_http

import (
	customer_usecase "go-invoice/internal/customer/usecase"

	"github.com/gin-gonic/gin"
)

func NewCustomerRoutes(router *gin.RouterGroup, usecase customer_usecase.CustomerUsecase) {
	customerUsecase := NewCustomerHandlers(usecase)
	route := router.Group("/customer")
	route.POST("/create", customerUsecase.CreateCustomer)
	route.GET("/customers", customerUsecase.FetchCustomers)
}
