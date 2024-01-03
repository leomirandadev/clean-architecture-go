package queue

import "context"

type QueueDoer interface {
	Subscribe(queue string, doer func(ctx context.Context, message []byte) error)
	Publish(ctx context.Context, queue string, message any) error
}
