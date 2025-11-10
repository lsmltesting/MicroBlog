package like

import (
	"context"

	"github.com/lsmltesting/MicroBlog/services/api/internal/models"
	"github.com/lsmltesting/MicroBlog/services/api/internal/repo/like"
	"github.com/lsmltesting/MicroBlog/services/api/internal/service/post"
	"github.com/lsmltesting/MicroBlog/services/api/internal/service/user"
)

type LikeService interface {
	Create(ctx context.Context, userID int, postID int) (int, error)
	FindById(ctx context.Context, likeID int) (*models.Like, error)
	GetAll(ctx context.Context) (map[int]*models.Like, error)
}

type likeService struct {
	repo like.LikeRepository

	userService user.UserService
	postService post.PostService
}

func NewLikeService(
	repo like.LikeRepository,
	userService user.UserService,
	postService post.PostService,
) LikeService {
	return &likeService{
		repo:        repo,
		userService: userService,
		postService: postService,
	}
}

func (l *likeService) Create(ctx context.Context, userID int, postID int) (int, error) {
	// Check if user exists
	_, err := l.userService.FindByID(ctx, userID)
	if err != nil {
		return 0, err
	}

	// Check if post exists
	_, err = l.postService.FindByID(ctx, postID)
	if err != nil {
		return 0, err
	}

	like := models.NewLike(userID, postID)

	// First save like in repo. After saving like will have actual ID
	likeID, err := l.repo.Create(ctx, like)
	if err != nil {
		return 0, err
	}

	return likeID, nil
}

func (l *likeService) FindById(ctx context.Context, likeID int) (*models.Like, error) {
	return l.repo.FindById(ctx, likeID)
}

func (l *likeService) GetAll(ctx context.Context) (map[int]*models.Like, error) {
	return l.repo.GetAll(ctx)
}
