package stream

import (
	"context"
	"slices"

	"github.com/redis/go-redis/v9"
)

type handler func([]redis.XStream) error
type values map[string]interface{}

type event interface {
	StreamName() string
	Values() values
}

func Emit(rdb *redis.Client, event event) error {
	_, err := rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream: event.StreamName(),
		Values: event.Values(),
	}).Result()

	if err != nil {
		return err
	}

	return nil
}

func Handle(rdb *redis.Client, handler handler, events ...event) error {
	ids := slices.Repeat([]string{"$"}, len(events))
	streams := []string{}

	for _, event := range events {
		streams = append(streams, event.StreamName())
	}

	res, err := rdb.XRead(context.Background(), &redis.XReadArgs{
		Streams: append(streams, ids...),
		Block:   0,
	}).Result()

	if err != nil {
		return err
	}

	err = handler(res)
	if err != nil {
		return err
	}

	return nil
}
