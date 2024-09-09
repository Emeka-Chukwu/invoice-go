package activitylog_repository

import (
	"fmt"
	"go-invoice/domain"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestFetchActivitieslog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	activityRepo := &activityRepository{Db: db}
	query := `select id, user_id, action, entity_type, created_at, updated_at from activities where user_id=\$1`
	rows := sqlmock.NewRows([]string{"id", "user_id", "action", "entity_type", "created_at", "updated_at"}).
		AddRow(1, 101, "CREATE", "USER", time.Now(), time.Now()).
		AddRow(2, 101, "UPDATE", "POST", time.Now(), time.Now())

	mock.ExpectQuery(query).
		WithArgs(101).
		WillReturnRows(rows)
	activities, err := activityRepo.FetchActivitieslog(101)
	require.NoError(t, err)
	require.Len(t, activities, 2)
	require.Equal(t, 101, activities[0].UserID)
	require.Equal(t, "CREATE", activities[0].Action)
	require.Equal(t, "UPDATE", activities[1].Action)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateActivity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	activityRepo := &activityRepository{Db: db}

	req := domain.CreateActivityDTO{
		UserID:     101,
		Action:     "CREATE",
		EntityType: "POST",
	}

	query := `insert into activities \(user_id, action, entity_type\) values \(\$1, \$2, \$3\)
	returning id, user_id, action, entity_type, created_at, updated_at`
	rows := sqlmock.NewRows([]string{"id", "user_id", "action", "entity_type", "created_at", "updated_at"}).
		AddRow(1, 101, "CREATE", "POST", time.Now(), time.Now())

	mock.ExpectQuery(query).
		WithArgs(req.UserID, req.Action, req.EntityType).
		WillReturnRows(rows)
	activity, err := activityRepo.CreateActivity(req)
	fmt.Println(activity, err)

	require.NoError(t, err)
	require.Equal(t, req.UserID, activity.UserID)
	require.Equal(t, req.Action, activity.Action)
	require.Equal(t, req.EntityType, activity.EntityType)
	require.NotZero(t, activity.ID)
	require.NotZero(t, activity.CreatedAt)
	require.NotZero(t, activity.UpdatedAt)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFetchActivityByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	activityRepo := &activityRepository{Db: db}
	query := `select id, user_id, action, entity_type, created_at, updated_at from activities where id=\$1`
	rows := sqlmock.NewRows([]string{"id", "user_id", "action", "entity_type", "created_at", "updated_at"}).
		AddRow(1, 101, "CREATE", "POST", time.Now(), time.Now())

	mock.ExpectQuery(query).
		WithArgs(1).
		WillReturnRows(rows)

	activity, err := activityRepo.FetchActivityByID(1)

	require.NoError(t, err)
	require.Equal(t, 101, activity.UserID)
	require.Equal(t, "CREATE", activity.Action)
	require.Equal(t, "POST", activity.EntityType)
	require.NotZero(t, activity.ID)
	require.NotZero(t, activity.CreatedAt)
	require.NotZero(t, activity.UpdatedAt)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
