package receiver

import (
	"encoding/json"
	"msgd/infra"
	"net/http"
)

type EnqueueRequest struct {
	Message string `json:"message"`
	// UserID string `json:"user_id"`
}

type EnqueueResponse struct {
	MessageID string `json:"message_id"`
}

type MsgQueuer interface {
	Enqueue(messageBody string, userID string) (messageID *string, err error)
}

func GetHandler(msgQueuer MsgQueuer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		p := EnqueueRequest{}
		err := json.NewDecoder(r.Body).Decode(&p)

		if err != nil {
			http.Error(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		userID, ok := infra.UserIDFrom(r.Context())
		if !ok {
			http.Error(w, "missing userid", http.StatusBadRequest)
			return
		}

		id, err := msgQueuer.Enqueue(p.Message, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // log
			return
		}

		if id == nil {
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
