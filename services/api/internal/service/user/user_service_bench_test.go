package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/lsmltesting/MicroBlog/db"
	"github.com/lsmltesting/MicroBlog/services/api/internal/logger"
	"github.com/lsmltesting/MicroBlog/services/api/internal/repo/user"
)

func BenchmarkCreateUser(b *testing.B) {
	lg := logger.NewLogger(
		logger.LoggerConfig{
			BufferSize: 100,
			Workers:    6,
		},
	)
	ctx := context.Background()
	pool, _ := db.NewPostgresPool(lg, ctx)

	userRepo := user.NewInMemoryUserRepo(pool)
	userService := NewUserService(userRepo)

	users := make(
		[]struct {
			username string
			email    string
			password string
		},
		b.N,
	)
	for i := range users {
		users[i].username = fmt.Sprintf("testUsername-%d", i)
		users[i].email = fmt.Sprintf("testemail-%d@gtest.ru", i)
		users[i].password = fmt.Sprintf("testPassword123qwerty-%d", i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userService.Create(ctx, users[i].username, users[i].email, users[i].password)
	}
}

func BenchmarkGetUserByID(b *testing.B) {
	lg := logger.NewLogger(
		logger.LoggerConfig{
			BufferSize: 100,
			Workers:    6,
		},
	)
	ctx := context.Background()

	pool, _ := db.NewPostgresPool(lg, ctx)
	userRepo := user.NewInMemoryUserRepo(pool)
	userService := NewUserService(userRepo)

	usersID := make([]int, b.N)

	for i := 0; i < b.N; i++ {
		userID, _ := userService.Create(
			ctx,
			fmt.Sprintf("testUsername-%d", i),
			fmt.Sprintf("testemail-%d@gtest.ru", i),
			fmt.Sprintf("testPassword123qwerty-%d", i),
		)
		usersID[i] = userID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userService.FindByID(ctx, usersID[i])
	}
}
