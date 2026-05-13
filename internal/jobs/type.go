package jobs

type RedisConfig struct {
	MaxActive int16
	MaxIdle   int16
	Wait      bool
	Dial      any
}
