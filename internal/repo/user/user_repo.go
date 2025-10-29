package user

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	customErrors "github.com/lsmltesting/MicroBlog/internal/errors"
	"github.com/lsmltesting/MicroBlog/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int, error)
	FindByID(ctx context.Context, ID int) (*models.User, error)
}

type inMemoryUserRepo struct {
	pool *pgxpool.Pool
	psql squirrel.StatementBuilderType
}

func NewInMemoryUserRepo(pool *pgxpool.Pool) UserRepository {
	return &inMemoryUserRepo{
		pool: pool,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *inMemoryUserRepo) Create(ctx context.Context, user *models.User) (int, error) {
	querySave, args, err := r.psql.
		Insert("users").
		Columns("username", "email", "password").
		Values(user.Username, user.Email, user.Password).
		Suffix("returning id").
		ToSql()
	if err != nil {
		return 0, customErrors.ErrFailedBuildQueryForPSQL
	}

	var id int
	if err := r.pool.QueryRow(ctx, querySave, args...).Scan(&id); err != nil {
		return 0, customErrors.ErrFailedCreateUserInPSQL
	}

	return id, nil
}

func (r *inMemoryUserRepo) FindByID(ctx context.Context, ID int) (*models.User, error) {

	queryFindUser, args, err := r.psql.
		Select("id", "username", "email", "password").
		From("users").
		Where(squirrel.Eq{"id": ID}).
		ToSql()
	if err != nil {
		return nil, customErrors.ErrFailedBuildQueryForPSQL
	}
	user := &models.User{}
	if err := r.pool.QueryRow(ctx, queryFindUser, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
	); err != nil {
		return nil, err
	}

	return user, nil
}
