package queue

import (
	"errors"
	"log"

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

func (r *RedisAdapter) Enqueue(taskName string, payloads map[string]any) error {
	_, err := r.enqueuer.Enqueue(taskName, payloads)

	log.Printf("Value come into Enqueue, name: %s \n", taskName)
	return err
}

func (r *RedisAdapter) GetPool() redis.Conn {
	return r.pool.Get()
}

func (r *RedisAdapter) Pool() *redis.Pool {
	return r.pool
}
