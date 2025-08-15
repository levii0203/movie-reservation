package rdb

import (
	"fmt"
	"context"
	"time"
	"strconv"
	"github.com/levii0203/user-service/pkg/utils/token"

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
		fmt.Println("[ERROR][user-service/pkg/utils/redis] Ping: %w", err)
		panic(err)
	}
}

func GetTokenAccessCount(token_str string) (int,error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	val, err:= Client.Get(ctx,token_str).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, ErrUserNotAuthorized
		}
		return 0,token.ErrInvalidToken
	}
	
	v , err := strconv.Atoi(val);
	if err !=nil {
		return 0, fmt.Errorf("runtime error")
	}
	return v, nil
}

func IncrementTokenAccessCount(token_str string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second);
	defer cancel()

	if err := Client.Incr(ctx, token_str).Err(); err != nil {
		if err == redis.Nil {
			return ErrUserNotAuthorized
		}
		return fmt.Errorf("something went wrong when authenticating user")
	}

	return nil
}

