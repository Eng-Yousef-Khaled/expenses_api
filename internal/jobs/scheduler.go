package jobs

import (
	"context"
	"errors"
	"fmt"

	repo "github.com/eng-yousef-khaled/expenses_api/internal/adapters/postgresql/sqlc"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type Service interface {
	NewEnqueuer(ctx context.Context, jobName string, data interface{}) error
	Log(ctx context.Context, job *work.Job, next work.NextMiddlewareFunc) error
}

type svr struct {
	cachingPool *redis.Pool
	repo        repo.Querier
	enqueuer    *work.Enqueuer
}

func (s *svr) Log(ctx context.Context, job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}

func (s *svr) NewEnqueuer(ctx context.Context, jobName string, data interface{}) error {
	_, err := s.enqueuer.Enqueue(jobName, work.Q{"data": data})
	return err
}

func CreateService(config RedisConfig, repo repo.Querier, namespace string) Service {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			dialFunc, ok := config.Dial.(func() (redis.Conn, error))
			if !ok {
				return nil, errors.New("invalid redis dial function type")
			}
			return dialFunc()
		},
		MaxIdle:   int(config.MaxActive),
		Wait:      config.Wait,
		MaxActive: int(config.MaxActive),
	}

	return &svr{
		cachingPool: pool,
		repo:        repo,
		enqueuer:    work.NewEnqueuer(namespace, pool),
	}
}
