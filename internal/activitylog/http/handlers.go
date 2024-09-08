package activity_http

import (
	"go-invoice/domain"
	activity_usecase "go-invoice/internal/activitylog/usecase"
	"go-invoice/security"
	"go-invoice/util"

	"github.com/gin-gonic/gin"
)

type ActivityHandler interface {
	FetchActivitieslog(ctx *gin.Context)
	FetchActivityByID(ctx *gin.Context)
}

type activityHandler struct {
	activityusecase activity_usecase.ActivityUsecase
}

func (ah *activityHandler) FetchActivitieslog(ctx *gin.Context) {
	authPayload := security.GetAuthsPayload(ctx)
	status, user, err := ah.activityusecase.FetchActivitieslog(authPayload.UserId)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, gin.H{"data": user})
}

func (ah *activityHandler) FetchActivityByID(ctx *gin.Context) {
	authPayload := security.GetAuthsPayload(ctx)
	payload := util.GetUrlParams[domain.IDParamPayload](ctx)
	status, user, err := ah.activityusecase.FetchActivityByID(payload.ID, authPayload.UserId)
	if err != nil {
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(status, gin.H{"data": user})
}

func NewActivityHandlers(activityusecase activity_usecase.ActivityUsecase) ActivityHandler {
	return &activityHandler{activityusecase: activityusecase}
}
