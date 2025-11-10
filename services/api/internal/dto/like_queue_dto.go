package dto

import "time"

type LikeQueuedResponse struct {
	Message   string    `json:"message"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	Timestamp time.Time `json:"timestamp"`
}
