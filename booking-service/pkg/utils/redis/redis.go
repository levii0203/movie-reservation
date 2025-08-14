package rdb

import (
	"fmt"
	"context"
	"time"

	"github.com/redis/go-redis/v9"

)

var (
	ErrUserNotAuthorized = fmt.Errorf("user not authorized")
)

const (
	REDIS_PORT = "5000"
)

var (
	Client *redis.Client = NewRedisClient(fmt.Sprintf("localhost:%s", REDIS_PORT))
)


func NewRedisClient(addr string) *redis.Client{
	cli:= redis.NewClient(&redis.Options{
		Addr: addr,
		DB: 0,
	})

	return cli
}

func Ping() {
	ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		fmt.Println("[ERROR][booking-service/pkg/utils/redis] Ping: %w", err)
		panic(err)
	}
}