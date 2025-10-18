package like

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	customErrors "github.com/lsmltesting/MicroBlog/internal/errors"
	"github.com/lsmltesting/MicroBlog/internal/models"
)

type LikeRepository interface {
	Save(ctx context.Context, like *models.Like) (int, error)
	FindLikeById(ctx context.Context, likeID int) (*models.Like, error)
	GetAllLikes(ctx context.Context) (map[int]*models.Like, error)
}

type inMemoryLikeRepo struct {
	pool *pgxpool.Pool
	psql squirrel.StatementBuilderType
}

func NewInMemoryLikeRepo(pool *pgxpool.Pool) LikeRepository {
	return &inMemoryLikeRepo{
		pool: pool,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (l *inMemoryLikeRepo) Save(ctx context.Context, like *models.Like) (int, error) {
	querySave, args, err := l.psql.
		Insert("likes").
		Columns("created_at", "post_id", "user_id").
		Values(like.CreatedAt, like.PostID, like.UserID).
		Suffix("returning id").
		ToSql()

	if err != nil {
		return 0, customErrors.ErrFailedBuildQueryForPSQL
	}

	var id int
	if err := l.pool.QueryRow(ctx, querySave, args...).Scan(&id); err != nil {
		return 0, customErrors.ErrFailedCreateLikeInPSQL
	}

	return id, nil
}

func (l *inMemoryLikeRepo) FindLikeById(ctx context.Context, likeID int) (*models.Like, error) {
	queryFind, args, err := l.psql.
		Select("id", "created_at", "post_id", "user_id").
		From("likes").
		Where(squirrel.Eq{"id": likeID}).
		ToSql()

	if err != nil {
		return nil, customErrors.ErrFailedBuildQueryForPSQL
	}

	like := &models.Like{}
	if err := l.pool.QueryRow(ctx, queryFind, args...).Scan(
		&like.ID,
		&like.CreatedAt,
		&like.PostID,
		&like.UserID,
	); err != nil {
		return nil, err
	}

	return like, nil
}

func (l *inMemoryLikeRepo) GetAllLikes(ctx context.Context) (map[int]*models.Like, error) {
	queryFind, args, err := l.psql.
		Select("id", "created_at", "post_id", "user_id").
		From("likes").
		ToSql()

	if err != nil {
		return nil, customErrors.ErrFailedBuildQueryForPSQL
	}

	rows, err := l.pool.Query(ctx, queryFind, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	likes := make(map[int]*models.Like)
	for rows.Next() {
		like := &models.Like{}

		err := rows.Scan(
			&like.ID,
			&like.CreatedAt,
			&like.PostID,
			&like.UserID,
		)

		if err != nil {
			return nil, err
		}

		likes[like.ID] = like
	}

	return likes, nil
}
