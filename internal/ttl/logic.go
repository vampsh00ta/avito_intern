package ttl

import (
	rep "avito/internal/db"
	"context"
	"strconv"
	"strings"
	"time"
)

//func (t *TTLMonitor) Collect(args *[]string, userId int, slug string, tm time.Time) {
//	key := strconv.Itoa(userId) + ":" + slug
//	value := tm.Format(time.RFC3339)
//	*args = append(*args, key)
//	*args = append(*args, value)
//}
//
//func (t *TTLMonitor) SetTTL(ctx context.Context, slugs ...string) error {
//
//	if slugs == nil {
//		return nil
//	}
//
//	if err := t.redis.SetTTL(ctx, slugs...); err != nil {
//		return err
//	}
//	return nil
//}
//func (t *TTLMonitor) DelUsersSegments(ctx context.Context, slugs ...string) error {
//	if err := t.redis.DelUsersSegments(ctx, slugs...); err != nil {
//		return err
//	}
//	return nil
//}
//func (t *TTLMonitor) Start(ctx context.Context, exit chan struct{}) {
//	for {
//		select {
//		case <-exit:
//			os.Exit(0)
//		default:
//
//			ttls, err := t.redis.GetTTls(ctx)
//			if err != nil {
//				t.logger.Errorw("redis error", "error", err)
//				continue
//			}
//			currTime := time.Now()
//			for keyUserIdSlug, timeValue := range ttls {
//				splitedKey := strings.Split(keyUserIdSlug, ":")
//				userId, err := strconv.Atoi(splitedKey[0])
//				if err != nil {
//					t.logger.Errorw("userId convert error", "error", err)
//					continue
//				}
//				slug := splitedKey[1]
//				ttlTime, err := time.Parse(time.RFC3339, timeValue)
//
//				if err != nil {
//					t.logger.Errorw("redis error", "error", err)
//					continue
//				}
//				if currTime.Compare(ttlTime) > 0 {
//					if err := t.rep.DeleteSegmentsFromUser(ctx, userId, rep.Segment{Slug: slug}); err != nil {
//						t.logger.Errorw("db err", "error", err)
//						continue
//					}
//					if err := t.redis.DelUsersSegments(ctx, keyUserIdSlug); err != nil {
//						t.logger.Errorw("redis error", "error", err)
//						continue
//					}
//					t.logger.Infow("deleted expired data", "user", userId, "slug", slug)
//				}
//
//			}
//
//		}
//		time.Sleep(t.cfg.TimeUpdate)
//	}
//}

//	func (t *TTLCache) Collect(args *[]Item, userId int, slug string, tm time.Time) {
//		key := strconv.Itoa(userId) + ":" + slug
//		*args = append(*args, Item{key: key, value: tm})
//	}
func (t *TTLCache) Set(userId int, segment rep.Segment, expireTime time.Time) {
	key := strconv.Itoa(userId) + ":" + segment.Slug
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
					if err := t.rep.DeleteSegmentsFromUser(context.Background(), userId, rep.Segment{Slug: slug}); err != nil {
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
