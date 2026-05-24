package queue

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/eng-yousef-khaled/expenses_api/internal/application/auth"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type RedisAdapter interface {
	Publish(ctx context.Context, job auth.Job) error
}
type redisAdapter struct {
	pool     *redis.Pool
	enqueuer *work.Enqueuer
}

func NewRedisAdapter(cfg RedisConfig, namespace string) RedisAdapter {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			dial, ok := cfg.Dial.(func() (redis.Conn, error))
			if !ok {
				return nil, errors.New("Failed to use config dial")
			}
			return dial()
		},
		MaxIdle:   int(cfg.MaxIdle),
		MaxActive: int(cfg.MaxActive),
		Wait:      cfg.Wait,
	}
	conn := pool.Get()

	defer conn.Close()

	_, err := conn.Do("PING")

	if err != nil {
		slog.Error("Redis is NOT connected", "error", err)
		os.Exit(0)
	} else {
		slog.Info("Redis connection is healthy")
	}
	return &redisAdapter{
		pool:     pool,
		enqueuer: work.NewEnqueuer(namespace, pool),
	}

}

func (r *redisAdapter) Publish(ctx context.Context, job auth.Job) error {
	_, err := r.enqueuer.Enqueue(job.Name, job.Payload)

	slog.Info("Value come into Enqueue, name: %s \n", "to", job.Name)
	return err
}
