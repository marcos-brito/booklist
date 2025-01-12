package conn

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func InitRedis(client *redis.Client) {
	Redis = client
}

func NewRedisClient() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
}
