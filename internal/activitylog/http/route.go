package activity_http

import (
	activity_usecase "go-invoice/internal/activitylog/usecase"

	"github.com/gin-gonic/gin"
)

func NewActivityRoutes(router *gin.RouterGroup, usecase activity_usecase.ActivityUsecase) {
	activityHandler := NewActivityHandlers(usecase)
	route := router.Group("/activities")
	route.POST("/logs", activityHandler.FetchActivitieslog)
	route.POST("/logs/:id", activityHandler.FetchActivityByID)
}
