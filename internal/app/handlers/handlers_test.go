package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hrapovd1/url-shortener/internal/app/storage"
)

func TestRootPost(t *testing.T) {
	type request struct {
		uri         string
		bodyURL     string
		contentType string
	}
	type want struct {
		status      int
		contentType string
	}
	stor := storage.NewMemStorage()
	bodyResponseLen := 31
	tests := []struct {
		name    string
		req     request
		want    want
		errored bool
	}{
		{"success", request{"/", "ya.ru", "text/plain"}, want{http.StatusCreated, "text/plain"}, false},
		{"wrong uri request", request{"/any", "ya.ru", "text/plain"}, want{http.StatusBadRequest, "text/plain; charset=utf-8"}, true},
		{"empty body url", request{"/", "", "text/plain"}, want{http.StatusCreated, "text/plain"}, false},
		{"wrong content type", request{"/", "ya.ru", ""}, want{http.StatusBadRequest, "text/plain; charset=utf-8"}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			postBody := bytes.NewBufferString(test.req.bodyURL)
			req := httptest.NewRequest(http.MethodPost, test.req.uri, postBody)
			req.Header.Set("content-type", test.req.contentType)
			rw := httptest.NewRecorder()

			rootPostHandler := rootPost(stor)
			rootPostHandler(rw, req)

			response := rw.Result()
			defer response.Body.Close()
			body, _ := io.ReadAll(response.Body)

			assert.Equal(t, test.want.status, response.StatusCode)
			assert.Equal(t, test.want.contentType, response.Header.Get("Content-Type"))
			if !test.errored {
				assert.Len(t, string(body), bodyResponseLen)
			}
		})
	}
}

func TestRootGet(t *testing.T) {
	type request struct {
		uri         string
		contentType string
	}
	type want struct {
		status   int
		location string
	}
	stor := storage.NewMemStorage()
	map[string]string(*stor)["yaShort"] = "ya.ru"
	map[string]string(*stor)["emptyShort"] = "empty"

	tests := []struct {
		name    string
		req     request
		want    want
		errored bool
	}{
		{"success", request{"/yaShort", "text/plain"}, want{http.StatusTemporaryRedirect, "ya.ru"}, false},
		{"wrong request uri", request{"/", "text/plain"}, want{http.StatusBadRequest, ""}, false},
		{"not existed request uri", request{"/NotExist", "text/plain"}, want{http.StatusBadRequest, ""}, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.req.uri, nil)
			rw := httptest.NewRecorder()

			rootGetHandler := rootGet(stor)
			rootGetHandler(rw, req)

			response := rw.Result()
			defer response.Body.Close()

			assert.Equal(t, test.want.status, response.StatusCode)
			assert.Equal(t, test.want.location, response.Header.Get("Location"))
		})
	}
}

func TestRootHandler(t *testing.T) {
	type request struct {
		method      string
		uri         string
		bodyURL     string
		contentType string
	}
	type want struct {
		status      int
		contentType string
		location    string
	}
	stor := storage.NewMemStorage()
	map[string]string(*stor)["yaShort"] = "ya.ru"
	bodyResponseLen := 31
	tests := []struct {
		name    string
		req     request
		want    want
		errored bool
	}{
		{"success Post", request{http.MethodPost, "/", "yandex.ru", "text/plain"}, want{http.StatusCreated, "text/plain", ""}, false},
		{"success Get", request{http.MethodGet, "/yaShort", "", "text/plain"}, want{http.StatusTemporaryRedirect, "", "ya.ru"}, false},
		{"wrong Method", request{http.MethodPut, "/yaShort", "", "text/plain"}, want{http.StatusMethodNotAllowed, "", "ya.ru"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			postBody := bytes.NewBufferString(test.req.bodyURL)
			req := httptest.NewRequest(test.req.method, test.req.uri, postBody)
			req.Header.Set("content-type", test.req.contentType)
			rw := httptest.NewRecorder()

			rootHandler := RootHandler(stor)
			rootHandler(rw, req)

			if test.req.method == http.MethodPost {
				response := rw.Result()
				defer response.Body.Close()
				body, _ := io.ReadAll(response.Body)

				assert.Equal(t, test.want.status, response.StatusCode)
				assert.Equal(t, test.want.contentType, response.Header.Get("Content-Type"))
				if !test.errored {
					assert.Len(t, string(body), bodyResponseLen)
				}
			} else if test.req.method == http.MethodGet {
				response := rw.Result()
				defer response.Body.Close()

				assert.Equal(t, test.want.status, response.StatusCode)
				assert.Equal(t, test.want.location, response.Header.Get("Location"))
			} else {
				response := rw.Result()
				defer response.Body.Close()

				assert.Equal(t, test.want.status, response.StatusCode)
			}
		})
	}
}
