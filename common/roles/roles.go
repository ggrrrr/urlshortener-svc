package roles

import (
	"context"

	"github.com/ggrrrr/urlshortener-svc/common/application"
)

type (
	ctxKeyT struct{}

	AuthenticatedUser struct {
		Username string
	}
)

var ctxKey ctxKeyT = ctxKeyT{}

func ExtractUser(ctx context.Context) (AuthenticatedUser, error) {
	v := ctx.Value(ctxKey)
	if v == nil {
		application.NewUnauthorized()
	}
	user, ok := v.(AuthenticatedUser)
	if ok {
		return user, nil
	}
	return AuthenticatedUser{}, application.NewSystemError("unable to cast user", nil)

}

func InjectUser(ctx context.Context, user AuthenticatedUser) context.Context {
	return context.WithValue(ctx, ctxKey, user)
}
