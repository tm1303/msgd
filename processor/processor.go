package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"msgd/domain"
	"time"
)

type pollerAction func(message *string, attributes map[string]interface{}) bool

type MsgPoller interface {
	poll(ctx context.Context, processorAction pollerAction, requestAttributes []string) (count int64)
}

func StartProcessor(ctx context.Context, poller MsgPoller, broadcastChan chan domain.MessageBody) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down processor...")
			return
		default:
			count := poller.poll(ctx,
				processMessageToQueue(broadcastChan),
				[]string{domain.UserIDAttributeName},
			)
			fmt.Printf("%d messages processed\n", count)
			if count == 0 {
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func processMessageToQueue(broadcastChan chan domain.MessageBody) pollerAction {

	return func(body *string, attributes map[string]interface{}) bool {
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

		var message domain.MessageBody

		err := json.Unmarshal([]byte(*body), &message)
		if err != nil {
			fmt.Printf("failed to unmarshal message body: %s", err)
			return false
		}

		broadcastChan <- message
		return true
	}
}
