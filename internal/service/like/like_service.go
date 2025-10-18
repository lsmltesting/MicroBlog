package like

import (
	"context"

	"github.com/lsmltesting/MicroBlog/internal/models"
	"github.com/lsmltesting/MicroBlog/internal/repo/like"
	"github.com/lsmltesting/MicroBlog/internal/service/post"
	"github.com/lsmltesting/MicroBlog/internal/service/user"
)

type LikeService interface {
	CreateLike(ctx context.Context, userID int, postID int) (int, error)
	GetLikeById(ctx context.Context, likeID int) (*models.Like, error)
	GetAllLikes(ctx context.Context) (map[int]*models.Like, error)
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

func (l *likeService) CreateLike(ctx context.Context, userID int, postID int) (int, error) {
	// Check if user exists
	_, err := l.userService.GetUserByID(ctx, userID)
	if err != nil {
		return 0, err
	}

	// Check if post exists
	_, err = l.postService.GetPostByID(ctx, postID)
	if err != nil {
		return 0, err
	}

	like := models.NewLike(userID, postID)

	// First save like in repo. After saving like will have actual ID
	likeID, err := l.repo.Save(ctx, like)
	if err != nil {
		return 0, err
	}

	return likeID, nil
}

func (l *likeService) GetLikeById(ctx context.Context, likeID int) (*models.Like, error) {
	return l.repo.FindLikeById(ctx, likeID)
}

func (l *likeService) GetAllLikes(ctx context.Context) (map[int]*models.Like, error) {
	return l.repo.GetAllLikes(ctx)
}
