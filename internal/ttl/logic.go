package ttl

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func (t *TTLMonitor) Collect(args *[]string, userId int, slug string, tm time.Time) {
	key := strconv.Itoa(userId) + ":" + slug
	value := tm.Format(time.RFC3339)
	fmt.Println(value)
	*args = append(*args, key)
	*args = append(*args, value)
}

func (t *TTLMonitor) Add(ctx context.Context, slugs ...any) error {

	if err := t.redis.SetTTL(ctx, slugs...); err != nil {
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
			for key, value := range ttls {
				splitedKey := strings.Split(key, ":")
				userId, err := strconv.Atoi(splitedKey[0])
				if err != nil {
					t.logger.Errorw("userId convert error", "error", err)
					continue
				}
				slug := splitedKey[1]
				ttlTime, err := time.Parse(time.RFC3339, value)

				if err != nil {
					t.logger.Errorw("redis error", "error", err)
					continue
				}
				if currTime.Compare(ttlTime) > 0 {
					if err := t.rep.DeleteSegmentsFromUser(ctx, userId, slug); err != nil {
						t.logger.Errorw("db err", "error", err)
						continue
					}
					if err := t.redis.DelUsersSegments(ctx, key); err != nil {
						t.logger.Errorw("redis error", "error", err)
						continue
					}
					t.logger.Infow("deleted expired data", "user", userId, "slug", slug)
				}

			}

		}
		time.Sleep(time.Second * 5)
	}
}
