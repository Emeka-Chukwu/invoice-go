package invoice_https

import (
	"fmt"
	"go-invoice/domain"
	invoice_usecase "go-invoice/internal/invoice/usecase"
	"go-invoice/security"
	"go-invoice/util"
	"go-invoice/worker"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

type InvoiceHandler interface {
	CreateInvoice(ctx *gin.Context)
	FetchInvoicesWithItems(ctx *gin.Context)
	FetchInvoices(ctx *gin.Context)
	FetchInvoiceWithItems(ctx *gin.Context)
	FetchInvoiceStats(ctx *gin.Context)
	DownloadInvoicePdf(ctx *gin.Context)
	UpdateInvoiceStatus(ctx *gin.Context)
}

type invoiceHandler struct {
	usecase invoice_usecase.InvoiceUsecase
	Worker  worker.TaskDistributor
}

// CreateInvoiceWithItems(req domain.CreateInvoiceRequestDTO) (int, int, error)

func (ch *invoiceHandler) CreateInvoice(ctx *gin.Context) {
	authPayload := security.GetAuthsPayload(ctx)
	payload := util.GetBody[domain.CreateInvoiceRequestDTO](ctx)
	payload.UserID = authPayload.UserId
	status, resp, err := ch.usecase.CreateInvoiceWithItems(payload)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	taskPayload := &domain.CreateActivityDTO{
		UserID:     payload.UserID,
		Action:     "You have created an invoice role for  ",
		EntityType: "Invoice",
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	ch.Worker.DistributeActivityLog(ctx, taskPayload, opts...)
	ctx.JSON(status, gin.H{"data": resp})
}

func (ch *invoiceHandler) FetchInvoicesWithItems(ctx *gin.Context) {
	authPayload := security.GetAuthsPayload(ctx)
	payload := util.GetUrlQueryParams[domain.PaginationDTO](ctx)
	pagination := ch.usecase.GetPagination(payload)
	status, resp, err := ch.usecase.FetchInvoicesWithItems(authPayload.UserId, int64(pagination.Limit), int64(pagination.Page))
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, gin.H{"data": resp})
}

func (ch *invoiceHandler) FetchInvoices(ctx *gin.Context) {
	authPayload := security.GetAuthsPayload(ctx)
	payload := util.GetUrlQueryParams[domain.PaginationDTO](ctx)
	pagination := ch.usecase.GetPagination(payload)
	status, resp, err := ch.usecase.FetchInvoices(authPayload.UserId, int64(pagination.Limit), int64(pagination.Page))
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, gin.H{"data": resp})
}

func (ch *invoiceHandler) FetchInvoiceWithItems(ctx *gin.Context) {
	authPayload := security.GetAuthsPayload(ctx)
	payload := util.GetUrlParams[domain.IDParamPayload](ctx)
	status, resp, err := ch.usecase.FetchInvoiceWithItems(payload.ID, authPayload.UserId)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, gin.H{"data": resp})
}

func (ch *invoiceHandler) FetchInvoiceStats(ctx *gin.Context) {
	authPayload := security.GetAuthsPayload(ctx)
	status, resp, err := ch.usecase.FetchInvoiceStats(authPayload.UserId)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, gin.H{"data": resp})
}

func NewInvoiceHandlers(usecase invoice_usecase.InvoiceUsecase, Worker worker.TaskDistributor) InvoiceHandler {
	return &invoiceHandler{usecase: usecase, Worker: Worker}
}

func (ch *invoiceHandler) DownloadInvoicePdf(ctx *gin.Context) {
	resp, err := ch.usecase.GenerateInvoicePDF()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fileName := "account statement " + fmt.Sprintf("%d", time.Now().UnixNano()) + ".pdf"
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Data(http.StatusOK, "application/pdf", resp)
}

func (ch *invoiceHandler) UpdateInvoiceStatus(ctx *gin.Context) {
	authPayload := security.GetAuthsPayload(ctx)
	payload := util.GetUrlParams[domain.IDParamPayload](ctx)
	status, err := ch.usecase.UpdateInvoiceStatus("paid", payload.ID, authPayload.UserId)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	taskPayload := &domain.CreateActivityDTO{
		UserID:     authPayload.UserId,
		Action:     "You have manually confirmed the payment",
		EntityType: "Invoice",
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	ch.Worker.DistributeActivityLog(ctx, taskPayload, opts...)
	ctx.JSON(status, gin.H{"data": "Manual Payment confirmed"})
}
