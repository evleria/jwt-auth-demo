package jwt

import (
	"github.com/evleria/jwt-auth-demo/pkg/common/config"
	"time"
)

type MakerConfig struct {
	AccessTokenSecret    string
	AccessTokenDuration  time.Duration
	RefreshTokenSecret   string
	RefreshTokenDuration time.Duration
}

func FromConfig() MakerConfig {
	return MakerConfig{
		AccessTokenSecret:    config.GetString("ACCESS_TOKEN_SECRET", "access_secret"),
		AccessTokenDuration:  config.GetDuration("ACCESS_TOKEN_DURATION", time.Minute*5),
		RefreshTokenSecret:   config.GetString("REFRESH_TOKEN_SECRET", "refresh_secret"),
		RefreshTokenDuration: config.GetDuration("REFRESH_TOKEN_DURATION", time.Hour*24*7),
	}
}
