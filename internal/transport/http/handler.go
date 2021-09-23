package http

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kiraplenkin/go-musthave-devops/internal/crypto"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/kiraplenkin/go-musthave-devops/internal/validator"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Handler stores pointers to service
type Handler struct {
	Router  *mux.Router
	Storage *storage.Store
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write ...
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipHandle handle which compress all handlers
func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			_, err := io.WriteString(w, err.Error())
			if err != nil {
				return
			}
			return
		}
		defer func(gz *gzip.Writer) {
			err := gz.Close()
			if err != nil {
				return
			}
		}(gz)

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

// NewHandler returns a pointer to Handler
func NewHandler(s storage.Store) *Handler {
	return &Handler{
		Storage: &s,
	}
}

// SetupRouters sets up all the routes for server
func (h *Handler) SetupRouters() {
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/{id}", h.GetStatsByID).Methods(http.MethodGet)
	h.Router.HandleFunc("/", h.GetAllStats).Methods(http.MethodGet)
	h.Router.HandleFunc("/", h.PostStat).Methods(http.MethodPost)
	h.Router.HandleFunc("/update/", h.PostStat).Methods(http.MethodPost)

	h.Router.HandleFunc("/health/", h.CheckHealth).Methods(http.MethodGet)
}

//GetStatsByID handler that return types.Stats by ID
func (h *Handler) GetStatsByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	i, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Unable to parse uint from id", http.StatusBadRequest)
		return
	}

	stat, err := h.Storage.GetStatsByID(i)
	if err != nil {
		http.Error(w, "Can't get stat by this ID", http.StatusBadRequest)
		return
	}

	_, err = fmt.Fprintf(w, "%+v", stat)
	if err != nil {
		http.Error(w, error.Error(err), http.StatusInternalServerError)
	}
}

//GetAllStats handler that return all values from storage.Store
func (h Handler) GetAllStats(w http.ResponseWriter, _ *http.Request) {
	allStats, err := h.Storage.GetAllStats()
	if err != nil {
		http.Error(w, "Can't get all stats", http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(w, "%v", *allStats)
	if err != nil {
		return
	}
}

//PostStat handler that save types.Stats to storage.Store
func (h Handler) PostStat(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	decodedBody, err := crypto.EncodeDecode(body, "decode")
	if err != nil {
		http.Error(w, "can't decode", http.StatusInternalServerError)
		return
	}

	decompressBody, err := crypto.Decompress(decodedBody)
	if err != nil {
		http.Error(w, "can't read body", http.StatusInternalServerError)
		return
	}

	var requestStats types.RequestStats
	err = json.Unmarshal(decompressBody, &requestStats)
	if err != nil {
		http.Error(w, "can't decode input json", http.StatusBadRequest)
		return
	}

	err = validator.Require(
		requestStats.ID,
		requestStats.TotalAlloc,
		requestStats.Sys,
		requestStats.Mallocs,
		requestStats.Frees,
		requestStats.LiveObjects,
		requestStats.NumGoroutine,
	)
	if err != nil {
		http.Error(w, "can't save stat", http.StatusBadRequest)
		return
	}

	newStat := types.Stats{
		TotalAlloc:   requestStats.TotalAlloc,
		Sys:          requestStats.Sys,
		Mallocs:      requestStats.Mallocs,
		Frees:        requestStats.Frees,
		LiveObjects:  requestStats.LiveObjects,
		NumGoroutine: requestStats.NumGoroutine,
	}

	err = h.Storage.SaveStats(requestStats.ID, newStat)
	if err != nil {
		http.Error(w, "can't save stat", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprintf(w, "%+v", newStat)
	if err != nil {
		return
	}
}

// CheckHealth handler to check health
func (h Handler) CheckHealth(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "alive!")
	if err != nil {
		return
	}
}
