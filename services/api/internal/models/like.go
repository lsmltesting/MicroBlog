package models

import "time"

type Like struct {
	CreatedAt time.Time
	ID        int
	UserID    int
	PostID    int
}

// TODO: add new param createdat
func NewLike(userID int, postID int) *Like {
	return &Like{
		CreatedAt: time.Now(),
		UserID:    userID,
		PostID:    postID,
	}
}
