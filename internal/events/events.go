package events

import (
	"os"
	"strings"
	"time"
)

type EventType string

const (
	EventTypeUserCreated EventType = "UserCreated"
	EventTypePostCreated EventType = "PostCreated"
	EventTypePostLiked   EventType = "PostLiked"
)

type Envelope struct {
	ID         string    `json:"id"`
	Type       EventType `json:"type"`
	Version    string    `json:"version"`
	OccurredAt time.Time `json:"occurred_at"`
	Payload    any       `json:"payload"`
	TraceID    string    `json:"trace_id,omitempty"`
}

// Payloads
type UserCreatedPayload struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type PostCreatedPayload struct {
	PostID int    `json:"post_id"`
	UserID int    `json:"user_id"`
	Text   string `json:"text"`
}

type PostLikedPayloadV1 struct {
	PostID int `json:"post_id"`
	UserID int `json:"user_id"`
}

// Config for producer
type Config struct {
	Brokers             []string
	TopicUserRegistered string
	TopicPostCreated    string
	TopicPostLiked      string
}

func ConfigFromEnv() Config {
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	var brokers []string
	if brokersEnv == "" {
		brokers = []string{"kafka:9092"}
	} else {
		brokers = strings.Split(brokersEnv, ",")
	}

	return Config{
		Brokers:             brokers,
		TopicUserRegistered: getenv("KAFKA_TOPIC_USER_REGISTERED", "user-registered"),
		TopicPostCreated:    getenv("KAFKA_TOPIC_POST_CREATED", "post-created"),
		TopicPostLiked:      getenv("KAFKA_TOPIC_POST_LIKED", "post-liked"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
