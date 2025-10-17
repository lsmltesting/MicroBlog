package post

import (
	"context"
	"fmt"
	"testing"

	"github.com/lsmltesting/MicroBlog/db"
	"github.com/lsmltesting/MicroBlog/internal/logger"
	"github.com/lsmltesting/MicroBlog/internal/models"
	"github.com/lsmltesting/MicroBlog/internal/repo/post"
)

type mockUserService struct {
	userID int
}

func (m *mockUserService) CreateUser(ctx context.Context, userName string, email string, password string) (int, error) {
	m.userID++
	return m.userID, nil
}

func (m *mockUserService) GetUserByID(ctx context.Context, ID int) (*models.User, error) {
	return &models.User{
		Username: "testUserName",
		Email:    "test@test.ru",
		Password: "testPassword",
		ID:       ID,
	}, nil
}

func BenchmarkCreatePost(b *testing.B) {
	lg := logger.NewLogger(
		logger.LoggerConfig{
			BufferSize: 100,
			Workers:    6,
		},
	)

	ctx := context.Background()
	pool, _ := db.NewPostgresPool(lg)
	postRepo := post.NewInMemoryPostRepo(pool)
	mockUserService := &mockUserService{}

	postService := NewPostService(postRepo, mockUserService)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		postService.CreatePost(ctx, i, fmt.Sprintf("testTextForPost-%d", i))
	}
}

func BenchmarkGetPostByID(b *testing.B) {
	lg := logger.NewLogger(
		logger.LoggerConfig{
			BufferSize: 100,
			Workers:    6,
		},
	)

	ctx := context.Background()
	pool, _ := db.NewPostgresPool(lg)
	postRepo := post.NewInMemoryPostRepo(pool)
	mockUserService := &mockUserService{}

	postService := NewPostService(postRepo, mockUserService)

	postsID := make([]int, b.N)
	// creating posts
	for i := 0; i < b.N; i++ {
		postID, _ := postService.CreatePost(ctx, i, fmt.Sprintf("testTextForPost-%d", i))
		postsID[i] = postID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		postService.GetPostByID(ctx, postsID[i])
	}
}

func BenchmarkGetAllPosts(b *testing.B) {
	lg := logger.NewLogger(
		logger.LoggerConfig{
			BufferSize: 100,
			Workers:    6,
		},
	)

	ctx := context.Background()
	pool, _ := db.NewPostgresPool(lg)
	postRepo := post.NewInMemoryPostRepo(pool)
	mockUserService := &mockUserService{}

	postService := NewPostService(postRepo, mockUserService)

	postsID := make([]int, b.N)
	// creating posts
	for i := 0; i < b.N; i++ {
		postID, _ := postService.CreatePost(ctx, i, fmt.Sprintf("testTextForPost-%d", i))
		postsID[i] = postID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		postService.GetAllPosts(ctx)
	}
}
