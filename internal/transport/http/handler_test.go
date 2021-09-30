package http

//import (
//	"bytes"
//	"fmt"
//	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
//	"github.com/kiraplenkin/go-musthave-devops/internal/types"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"io/ioutil"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//var (
//	contentType = "text/plain; charset=utf-8"
//	endpoint    = "/"
//	testCfg     = types.ServerConfig{
//		ServerAddress:   "localhost:8080",
//		FileStoragePath: "test_file",
//	}
//	testStore, _ = storage.NewStorage(&testCfg)
//	handler      = NewHandler(*testStore)
//)
//
//func TestGetAllStats(t *testing.T) {
//	type want struct {
//		code        int
//		contentType string
//		text        string
//	}
//	tests := []struct {
//		name     string
//		endpoint string
//		want     want
//	}{
//		{
//			name:     "Positive test",
//			endpoint: endpoint,
//			want: want{
//				code:        http.StatusOK,
//				contentType: contentType,
//				text:        "map[]",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			handler.SetupRouters()
//
//			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
//			require.NoError(t, err)
//			rec := httptest.NewRecorder()
//
//			handler.Router.HandleFunc(endpoint, handler.GetAllStats)
//			handler.Router.ServeHTTP(rec, req)
//
//			res := rec.Result()
//			assert.Equal(t, tt.want.code, res.StatusCode)
//			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
//			b, err := ioutil.ReadAll(res.Body)
//			require.NoError(t, err)
//			assert.Equal(t, tt.want.text, string(b))
//			err = res.Body.Close()
//			require.NoError(t, err)
//		})
//	}
//}
//
//func TestGetStats(t *testing.T) {
//	type want struct {
//		code        int
//		contentType string
//		text        string
//	}
//	tests := []struct {
//		name string
//		id   string
//		want want
//	}{
//		{
//			name: "Not existed ID",
//			id:   "1",
//			want: want{
//				code:        http.StatusBadRequest,
//				contentType: contentType,
//				text:        "Can't get stat by this ID\n",
//			},
//		},
//		{
//			name: "Bad id",
//			id:   "test_id",
//			want: want{
//				code:        http.StatusBadRequest,
//				contentType: contentType,
//				text:        "Unable to parse uint from id\n",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			handler.SetupRouters()
//			path := fmt.Sprintf("%s%s", endpoint, tt.id)
//
//			req, err := http.NewRequest(http.MethodGet, path, nil)
//			require.NoError(t, err)
//			rec := httptest.NewRecorder()
//
//			handler.Router.HandleFunc("/{id}", handler.GetStatsByID)
//			handler.Router.ServeHTTP(rec, req)
//
//			res := rec.Result()
//			assert.Equal(t, tt.want.code, res.StatusCode)
//			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
//			b, err := ioutil.ReadAll(res.Body)
//			require.NoError(t, err)
//			assert.Equal(t, tt.want.text, string(b))
//			err = res.Body.Close()
//			require.NoError(t, err)
//		})
//	}
//}
//
//func TestPostStat(t *testing.T) {
//	type want struct {
//		code int
//	}
//
//	tests := []struct {
//		name string
//		data []byte
//		want want
//	}{
//		{
//			name: "Positive data",
//			data: []byte(`{"id":1,"totalAlloc":213880,"sys":8735760,"mallocs":760,"frees":11,"liveObjects":749,"numGoroutine":1}`),
//			want: want{
//				code: http.StatusCreated,
//			},
//		},
//		{
//			name: "Empty post data",
//			data: []byte("{}"),
//			want: want{
//				code: http.StatusBadRequest,
//			},
//		},
//		{
//			name: "Empty Id",
//			data: []byte(`{"totalAlloc":213880,"sys":8735760,"mallocs":760,"frees":11,"liveObjects":749,"numGoroutine":1}`),
//			want: want{
//				code: http.StatusBadRequest,
//			},
//		},
//		{
//			name: "Bad value",
//			data: []byte(`{"id":1,"totalAlloc":"213880","sys":8735760,"mallocs":760,"frees":11,"liveObjects":749,"numGoroutine":1}`),
//			want: want{
//				code: http.StatusBadRequest,
//			},
//		},
//		{
//			name: "Empty value",
//			data: []byte(`{"id":1,"sys":8735760,"mallocs":760,"frees":11,"liveObjects":749,"numGoroutine":1}`),
//			want: want{
//				code: http.StatusBadRequest,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			handler.SetupRouters()
//			req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(string(tt.data)))
//			req.Header.Set("Content-Type", "application/json")
//			w := httptest.NewRecorder()
//			h := http.HandlerFunc(handler.PostStat)
//			h.ServeHTTP(w, req)
//			res := w.Result()
//			assert.Equal(t, tt.want.code, res.StatusCode)
//			err := res.Body.Close()
//			require.NoError(t, err)
//		})
//	}
//}
//
//func TestCheckHealth(t *testing.T) {
//	type want struct {
//		code        int
//		contentType string
//		text        string
//	}
//	tests := []struct {
//		name string
//		want want
//	}{
//		{
//			name: "Check health",
//			want: want{
//				code:        http.StatusOK,
//				contentType: contentType,
//				text:        "alive!",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			handler.SetupRouters()
//
//			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
//			require.NoError(t, err)
//			rec := httptest.NewRecorder()
//
//			handler.Router.HandleFunc(endpoint, handler.CheckHealth)
//			handler.Router.ServeHTTP(rec, req)
//
//			res := rec.Result()
//			assert.Equal(t, tt.want.code, res.StatusCode)
//			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
//			b, err := ioutil.ReadAll(res.Body)
//			require.NoError(t, err)
//			assert.Equal(t, tt.want.text, string(b))
//			err = res.Body.Close()
//			require.NoError(t, err)
//		})
//	}
//}
