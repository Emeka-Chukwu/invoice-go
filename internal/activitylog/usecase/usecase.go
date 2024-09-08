package activity_usecase

import (
	"errors"
	"go-invoice/domain"
	activitylog_repository "go-invoice/internal/activitylog/repository"
	"net/http"
)

type ActivityUsecase interface {
	FetchActivityByID(id, userId int) (int, domain.Activity, error)
	FetchActivitieslog(userId int) (int, []domain.Activity, error)
}

type activityUsecase struct {
	Repo activitylog_repository.ActivityRepository
}

func NewActivityUsecase(Repo activitylog_repository.ActivityRepository) ActivityUsecase {
	return activityUsecase{Repo: Repo}
}

// FetchActivitieslog implements ActivityUsecase.
func (a activityUsecase) FetchActivitieslog(userId int) (int, []domain.Activity, error) {
	resp, err := a.Repo.FetchActivitieslog(userId)
	if err != nil {
		return http.StatusInternalServerError, []domain.Activity{}, err
	}
	return http.StatusOK, resp, nil
}

// FetchActivityByID implements ActivityUsecase.
func (a activityUsecase) FetchActivityByID(id, userId int) (int, domain.Activity, error) {
	resp, err := a.Repo.FetchActivityByID(id)
	if err != nil {
		return http.StatusInternalServerError, domain.Activity{}, err
	}
	if resp.UserID != userId {
		return http.StatusForbidden, domain.Activity{}, errors.New("resources forbidden")
	}
	return http.StatusOK, resp, nil
}
