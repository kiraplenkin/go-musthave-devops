package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/kiraplenkin/go-musthave-devops/internal/validator"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Handler stores pointers to service
type Handler struct {
	Router  *mux.Router
	Storage *storage.Store
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

	h.Router.HandleFunc("/stats/{id}", h.GetStatsByID).Methods(http.MethodGet)
	h.Router.HandleFunc("/stats/", h.GetAllStats).Methods(http.MethodGet)
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
	var requestStats types.RequestStats

	err = json.Unmarshal(body, &requestStats)
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
