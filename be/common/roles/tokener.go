package roles

import "context"

type (
	TokenGenerator interface {
		Generate(ctx context.Context, username string) (string, error)
	}

	TokenVerifier interface {
		Verify(ctx context.Context, token string) (AuthenticatedUser, error)
	}
)
