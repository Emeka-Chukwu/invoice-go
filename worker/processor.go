package worker

import (
	"context"
	activitylog_repository "go-invoice/internal/activitylog/repository"
	payment_repository "go-invoice/internal/payments/repository"
	"go-invoice/util"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessActivityLog(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server          *asynq.Server
	paymentStore    payment_repository.PaymentRepo
	activitiesStore activitylog_repository.ActivityRepository
	config          util.Config
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, storepayment payment_repository.PaymentRepo, activitiesStore activitylog_repository.ActivityRepository, config util.Config) TaskProcessor {

	server := asynq.NewServer(
		redisOpt, asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: NewLogger(),
		},
	)
	return &RedisTaskProcessor{
		server:          server,
		paymentStore:    storepayment,
		activitiesStore: activitiesStore,
		config:          config,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskCreateActivityLog, processor.ProcessActivityLog)
	// mux.HandleFunc(TaskResendVerifyEmail, processor.ProcessTaskResendVerifyEmail)
	// mux.HandleFunc(TaskLinkEventTicket, processor.ProcessChainingOfEventTicket)
	// mux.HandleFunc(TaskVerifyPaymentTicket, processor.ProcessVerifyPaymentEmail)
	return processor.server.Start(mux)

}
