package app

import (
	"database/sql"
	"fmt"
	activity_http "go-invoice/internal/activitylog/http"
	activitylog_repository "go-invoice/internal/activitylog/repository"
	activity_usecase "go-invoice/internal/activitylog/usecase"
	auth_http "go-invoice/internal/auths/http"
	auth_repository "go-invoice/internal/auths/repository"
	auth_usecase "go-invoice/internal/auths/usecase"
	bankinfo_http "go-invoice/internal/bank_info/http"
	bankinfo_repository "go-invoice/internal/bank_info/repository"
	bankinfo_usecase "go-invoice/internal/bank_info/usecase"
	customer_http "go-invoice/internal/customer/http"
	customer_repository "go-invoice/internal/customer/repository"
	customer_usecase "go-invoice/internal/customer/usecase"
	invoice_https "go-invoice/internal/invoice/https"
	invoice_repository "go-invoice/internal/invoice/repository"
	invoice_usecase "go-invoice/internal/invoice/usecase"
	"go-invoice/middleware"
	"go-invoice/security"
	"go-invoice/util"
	"go-invoice/worker"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Config          util.Config
	conn            *sql.DB
	Router          *gin.Engine
	TokenMaker      security.Maker
	taskDistributor worker.TaskDistributor
}

//// server serves out http request for our backend service

func NewServer(config util.Config, conn *sql.DB, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := security.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create initiating security: %w", err)
	}
	server := &Server{TokenMaker: tokenMaker, Config: config, conn: conn, taskDistributor: taskDistributor}
	server.setupRouter()
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Personal Budget app ruuning at %s", server.Config.HTTPServerAddress),
		})
	})
	groupRouter := router.Group("/api/v1")

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":       "Resource not found",
			"route":       c.Request.URL.Path,
			"status_code": 404,
		})
	})

	//////auths
	authRepo := auth_repository.NewAuthRepository(server.conn)
	authusecase := auth_usecase.NewAuthUsecase(authRepo, server.TokenMaker, server.Config)
	auth_http.NewAuthRoutes(groupRouter, authusecase)

	//////middleware
	groupRouter.Use(middleware.AuthMiddleware(server.TokenMaker))

	//////activity
	activityRepo := activitylog_repository.NewAuthRepository(server.conn)
	activityusecase := activity_usecase.NewActivityUsecase(activityRepo)
	activity_http.NewActivityRoutes(groupRouter, activityusecase)

	//////bank
	bankRepo := bankinfo_repository.NewBankInfoRepository(server.conn)
	bankusecase := bankinfo_usecase.NewBankInfoUsecase(bankRepo)
	bankinfo_http.NewBankInfoRoutes(groupRouter, bankusecase)

	//////customer
	customerRepo := customer_repository.NewCustomerRepository(server.conn)
	customerusecase := customer_usecase.NewCustomerUsecase(customerRepo)
	customer_http.NewCustomerRoutes(groupRouter, customerusecase)

	//////invoices
	invoiceRepo := invoice_repository.NewInvoiceWithItems(server.conn)
	invoiceusecase := invoice_usecase.NewInvoiceUsecase(invoiceRepo)
	invoice_https.NewInvoiceRoutes(groupRouter, invoiceusecase, server.taskDistributor)

	server.Router = router
}
