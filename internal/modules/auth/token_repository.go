package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type TokenRepository interface {
	Blacklist(userId int, t time.Time, ttl time.Duration) error
	IsBlacklisted(userId int) (time.Time, bool, error)
}

type tokenRepository struct {
	redis *redis.Client
}

func NewTokenRepository(redis *redis.Client) TokenRepository {
	return &tokenRepository{
		redis: redis,
	}
}

func (r *tokenRepository) Blacklist(userId int, t time.Time, ttl time.Duration) error {
	key := getBlacklistKey(userId)
	return r.redis.Set(context.TODO(), key, t, ttl).Err()
}

func (r *tokenRepository) IsBlacklisted(userId int) (time.Time, bool, error) {
	key := getBlacklistKey(userId)
	s, err := r.redis.Get(context.TODO(), key).Result()
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

func getBlacklistKey(userId int) string {
	return fmt.Sprintf("blacklist::%d", userId)
}
