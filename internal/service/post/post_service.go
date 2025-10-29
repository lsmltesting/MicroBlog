package post

import (
	"context"

	"github.com/lsmltesting/MicroBlog/internal/models"
	"github.com/lsmltesting/MicroBlog/internal/repo/post"
	"github.com/lsmltesting/MicroBlog/internal/service/user"
)

type PostService interface {
	Create(ctx context.Context, user int, text string) (int, error)
	FindByID(ctx context.Context, postID int) (*models.Post, error)
	GetAll(ctx context.Context) (map[int]*models.Post, error)
}

type postService struct {
	repo        post.PostRepository
	userService user.UserService
}

func NewPostService(repo post.PostRepository, userService user.UserService) PostService {
	return &postService{
		repo:        repo,
		userService: userService,
	}
}

func (s *postService) Create(ctx context.Context, userID int, text string) (int, error) {
	// Check if user with shared userId is exists
	_, err := s.userService.FindByID(ctx, userID)
	if err != nil {
		return 0, err
	}

	post, err := models.NewPost(userID, text)
	if err != nil {
		return 0, err
	}

	// After creating post update user's posthistory map
	postID, err := s.repo.Create(ctx, post)
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (s *postService) FindByID(ctx context.Context, postID int) (*models.Post, error) {
	return s.repo.FindByID(ctx, postID)
}

func (s *postService) GetAll(ctx context.Context) (map[int]*models.Post, error) {
	return s.repo.GetAll(ctx)
}
