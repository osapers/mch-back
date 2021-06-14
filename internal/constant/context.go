// Package constant contains common variables
package constant

import (
	"context"
)

var (
	// UserIDCtxKey is a key name of customerID in ctx
	UserIDCtxKey = &contextKey{"userID"}
	// ErrorCtxKey is a key name of auth error context in ctx
	ErrorCtxKey = &contextKey{"error"}
	// DBSecretCtxKey is a key name of secret for db
	DBSecretCtxKey = &contextKey{"dbSecret"}
)

type contextKey struct {
	name string
}

func (c contextKey) String() string {
	return c.name
}

// GetUserIDFromCtx returns userID from context
func GetUserIDFromCtx(ctx context.Context) string {
	val := ctx.Value(UserIDCtxKey)
	if val == nil {
		return ""
	}

	return val.(string)
}
