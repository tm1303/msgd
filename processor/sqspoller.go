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

func (r sqsPoller) poll(ctx context.Context, action pollerAction, requestAttributes []string) int64 {

	pollParams := &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: r.batchSize,
		QueueUrl:            r.queue,
		MessageAttributeNames: aws.StringSlice(requestAttributes),
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
		a := convertMessageAttributes(v.MessageAttributes)
		if action(b, a) { // if the message has been processed delete it from the queue
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

func convertMessageAttributes(attrs map[string]*sqs.MessageAttributeValue) map[string]interface{} {
    result := make(map[string]interface{})
    for key, attr := range attrs {
        switch *attr.DataType {
        case "String":
            result[key] = *attr.StringValue
        case "Number":
            result[key] = *attr.StringValue // bit annoying :/
        case "Binary":
            result[key] = attr.BinaryValue
        default:
            result[key] = nil // or handle unknown types
        }
    }
    return result
}