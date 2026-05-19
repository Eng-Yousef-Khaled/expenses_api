package queue

import (
	"context"

	"github.com/eng-yousef-khaled/expenses_api/internal/domain/auth"
)

type TaskQueue interface {
	Publish(ctx context.Context, job auth.Job) error
}

type RedisConfig struct {
	MaxActive int16
	MaxIdle   int16
	Wait      bool
	Dial      any
}
