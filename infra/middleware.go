package infra

import (
	"net/http"
)

const userHeader string = "X-User-ID"

func UserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userIDVal := r.Header.Get(userHeader)

		if userIDVal != "" {
			r = r.WithContext(ContextWithUserID(r.Context(), userIDVal))
		}

		next.ServeHTTP(w, r)
	})
}
