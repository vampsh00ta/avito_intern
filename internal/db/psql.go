package repository

import (
	"avito/pkg/client"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type Db struct {
	client postgresql.Client
	log    *zap.Logger
}

func (d *Db) PgError(err error) error {
	if pgError, ok := err.(*pgconn.PgError); ok {
		d.log.Debug(pgError.Error())
		return errors.New(dbError)
	}
	return nil
}

const (
	dbError = "smth went wrong with db"
)

func (d *Db) GetUsersTags(ctx context.Context, userId int) ([]Segment, error) {
	q := `select slug  from  segments  as s
    inner join  user_segment as us on s.id=us.segment_id
	inner join  users as u on u.id=us.user_id
	where u.id = $1
`

	rows, err := d.client.Query(ctx, q, userId)

	if err != nil {
		return nil, d.PgError(err)
	}
	segments, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, d.PgError(err)
	}
	fmt.Println(segments)
	if err := d.client.QueryRow(ctx, q, userId).Scan(); err != nil {
		return nil, d.PgError(err)
	}
	return nil, nil

}
func (d *Db) GetSegmentsIds(ctx context.Context, tx interface{}, slugs ...any) (ids []int, err error) {
	q := `select id from  segments where slug in ( `
	for i, _ := range slugs {
		toAdd := fmt.Sprintf(`$%d,`, i+1)
		q += toAdd
	}
	txOk, ok := tx.(pgx.Tx)

	q = q[0 : len(q)-1]
	q += ")"
	var rows pgx.Rows

	if ok {
		rows, err = txOk.Query(ctx, q, slugs...)

	} else {
		rows, err = d.client.Query(ctx, q, slugs...)

	}

	if err != nil {
		return nil, d.PgError(err)
	}
	ids, err = pgx.CollectRows(rows, pgx.RowTo[int])
	if err != nil {
		return nil, d.PgError(err)
	}

	return ids, nil
}
func (d *Db) AddTagsToUser(ctx context.Context, userId int, slugs ...any) (err error) {
	tx, err := d.client.Begin(ctx)
	defer tx.Commit(ctx)
	if err != nil {
		return d.PgError(err)
	}
	slugsIds, err := d.GetSegmentsIds(ctx, tx, slugs...)
	if err != nil {
		return d.PgError(err)
	}

	q := `insert into user_segment (user_id,segment_id) values  `
	for _, id := range slugsIds {
		toAdd := fmt.Sprintf(` ($1,%d),`, id)
		q += toAdd

	}
	q = q[0 : len(q)-1]

	if err := tx.QueryRow(ctx, q, userId).Scan(); err != nil {
		return d.PgError(err)
	}

	return nil
}

func (d *Db) CreateUser(ctx context.Context, username string) error {
	q := `insert into users (username) values ($1)`
	if err := d.client.QueryRow(ctx, q, username).Scan(); err != nil {
		return d.PgError(err)
	}
	return nil
}

func (d *Db) DeleteUser(ctx context.Context, username string) error {
	q := `delete from users where username = $1 `
	if err := d.client.QueryRow(ctx, q, username).Scan(); err != nil {
		return d.PgError(err)
	}
	return nil
}

func (d *Db) CreateSegment(ctx context.Context, slug string) error {
	q := `insert into segments (slug) values ($1)`
	if err := d.client.QueryRow(ctx, q, slug).Scan(); err != nil {
		return d.PgError(err)
	}
	return nil
}

func (d *Db) DeleteSegment(ctx context.Context, slug string) error {
	q := `delete from segments where slug = $1 `
	if err := d.client.QueryRow(ctx, q, slug).Scan(); err != nil {
		return d.PgError(err)
	}
	return nil
}

func New(client postgresql.Client, logger *zap.Logger) Repository {
	return &Db{
		client,
		logger,
	}
}
