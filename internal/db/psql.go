package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

func (d *Db) PgError(err error) error {
	pgError, ok := err.(*pgconn.PgError)
	if !ok {
		return nil
	}
	switch pgError.Code {
	case "23505":
		return errors.New("already exists")
	case notError:
		return nil
	case emptyValue:
		return errors.New(emptyValue)
	default:
		return err
	}
	return nil
}

const (
	notError   = "no rows in result set"
	emptyValue = "empty value"
)

func (d *Db) GetUserIds(ctx context.Context) ([]int, error) {
	q := `select id  from  users`
	rows, err := d.client.Query(ctx, q)
	if d.PgError(err) != nil {
		return nil, d.PgError(err)
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[int])
	return ids, nil

}
func (d *Db) GetUsersSegments(ctx context.Context, userId int) ([]Segment, error) {
	q := `select slug  from  segments  as s
    inner join  user_segment as us on s.id=us.segment_id
	inner join  users as u on u.id=us.user_id
	where u.id = $1
`

	rows, err := d.client.Query(ctx, q, userId)

	if err != nil {
		return nil, d.PgError(err)
	}
	segments, err := pgx.CollectRows(rows, pgx.RowToStructByName[Segment])
	if err != nil {
		return nil, d.PgError(err)
	}
	return segments, nil

}

func (d *Db) DeleteSegmentsFromUser(ctx context.Context, userId int, slugs ...any) (err error) {
	tx, err := d.client.Begin(ctx)
	defer tx.Commit(ctx)

	q := `delete from user_segment as us
       where us.user_id = $1 
         and 
       us.segment_id in (select id from segments  as s where s.slug in ( `
	for _, slug := range slugs {
		toAdd := fmt.Sprintf(`'%s',`, slug)
		q += toAdd

	}
	q = q[0:len(q)-1] + ")) "
	if err := tx.QueryRow(ctx, q, userId).Scan(); d.PgError(err) != nil {
		return d.PgError(err)
	}
	if err := d.AddToHistoryUserSlugs(ctx, tx, userId, false, slugs...); d.PgError(err) != nil {
		if errtx := tx.Rollback(ctx); errtx != nil {
			return errors.New("tx error")
		}
		return d.PgError(err)
	}
	return nil
}
func (d *Db) AddSegmentsToUser(ctx context.Context, userId int, slugs ...any) (err error) {
	tx, err := d.client.Begin(ctx)
	defer tx.Commit(ctx)
	if err != nil {
		return d.PgError(err)
	}
	q := `insert into user_segment (user_id,segment_id) select $1,segments.id from segments where slug in ( `
	var args = []any{userId}

	for i, slug := range slugs {
		toAdd := fmt.Sprintf(`$%d,`, i+2)
		q += toAdd
		args = append(args, slug)

	}
	q = q[0:len(q)-1] + ")"

	if err := tx.QueryRow(ctx, q, args...).Scan(); d.PgError(err) != nil {
		return d.PgError(err)
	}
	if err := d.AddToHistoryUserSlugs(ctx, tx, userId, true, slugs...); d.PgError(err) != nil {
		if errtx := tx.Rollback(ctx); errtx != nil {
			return errors.New("tx error")
		}
		return d.PgError(err)
	}
	return nil
}
func (d *Db) AddSlugIdToUsers(ctx context.Context, slugId int, ids ...int) (err error) {
	tx, err := d.client.Begin(ctx)
	defer tx.Commit(ctx)
	if err != nil {
		return d.PgError(err)
	}

	q := `insert into user_segment (user_id,segment_id) values `
	var args = []any{slugId}

	for i, id := range ids {
		toAdd := fmt.Sprintf(`($%d,$1),`, i+2)
		q += toAdd
		args = append(args, id)

	}
	q = q[0 : len(q)-1]

	if err := tx.QueryRow(ctx, q, args...).Scan(); d.PgError(err) != nil {
		return d.PgError(err)
	}
	if err := d.AddToHistorySlugUsers(ctx, tx, slugId, true, ids...); d.PgError(err) != nil {
		if errtx := tx.Rollback(ctx); errtx != nil {
			return errors.New("tx error")
		}
		return d.PgError(err)
	}
	return nil
}
func (d *Db) CreateUser(ctx context.Context, username string) error {
	q := `insert into users (username) values ($1) `
	if err := d.client.QueryRow(ctx, q, username).Scan(); err != nil {
		return d.PgError(err)
	}
	return nil
}

func (d *Db) DeleteUser(ctx context.Context, id int) error {
	q := `delete from users where id = $1 `
	if err := d.client.QueryRow(ctx, q, id).Scan(); err != nil {
		return d.PgError(err)
	}
	return nil
}

func (d *Db) CreateSegment(ctx context.Context, slug string) (int, error) {
	q := `insert into segments (slug) values ($1) returning id`
	var id int
	if err := d.client.QueryRow(ctx, q, slug).Scan(&id); err != nil {
		return -1, d.PgError(err)
	}
	return id, nil
}

func (d *Db) DeleteSegment(ctx context.Context, slug string) error {
	q := `delete from segments where slug = $1 `
	if err := d.client.QueryRow(ctx, q, slug).Scan(); err != nil {
		return d.PgError(err)
	}
	return nil
}

func (d *Db) AddToHistorySlugUsers(ctx context.Context, tx pgx.Tx, slugId int, operationType bool, ids ...int) error {

	var operationStr string
	if operationType {
		operationStr = "insert"
	} else {
		operationStr = "delete"
	}
	var args []any = []any{slugId, operationStr, time.Now()}

	q := `insert into history (user_id,slug,operation,update_time) values `

	for i, id := range ids {
		toAdd := fmt.Sprintf(` ($%d,$1,$2,$3),`, i+4)
		q += toAdd
		args = append(args, id)
	}
	q = q[0 : len(q)-1]
	fmt.Println(q, ids)
	if err := tx.QueryRow(ctx, q, args...).Scan(); d.PgError(err) != nil {
		fmt.Println(err)
		return d.PgError(err)
	}
	return nil
}

func (d *Db) AddToHistoryUserSlugs(ctx context.Context, tx pgx.Tx, userId int, operationType bool, slugs ...any) error {
	var operationStr string
	if operationType {
		operationStr = "insert"
	} else {
		operationStr = "delete"
	}
	var args []any = []any{userId, operationStr, time.Now()}

	q := `insert into history (user_id,slug,operation,update_time) values `

	for i, slug := range slugs {
		toAdd := fmt.Sprintf(` ($1,$%d,$2,$3),`, i+4)
		q += toAdd
		args = append(args, slug)
	}
	q = q[0 : len(q)-1]

	if err := tx.QueryRow(ctx, q, args...).Scan(); d.PgError(err) != nil {
		return d.PgError(err)
	}
	return nil

}
func (d *Db) GetHistory(ctx context.Context, userId int, year, month int) ([]HistoryRow, error) {
	q := `select user_id,slug,operation,update_time from  history
		  where user_id = $1 and
		  date_part('year', update_time) = $2 and 
		   date_part('month', update_time)  = $3`

	rows, err := d.client.Query(ctx, q, userId, year, month)

	if d.PgError(err) != nil {
		return nil, d.PgError(err)
	}
	history, err := pgx.CollectRows(rows, pgx.RowToStructByName[HistoryRow])
	return history, nil
}
