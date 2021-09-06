package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/kiraplenkin/go-musthave-devops/internal/validator"
	"net/http"
	"strconv"
)

// Handler - stores pointers to service
type Handler struct {
	Router  *mux.Router
	Service *storage.Store
}

// NewHandler - returns a pointer to Handler
func NewHandler(s storage.Store) *Handler {
	return &Handler{
		Service: &s,
	}
}

// SetupRouters - sets up all the routes for App
func (h *Handler) SetupRouters() {
	fmt.Println("Setting Up Routers")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/stat/{id}", h.GetStats).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/stat/", h.GetAllStats).Methods(http.MethodGet)
	h.Router.HandleFunc("/api/stat/", h.PostStat).Methods(http.MethodPost)

	h.Router.HandleFunc("/api/health/", h.CheckHealth).Methods(http.MethodGet)
}

//GetStats - handler that return types.Stats by ID
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// TODO if id null
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Unable to parse uint from id", http.StatusBadRequest)
		return
	}

	stat, err := h.Service.GetStatsByID(uint(i))
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Can't get stat by this ID", http.StatusBadRequest)
		return
	}

	_, err = fmt.Fprintf(w, "%+v", stat)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

//GetAllStats - handler that return all values from storage.Store
func (h Handler) GetAllStats(w http.ResponseWriter, _ *http.Request) {
	allStats, err := h.Service.GetAllStats()
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Can't get all stats", http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(w, "%v", *allStats)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

//PostStat - handler that save types.Stats to storage.Store
func (h Handler) PostStat(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestStats types.RequestStats

	err := decoder.Decode(&requestStats)
	if err != nil {
		http.Error(w, "Can't decode input json", http.StatusInternalServerError)
	}

	id := requestStats.ID
	statsType := requestStats.Type
	statsValue := requestStats.Value
	err = validator.Require(id, statsType, statsValue)
	if err != nil {
		http.Error(w, error.Error(err), http.StatusBadRequest)
	}

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Unable to parse uint from id", http.StatusBadRequest)
		return
	}

	newStat := types.Stats{StatsType: statsType, StatsValue: statsValue}
	err = h.Service.SaveStats(uint(i), newStat)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Can't save stat", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprintf(w, "%+v", newStat)
	if err != nil {
		return
	}
}

// CheckHealth - handler to check health
func (h Handler) CheckHealth(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "alive!")
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}
