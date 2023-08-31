package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

func (d *Db) isError(err error) error {
	_, ok := err.(*pgconn.PgError)

	if ok {
		return err
	}

	return nil
}

const (
	notError      = "no rows in result set"
	emptyValue    = "empty value"
	alreadyExists = "already exists"
)

func (d *Db) GetUserIds(ctx context.Context) ([]int, error) {
	q := `select id  from  users`
	rows, err := d.client.Query(ctx, q)
	if d.isError(err) != nil {
		return nil, err
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[int])
	return ids, nil

}
func (d *Db) GetUsersSegments(ctx context.Context, userId int) (*[]Segment, error) {
	q := `select slug  from  segments  as s
    inner join  user_segment as us on s.id=us.segment_id
	inner join  users as u on u.id=us.user_id
	where u.id = $1
`

	rows, err := d.client.Query(ctx, q, userId)
	fmt.Println(rows)
	if d.isError(err) != nil {
		return nil, err
	}
	segments, err := pgx.CollectRows(rows, pgx.RowToStructByName[Segment])
	if err != nil {
		return nil, d.isError(err)
	}
	return &segments, nil

}

func (d *Db) DeleteSegmentsFromUser(ctx context.Context, userId int, segments ...*Segment) (err error) {
	tx, err := d.client.Begin(ctx)
	defer tx.Commit(ctx)
	var args []any = []any{userId}
	q := `delete from user_segment as us
       where us.user_id = $1 
         and 
       us.segment_id in (select id from segments  as s where s.slug in ( `
	for i, segment := range segments {
		toAdd := fmt.Sprintf(`$%d,`, i+2)
		q += toAdd
		args = append(args, segment.Slug)

	}
	q = q[0:len(q)-1] + ")) "
	if err := tx.QueryRow(ctx, q, args...).Scan(); d.isError(err) != nil {
		return err
	}
	if err := d.AddToHistoryUserSlugs(ctx, tx, userId, false, segments...); d.isError(err) != nil {
		if errtx := tx.Rollback(ctx); d.isError(errtx) != nil {
			return err
		}
		return err
	}
	return nil
}
func (d *Db) AddSegmentsToUser(ctx context.Context, userId int, segments ...*Segment) (err error) {
	tx, err := d.client.Begin(ctx)
	defer tx.Commit(ctx)
	if err != nil {
		return d.isError(err)
	}
	q := `insert into user_segment (user_id,segment_id) select $1,segments.id from segments where slug in ( `
	var args = []any{userId}

	for i, segment := range segments {
		toAdd := fmt.Sprintf(`$%d,`, i+2)
		q += toAdd
		args = append(args, segment.Slug)

	}
	q = q[0:len(q)-1] + ")"
	if err := tx.QueryRow(ctx, q, args...).Scan(); d.isError(err) != nil {
		return err
	}
	if err := d.AddToHistoryUserSlugs(ctx, tx, userId, true, segments...); d.isError(err) != nil {
		if errtx := tx.Rollback(ctx); d.isError(errtx) != nil {
			return err
		}
		return err
	}
	return nil
}
func (d *Db) AddSlugIdToUsers(ctx context.Context, segment Segment, ids ...int) (err error) {
	tx, err := d.client.Begin(ctx)
	defer tx.Commit(ctx)
	if d.isError(err) != nil {
		return err
	}

	q := `insert into user_segment (user_id,segment_id) values `
	var args = []any{segment.Id}
	for i, id := range ids {
		toAdd := fmt.Sprintf(`($%d,$1),`, i+2)
		q += toAdd
		args = append(args, id)
	}
	q = q[0 : len(q)-1]

	if err := tx.QueryRow(ctx, q, args...).Scan(); d.isError(err) != nil {
		return err
	}
	if err := d.AddToHistorySlugUsers(ctx, tx, segment, true, ids...); d.isError(err) != nil {
		if errtx := tx.Rollback(ctx); d.isError(errtx) != nil {
			return err
		}
		return err
	}
	return nil
}
func (d *Db) CreateUser(ctx context.Context, username string) error {
	q := `insert into users (username) values ($1)`

	if err := d.client.QueryRow(ctx, q, username).Scan(); d.isError(err) != nil {
		return err
	}
	return nil
}

func (d *Db) DeleteUser(ctx context.Context, id int) error {
	q := `delete from users where id = $1 `
	if err := d.client.QueryRow(ctx, q, id).Scan(); d.isError(err) != nil {
		return err
	}
	return nil
}

func (d *Db) CreateSegment(ctx context.Context, segment Segment) (int, error) {
	q := `insert into segments (slug) values ($1) returning id`
	var id int
	if err := d.client.QueryRow(ctx, q, segment.Slug).Scan(&id); d.isError(err) != nil {
		return -1, err
	}
	return id, nil
}

func (d *Db) DeleteSegment(ctx context.Context, segment Segment) error {
	q := `delete from segments where slug = $1 `
	if err := d.client.QueryRow(ctx, q, segment.Slug).Scan(); d.isError(err) != nil {
		return err
	}
	return nil
}

func (d *Db) AddToHistorySlugUsers(ctx context.Context, tx pgx.Tx, segment Segment, operationType bool, ids ...int) error {

	var operationStr string
	if operationType {
		operationStr = "insert"
	} else {
		operationStr = "delete"
	}
	var args []any = []any{segment.Slug, operationStr, time.Now()}

	q := `insert into history (user_id,slug,operation,update_time) values `

	for i, id := range ids {
		toAdd := fmt.Sprintf(` ($%d,$1,$2,$3),`, i+4)
		q += toAdd
		args = append(args, id)
	}
	q = q[0 : len(q)-1]
	if err := tx.QueryRow(ctx, q, args...).Scan(); d.isError(err) != nil {

		return err
	}
	return nil
}

func (d *Db) AddToHistoryUserSlugs(ctx context.Context, tx pgx.Tx, userId int, operationType bool, segments ...*Segment) error {
	var operationStr string
	if operationType {
		operationStr = "insert"
	} else {
		operationStr = "delete"
	}
	var args []any = []any{userId, operationStr, time.Now()}

	q := `insert into history (user_id,slug,operation,update_time) values `

	for i, segment := range segments {
		toAdd := fmt.Sprintf(` ($1,$%d,$2,$3),`, i+4)
		q += toAdd
		args = append(args, segment.Slug)
	}
	q = q[0 : len(q)-1]

	if err := tx.QueryRow(ctx, q, args...).Scan(); d.isError(err) != nil {
		return err
	}
	return nil

}
func (d *Db) GetHistoryById(ctx context.Context, userId int, year, month int) (*[]HistoryRow, error) {
	q := `select user_id,slug,operation,update_time from  history
		  where user_id = $1 and
		  date_part('year', update_time) = $2 and 
		   date_part('month', update_time)  = $3`

	rows, err := d.client.Query(ctx, q, userId, year, month)

	if d.isError(err) != nil {
		return nil, d.isError(err)
	}
	history, err := pgx.CollectRows(rows, pgx.RowToStructByName[HistoryRow])
	return &history, nil
}

func (d *Db) GetHistoryAll(ctx context.Context, year, month int) (*[]HistoryRow, error) {
	q := `select user_id,slug,operation,update_time from  history
		  where 
		  date_part('year', update_time) = $1 and 
		   date_part('month', update_time)  = $2`

	rows, err := d.client.Query(ctx, q, year, month)

	if d.isError(err) != nil {
		return nil, err
	}
	history, err := pgx.CollectRows(rows, pgx.RowToStructByName[HistoryRow])
	return &history, nil
}
