package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

func (distributor *RedisTaskDistributer) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	options ...asynq.Option,

) error {
	json_payload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshall payload %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, json_payload, options...)
	taskInfo, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueu task %w", err)
		
	}

	slog.Info("task enqueued",
			slog.String("task_type", task.Type()),               // Task type
			slog.String("task_payload", string(task.Payload())), // Task payload
			slog.String("task_queue", taskInfo.Queue),           // Task queue
			slog.Int("task_max_retry", taskInfo.MaxRetry),       // Task max retry
			)


			return nil

}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {

	var payload PayloadSendVerifyEmail

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("error unmarshalling payload %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user Not Found %w", asynq.SkipRetry)
		}

		return fmt.Errorf("Internal server error %w", err)
		
	}

	slog.Info(
		"processed task",
		slog.String("task type", task.Type()),
		slog.String("email", user.Email),
	)


	return nil


}
