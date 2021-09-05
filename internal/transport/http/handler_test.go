package http

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetAllStats(t *testing.T) {
	type fields struct {
		Router  *mux.Router
		Service *storage.Store
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
				Service: storage.NewStorage(),
			},
			endpoint: "/",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				text:        "{map[]}",
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
			require.NoError(t, err)
			assert.Equal(t, tt.want.text, string(b))
			err = res.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestGetStats(t *testing.T) {
	type fields struct {
		Router  *mux.Router
		Service *storage.Store
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
		//{
		//	name: "Not existed ID",
		//	fields: fields{
		//		Router:  mux.NewRouter(),
		//		Service: stats.NewService(storage.New()),
		//	},
		//	endpoint: "/1",
		//	want: want{
		//		code:        http.StatusBadRequest,
		//		contentType: "text/plain; charset=utf-8",
		//		text:        "Can't get stat by this ID\n",
		//	},
		//},
		{
			name: "Bad id",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: storage.NewStorage(),
			},
			endpoint: "/test",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				text:        "Unable to parse uint from id\n",
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
			require.NoError(t, err)
			assert.Equal(t, tt.want.text, string(b))
			err = res.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestPostStat(t *testing.T) {
	type fields struct {
		Router  *mux.Router
		Service *storage.Store
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
		data     url.Values
		want     want
	}{
		{
			name: "Positive data",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: storage.NewStorage(),
			},
			endpoint: "/",
			data:     positiveData,
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name: "Empty post data",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: storage.NewStorage(),
			},
			endpoint: "/",
			data:     emptyData,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Bad data",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: storage.NewStorage(),
			},
			endpoint: "/",
			data:     negativeData,
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
			err := res.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestCheckHealth(t *testing.T) {
	type fields struct {
		Router  *mux.Router
		Service *storage.Store
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
				Service: storage.NewStorage(),
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
			require.NoError(t, err)
			assert.Equal(t, tt.want.text, string(b))
			err = res.Body.Close()
			require.NoError(t, err)
		})
	}
}
