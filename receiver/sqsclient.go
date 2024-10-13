package receiver

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type sqsQueuer struct {
	sqs   *sqs.SQS
	queue *string
}

func NewSqsQueuer(awsSession *session.Session, queue string) MsgQueuer {
	return sqsQueuer{
		sqs:   sqs.New(awsSession),
		queue: aws.String(queue),
	}
}

func (r sqsQueuer) Enqueue(messageBody string) (*string, error) {

	sendParams := &sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    r.queue,
	}

	result, err := r.sqs.SendMessage(sendParams)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return result.MessageId, nil
}
