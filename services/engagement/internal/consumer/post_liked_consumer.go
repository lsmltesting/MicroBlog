package consumer

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/segmentio/kafka-go"

	"github.com/lsmltesting/MicroBlog/services/api/internal/events"
	"github.com/lsmltesting/MicroBlog/services/api/internal/logger"
	"github.com/lsmltesting/MicroBlog/services/engagement/internal/repository"
)

type PostLikedConsumer struct {
	reader kafka.Reader
	repo   repository.StatsRepository
	lg     logger.Logger
}

func NewPostLikedConsumer(cfg events.Config, repo repository.StatsRepository, lg logger.Logger) *PostLikedConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: "engagement-post-liked",
		Topic:   cfg.TopicPostLiked,
	})

	return &PostLikedConsumer{
		reader: *reader,
		repo:   repo,
		lg:     lg,
	}
}

func (c *PostLikedConsumer) Run(ctx context.Context) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			c.lg.AddLog(
				logger.LevelError,
				logger.SourceConsumer,
				map[string]string{},
				err.Error(),
			)
			continue
		}

		var env events.Envelope
		if err := json.Unmarshal(m.Value, &env); err != nil {
			c.lg.AddLog(
				logger.LevelError,
				logger.SourceConsumer,
				map[string]string{},
				"failed to unmarshal envelope: "+err.Error(),
			)
			continue
		}

		payloadBytes, err := json.Marshal(env.Payload)
		if err != nil {
			c.lg.AddLog(
				logger.LevelError,
				logger.SourceConsumer,
				map[string]string{},
				"failed to re-marshal payload: "+err.Error(),
			)
			continue
		}

		var payload events.PostLikedPayloadV1
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			c.lg.AddLog(
				logger.LevelError,
				logger.SourceConsumer,
				map[string]string{},
				"failed to unmarshal PostLiked payload: "+err.Error(),
			)
			continue
		}

		if err := c.repo.IncrementPostLikes(ctx, payload.PostID, env.ID); err != nil {
			c.lg.AddLog(
				logger.LevelError,
				logger.SourceConsumer,
				map[string]string{
					"post_id": strconv.Itoa(payload.PostID),
					"user_id": strconv.Itoa(payload.UserID),
				},
				"failed to increment like counter: "+err.Error(),
			)
		}
	}
}

func (c *PostLikedConsumer) Close() error {
	return c.reader.Close()
}
