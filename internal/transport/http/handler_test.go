package http

import (
	"github.com/gorilla/mux"
	"github.com/kiraplenkin/go-musthave-devops/internal/stats"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
			name: "test",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: stats.NewService(storage.New()),
			},
			endpoint: "/api/stat/1",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				text:        "This Id doesn''t exists\n",
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
			endpoint: "/api/health/",
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
