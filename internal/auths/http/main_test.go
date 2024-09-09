package auth_http_test

import (
	"go-invoice/app"
	"go-invoice/domain"
	auth_http "go-invoice/internal/auths/http"
	auth_usecase "go-invoice/internal/auths/usecase"
	"go-invoice/middleware"
	"go-invoice/util"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store auth_usecase.AuthUsecase) *app.Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	server, err := app.NewServer(config, db, nil)
	require.NoError(t, err)
	err = SetupRouter(server, store)
	if err != nil {
		return &app.Server{}
	}
	return server
}

// func TestMain(m *testing.M) {
// 	gin.SetMode(gin.TestMode)
// 	os.Exit(m.Run())
// }

// func newTestServers(t *testing.T, store *gorm.DB, usecase user_services.AccountUsecase) *serverpkg.Server {
// 	config := util.Config{
// 		TokenSymmetricKey:   util.RandomString(32),
// 		AccessTokenDuration: time.Minute,
// 	}
// 	server, err := serverpkg.NewServer(config, store)
// 	err = SetupRouter(server, usecase)
// 	if err != nil {
// 		return &serverpkg.Server{}
// 	}
// 	require.NoError(t, err)
// 	return server
// }

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func SetupRouter(server *app.Server, usecase auth_usecase.AuthUsecase) error {
	router := gin.Default()
	server.Router = router
	groupRouter := router.Group("/api/v1")
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "app is unning fine at" + server.Config.HTTPServerAddress})
	})
	NewUserHandlers(groupRouter, usecase)
	return nil
}

func NewUserHandlers(router *gin.RouterGroup, store auth_usecase.AuthUsecase) {
	authsHandler := auth_http.NewAuthsHandlers(store)
	route := router.Group("/auths")
	route.POST("/register", middleware.ValidatorMiddleware[domain.CreateUserRequestDto], authsHandler.RegisterUser)
	route.POST("/login", middleware.ValidatorMiddleware[domain.LoginRequestDto], authsHandler.LoginUser)
}

// func addAuthorization(
// 	t *testing.T,
// 	request *http.Request,
// 	tokenMaker token.Maker,
// 	authorizationType string,
// 	userId uuid.UUID,
// 	duration time.Duration,
// ) {
// 	token, payload, err := tokenMaker.CreateToken(userId, duration)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, payload)
// 	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
// 	request.Header.Set(authorizationHeaderKey, authorizationHeader)
// }
