package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
			handler.SetupRouters()
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
		id       string
		want     want
	}{
		{
			name: "Not existed ID",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: storage.NewStorage(),
			},
			endpoint: "/",
			id:       "1",
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
				Service: storage.NewStorage(),
			},
			endpoint: "/",
			id:       "test_id",
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
			path := fmt.Sprintf("/%s", tt.id)
			req, err := http.NewRequest(http.MethodGet, path, nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()
			handler.SetupRouters()
			handler.Router.HandleFunc("/{id}", handler.GetStats)
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
	type fields struct {
		Router  *mux.Router
		Service *storage.Store
	}
	type want struct {
		code int
	}

	positiveRequest := types.RequestStats{
		ID:    "1",
		Type:  "test_type",
		Value: "test_value",
	}

	emptyIdRequest := types.RequestStats{
		ID:    "",
		Type:  "test_type",
		Value: "test_value",
	}

	badIdRequest := types.RequestStats{
		ID:    "test_ID",
		Type:  "test_type",
		Value: "test_value",
	}

	emptyTypeRequest := types.RequestStats{
		ID:    "1",
		Type:  "",
		Value: "test_value",
	}

	emptyValueRequest := types.RequestStats{
		ID:    "1",
		Type:  "test_type",
		Value: "",
	}

	tests := []struct {
		name     string
		fields   fields
		endpoint string
		data     types.RequestStats
		want     want
	}{
		{
			name: "Positive data",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: storage.NewStorage(),
			},
			endpoint: "/",
			data:     positiveRequest,
			want: want{
				code: http.StatusCreated,
			},
		},
		{
			name: "Bad ID",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: storage.NewStorage(),
			},
			endpoint: "/",
			data:     badIdRequest,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Empty ID",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: storage.NewStorage(),
			},
			endpoint: "/",
			data:     emptyIdRequest,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Empty Type",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: storage.NewStorage(),
			},
			endpoint: "/",
			data:     emptyTypeRequest,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Empty Value",
			fields: fields{
				Router:  mux.NewRouter(),
				Service: storage.NewStorage(),
			},
			endpoint: "/",
			data:     emptyValueRequest,
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
			r, err := json.Marshal(tt.data)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, tt.endpoint, bytes.NewBufferString(string(r)))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.PostStat)
			h.ServeHTTP(w, req)
			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			err = res.Body.Close()
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
