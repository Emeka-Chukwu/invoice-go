package invoice_https_test

import (
	"fmt"
	"go-invoice/app"
	"go-invoice/domain"
	invoice_https "go-invoice/internal/invoice/https"
	invoice_usecase "go-invoice/internal/invoice/usecase"
	"go-invoice/middleware"
	"go-invoice/security"
	"go-invoice/util"
	"go-invoice/worker"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func newTestServer(t *testing.T, store invoice_usecase.InvoiceUsecase) *app.Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
		RedisAddress:        "0.0.0.0:6379",
	}
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	server, err := app.NewServer(config, db, taskDistributor)
	require.NoError(t, err)
	err = SetupRouter(server, store, taskDistributor)
	if err != nil {
		return &app.Server{}
	}
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func SetupRouter(server *app.Server, usecase invoice_usecase.InvoiceUsecase, work worker.TaskDistributor) error {
	router := gin.Default()
	server.Router = router
	groupRouter := router.Group("/api/v1")
	groupRouter.Use(middleware.AuthMiddleware(server.TokenMaker))
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "app is unning fine at" + server.Config.HTTPServerAddress})
	})
	NewInvoiceHandlers(groupRouter, usecase, work)
	return nil
}

func NewInvoiceHandlers(router *gin.RouterGroup, store invoice_usecase.InvoiceUsecase, work worker.TaskDistributor) {
	handler := invoice_https.NewInvoiceHandlers(store, work)
	route := router.Group("/invoice")

	route.POST("/create", middleware.ValidatorMiddleware[domain.CreateInvoiceRequestDTO], handler.CreateInvoice)
	route.GET("/invoices", handler.FetchInvoices)
	route.PATCH("/update/:id", handler.UpdateInvoiceStatus)
	route.GET("/invoices/stats", handler.FetchInvoiceStats)
	// route.POST("/login", middleware.ValidatorMiddleware[domain.LoginRequestDto], authsHandler.LoginUser)
}

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker security.Maker,
	authorizationType string,
	userId int,
	duration time.Duration,
) {
	token, payload, err := tokenMaker.CreateToken(userId, duration, true)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}
