package logger

import (
	"log"
	"time"

	customErrors "github.com/lsmltesting/MicroBlog/internal/errors"
)

type LoggerConfig struct {
	Workers    int
	BufferSize int
}

type Logger interface {
	Close()
	AddLog(
		level Level,
		source Source,
		fields map[string]string,
		msg string,
	) error
}

type logger struct {
	messages chan Message
	stop     chan struct{}
	workers  int
	writer   Writer
}

func NewLogger(config LoggerConfig) Logger {
	logger := &logger{
		messages: make(chan Message, config.BufferSize),
		stop:     make(chan struct{}),
		workers:  config.Workers,
		writer:   &JSONWriter{},
	}
	logger.startWorkers()
	return logger
}

func (l *logger) AddLog(
	level Level,
	source Source,
	fields map[string]string,
	msg string,
) error {
	message := Message{
		Timestamp: time.Now(),
		Level:     level,
		Source:    source,
		Message:   msg,
	}

	select {
	case l.messages <- message:
		return nil
	case <-l.stop:
		return customErrors.ErrLoggerChanClosed
	}
}

func (l *logger) Close() {
	close(l.messages)
	close(l.stop)
}

func (l *logger) startWorkers() {
	for i := 0; i < l.workers; i++ {
		go l.worker(i)
	}
}

func (l *logger) worker(workerID int) {
	log.Printf("[channel for logging] Starting handle with workerID:%d", workerID)

	for {
		log.Printf("[channel for logging] Worker %d WAINTING for task...", workerID)

		select {
		case msg, ok := <-l.messages:
			if !ok {
				log.Printf("[channel for logging] Worker: %d, channel is closed", workerID)
				return
			}

			log.Printf("[channel for logging] Worker %d PROCESSING task", workerID)
			l.writer.Write(msg)
		case <-l.stop:
			log.Printf("[channel for logging] Catch stop signal with worker: %d", workerID)
			return
		}
	}
}
