package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsRepository interface {
	IncrementPostLikes(ctx context.Context, postID int, eventID string) error
	GetPostLikeCount(ctx context.Context, postID int) (int64, error)
}

type statsRepository struct {
	pool *pgxpool.Pool
	psql squirrel.StatementBuilderType
}

func NewStatsRepository(pool *pgxpool.Pool) StatsRepository {
	return &statsRepository{
		pool: pool,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *statsRepository) IncrementPostLikes(ctx context.Context, postID int, eventID string) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Idempotency: we record the event_id, if it already exists, we simply exit
	q1, args1, err := r.psql.
		Insert("processed_events").
		Columns("event_id").
		Values(eventID).
		Suffix("ON CONFLICT (event_id) DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}

	tag, err := tx.Exec(ctx, q1, args1...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		// already handled
		return nil
	}

	q2, args2, err := r.psql.
		Insert("post_like_counters").
		Columns("post_id", "like_count").
		Values(postID, 1).
		Suffix(`
			ON CONFLICT (post_id)
			DO UPDATE SET like_count = post_like_counters.like_count + EXCLUDED.like_count
		`).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, q2, args2...); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *statsRepository) GetPostLikeCount(ctx context.Context, postID int) (int64, error) {
	q, args, err := r.psql.
		Select("like_count").
		From("post_like_counters").
		Where(squirrel.Eq{"post_id": postID}).
		ToSql()
	if err != nil {
		return 0, err
	}

	var count int64
	err = r.pool.QueryRow(ctx, q, args...).Scan(&count)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}
