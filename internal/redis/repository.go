package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Repository interface {
	SetTTL(ctx context.Context, fields ...string) error
	GetTTls(ctx context.Context) (map[string]string, error)
	DelUsersSegments(tx context.Context, fields ...string) error
}

type Redis struct {
	client *redis.Client
	logger *zap.SugaredLogger
}

func New(client *redis.Client, logger *zap.SugaredLogger) Repository {
	return &Redis{
		client: client,
		logger: logger,
	}
}
