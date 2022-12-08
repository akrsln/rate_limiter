package main

import (
	"context"
	"github.com/go-redis/redis/v9"
	"time"
)

var ctx = context.Background()

type Bucket struct {
	limit       int64         // maximum number of requests allowed in the given duration
	duration    time.Duration // time window for the rate limiter
	redisClient *redis.Client // redis client instance
}

func NewBucket(limit int64, duration time.Duration, redisClient *redis.Client) *Bucket {
	return &Bucket{
		limit:       limit,
		duration:    duration,
		redisClient: redisClient,
	}
}

// Allow method checks whether the given flag is permitted to make a request
// flag can be username, IP etc.

func (b *Bucket) Allow(flag string) bool {

	key := flag + "_request_count"

	reqCount, err := b.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return false
	}

	if reqCount > b.limit {
		return false
	}

	b.redisClient.Expire(ctx, key, b.duration)

	return true
}

func (b *Bucket) Remaining(flag string) int64 {

	key := flag + "_request_count"
	reqCount, err := b.redisClient.Get(ctx, key).Int64()
	if err != nil {
		return 0
	}

	return b.limit - reqCount
}

func (b *Bucket) NextAllowed(flag string) time.Duration {

	key := flag + "_request_count"
	_, err := b.redisClient.Get(ctx, key).Int()
	if err != nil {
		return b.duration
	}

	remaining, _ := b.redisClient.TTL(ctx, key).Result()
	if err != nil {
		return b.duration
	}

	return remaining
}
