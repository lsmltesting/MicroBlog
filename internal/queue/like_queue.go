package queue

import (
	"context"
	"time"

	customErrors "github.com/lsmltesting/MicroBlog/internal/errors"
	"github.com/lsmltesting/MicroBlog/internal/service/like"
)

type LikeQueue interface {
	AddLike(userID int, postID int) error
	Close()
}

type likeQueueImplement struct {
	qLike       chan LikeForChan
	stop        chan struct{}
	workers     int
	likeService like.LikeService
}

// struct for creating Like model from channel
type LikeForChan struct {
	CreatedAt time.Time
	UserID    int
	PostID    int
}

// config for q
type LikeQueueConfig struct {
	BufferSize int
	Workers    int
}

func NewLikeQueue(ctx context.Context, config LikeQueueConfig, likeService like.LikeService) LikeQueue {
	q := &likeQueueImplement{
		qLike:       make(chan LikeForChan, config.BufferSize),
		workers:     config.Workers,
		stop:        make(chan struct{}),
		likeService: likeService,
	}
	q.startWorkers(ctx)

	return q
}

func (l *likeQueueImplement) AddLike(userID int, postID int) error {
	likeForChan := LikeForChan{
		CreatedAt: time.Now(),
		UserID:    userID,
		PostID:    postID,
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
