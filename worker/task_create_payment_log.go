package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"go-invoice/domain"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskCreatePaymentLog = "task:task_create_payment_log"

func (distributor *RedisTaskDistributor) DistributeCreatePayment(
	ctx context.Context,
	payload *domain.PaymentDto,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskCreatePaymentLog, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueuedn task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessCreatePayment(ctx context.Context, task *asynq.Task) error {
	var payload domain.PaymentDto
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}
	_, err := processor.paymentStore.CreatePayment(payload)
	if err != nil {
		return err
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Msg("Activity log recorded")
	return nil

}
