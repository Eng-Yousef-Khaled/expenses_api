package queue

type TaskQueue interface {
	Enqueue(taskName string, payload map[string]any) error
}

type RedisConfig struct {
	MaxActive int16
	MaxIdle   int16
	Wait      bool
	Dial      any
}
