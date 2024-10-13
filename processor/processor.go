package processor

import (
	"context"
	"fmt"
	"msgd/domain"
	"time"
)

type pollerAction func(message *string, attributes map[string]interface{}) bool

type MsgPoller interface {
	poll(ctx context.Context, processorAction pollerAction, requestAttributes []string) (count int64)
}

func StartProcessor(ctx context.Context, poller MsgPoller) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down processor...")
			return
		default:
			if err := pollMessages(ctx, poller); err != nil {
				fmt.Printf("Error during polling: %v\n", err)
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func pollMessages(ctx context.Context, poller MsgPoller) error {
	fmt.Println("start message polling")
	count := poller.poll(ctx, processMessage, []string{domain.UserIDAttributeName})
	fmt.Printf("%d messages processed\n", count)

	return nil // Return nil for successful processing; add error handling if needed
}

func processMessage(body *string, attributes map[string]interface{}) bool {
	//TODO: do something fun with our message!
	if body == nil || attributes[domain.UserIDAttributeName] == nil {
		return false // log
	}
	userID, ok := attributes[domain.UserIDAttributeName].(string)
	if !ok {
		return false // log
	}
	if _, err := fmt.Printf("Message from user `%s` dequeued: %s\n", userID, *body); err != nil {
		return false // log
	}
	return true
}
