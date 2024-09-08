package invoice_https

import (
	"go-invoice/domain"
	invoice_usecase "go-invoice/internal/invoice/usecase"
	"go-invoice/security"
	"go-invoice/util"
	"go-invoice/worker"
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

// FetchInvoicesWithItems(userId int, limit, offset int64) (int, []domain.InvoiceResponse, error)
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

// FetchInvoices(userId int, page, limit int64) (int, []domain.InvoiceResponse, error)
// FetchInvoiceWithItems(invoiceId int) (int, domain.InvoiceResponse, error)
// DeleteInvoiceItems(itemIds []int) (int, error)
// FetchInvoiceStats(userId int) (int, map[string]domain.InvoiceStats, error)

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

// FetchInvoiceWithItems(invoiceId int) (int, domain.InvoiceResponse, error)
// DeleteInvoiceItems(itemIds []int) (int, error)
// FetchInvoiceStats(userId int) (int, map[string]domain.InvoiceStats, error)

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

// DeleteInvoiceItems(itemIds []int) (int, error)
// FetchInvoiceStats(userId int) (int, map[string]domain.InvoiceStats, error)

func (ch *invoiceHandler) FetchInvoiceStats(ctx *gin.Context) {
	authPayload := security.GetAuthsPayload(ctx)
	status, resp, err := ch.usecase.FetchInvoiceStats(authPayload.UserId)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, gin.H{"data": resp})
}

func NewInvoiceHandlers(usecase invoice_usecase.InvoiceUsecase) InvoiceHandler {
	return &invoiceHandler{usecase: usecase}
}
