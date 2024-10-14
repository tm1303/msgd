package broadcaster

import (
	"context"
	"encoding/json"
	"fmt"
	"msgd/domain"

	"github.com/gorilla/websocket"
)

func StartBroadcaster(ctx context.Context, broadcast chan domain.MessageBody) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down processor...")
			return
		case m := <-broadcast:
			fmt.Println("Broadcast...")

			json, err := json.Marshal(m)
			if err != nil {
				fmt.Println("Error marshling for broadcast...")
			}
			for c := range connections {
				if err := c.WriteMessage(websocket.TextMessage, json); err != nil {
					fmt.Println("error writing message: %w", err)
					c.Close()
					delete(connections, c)
				}
			}
		}
	}
}
