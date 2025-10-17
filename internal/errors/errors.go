package errors

import "errors"

var (
	ErrWrongUserName     = errors.New("wrong user name")
	ErrWrongEmailAddress = errors.New("wrong email address")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrNotFindUser       = errors.New("not find user")

	ErrEmptyPostText          = errors.New("empty text for post")
	ErrPostLikeAlreadyCreated = errors.New("like is already created")
	ErrNotFindPost            = errors.New("not find post")
	ErrNotAnyPostExists       = errors.New("no one post exists")

	ErrNotFindLike      = errors.New("like not found")
	ErrNotAnyLikeExists = errors.New("no one like exists")

	ErrQueueClosed = errors.New("queue is closed")

	ErrLoggerChanClosed = errors.New("channel message of logs is closed")

	ErrFailedBuildQueryForPSQL = errors.New("build insert query is failed")
	ErrFailedCreateUserInPSQL  = errors.New("create new user is failed")
	ErrFailedCreatePostInPSQL  = errors.New("create new post is failed")
)
