package events

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type KafkaConfig struct {
	Brokers             []string
	TopicUserRegistered string
	TopicPostCreated    string
	TopicPostLiked      string
}

func FromEnv() KafkaConfig {
	wd, _ := os.Getwd()
	envPath := filepath.Join(wd, "..", ".env")
	_ = godotenv.Load(envPath)

	return KafkaConfig{
		Brokers: []string{
			os.Getenv("KAFKA_BROKER"),
		},
		TopicUserRegistered: getenvDefault("KAFKA_TOPIC_USER_REGISTERED", "user-registered"),
		TopicPostCreated:    getenvDefault("KAFKA_TOPIC_POST_CREATED", "post-created"),
		TopicPostLiked:      getenvDefault("KAFKA_TOPIC_POST_LIKED", "post-liked"),
	}
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
