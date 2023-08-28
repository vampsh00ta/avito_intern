package ttl

import (
	"avito/config"
	rep "avito/internal/db"
	"avito/internal/redis"
	"context"
	"go.uber.org/zap"
	"time"
)

type TTLMonitor struct {
	rep    rep.Repository
	logger *zap.SugaredLogger
	redis  redis.Repository
	cfg    *config.Config
}
type TTL interface {
	Collect(args *[]string, userId int, slug string, time time.Time)
	SetTTL(ctx context.Context, slugs ...string) error
	DelUsersSegments(ctx context.Context, slugs ...string) error
}

func NewTTL(rep rep.Repository, logger *zap.SugaredLogger, redis redis.Repository, cfg *config.Config) *TTLMonitor {
	return &TTLMonitor{
		rep:    rep,
		logger: logger,
		redis:  redis,
		cfg:    cfg,
	}
}
