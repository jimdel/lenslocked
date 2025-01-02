package context

import (
	"context"

	"github.com/jimdel/lenslocked/models"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	value := ctx.Value(userKey)
	user, ok := value.(*models.User)
	if !ok {
		// Most likely occurs if no data was stored in the ctx
		// so it does not have the type *models.User
		return nil
	}
	return user
}
