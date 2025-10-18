package like

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lsmltesting/MicroBlog/db"
	"github.com/lsmltesting/MicroBlog/internal/logger"
	"github.com/lsmltesting/MicroBlog/internal/models"
	"github.com/lsmltesting/MicroBlog/internal/repo/like"
)

type mockUserService struct {
	userID int
}

func (mu *mockUserService) CreateUser(ctx context.Context, userName string, email string, password string) (int, error) {
	mu.userID++
	return mu.userID, nil
}

func (mu *mockUserService) GetUserByID(ctx context.Context, ID int) (*models.User, error) {
	return &models.User{
		Username: "testUserName",
		Email:    "test@test.ru",
		Password: "testPassword",
		ID:       ID,
	}, nil
}

type mockPostService struct {
	postID int
}

func (mp *mockPostService) CreatePost(ctx context.Context, userID int, text string) (int, error) {
	mp.postID++
	return mp.postID, nil
}

func (mp *mockPostService) GetPostByID(ctx context.Context, postID int) (*models.Post, error) {
	return &models.Post{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Text:      "testPostText",
		UserID:    1,
		ID:        postID,
	}, nil
}

func (mp *mockPostService) GetAllPosts(ctx context.Context) (map[int]*models.Post, error) {
	posts := make(map[int]*models.Post)
	for i := 0; i < 10; i++ {
		posts[i] = &models.Post{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Text:      fmt.Sprintf("testPostText-%d", i),
			UserID:    i + 10,
			ID:        i,
		}
	}

	return posts, nil
}

func BenchmarkCreateLike(b *testing.B) {
	lg := logger.NewLogger(
		logger.LoggerConfig{
			BufferSize: 100,
			Workers:    6,
		},
	)

	ctx := context.Background()
	pool, _ := db.NewPostgresPool(lg)
	likeRepo := like.NewInMemoryLikeRepo(pool)

	mockUserService := &mockUserService{}
	mockPostService := &mockPostService{}

	likeService := NewLikeService(likeRepo, mockUserService, mockPostService)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		likeService.CreateLike(ctx, i, i)
	}
}

func BenchmarkGetLikeById(b *testing.B) {
	lg := logger.NewLogger(
		logger.LoggerConfig{
			BufferSize: 100,
			Workers:    6,
		},
	)

	ctx := context.Background()
	pool, _ := db.NewPostgresPool(lg)
	likeRepo := like.NewInMemoryLikeRepo(pool)

	mockUserService := &mockUserService{}
	mockPostService := &mockPostService{}

	likeService := NewLikeService(likeRepo, mockUserService, mockPostService)

	likesID := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		likeID, _ := likeService.CreateLike(ctx, i, i)
		likesID[i] = likeID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		likeService.GetLikeById(ctx, likesID[i])
	}
}

func BenchmarkGetAllLikes(b *testing.B) {
	lg := logger.NewLogger(
		logger.LoggerConfig{
			BufferSize: 100,
			Workers:    6,
		},
	)

	ctx := context.Background()
	pool, _ := db.NewPostgresPool(lg)
	likeRepo := like.NewInMemoryLikeRepo(pool)

	mockUserService := &mockUserService{}
	mockPostService := &mockPostService{}

	likeService := NewLikeService(likeRepo, mockUserService, mockPostService)

	likesID := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		likeID, _ := likeService.CreateLike(ctx, i, i)
		likesID[i] = likeID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		likeService.GetAllLikes(ctx)
	}
}
