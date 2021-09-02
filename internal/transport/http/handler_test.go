package http

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/kiraplenkin/go-musthave-devops/internal/stats"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetAllStats(t *testing.T) {
	type fields struct {
		Router  *mux.Router
		Service *stats.Service
	}
	type want struct {
		code        int
		contentType string
		text        string
	}
	tests := []struct {
		name     string
		fields   fields
		endpoint string
		want     want
	}{
		{
			name: "OK",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: stats.NewService(storage.New()),
			},
			endpoint: "/",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				text:        "{Storage:map[]}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &Handler{
				Router:  tt.fields.Router,
				Service: tt.fields.Service,
			}
			req := httptest.NewRequest(http.MethodGet, tt.endpoint, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.GetAllStats)
			h.ServeHTTP(w, req)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf(err.Error())
			}
			assert.Equal(t, tt.want.text, string(b))
			err = res.Body.Close()
			if err != nil {
				// TODO return error
				return
			}
		})
	}
}

func TestGetStats(t *testing.T) {
	type fields struct {
		Router  *mux.Router
		Service *stats.Service
	}
	type want struct {
		code        int
		contentType string
		text        string
	}
	tests := []struct {
		name     string
		fields   fields
		endpoint string
		want     want
	}{
		{
			name: "Not existed ID",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: stats.NewService(storage.New()),
			},
			endpoint: "/?id=1",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				text:        "Can't get stat by this ID\n",
			},
		},
		{
			name: "Bad id",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: stats.NewService(storage.New()),
			},
			endpoint: "/?id=test",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				text:        "Unable to parse uint from id\n",
			},
		},
		{
			name: "Not id param in request",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: stats.NewService(storage.New()),
			},
			endpoint: "/?i=test",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				text:        "URL param id is missing\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &Handler{
				Router:  tt.fields.Router,
				Service: tt.fields.Service,
			}
			req := httptest.NewRequest(http.MethodGet, tt.endpoint, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.GetStats)
			h.ServeHTTP(w, req)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf(err.Error())
			}
			assert.Equal(t, tt.want.text, string(b))
			err = res.Body.Close()
			if err != nil {
				// TODO return error
				return
			}
		})
	}
}

func TestPostStat(t *testing.T) {
	type fields struct {
		Router  *mux.Router
		Service *stats.Service
	}
	type want struct {
		code int
	}

	positiveData := url.Values{}
	positiveData.Set("id", "1")
	positiveData.Set("type", "test")
	positiveData.Set("value", "1")

	emptyData := url.Values{}

	negativeData := url.Values{}
	negativeData.Set("id", "")
	negativeData.Set("type", "")
	negativeData.Set("value", "")

	tests := []struct {
		name     string
		fields   fields
		endpoint string
		data url.Values
		want     want
	}{
		{
			name: "Positive data",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: stats.NewService(storage.New()),
			},
			endpoint: "/",
			data: positiveData,
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name: "Empty post data",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: stats.NewService(storage.New()),
			},
			endpoint: "/",
			data: emptyData,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Bad data",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: stats.NewService(storage.New()),
			},
			endpoint: "/",
			data: negativeData,
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &Handler{
				Router:  tt.fields.Router,
				Service: tt.fields.Service,
			}
			req := httptest.NewRequest(http.MethodPost, tt.endpoint, bytes.NewBufferString(tt.data.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.PostStat)
			h.ServeHTTP(w, req)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func TestCheckHealth(t *testing.T) {
	type fields struct {
		Router  *mux.Router
		Service *stats.Service
	}
	type want struct {
		code int
		text string
	}
	tests := []struct {
		name     string
		fields   fields
		endpoint string
		want     want
	}{
		{
			name: "Check health",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: stats.NewService(storage.New()),
			},
			endpoint: "/",
			want: want{
				code: http.StatusOK,
				text: "alive!",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := Handler{
				Router:  tt.fields.Router,
				Service: tt.fields.Service,
			}
			req := httptest.NewRequest(http.MethodGet, tt.endpoint, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.CheckHealth)
			h.ServeHTTP(w, req)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf(err.Error())
			}
			assert.Equal(t, tt.want.text, string(b))
			err = res.Body.Close()
			if err != nil {
				// TODO return error
				return
			}
		})
	}
}
