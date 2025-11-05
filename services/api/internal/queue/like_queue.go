package queue

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	customErrors "github.com/lsmltesting/MicroBlog/services/api/internal/errors"
	"github.com/lsmltesting/MicroBlog/services/api/internal/events"
	"github.com/lsmltesting/MicroBlog/services/api/internal/logger"
	"github.com/lsmltesting/MicroBlog/services/api/internal/service/like"
	"github.com/lsmltesting/MicroBlog/services/api/internal/trace"
)

type LikeQueue interface {
	AddLike(ctx context.Context, userID int, postID int) error
	Close()
}

type likeQueueImplement struct {
	qLike       chan LikeForChan
	stop        chan struct{}
	workers     int
	likeService like.LikeService
	producer    events.Producer
	lg          logger.Logger
}

// struct for creating Like model from channel
type LikeForChan struct {
	CreatedAt time.Time
	UserID    int
	PostID    int
	TraceID   string
}

// config for q
type LikeQueueConfig struct {
	BufferSize int
	Workers    int
}

func NewLikeQueue(
	ctx context.Context,
	config LikeQueueConfig,
	likeService like.LikeService,
	producer events.Producer,
	lg logger.Logger,
) LikeQueue {
	q := &likeQueueImplement{
		qLike:       make(chan LikeForChan, config.BufferSize),
		workers:     config.Workers,
		stop:        make(chan struct{}),
		likeService: likeService,
		producer:    producer,
		lg:          lg,
	}
	q.startWorkers(ctx)
	return q
}

func (l *likeQueueImplement) AddLike(ctx context.Context, userID int, postID int) error {
	likeForChan := LikeForChan{
		CreatedAt: time.Now(),
		UserID:    userID,
		PostID:    postID,
		TraceID:   trace.IDFromContext(ctx),
	}

	select {
	case l.qLike <- likeForChan:
		return nil
	case <-l.stop:
		return customErrors.ErrQueueClosed
	}
}

func (l *likeQueueImplement) Close() {
	close(l.qLike)
	close(l.stop)
}

func (l *likeQueueImplement) startWorkers(ctx context.Context) {
	for i := 0; i < l.workers; i++ {
		workerID := i

		go func() {
			for {
				select {

				case <-ctx.Done():
					return

				case msg, ok := <-l.qLike:
					if !ok {
						return
					}

					likeCtx := trace.WithTraceID(ctx, msg.TraceID)

					likeID, err := l.likeService.Create(likeCtx, msg.UserID, msg.PostID)
					if err != nil {
						l.lg.AddLog(
							logger.LevelError,
							logger.SourceDB,
							map[string]string{
								"worker":  strconv.Itoa(workerID),
								"user_id": strconv.Itoa(msg.UserID),
								"post_id": strconv.Itoa(msg.PostID),
							},
							err.Error(),
						)
						continue
					}

					env := events.Envelope{
						ID:         uuid.NewString(),
						OccurredAt: msg.CreatedAt,
						TraceID:    msg.TraceID,
					}

					if err := l.producer.PublishPostLiked(likeCtx, env, events.PostLikedPayloadV1{
						PostID: msg.PostID,
						UserID: msg.UserID,
					}); err != nil {
						l.lg.AddLog(
							logger.LevelError,
							logger.SourceQueue,
							map[string]string{
								"worker":  strconv.Itoa(workerID),
								"user_id": strconv.Itoa(msg.UserID),
								"post_id": strconv.Itoa(msg.PostID),
								"like_id": strconv.Itoa(likeID),
							},
							err.Error(),
						)
					}
				}
			}
		}()
	}
}
