package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/urlshortener-svc/be/common/web"
	"github.com/ggrrrr/urlshortener-svc/be/shorturl/models"
)

func (s server) handleForward(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	toURL, err := s.app.GetLongURL(r.Context(), key)
	if err != nil {
		slog.Error("handleForward.GetLongURL", slog.Any("error", err))
		web.SendJSONError(w, err)
		return
	}

	web.Redirect(w, r, toURL)
}

func (s server) handleCreateShortURL(w http.ResponseWriter, r *http.Request) {
	req := models.CreateShortURL{}
	err := web.UnmarshalJSON(r, &req)
	if err != nil {
		slog.Error("handleCreateShortURL", slog.Any("error", err))
		web.SendJSONError(w, err)
		return
	}

	shortURL, err := s.app.Create(r.Context(), req)
	if err != nil {
		slog.Error("handleCreateShortURL", slog.Any("error", err))
		web.SendJSONError(w, err)
		return
	}

	web.SendJSONPayload(w, shortURL)
}

func (s server) handleDeleteShortURL(w http.ResponseWriter, r *http.Request) {
	req := models.DeleteShortURL{}
	err := web.UnmarshalJSON(r, &req)
	if err != nil {
		slog.Error("handleDeleteShortURL", slog.Any("error", err))
		web.SendJSONError(w, err)
		return
	}

	err = s.app.Delete(r.Context(), req)
	if err != nil {
		slog.Error("handleDeleteShortURL", slog.Any("error", err))
		web.SendJSONError(w, err)
		return
	}

	web.SendJSONPayload(w, nil)
}

func (s server) handleUpdateShortURL(w http.ResponseWriter, r *http.Request) {

	req := models.UpdateShortURL{}
	err := web.UnmarshalJSON(r, &req)
	if err != nil {
		slog.Error("handleUpdateShortURL", slog.Any("error", err))
		web.SendJSONError(w, err)
		return
	}

	err = s.app.Update(r.Context(), req)
	if err != nil {
		slog.Error("handleUpdateShortURL", slog.Any("error", err))
		web.SendJSONError(w, err)
		return
	}

	web.SendJSONPayload(w, nil)

}

func (s server) handleListShortURL(w http.ResponseWriter, r *http.Request) {
	list, err := s.app.ListForOwner(r.Context())
	if err != nil {
		slog.Error("handleListShortURL", slog.Any("error", err))
		web.SendJSONError(w, err)
		return
	}

	web.SendJSONPayload(w, list)

}
