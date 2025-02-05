package web

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/ggrrrr/urlshortener-svc/be/common/roles"
)

func (l *Listener) httpHandlerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			next.ServeHTTP(w, r)
		}()

		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			return
		}

		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) != 2 {
			slog.ErrorContext(r.Context(), "httpHandlerAuth.NoHeader.Bearer")
			return
		}
		slog.InfoContext(r.Context(), "httpHandlerAuth", slog.Any("Authorization", splitToken[1]))

		authInfo, err := l.verifier.Verify(r.Context(), splitToken[1])
		if err != nil {
			slog.ErrorContext(r.Context(), "httpHandlerAuth", slog.Any("error", err))
			return
		}

		slog.InfoContext(r.Context(), "httpHandlerAuth", slog.Any("authInfo", authInfo))
		ctx := roles.InjectUser(r.Context(), authInfo)
		r = r.WithContext(ctx)
	})
}
