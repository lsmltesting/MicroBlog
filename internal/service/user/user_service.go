package user

import (
	"context"

	"github.com/lsmltesting/MicroBlog/internal/models"
	"github.com/lsmltesting/MicroBlog/internal/repo/user"
)

type UserService interface {
	CreateUser(ctx context.Context, username string, email string, password string) (int, error)
	GetUserByID(ctx context.Context, ID int) (*models.User, error)
}

type userService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, username string, email string, password string) (int, error) {
	user, err := models.NewUser(username, email, password)
	if err != nil {
		return 0, err
	}
	return s.repo.Save(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, ID int) (*models.User, error) {
	return s.repo.FindUserByID(ctx, ID)
}
