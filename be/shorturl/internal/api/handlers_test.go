package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/urlshortener-svc/be/common/application"
	"github.com/ggrrrr/urlshortener-svc/be/shorturl/models"
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
				mockApp.On("GetLongURL", "shortKey").Return(srvRedirect.URL, nil).Once()

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
				mockApp.On("GetLongURL", "notfound").Return("", application.NewNotFound()).Once()

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
				mockApp.On("GetLongURL", "syserr").Return("", application.NewSystemError("some error", fmt.Errorf("some error"))).Once()

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
				mockApp.On("Create", models.CreateShortURL{LongURL: "long_url"}).Return(&models.Key{Key: "key1"}, nil).Once()

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				testResp(t, resp, 200, `{"key":"key1"}`)
			},
		},
		{
			name: "err create 500",
			prepFunc: func(t *testing.T) {
				req, err := http.NewRequest("POST", srv.URL+"/admin/v1", strings.NewReader(`{"long_url":"long_error"}`))
				require.NoError(t, err)
				mockApp.On("Create", models.CreateShortURL{LongURL: "long_error"}).Return(
					nil,
					application.NewSystemError("sys", fmt.Errorf("error")),
				).Once()

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				testResp(t, resp, 500, `{"code":500,"error":"sys"}`)
			},
		},
		{
			name: "err create 400",
			prepFunc: func(t *testing.T) {
				req, err := http.NewRequest("POST", srv.URL+"/admin/v1", strings.NewReader(`{long_url":"long_error"}`))
				require.NoError(t, err)

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				testResp(t, resp, 400, `{"code":400,"error":"unable to parse request json"}`)
			},
		},
		{
			name: "ok delete",
			prepFunc: func(t *testing.T) {
				req, err := http.NewRequest("DELETE", srv.URL+"/admin/v1", strings.NewReader(`{"key":"key_ok"}`))
				require.NoError(t, err)
				mockApp.On("Delete", "key_ok").Return(nil).Once()

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				testResp(t, resp, 204, ``)
			},
		},
		{
			name: "err delete 500",
			prepFunc: func(t *testing.T) {
				req, err := http.NewRequest("DELETE", srv.URL+"/admin/v1", strings.NewReader(`{"key":"key_err"}`))
				require.NoError(t, err)
				mockApp.On("Delete", "key_err").Return(application.NewSystemError("msg", fmt.Errorf("err"))).Once()

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				testResp(t, resp, 500, `{"code":500,"error":"msg"}`)
			},
		},
		{
			name: "ok update",
			prepFunc: func(t *testing.T) {
				req, err := http.NewRequest("PUT", srv.URL+"/admin/v1", strings.NewReader(`{"key":"key_ok"}`))
				require.NoError(t, err)
				mockApp.On("Update", "key_ok").Return(nil).Once()

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				testResp(t, resp, 204, ``)
			},
		},
		{
			name: "list ok",
			prepFunc: func(t *testing.T) {
				req, err := http.NewRequest("GET", srv.URL+"/admin/v1", nil)
				require.NoError(t, err)
				mockApp.On("ListForOwner").Return([]*models.ShortURLRecord{}, nil).Once()

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				testResp(t, resp, 200, `[]`)
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
