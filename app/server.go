package app

import (
	"database/sql"
	"fmt"
	"go-invoice/security"
	"go-invoice/util"
	"go-invoice/worker"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config          util.Config
	conn            *sql.DB
	router          *gin.Engine
	tokenMaker      security.Maker
	taskDistributor worker.TaskDistributor
}

//// server serves out http request for our backend service

func NewServer(config util.Config, conn *sql.DB, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := security.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create initiating security: %w", err)
	}
	server := &Server{tokenMaker: tokenMaker, config: config, conn: conn, taskDistributor: taskDistributor}
	server.setupRouter()
	return server, nil
}

func errorrResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Personal Budget app ruuning at %s", server.config.HTTPServerAddress),
		})
	})
	// groupRouter := router.Group("/api/v1")

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":       "Resource not found",
			"route":       c.Request.URL.Path,
			"status_code": 404,
		})
	})

	// //////user
	// userRepo := repositories_users.NewUserAuths(server.conn)
	// userCase := usecase_user.NewUsecaseUser(server.config, server.tokenMaker, userRepo)
	// users_v1.NewUserRoutes(groupRouter, userCase)

	server.router = router
}
