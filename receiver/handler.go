package receiver

import (
	"encoding/json"
	"net/http"
)

type EnqueuePayload struct {
	Message string `json:"message"`
}

type MsgClient interface {
	Enqueue(messageBody string) *string
}

func GetHandler(msgClient MsgClient) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		p := EnqueuePayload{}
		err := json.NewDecoder(r.Body).Decode(&p)

		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		msgClient.Enqueue(p.Message)
	}
}
