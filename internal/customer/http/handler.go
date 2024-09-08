package customer_http

import (
	"go-invoice/domain"
	customer_usecase "go-invoice/internal/customer/usecase"
	"go-invoice/security"
	"go-invoice/util"

	"github.com/gin-gonic/gin"
)

type CustomerHandler interface {
	CreateCustomer(ctx *gin.Context)
	FetchCustomers(ctx *gin.Context)
}

type customerHandler struct {
	usecase customer_usecase.CustomerUsecase
}

func (ch *customerHandler) CreateCustomer(ctx *gin.Context) {
	authPayload := security.GetAuthsPayload(ctx)
	payload := util.GetBody[domain.CustomerDTO](ctx)
	payload.UserID = authPayload.UserId
	status, resp, err := ch.usecase.CreateCustomer(payload)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, gin.H{"data": resp})
}

func (ch *customerHandler) FetchCustomers(ctx *gin.Context) {
	authPayload := security.GetAuthsPayload(ctx)
	status, resp, err := ch.usecase.FetchCustomers(authPayload.UserId)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, gin.H{"data": resp})
}

func NewCustomerHandlers(usecase customer_usecase.CustomerUsecase) CustomerHandler {
	return &customerHandler{usecase: usecase}
}
