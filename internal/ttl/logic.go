package ttl

import (
	rep "avito/internal/db"
	"context"
	"os"
	"strconv"
	"strings"
	"time"
)

func (t *TTLMonitor) Collect(args *[]string, userId int, slug string, tm time.Time) {
	key := strconv.Itoa(userId) + ":" + slug
	value := tm.Format(time.RFC3339)
	*args = append(*args, key)
	*args = append(*args, value)
}

func (t *TTLMonitor) SetTTL(ctx context.Context, slugs ...string) error {

	if slugs == nil {
		return nil
	}

	if err := t.redis.SetTTL(ctx, slugs...); err != nil {
		return err
	}
	return nil
}
func (t *TTLMonitor) DelUsersSegments(ctx context.Context, slugs ...string) error {
	if err := t.redis.DelUsersSegments(ctx, slugs...); err != nil {
		return err
	}
	return nil
}
func (t *TTLMonitor) Start(ctx context.Context, exit chan struct{}) {
	for {
		select {
		case <-exit:
			os.Exit(0)
		default:

			ttls, err := t.redis.GetTTls(ctx)
			if err != nil {
				t.logger.Errorw("redis error", "error", err)
				continue
			}
			currTime := time.Now()
			for keyUserIdSlug, timeValue := range ttls {
				splitedKey := strings.Split(keyUserIdSlug, ":")
				userId, err := strconv.Atoi(splitedKey[0])
				if err != nil {
					t.logger.Errorw("userId convert error", "error", err)
					continue
				}
				slug := splitedKey[1]
				ttlTime, err := time.Parse(time.RFC3339, timeValue)

				if err != nil {
					t.logger.Errorw("redis error", "error", err)
					continue
				}
				if currTime.Compare(ttlTime) > 0 {
					if err := t.rep.DeleteSegmentsFromUser(ctx, userId, rep.Segment{Slug: slug}); err != nil {
						t.logger.Errorw("db err", "error", err)
						continue
					}
					if err := t.redis.DelUsersSegments(ctx, keyUserIdSlug); err != nil {
						t.logger.Errorw("redis error", "error", err)
						continue
					}
					t.logger.Infow("deleted expired data", "user", userId, "slug", slug)
				}

			}

		}
		time.Sleep(t.cfg.TimeUpdate)
	}
}
