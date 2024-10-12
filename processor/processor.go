package processor

import (
	"net/http"
	"time"
)

type MsgPoller interface {
	Poll() *string
}

func StartProcessor(poller MsgPoller) func(w http.ResponseWriter, r *http.Request) {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop() // Stop the ticker when the function returns

	for {
		select {
		case <-ticker.C:
			poller.Poll()
		}
	}
}
