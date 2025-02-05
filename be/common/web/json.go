package web

import (
	"encoding/json"
	"net/http"

	"github.com/ggrrrr/urlshortener-svc/be/common/application"
)

func UnmarshalJSON(r *http.Request, v any) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return application.NewBadRequest("unable to parse request json", err)
	}
	defer r.Body.Close()

	return nil
}
