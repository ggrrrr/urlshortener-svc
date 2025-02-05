package web

import (
	"log/slog"
	"net/http"
)

func (l *Listener) httpHandlerCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)
		if l.cfg.CORSHosts == "" {
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", l.cfg.CORSHosts)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Authorization")

		if r.Method == http.MethodOptions {
			out := "."
			w.WriteHeader(200)
			_, err := w.Write([]byte(out))
			if err != nil {
				slog.ErrorContext(r.Context(), "cant write body")
			}
			return
		}
	})
}
