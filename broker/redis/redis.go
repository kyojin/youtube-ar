package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/yansal/youtube-ar/log"
)

// New returns a new redis client.
func New(logger log.Logger) (*redis.Client, error) {
	redisURL := "redis://"
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)

	client.WrapProcess(func(old func(redis.Cmder) error) func(redis.Cmder) error {
		return func(cmd redis.Cmder) error {
			start := time.Now()
			err := old(cmd)

			fields := []log.Field{
				log.String("command", cmd.Name()),
				log.Stringer("duration", time.Since(start)),
			}

			var msg string
			if stringer, ok := cmd.(fmt.Stringer); ok {
				msg = stringer.String()
			} else {
				msg = cmderString(cmd)
			}
			// TODO: get context from cmd
			logger.Log(context.Background(), "redis: "+msg, fields...)

			return err
		}
	})

	return client, nil
}

func cmderString(cmd redis.Cmder) string {
	var ss []string
	for _, arg := range cmd.Args() {
		ss = append(ss, fmt.Sprint(arg))
	}
	s := strings.Join(ss, " ")
	if err := cmd.Err(); err != nil {
		return s + ": " + err.Error()
	}
	return s
}
