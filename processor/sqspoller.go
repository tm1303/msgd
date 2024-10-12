package processor

import (
	"fmt"
	"os"

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

func (r sqsPoller) Poll() *string {

	pollParams := &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: r.batchSize,
		QueueUrl:            r.queue,
	}

	result, err := r.sqs.ReceiveMessage(pollParams)
	if err != nil {
		fmt.Printf("Failed to receive message: %s\n", err)
		os.Exit(1) // TODO
	}

	fmt.Printf("Message sent successfully, Message[0]: %s\n", *result.Messages[0].Body)
	return result.Messages[0].ReceiptHandle
}
