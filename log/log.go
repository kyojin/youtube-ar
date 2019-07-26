package log

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"
)

// Logger is a logger.
type Logger interface {
	Log(ctx context.Context, msg string, fields ...Field)
}

// New returns a new logger.
func New() Logger {
	return &logger{logger: log.New(os.Stdout, "", 0)}
}

type logger struct {
	logger *log.Logger
}

func (l *logger) Log(ctx context.Context, msg string, fields ...Field) {
	event := map[string]interface{}{
		"msg":       msg,
		"timestamp": time.Now(),
	}
	for _, field := range fields {
		event[field.key()] = field.value()
	}
	b, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	l.logger.Printf("%s\n", b)
}
