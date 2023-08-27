package ttl

import (
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
}
type TTL interface {
	Collect(args *[]string, userId int, slug string, time time.Time)
	SetTTL(ctx context.Context, slugs ...string) error
	DelUsersSegments(ctx context.Context, slugs ...string) error
}

func NewTTL(rep rep.Repository, logger *zap.SugaredLogger, redis redis.Repository) *TTLMonitor {
	return &TTLMonitor{
		rep:    rep,
		logger: logger,
		redis:  redis,
	}
}
