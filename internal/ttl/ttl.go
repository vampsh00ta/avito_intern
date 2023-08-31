package ttl

import (
	"avito/config"
	rep "avito/internal/db"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Cache struct {
	sync.Mutex
	storage map[string]time.Time
}
type TTLCache struct {
	rep    rep.Repository
	logger *zap.SugaredLogger
	cfg    *config.Config
	cache  *Cache
}
type TTL interface {
	Set(userId int, slug string, expireTime time.Time)
	GetAll() map[string]time.Time
	Delete(keys ...string)
	DeleteAll()
	Start(exit chan struct{})
}

func New(rep rep.Repository, logger *zap.SugaredLogger, cfg *config.Config) TTL {
	storage := make(map[string]time.Time)
	return &TTLCache{
		rep:    rep,
		logger: logger,
		cfg:    cfg,
		cache:  &Cache{storage: storage},
	}
}
