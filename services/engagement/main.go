package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/lsmltesting/MicroBlog/db"
	"github.com/lsmltesting/MicroBlog/services/api/internal/events"
	"github.com/lsmltesting/MicroBlog/services/api/internal/logger"
	"github.com/lsmltesting/MicroBlog/services/api/internal/trace"
	"github.com/lsmltesting/MicroBlog/services/engagement/internal/consumer"
	"github.com/lsmltesting/MicroBlog/services/engagement/internal/handlers"
	"github.com/lsmltesting/MicroBlog/services/engagement/internal/repository"
)

func main() {
	lg := logger.NewLogger(
		logger.LoggerConfig{
			BufferSize: 100,
			Workers:    4,
		},
	)

	rootCtx := context.Background()
	ctx, stop := signal.NotifyContext(rootCtx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPostgresPool(lg, ctx)
	if err != nil {
		lg.AddLog(
			logger.LevelError,
			logger.SourceMain,
			map[string]string{},
			"failed to connect to db: "+err.Error(),
		)
		return
	}
	defer pool.Close()

	statsRepo := repository.NewStatsRepository(pool)

	eventsCfg := events.ConfigFromEnv()
	postLikedConsumer := consumer.NewPostLikedConsumer(eventsCfg, statsRepo, lg)
	defer postLikedConsumer.Close()

	// запускаем consumer
	go func() {
		if err := postLikedConsumer.Run(ctx); err != nil {
			lg.AddLog(
				logger.LevelError,
				logger.SourceConsumer,
				map[string]string{},
				"consumer stopped: "+err.Error(),
			)
		}
	}()

	// HTTP
	r := mux.NewRouter()
	r.Use(trace.Middleware)

	statsHandler := handlers.NewStatsHandler(statsRepo)
	r.HandleFunc("/stats/posts/{id}", statsHandler.GetPostStats).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:         ":8090",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		lg.AddLog(
			logger.LevelInfo,
			logger.SourceMain,
			map[string]string{},
			"engagement http server started on :8090",
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.AddLog(
				logger.LevelError,
				logger.SourceMain,
				map[string]string{},
				"engagement http server error: "+err.Error(),
			)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		lg.AddLog(
			logger.LevelError,
			logger.SourceMain,
			map[string]string{},
			"engagement http shutdown error: "+err.Error(),
		)
	}

	lg.Close()
}
