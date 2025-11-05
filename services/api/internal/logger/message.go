package logger

import "time"

type Message struct {
	Timestamp time.Time
	Level     Level
	Source    Source            // who generate log: handler/service/repo
	Fields    map[string]string // key-value (userID, postID, etc)
	Message   string
}
