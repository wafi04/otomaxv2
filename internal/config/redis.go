package config

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)
type RedisConnection struct {
	Config *RedisConfig
	Client *redis.Client
}



func NewRedisConnection(cfg *RedisConfig) *RedisConnection{
	return &RedisConnection{
		Config: cfg,
	}
}



func (c *RedisConnection) Connect()*redis.Client{
	client := redis.NewClient(&redis.Options{
		Addr:	  fmt.Sprintf("%s:%s",c.Config.Host,c.Config.Port),
        Password: "", 
        DB:		  0,  
        Protocol: 2,
		MaxRetries: c.Config.MaxRetries,
		PoolTimeout: c.Config.PoolTimeout,
		PoolSize:  c.Config.PoolSize,
		DialTimeout: c.Config.IdleTimeout,  
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Printf("❌ Redis connection failed: %v\n", err)
		return nil
	}

	fmt.Println("✅ Redis connected successfully")
	return client
}

func (c *RedisConnection) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := c.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

func (c *RedisConnection) Get(ctx context.Context, key string) (string, error) {
	val, err := c.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s does not exist", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

func (c *RedisConnection) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	return c.Client.Set(ctx, key, data, expiration).Err()
}

func (c *RedisConnection) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := c.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key %s does not exist", key)
	} else if err != nil {
		return fmt.Errorf("failed to get key %s: %w", key, err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}
	return nil
}