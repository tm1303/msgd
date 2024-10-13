package receiver

import (
	"encoding/json"
	"net/http"
)

type EnqueueRequest struct {
	Message string `json:"message"`
}

type EnqueueResponse struct {
	MessageID string `json:"message_id"`
}

type MsgQueuer interface {
	Enqueue(messageBody string) (messageID *string, err error)
}

func GetHandler(msgQueuer MsgQueuer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		p := EnqueueRequest{}
		err := json.NewDecoder(r.Body).Decode(&p)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		id, err := msgQueuer.Enqueue(p.Message)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // log
			return
		}

		if id==nil{
			w.WriteHeader(http.StatusInternalServerError) // log
			return
		}

		resp := EnqueueResponse{
			MessageID: *id,
		}

		if err = json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError) // log
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}
