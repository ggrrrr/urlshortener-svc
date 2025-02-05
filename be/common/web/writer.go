package web

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ggrrrr/urlshortener-svc/be/common/application"
)

type errorResponse struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
}

func Redirect(w http.ResponseWriter, r *http.Request, toURL string) {
	http.Redirect(w, r, toURL, http.StatusMovedPermanently)
}

func SendJSONError(w http.ResponseWriter, err error) {
	appError, ok := err.(*application.AppError)
	if ok {
		sendJSONPayload(w, appError.Code(), errorResponse{
			Code:  appError.Code(),
			Error: appError.Message(),
		})
		return
	}

	sendJSONPayload(w, http.StatusInternalServerError, errorResponse{
		Code:  http.StatusInternalServerError,
		Error: err.Error(),
	})
}

func SendJSONPayload(w http.ResponseWriter, payload any) {
	if payload == nil {
		sendEmpty(w, http.StatusNoContent)
	}
	sendJSONPayload(w, http.StatusOK, payload)
}

func sendJSONPayload(w http.ResponseWriter, code int, payload any) {
	w.Header().Add("Content-Type", "application/json")
	b, err := json.Marshal(payload)
	if err != nil {
		slog.Error("unable to write response",
			slog.Any("http_code", code),
			slog.Any("payload", payload),
			slog.Any("error", err),
		)
	}
	w.WriteHeader(code)
	_, err = w.Write(b)
	if err != nil {
		slog.Error("unable to write response",
			slog.Any("error", err),
		)
	}
}

func sendEmpty(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
