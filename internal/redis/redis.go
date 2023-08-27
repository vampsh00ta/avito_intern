package redis

import (
	"context"
)

func (r *Redis) SetTTL(ctx context.Context, fields ...string) error {
	args := make([]interface{}, len(fields))
	for i, v := range fields {
		args[i] = v
	}
	if err := r.client.HSet(ctx, "ttl", args...).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetTTls(ctx context.Context) (map[string]string, error) {
	ttls, err := r.client.HGetAll(ctx, "ttl").Result()

	if err != nil {
		return nil, err
	}
	return ttls, nil

}

func (r *Redis) DelUsersSegments(ctx context.Context, fields ...string) error {
	if err := r.client.HDel(ctx, "ttl", fields...).Err(); err != nil {
		return err
	}
	return nil
}
