package api

import (
	"log/slog"
	"net/http"

	"github.com/ggrrrr/urlshortener-svc/be/common/web"
	"github.com/ggrrrr/urlshortener-svc/be/login/models"
)

func (s server) handleLogin(w http.ResponseWriter, r *http.Request) {
	req := models.UserPasswordRequest{}
	err := web.UnmarshalJSON(r, &req)
	if err != nil {
		slog.Error("handleLogin", slog.Any("error", err))
		web.SendJSONError(w, err)
		return
	}

	token, err := s.app.Login(r.Context(), req)
	if err != nil {
		slog.Error("handleLogin.Login", slog.Any("error", err))
		web.SendJSONError(w, err)
		return
	}

	web.SendJSONPayload(w, models.UserPasswordResponse{Token: token})
}
