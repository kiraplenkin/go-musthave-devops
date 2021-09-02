package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kiraplenkin/go-musthave-devops/internal/stats"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"net/http"
	"strconv"
)

// Handler - stores pointers to our service
type Handler struct {
	Router  *mux.Router
	Service *stats.Service
}

// NewHandler - returns a pointer to HAndler
func NewHandler(s *stats.Service) *Handler {
	return &Handler{
		Service: s,
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

//GetStats - ...
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Unable to parse uint from id", http.StatusBadRequest)
		return
	}

	stat, err := h.Service.GetStats(uint(i))
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

//GetAllStats - ...
func (h Handler) GetAllStats(w http.ResponseWriter, r *http.Request) {
	allStats, err := h.Service.GetAllStats()
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Can't get all stats", http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(w, "%+v", allStats)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

//PostStat - ...
func (h Handler) PostStat(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		_, err := fmt.Fprintf(w, "Can't parse form")
		if err != nil {
			return
		}
	}
	id := r.Form.Get("id")
	if id == "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Id is empty", http.StatusBadRequest)
		return
	}
	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Unable to parse uint from id", http.StatusBadRequest)
		return
	}

	statsType := r.Form.Get("type")
	if statsType == "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Stats type is empty", http.StatusBadRequest)
		return
	}

	statsValue := r.Form.Get("value")
	if statsValue == "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Stats type is empty", http.StatusBadRequest)
		return
	}

	newStat := storage.Stats{StatsType: statsType, StatsValue: statsValue}
	err = h.Service.PostStats(uint(i), newStat)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Can't save stat", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// CheckHealth - endpoint to check health
func (h Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "alive!")
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}
