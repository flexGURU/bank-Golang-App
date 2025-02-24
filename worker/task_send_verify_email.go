package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/hibiken/asynq"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

func (distributor *RedisTaskDsitrbuter) DistributeTaskSendVerifyEmail(
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info(taskInfo.ID)

	return nil
}