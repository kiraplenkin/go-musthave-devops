package http

import (
	"bytes"
	"fmt"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	contentType = "text/plain; charset=utf-8"
	endpoint    = "/"
	handler     = NewHandler(*storage.NewStorage())
)

func TestGetAllStats(t *testing.T) {
	type want struct {
		code        int
		contentType string
		text        string
	}
	tests := []struct {
		name     string
		endpoint string
		want     want
	}{
		{
			name:     "Positive test",
			endpoint: endpoint,
			want: want{
				code:        http.StatusOK,
				contentType: contentType,
				text:        "{map[]}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.SetupRouters()
			//path := fmt.Sprintf("%s", endpoint)
			//path := endpoint

			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()

			handler.Router.HandleFunc(endpoint, handler.GetAllStats)
			handler.Router.ServeHTTP(rec, req)

			res := rec.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			b, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.text, string(b))
			err = res.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestGetStats(t *testing.T) {
	type want struct {
		code        int
		contentType string
		text        string
	}
	tests := []struct {
		name string
		id   string
		want want
	}{
		{
			name: "Not existed ID",
			id:   "1",
			want: want{
				code:        http.StatusBadRequest,
				contentType: contentType,
				text:        "Can't get stat by this ID\n",
			},
		},
		{
			name: "Bad id",
			id:   "test_id",
			want: want{
				code:        http.StatusBadRequest,
				contentType: contentType,
				text:        "Unable to parse uint from id\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.SetupRouters()
			path := fmt.Sprintf("%s%s", endpoint, tt.id)

			req, err := http.NewRequest(http.MethodGet, path, nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()

			handler.Router.HandleFunc("/{id}", handler.GetStatsByID)
			handler.Router.ServeHTTP(rec, req)

			res := rec.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			b, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.text, string(b))
			err = res.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestPostStat(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		name string
		data url.Values
		want want
	}{
		{
			name: "Positive data",
			data: func() url.Values {
				v := url.Values{}
				v.Set("id", "1")
				v.Set("Alloc", "100")
				v.Set("TotalAlloc", "101")
				v.Set("Sys", "102")
				v.Set("Mallocs", "103")
				v.Set("Frees", "104")
				v.Set("LiveObjects", "105")
				v.Set("PauseTotalNs", "106")
				v.Set("NumGC", "107")
				v.Set("NumGoroutine", "108")
				return v
			}(),
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name: "Empty post data",
			data: func() url.Values {
				v := url.Values{}
				return v
			}(),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Empty Id",
			data: func() url.Values {
				v := url.Values{}
				v.Set("id", "")
				v.Set("Alloc", "100")
				v.Set("TotalAlloc", "100")
				v.Set("Sys", "100")
				v.Set("Mallocs", "100")
				v.Set("Frees", "100")
				v.Set("LiveObjects", "100")
				v.Set("PauseTotalNs", " 100")
				v.Set("NumGC", "100")
				v.Set("NumGoroutine", "100")
				return v
			}(),
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.SetupRouters()
			req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(tt.data.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.PostStat)
			h.ServeHTTP(w, req)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			err := res.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestCheckHealth(t *testing.T) {
	type want struct {
		code        int
		contentType string
		text        string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Check health",
			want: want{
				code:        http.StatusOK,
				contentType: contentType,
				text:        "alive!",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.SetupRouters()
			//path := fmt.Sprintf("%s", endpoint)

			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()

			handler.Router.HandleFunc(endpoint, handler.CheckHealth)
			handler.Router.ServeHTTP(rec, req)

			res := rec.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			b, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.text, string(b))
			err = res.Body.Close()
			require.NoError(t, err)
		})
	}
}
