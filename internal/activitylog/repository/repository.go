package activitylog_repository

import (
	"context"
	"database/sql"
	"go-invoice/domain"
	"go-invoice/util"
)

type activityRepository struct {
	Db *sql.DB
}

func (a *activityRepository) FetchActivitieslog(userId int) ([]domain.Activity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), util.DbTimeout)
	defer cancel()
	stmt := `select id, user_id, action, entity_type, created_at, updated_at from activities where user_id=$1`
	rows, err := a.Db.QueryContext(ctx, stmt, userId)
	results := make([]domain.Activity, 0)
	for rows.Next() {
		var model domain.Activity
		rows.Scan(&model.ID, &model.UserID, &model.Action, &model.EntityID, &model.CreatedAt, &model.UpdatedAt)
		results = append(results, model)
	}
	return results, err
}

func (a *activityRepository) CreateActivity(req domain.CreateActivityDTO) (domain.Activity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), util.DbTimeout)
	defer cancel()
	stmt := `insert into activities (user_id, action, entity_type ) values ($1, $2, $3 )
	returning id, user_id, action, entity_type, created_at, updated_at`
	var response domain.Activity
	err := a.Db.QueryRowContext(ctx, stmt, req.UserID, req.Action, req.EntityType).
		Scan(&response.ID, &response.UserID, &response.Action, &response.EntityID, &response.CreatedAt, &response.UpdatedAt)
	return response, err
}

func (a *activityRepository) FetchActivityByID(id int) (domain.Activity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), util.DbTimeout)
	defer cancel()
	stmt := `select id, user_id, action, entity_type, created_at, updated_at from activities where id=$1`
	var model domain.Activity
	err := a.Db.QueryRowContext(ctx, stmt, id).
		Scan(&model.ID, &model.UserID, &model.Action, &model.EntityID, &model.CreatedAt, &model.UpdatedAt)
	return model, err
}

type ActivityRepository interface {
	FetchActivityByID(id int) (domain.Activity, error)
	FetchActivitieslog(userId int) ([]domain.Activity, error)
	CreateActivity(req domain.CreateActivityDTO) (domain.Activity, error)
}

func NewAuthRepository(db *sql.DB) ActivityRepository {
	return &activityRepository{Db: db}
}
