package post

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	customErrors "github.com/lsmltesting/MicroBlog/internal/errors"
	"github.com/lsmltesting/MicroBlog/internal/models"
)

type PostRepository interface {
	Save(ctx context.Context, post *models.Post) (int, error)
	FindPostByID(ctx context.Context, postID int) (*models.Post, error)
	GetAllPosts(ctx context.Context) (map[int]*models.Post, error)
}

type inMemoryPostRepo struct {
	pool *pgxpool.Pool
	psql squirrel.StatementBuilderType
}

func NewInMemoryPostRepo(pool *pgxpool.Pool) PostRepository {
	return &inMemoryPostRepo{
		pool: pool,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *inMemoryPostRepo) Save(ctx context.Context, post *models.Post) (int, error) {
	querySave, args, err := r.psql.
		Insert("posts").
		Columns("created_at", "udpate_at", "text", "user_id").
		Values(post.CreatedAt, post.UpdatedAt, post.Text, post.UserID).
		Suffix("returning id").
		ToSql()

	if err != nil {
		return 0, customErrors.ErrFailedBuildQueryForPSQL
	}

	var id int
	if err := r.pool.QueryRow(ctx, querySave, args...).Scan(&id); err != nil {
		return 0, customErrors.ErrFailedCreatePostInPSQL
	}

	return id, nil
}

func (r *inMemoryPostRepo) FindPostByID(ctx context.Context, postID int) (*models.Post, error) {
	queryFind, args, err := r.psql.
		Select("created_at", "udpate_at", "text", "user_id", "id").
		From("posts").
		Where(squirrel.Eq{"id": postID}).
		ToSql()

	if err != nil {
		return nil, customErrors.ErrFailedBuildQueryForPSQL
	}

	post := &models.Post{}
	err = r.pool.QueryRow(ctx, queryFind, args...).Scan(
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Text,
		&post.UserID,
		&post.ID,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *inMemoryPostRepo) GetAllPosts(ctx context.Context) (map[int]*models.Post, error) {
	queryAll, args, err := r.psql.
		Select("created_at", "updated_at", "text", "user_id", "id").
		From("posts").
		ToSql()

	if err != nil {
		return nil, customErrors.ErrFailedBuildQueryForPSQL
	}

	rows, err := r.pool.Query(ctx, queryAll, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make(map[int]*models.Post)
	for rows.Next() {
		post := &models.Post{}

		err = rows.Scan(
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Text,
			&post.UserID,
			&post.ID,
		)
		if err != nil {
			return nil, err
		}

		posts[post.ID] = post
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
