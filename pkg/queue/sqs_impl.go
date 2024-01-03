package queue

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

/*
NewSQS creates a instance that provides a way to publish and subscribe on sqs queues.
"key" and "secret" are the credentials while "unqueueTimeoutSeconds" is the maximum time (only needed to subscriptions).
that one message will take to be unqueued.
*/
func NewSQS(key, secret, region, endpoint string, unqueueTimeoutSeconds int32) QueueDoer {

	options := sqs.Options{
		Credentials: credentials.NewStaticCredentialsProvider(key, secret, ""),
		Region:      region,
	}

	if endpoint != "" {
		options.BaseEndpoint = &endpoint
	}

	return &sqsImpl{
		client:                sqs.New(options),
		unqueueTimeoutSeconds: unqueueTimeoutSeconds,
		maxParallelMessages:   10,
	}
}

type sqsImpl struct {
	client                *sqs.Client
	unqueueTimeoutSeconds int32
	maxParallelMessages   int32
}

func (s sqsImpl) Subscribe(queue string, doer func(ctx context.Context, message []byte) error) {
	ctx := context.Background()

	queueUrlResult, err := s.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &queue,
	})
	if err != nil {
		slog.Error("sqs: get queue url fails", "queue", queue, "err", err)
		return
	}

	gMInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            queueUrlResult.QueueUrl,
		MaxNumberOfMessages: s.maxParallelMessages,
		VisibilityTimeout:   s.unqueueTimeoutSeconds,
	}

	for {
		messagesResult, err := s.client.ReceiveMessage(ctx, gMInput)
		if err != nil {
			slog.Error("sqs: receive message fails", "queue", queue, "err", err)
		}

		for _, message := range messagesResult.Messages {
			go func(msg types.Message) {
				slog.Info("processing new message")

				var msgByte []byte
				if msg.Body != nil {
					msgByte = []byte(*msg.Body)
				}

				if err := doer(ctx, msgByte); err != nil {
					slog.Error("sqs: receive message fails", "queue", queue, "err", err)
				} else {
					// remove message from the queue
					s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
						QueueUrl:      queueUrlResult.QueueUrl,
						ReceiptHandle: msg.ReceiptHandle,
					})
				}

			}(message)
		}

	}
}

func (s sqsImpl) Publish(ctx context.Context, queue string, data any) error {
	parsedData, err := json.Marshal(data)
	if err != nil {
		slog.Warn("sqs: marshal fails", "queue", queue, "err", err)
		return err
	}

	dataStr := string(parsedData)

	queueUrlResult, err := s.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &queue,
	})
	if err != nil {
		slog.Error("sqs: get queue url fails", "queue", queue, "err", err)
		return err
	}

	message := &sqs.SendMessageInput{
		MessageBody: &(dataStr),
		QueueUrl:    queueUrlResult.QueueUrl,
	}

	if _, err = s.client.SendMessage(ctx, message); err != nil {
		slog.Warn("sqs: send message fails", "queue", queue, "err", err)
		return err
	}

	return nil
}
