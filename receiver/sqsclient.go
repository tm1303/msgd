package receiver

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type sqsClient struct {
	sqs   *sqs.SQS
	queue *string
}

func NewSqsClient(awsSession *session.Session, queue string) MsgClient {
	return sqsClient{
		sqs:   sqs.New(awsSession),
		queue: aws.String(queue),
	}
}

func (r sqsClient) Enqueue(messageBody string) *string {

	sendParams := &sqs.SendMessageInput{
		MessageBody: aws.String(messageBody),
		QueueUrl:    r.queue,
	}

	result, err := r.sqs.SendMessage(sendParams)
	if err != nil {
		fmt.Printf("Failed to send message: %s\n", err)
		os.Exit(1) // TODO
	}

	fmt.Printf("Message sent successfully, Message ID: %s\n", *result.MessageId)
	return result.MessageId
}
