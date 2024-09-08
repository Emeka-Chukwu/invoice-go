package bankinfo_http

import (
	"go-invoice/domain"
	bankinfo_usecase "go-invoice/internal/bank_info/usecase"
	"go-invoice/security"
	"go-invoice/util"

	"github.com/gin-gonic/gin"
)

type BankInfoHandler interface {
	CreateBankInfo(ctx *gin.Context)
	FetchBankInfo(ctx *gin.Context)
}

type bankinfoHandler struct {
	bankusecase bankinfo_usecase.BankInfoUsecase
}

func (bh *bankinfoHandler) CreateBankInfo(ctx *gin.Context) {
	payload := util.GetBody[domain.CreateBankInformationDTO](ctx)
	status, user, err := bh.bankusecase.CreateBankInformation(payload)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, gin.H{"data": user})
}

func (bh *bankinfoHandler) FetchBankInfo(ctx *gin.Context) {
	authPayload := security.GetAuthsPayload(ctx)
	status, user, err := bh.bankusecase.FetchBankInformation(authPayload.UserId)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, gin.H{"data": user})
}

func NewBankInfoHandlers(authusecase bankinfo_usecase.BankInfoUsecase) BankInfoHandler {
	return &bankinfoHandler{bankusecase: authusecase}
}
