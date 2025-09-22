package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// RedisClient 封装 go-redis 客户端
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient 创建 Redis 客户端
func NewRedisClient(addr, password string, dbNum int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbNum,
	})
	return &RedisClient{client: rdb}
}

// Ping 检查 Redis 连接
func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close 关闭 Redis 连接
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// GetClient 获取原生 go-redis 客户端
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// Shutdown 实现 app.Shutdownable 接口
func (r *RedisClient) Shutdown(ctx context.Context) error {
	return r.Close()
}

// Get 获取键值
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set 设置键值
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Del 删除键
func (r *RedisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Del(ctx, keys...).Result()
}

// Exists 检查键是否存在
func (r *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}
