package user

import (
	"context"

	"github.com/lsmltesting/MicroBlog/services/api/internal/models"
	"github.com/lsmltesting/MicroBlog/services/api/internal/repo/user"
)

type UserService interface {
	Create(ctx context.Context, username string, email string, password string) (int, error)
	FindByID(ctx context.Context, ID int) (*models.User, error)
}

type userService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Create(ctx context.Context, username string, email string, password string) (int, error) {
	user, err := models.NewUser(username, email, password)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, user)
}

func (s *userService) FindByID(ctx context.Context, ID int) (*models.User, error) {
	return s.repo.FindByID(ctx, ID)
}
