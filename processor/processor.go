package processor

import (
	"context"
	"fmt"
	"time"
)

type MsgPoller interface {
	Poll(ctx context.Context, processorAction func(Message *string) bool) (count int64)
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
	fmt.Println("start messages process")
	count := poller.Poll(ctx, processMessage)
	fmt.Printf("%d messages processed\n", count)

	return nil // Return nil for successful processing; add error handling if needed
}

func processMessage(Message *string) bool {
	if _, err := fmt.Printf("Message dequeued: %s\n", *Message); err != nil {
		return false
	}
	return true
}
