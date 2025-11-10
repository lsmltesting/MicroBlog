package events

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishUserRegistered(ctx context.Context, env Envelope, payload UserCreatedPayload) error
	PublishPostCreated(ctx context.Context, env Envelope, payload PostCreatedPayload) error
	PublishPostLiked(ctx context.Context, env Envelope, payload PostLikedPayloadV1) error
	Close() error
}

type kafkaProducer struct {
	cfg     Config
	writers map[string]*kafka.Writer
}

func NewKafkaProducer(cfg Config) Producer {
	writers := make(map[string]*kafka.Writer)

	mkWriter := func(topic string) *kafka.Writer {
		return &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		}
	}

	writers[cfg.TopicUserRegistered] = mkWriter(cfg.TopicUserRegistered)
	writers[cfg.TopicPostCreated] = mkWriter(cfg.TopicPostCreated)
	writers[cfg.TopicPostLiked] = mkWriter(cfg.TopicPostLiked)

	return &kafkaProducer{
		cfg:     cfg,
		writers: writers,
	}
}

func (p *kafkaProducer) Close() error {
	var lastErr error
	for _, w := range p.writers {
		if err := w.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (p *kafkaProducer) publish(ctx context.Context, topic string, env Envelope, payload any) error {
	env.Payload = payload

	data, err := json.Marshal(env)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	w := p.writers[topic]

	var lastErr error
	backoff := 100 * time.Millisecond

	for i := 0; i < 3; i++ {
		err = w.WriteMessages(ctx, kafka.Message{
			Key:   []byte(env.ID),
			Value: data,
		})
		if err == nil {
			return nil
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		lastErr = err
		time.Sleep(backoff)
		backoff *= 2
	}

	return lastErr
}

func (p *kafkaProducer) PublishUserRegistered(ctx context.Context, env Envelope, payload UserCreatedPayload) error {
	env.Type = EventTypeUserCreated
	env.Version = "v1"
	return p.publish(ctx, p.cfg.TopicUserRegistered, env, payload)
}

func (p *kafkaProducer) PublishPostCreated(ctx context.Context, env Envelope, payload PostCreatedPayload) error {
	env.Type = EventTypePostCreated
	env.Version = "v1"
	return p.publish(ctx, p.cfg.TopicPostCreated, env, payload)
}

func (p *kafkaProducer) PublishPostLiked(ctx context.Context, env Envelope, payload PostLikedPayloadV1) error {
	env.Type = EventTypePostLiked
	env.Version = "v1"
	return p.publish(ctx, p.cfg.TopicPostLiked, env, payload)
}
