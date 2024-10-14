package broadcaster

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connections = make(map[*websocket.Conn]bool)
var keepalive []byte = []byte(`{"keep_alive":true}`)

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error during connection upgrade: %w", err)
		return
	}

	defer func() {
		conn.Close()
		delete(connections, conn)
	}()

	connections[conn] = true
	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down broadcaster...")
			return
		default:
			if err := conn.WriteMessage(websocket.TextMessage, keepalive); err != nil {
				fmt.Println("error writing message: %w", err)
				conn.Close()
				delete(connections, conn)
			}
			time.Sleep(5 * time.Second)
		}
	}
}
