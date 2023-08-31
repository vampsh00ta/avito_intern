package ttl

import (
	rep "avito/internal/db"
	"context"
	"strconv"
	"strings"
	"time"
)

func (t *TTLCache) Set(userId int, slug string, expireTime time.Time) {
	key := strconv.Itoa(userId) + ":" + slug
	t.cache.Lock()
	t.cache.storage[key] = expireTime
	t.cache.Unlock()
}
func (t *TTLCache) Delete(keys ...string) {
	t.cache.Lock()
	for _, key := range keys {
		delete(t.cache.storage, key)
	}
	t.cache.Unlock()
}
func (t *TTLCache) GetAll() map[string]time.Time {
	return t.cache.storage
}
func (t *TTLCache) DeleteAll() {
	t.cache.Lock()
	for key, _ := range t.cache.storage {
		delete(t.cache.storage, key)
	}
	t.cache.Unlock()
}
func (t *TTLCache) Start(exit chan struct{}) {
	for {
		ticker := time.NewTicker(t.cfg.TimeUpdate)

		select {
		case <-exit:
			t.DeleteAll()
		case <-ticker.C:
			ttls := t.GetAll()
			currTime := time.Now()
			for keyUserIdSlug, expireTime := range ttls {
				if currTime.Compare(expireTime) > 0 {
					splitedKey := strings.Split(keyUserIdSlug, ":")
					userId, err := strconv.Atoi(splitedKey[0])
					if err != nil {
						t.logger.Errorw("userId convert error", "error", err)
						continue
					}
					slug := splitedKey[1]
					if err := t.rep.DeleteSegmentsFromUser(context.Background(), userId, &rep.Segment{Slug: slug}); err != nil {
						t.logger.Errorw("db err", "error", err)
						continue
					}
					t.Delete(keyUserIdSlug)
					t.logger.Infow("deleted expired data", "user", userId, "slug", slug)
				}

			}

		}

	}
}
