package queue

import (
	"context"
	"errors"
	"log"

	"github.com/eng-yousef-khaled/expenses_api/internal/domain/auth"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type RedisAdapter struct {
	pool     *redis.Pool
	enqueuer *work.Enqueuer
}

func NewRedisAdapter(cfg RedisConfig, namespace string) *RedisAdapter {
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
	return &RedisAdapter{
		pool:     pool,
		enqueuer: work.NewEnqueuer(namespace, pool),
	}

}

func (r *RedisAdapter) Publish(ctx context.Context, job auth.Job) error {
	_, err := r.enqueuer.Enqueue(job.Name, job.Payload)

	log.Printf("Value come into Enqueue, name: %s \n", job.Name)
	return err
}

func (r *RedisAdapter) GetPool() redis.Conn {
	return r.pool.Get()
}

func (r *RedisAdapter) Pool() *redis.Pool {
	return r.pool
}
