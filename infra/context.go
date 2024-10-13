package infra

import (
	"context"
)

type contextKey string

const userContextKey contextKey = "UserID"

func ContextWithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userContextKey, id)
}

func UserIDFrom(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userContextKey).(string)
	return id, ok
}
