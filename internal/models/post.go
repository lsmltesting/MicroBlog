package models

import (
	"time"

	"github.com/lsmltesting/MicroBlog/internal/errors"
)

type Post struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Text      string
	UserID    int
	ID        int
}

func NewPost(userID int, text string) (*Post, error) {
	post := &Post{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
	}

	// Set text for post after validating
	if err := post.SetText(text); err != nil {
		return nil, err
	}

	return post, nil
}

func (post *Post) SetText(text string) error {
	if text == "" {
		return errors.ErrEmptyPostText
	}

	post.Text = text
	post.UpdatedAt = time.Now()
	return nil
}

// func (post *Post) SetLike(userID int, likeID int) error {

// 	// Check if like from userId is already created
// 	if _, ok := post.HistoryLikes[userID]; !ok {
// 		return errors.ErrPostLikeAlreadyCreated
// 	}

// 	post.HistoryLikes[userID] = struct{}{}
// 	return nil
// }
