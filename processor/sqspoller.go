package processor

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type sqsPoller struct {
	sqs       *sqs.SQS
	queue     *string
	batchSize *int64
}

func NewSqsPoller(awsSession *session.Session, queue string, batchSize int64) MsgPoller {
	return sqsPoller{
		sqs:       sqs.New(awsSession),
		queue:     aws.String(queue),
		batchSize: aws.Int64(batchSize),
	}
}

func (r sqsPoller) Poll(ctx context.Context, action func(message *string) bool) int64 {

	pollParams := &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: r.batchSize,
		QueueUrl:            r.queue,
	}

	result, err := r.sqs.ReceiveMessage(pollParams)
	if err != nil {
		panic(fmt.Errorf("failed to read message queue: %w", err)) 
	}

	if len(result.Messages) == 0 {
		return 0
	}

	msgCount := int64(0)
	for _, v := range result.Messages {
		b := v.Body
		if action(b) { // if the message has been procssed delete it from the queue
			deleteParams := &sqs.DeleteMessageInput{
				QueueUrl:      r.queue,
				ReceiptHandle: v.ReceiptHandle,
			}
			_, err = r.sqs.DeleteMessage(deleteParams)
			if err != nil {
				panic(fmt.Errorf("failed to delete from queue: %w", err)) // unknown state
			}
			msgCount++
		}
	}

	return int64(msgCount)
}
