package invoice_https

import (
	invoice_usecase "go-invoice/internal/invoice/usecase"
	"go-invoice/worker"

	"github.com/gin-gonic/gin"
)

func NewInvoiceRoutes(router *gin.RouterGroup, usecase invoice_usecase.InvoiceUsecase, worker worker.TaskDistributor) {
	invoiceHandler := NewInvoiceHandlers(usecase, worker)
	route := router.Group("/invoice")
	route.POST("/create", invoiceHandler.CreateInvoice)
	route.GET("/all", invoiceHandler.FetchInvoicesWithItems)
	route.GET("/invoices", invoiceHandler.FetchInvoices)
	route.GET("/invoices/:id", invoiceHandler.FetchInvoiceWithItems)
	route.GET("/invoices/stats", invoiceHandler.FetchInvoiceStats)
	route.GET("/download", invoiceHandler.DownloadInvoicePdf)
	route.PATCH("/update/:id", invoiceHandler.UpdateInvoiceStatus)
	route.GET("/:id/download/:customer_id", invoiceHandler.UpdateInvoiceStatus)
}
