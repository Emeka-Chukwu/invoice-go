package bankinfo_http

import (
	bankinfo_usecase "go-invoice/internal/bank_info/usecase"

	"github.com/gin-gonic/gin"
)

func NewBankInfoRoutes(router *gin.RouterGroup, usecase bankinfo_usecase.BankInfoUsecase) {
	bankInfo := NewBankInfoHandlers(usecase)
	route := router.Group("/bank-info")
	route.POST("/create", bankInfo.CreateBankInfo)
	route.GET("/info", bankInfo.FetchBankInfo)
}
