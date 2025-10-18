package queue

import (
	"context"
	"log"
)

func (l *likeQueueImplement) startWorkers(ctx context.Context) {
	for i := 0; i < l.workers; i++ {
		go l.worker(ctx, i)
	}
}

func (l *likeQueueImplement) worker(ctx context.Context, workerID int) {
	log.Println("Starting handle queue with workerID:", workerID)
	for {
		select {
		case likeForChan, ok := <-l.qLike:
			if !ok {
				log.Printf("Worker: %d, channel is closed.", workerID)
				return
			}

			l.handleQueue(ctx, likeForChan)

		case <-l.stop:
			log.Printf("Catch stop signal with worker: %d", workerID)
			return

		case <-ctx.Done():
			log.Printf("Context done for worker: %d, stopping", workerID)

		}

	}
}

func (l *likeQueueImplement) handleQueue(ctx context.Context, likeForChan LikeForChan) (int, error) {
	likeID, err := l.likeService.CreateLike(ctx, likeForChan.UserID, likeForChan.PostID)
	if err != nil {
		return 0, err
	}
	return likeID, err
}
