package post

import (
	"context"

	"github.com/lsmltesting/MicroBlog/internal/models"
	"github.com/lsmltesting/MicroBlog/internal/repo/post"
	"github.com/lsmltesting/MicroBlog/internal/service/user"
)

type PostService interface {
	CreatePost(ctx context.Context, user int, text string) (int, error)
	GetPostByID(ctx context.Context, postID int) (*models.Post, error)
	GetAllPosts(ctx context.Context) (map[int]*models.Post, error)
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

func (s *postService) CreatePost(ctx context.Context, userID int, text string) (int, error) {
	// Check if user with shared userId is exists
	_, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return 0, err
	}

	post, err := models.NewPost(userID, text)
	if err != nil {
		return 0, err
	}

	// After creating post update user's posthistory map
	postID, err := s.repo.Save(ctx, post)
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func (s *postService) GetPostByID(ctx context.Context, postID int) (*models.Post, error) {
	return s.repo.FindPostByID(ctx, postID)
}

func (s *postService) GetAllPosts(ctx context.Context) (map[int]*models.Post, error) {
	return s.repo.GetAllPosts(ctx)
}
