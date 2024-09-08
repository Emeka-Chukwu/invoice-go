package worker

import (
	"context"
	"go-invoice/domain"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeActivityLog(ctx context.Context, payload *domain.CreateActivityDTO, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
