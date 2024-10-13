package receiver

import (
	"encoding/json"
	"fmt"
	"msgd/domain"
	"time"

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

func (r sqsQueuer) Enqueue(messageBody string, userID string) (*string, error) {
	body := domain.MessageBody{
        Message: messageBody,
        Date:   time.Now(),
		UserID: userID,
    }

    bodyJSON, err := json.Marshal(body)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal message body: %w", err)
    }

    messageAttributes := map[string]*sqs.MessageAttributeValue{
        "UserID": {
            DataType:    aws.String("String"),
            StringValue: aws.String(userID),
        },
    }

    sendParams := &sqs.SendMessageInput{
        MessageBody:       aws.String(string(bodyJSON)),
        QueueUrl:          r.queue,
        MessageAttributes: messageAttributes,
    }
	result, err := r.sqs.SendMessage(sendParams)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return result.MessageId, nil
}
