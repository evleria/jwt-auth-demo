// Package repository encapsulates work with databases
package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// Token contains methods to work with token-related data
type Token interface {
	Blacklist(ctx context.Context, userID int, t time.Time, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, userID int) (time.Time, bool, error)
}

type token struct {
	redis *redis.Client
}

// NewTokenRepository creates token repository
func NewTokenRepository(redisClient *redis.Client) Token {
	return &token{
		redis: redisClient,
	}
}

func (r *token) Blacklist(ctx context.Context, userID int, t time.Time, ttl time.Duration) error {
	key := getBlacklistKey(userID)
	return r.redis.Set(ctx, key, t.UnixNano(), ttl).Err()
}

func (r *token) IsBlacklisted(ctx context.Context, userID int) (time.Time, bool, error) {
	key := getBlacklistKey(userID)
	s, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return time.Time{}, false, nil
		}
		return time.Time{}, false, err
	}
	if s == "" {
		return time.Time{}, false, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return time.Time{}, false, err
	}
	t := time.Unix(0, int64(i))
	return t, true, nil
}

func getBlacklistKey(userID int) string {
	return fmt.Sprintf("blacklist::%d", userID)
}
