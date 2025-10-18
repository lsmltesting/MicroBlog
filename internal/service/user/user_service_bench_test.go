package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/lsmltesting/MicroBlog/db"
	"github.com/lsmltesting/MicroBlog/internal/logger"
	"github.com/lsmltesting/MicroBlog/internal/repo/user"
)

func BenchmarkCreateUser(b *testing.B) {
	lg := logger.NewLogger(
		logger.LoggerConfig{
			BufferSize: 100,
			Workers:    6,
		},
	)

	pool, _ := db.NewPostgresPool(lg)

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

	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userService.CreateUser(ctx, users[i].username, users[i].email, users[i].password)
	}
}

func BenchmarkGetUserByID(b *testing.B) {
	lg := logger.NewLogger(
		logger.LoggerConfig{
			BufferSize: 100,
			Workers:    6,
		},
	)

	pool, _ := db.NewPostgresPool(lg)
	userRepo := user.NewInMemoryUserRepo(pool)
	userService := NewUserService(userRepo)

	usersID := make([]int, b.N)
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		userID, _ := userService.CreateUser(
			ctx,
			fmt.Sprintf("testUsername-%d", i),
			fmt.Sprintf("testemail-%d@gtest.ru", i),
			fmt.Sprintf("testPassword123qwerty-%d", i),
		)
		usersID[i] = userID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userService.GetUserByID(ctx, usersID[i])
	}
}
