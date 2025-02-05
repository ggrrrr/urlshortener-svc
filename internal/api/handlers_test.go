package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/urlshortener-svc/common/application"
	"github.com/ggrrrr/urlshortener-svc/models"
)

func TestRouter(t *testing.T) {
	mockApp := new(MockApp)

	srv := httptest.NewServer(CreateRouter(mockApp))
	defer srv.Close()

	tests := []struct {
		name     string
		prepFunc func(t *testing.T)
	}{
		{
			name: "ok forward",
			prepFunc: func(t *testing.T) {
				srvRedirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
					_, _ = w.Write([]byte("ok"))
				}))
				defer srvRedirect.Close()

				req, err := http.NewRequest("GET", srv.URL+"/shortKey", nil)
				require.NoError(t, err)
				mockApp.On("GetLongURL", "shortKey").Return(srvRedirect.URL, nil)

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)

				testResp(t, resp, 200, "ok")

			},
		},
		{
			name: "404 forward",
			prepFunc: func(t *testing.T) {
				req, err := http.NewRequest("GET", srv.URL+"/notfound", nil)
				require.NoError(t, err)
				mockApp.On("GetLongURL", "notfound").Return("", application.NewNotFound())

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				assert.Equal(t, 404, resp.StatusCode)

				testResp(t, resp, 404, `{"code":404}`)
			},
		},
		{
			name: "500 forward",
			prepFunc: func(t *testing.T) {
				req, err := http.NewRequest("GET", srv.URL+"/syserr", nil)
				require.NoError(t, err)
				mockApp.On("GetLongURL", "syserr").Return("", application.NewSystemError("some error", fmt.Errorf("some error")))

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				assert.Equal(t, 500, resp.StatusCode)
				testResp(t, resp, 500, `{"code":500,"error":"some error"}`)
			},
		},
		{
			name: "ok create",
			prepFunc: func(t *testing.T) {
				req, err := http.NewRequest("POST", srv.URL+"/admin/v1", strings.NewReader(`{"long_url":"long_url"}`))
				require.NoError(t, err)
				mockApp.On("Create", "long_url").Return(
					models.ShortURLRecord{
						Key:       "new_key",
						LongURL:   "long_url",
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					},
					nil,
				)

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				testResp(t, resp, 200, `{"key":"new_key","long_url":"long_url","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`)
			},
		},
		{
			name: "err create",
			prepFunc: func(t *testing.T) {
				req, err := http.NewRequest("POST", srv.URL+"/admin/v1", strings.NewReader(`{"long_url":"long_error"}`))
				require.NoError(t, err)
				mockApp.On("Create", "long_error").Return(
					models.ShortURLRecord{},
					application.NewSystemError("sys", fmt.Errorf("error")),
				)

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				testResp(t, resp, 500, `{"code":500,"error":"sys"}`)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, tc.prepFunc)

	}
}

func testResp(t *testing.T, resp *http.Response, code int, body string) {
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, code, resp.StatusCode)
	assert.Equal(t, body, string(b))
}
